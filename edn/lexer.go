// Copyright 2018-2019 Workiva Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package edn

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type PrimitiveType int

const (
	LiteralPrimitive PrimitiveType = iota
	IntegerPrimitive
	FloatPrimitive
	CharacterPrimitive
	SymbolPrimitive
	StringPrimitive
	lastPrimitivePriority
)

type PrimitiveProcessor func(tag string, tokenValue string) (Element, error)
type CollectionProcessor func(tag string, elements []Element) (el Element, e error)

type collProcDef struct {
	start     string
	end       string
	processor CollectionProcessor
}

type tokenType string

const (

	// because all blanks are skipped, all these token types are # of blanks.
	skipToken    tokenType = " "
	elementToken tokenType = "  "
)

func (tt tokenType) String() string {
	switch tt {
	case skipToken:
		return "[Skip Token]"
	case elementToken:
		return "[Element]"
	default:
		return string(tt)
	}
}

func (tt tokenType) Is(this string) bool {
	return string(tt) == this
}

// Lexer defines the lexical analyser for the
type Lexer interface {

	// AddPattern will take a pattern and attach the processor for that pattern.
	AddPattern(priority PrimitiveType, pattern string, processor PrimitiveProcessor)

	AddCollectionPattern(start string, end string, processor CollectionProcessor)

	Parse(data io.Reader) (Element, error)
}

func splitTag(data []byte, possible string) (tag string, value string) {

	// Special case, if the #{ appears then ignore the splitting and just return the value.
	if full := string(data); !strings.HasPrefix(full, SetStartLiteral) && strings.HasPrefix(full, TagPrefix) {
		parts := strings.Fields(full)
		tag = parts[0]
		value = strings.TrimPrefix(full, tag)
		tag = strings.TrimPrefix(tag, TagPrefix)
		value = strings.TrimSpace(value)

		if len(possible) > 0 && strings.HasSuffix(tag, possible) {
			tag = strings.TrimSuffix(tag, possible)
		}
	} else {
		value = full
	}

	return tag, value
}

func buildTagPattern(pattern string, mustHasSpace bool) []byte {

	subPattern := "*"
	if mustHasSpace {
		subPattern = "+"
	}

	return []byte(fmt.Sprintf("(%s[A-Za-z][-A-Za-z0-9_/.]*(\\s)%s)?%s", TagPrefix, subPattern, pattern))
}

func runScanner(scanner *lexmachine.Scanner) (tokType tokenType, elems []Element, err error) {
	var t interface{}

	var eos bool
	for t, err, eos = scanner.Next(); !eos && err == nil; t, err, eos = scanner.Next() {
		switch v := t.(type) {
		case Element:
			elems = append(elems, v)
			tokType = elementToken
		case tokenType:
			tokType = v
		}

		if tokType != elementToken && tokType != skipToken {
			break
		}
	}

	if err != nil {
		switch v := err.(type) {
		case *machines.UnconsumedInput:
			err = MakeError(ErrParserError, struct {
				message string
				elem    []Element
			}{
				v.Error(),
				elems,
			})
		}
	}

	return tokType, elems, err
}

///// ----------------------------------------------

type lexerImpl struct {
	primitivePatterns  map[PrimitiveType]map[string]PrimitiveProcessor
	collectionPatterns map[string]*collProcDef
	lex                *lexmachine.Lexer
	built              bool
	patternAdder       func([]byte, lexmachine.Action)
}

// newLexer will create a new lexer.
func newLexer() (lexer Lexer, err error) {
	lexer = &lexerImpl{
		primitivePatterns:  map[PrimitiveType]map[string]PrimitiveProcessor{},
		collectionPatterns: map[string]*collProcDef{},
		lex:                lexmachine.NewLexer(),
		built:              false,
	}

	return lexer, err
}

// completeStartup of the lexer
func (lexer *lexerImpl) completeStartup() error {

	if !lexer.built {

		for i := PrimitiveType(0); i < lastPrimitivePriority; i++ {
			if v, has := lexer.primitivePatterns[i]; has {

				for pattern, p := range v {
					processor := p // this is required as the processor needs to have a local reference... yay golang oddities! :(
					lexer.addPattern(buildTagPattern(pattern, true), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
						tag, value := splitTag(match.Bytes, "")
						return processor(tag, value)
					})
				}
			}
		}

		endPatterns := map[string]bool{}

		lexSpecialChars := []string{
			"\\", "[", "]", "{", "}", "(", ")",
		}

		for _, def := range lexer.collectionPatterns {
			processor := def.processor // and again boo - golang oddities! :(
			end := def.end
			start := def.start

			for _, c := range lexSpecialChars {
				start = strings.Replace(start, c, "\\"+c, -1)
				end = strings.Replace(end, c, "\\"+c, -1)
			}

			startRaw := def.start

			if _, has := endPatterns[end]; !has {
				lexer.addPattern([]byte(end), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
					return tokenType(end), nil
				})
				endPatterns[end] = true
			}

			// Add the non tagged items.
			lexer.addPattern(buildTagPattern(start, false), func(scan *lexmachine.Scanner, match *machines.Match) (v interface{}, e error) {
				tag, _ := splitTag(match.Bytes, startRaw)

				var tt tokenType
				var children []Element
				var c []Element

				for tt, c, e = runScanner(scan); ; tt, c, e = runScanner(scan) {
					stop := true
					if e == nil {
						children = append(children, c...)
						switch {
						case tt == elementToken:
							stop = false
						case tt.Is(end):
						default:
							e = MakeErrorWithFormat(ErrParserError, "Unexpected end token: '%s' instead of '%s'", tt.String(), end)
						}
					}

					if stop {
						break
					}
				}

				if e == nil {
					v, e = processor(tag, children)
				}

				return v, e
			})
		}

		lexer.addPattern([]byte("(\\s|,)+"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			return skipToken, nil
		})

		lexer.addPattern([]byte(";[^\\n]*(\\n)?"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			return skipToken, nil
		})

		compile := lexer.lex.CompileNFA
		if err := compile(); err != nil {
			return err
		}
		lexer.built = true
	}

	return nil
}

// Parse the value
func (lexer *lexerImpl) Parse(data io.Reader) (_ Element, err error) {

	if data == nil {
		return nil, MakeErrorWithFormat(ErrParserError, "parse input was nil")
	}

	var bytes []byte
	if bytes, err = ioutil.ReadAll(data); err != nil {
		return nil, MakeErrorWithFormat(ErrParserError, "parse input error: %s", err.Error())
	}

	if err = lexer.completeStartup(); err != nil {
		return nil, err
	}

	var scanner *lexmachine.Scanner
	if scanner, err = lexer.lex.Scanner(bytes); err != nil {
		return nil, err
	}

	var elems []Element
	if _, elems, err = runScanner(scanner); err != nil {
		return nil, err
	}

	switch {
	case len(elems) == 1:
		return elems[0], nil
	default:
		return nil, MakeErrorWithFormat(ErrParserError, "Expected one result, got: %d", len(elems))
	}
}

// AddPattern will add a pattern to the lexer
func (lexer *lexerImpl) AddPattern(priority PrimitiveType, pattern string, processor PrimitiveProcessor) {

	//  map[PrimitiveType]map[string]PrimitiveProcessor
	if _, has := lexer.primitivePatterns[priority]; !has {
		lexer.primitivePatterns[priority] = map[string]PrimitiveProcessor{}
	}

	if _, has := lexer.primitivePatterns[priority][pattern]; !has {
		lexer.primitivePatterns[priority][pattern] = processor
	}
}

// AddCollectionPattern will add the collection pattern to this one.
func (lexer *lexerImpl) AddCollectionPattern(start string, end string, processor CollectionProcessor) {
	pattern := start + end
	if _, has := lexer.collectionPatterns[pattern]; !has {
		lexer.collectionPatterns[pattern] = &collProcDef{
			start:     start,
			end:       end,
			processor: processor,
		}
	}
}

func (lexer *lexerImpl) addPattern(pattern []byte, action lexmachine.Action) {
	if lexer.patternAdder != nil {
		lexer.patternAdder(pattern, action)
		return
	}

	lexer.lex.Add(pattern, action)
}

func (lexer *lexerImpl) addPatternDebugger(pattern []byte, action lexmachine.Action) {
	lexer.lex.Add(pattern, action)

	fmt.Println("Pattern: ", string(pattern))
	lex := lexmachine.NewLexer()

	lex.Add([]byte(pattern), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		fmt.Println("  Matches empty: ", string(pattern))
		return true, nil
	})

	compile := lex.CompileNFA

	var err error
	if err = compile(); err == nil {

		var s *lexmachine.Scanner
		if s, err = lex.Scanner([]byte("")); err == nil {

			var tok interface{}
			var end bool

			if tok, err, end = s.Next(); err == nil {
				if !end {
					err = errors.New("expected an end of string")
				} else {
					if tok != nil {
						switch v := tok.(type) {
						case bool:
							if v {
								err = errors.New("expected false")
							}
						default:
							err = errors.New(fmt.Sprintf("expected a bool: %#v", tok))
						}
					}
				}
			}
		}
	}

	if err != nil {
		fmt.Println("  Error: ", err)
	}
}
