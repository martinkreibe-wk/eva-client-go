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

type PrimitiveProcessor func(value string) (Element, error)
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

// ElementTypeFactory defines the factory for an element.
type ElementTypeFactory func(interface{}) (Element, error)

// Lexer defines the lexical analyser for the
type Lexer interface {
	GetFactory(elementType ElementType, tag string) (ElementTypeFactory, bool)
	RemoveFactory(elementType ElementType, tag string)

	// AddPattern will take a pattern and attach the processor for that pattern.
	AddPrimitiveFactory(priority PrimitiveType, elemType ElementType, tag string, Element ElementTypeFactory, processor PrimitiveProcessor, patterns ...string) error

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
	primitives         map[PrimitiveType][]*primitiveDef
	collectionPatterns map[string]*collProcDef
	factories          map[ElementType]map[string]ElementTypeFactory
	lex                *lexmachine.Lexer
}

type elementDefinition struct {
	elemType    ElementType
	initializer func(lexer Lexer) error
}

// newLexer will create a new lexer.
func newLexer() (Lexer, error) {
	lexer := &lexerImpl{}

	// typeDefinitions holds the type to name/initializer mappings
	// NOTE: ORDER MATTERS!!
	var typeDefinitions = []*elementDefinition{
		{UnknownType, nil},
		{NilType, initNil},
		{BooleanType, initBoolean},
		{StringType, initString},
		{CharacterType, initCharacter},
		{SymbolType, initSymbol},
		{KeywordType, initKeyword},
		{IntegerType, initInteger},
		{FloatType, initFloat},
		{InstantType, initInstant},
		{UUIDType, initUUID},
		{ListType, initList},
		{VectorType, initVector},
		{MapType, initMap},
		{SetType, initSet},

		// TODO
		{URIType, nil},
		{BytesType, nil},
		{BigIntType, nil},
		{BigDecType, nil},
		{DoubleType, nil},
		{RefType, nil},
	}

	for _, def := range typeDefinitions {
		if def.initializer != nil {
			if err := def.initializer(lexer); err != nil {
				return nil, err
			}
		}
	}

	lex := lexmachine.NewLexer()

	if lexer.primitives != nil {
		for i := PrimitiveType(0); i < lastPrimitivePriority; i++ {
			if defs, has := lexer.primitives[i]; has {

				var orderedDefs []*primitiveDef
				var last []*primitiveDef

				for _, def := range defs {
					if def.tag == "" {
						last = append(last, def)
					} else {
						orderedDefs = append(orderedDefs, def)
					}
				}

				if last != nil {
					orderedDefs = append(orderedDefs, last...)
				}

				for _, def := range orderedDefs {
					for _, pattern := range def.patterns {
						proc := def.processor

						var realPattern []byte
						if len(def.tag) == 0 {
							realPattern = buildTagPattern(pattern, true)
						} else {
							realPattern = []byte(fmt.Sprintf("%s%s(\\s)+%s", TagPrefix, def.tag, pattern))
						}

						lex.Add(realPattern, func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
							tag, value := splitTag(match.Bytes, "")

							elem, err := proc(value)
							if err != nil {
								return nil, err
							}

							err = elem.SetTag(tag)
							if err != nil {
								return nil, err
							}

							return elem, nil
						})
					}
				}
			}
		}
	}

	endPatterns := map[string]bool{}

	lexSpecialChars := []string{
		"\\", "[", "]", "{", "}", "(", ")",
	}

	if lexer.collectionPatterns != nil {
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
				lex.Add([]byte(end), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
					return tokenType(end), nil
				})
				endPatterns[end] = true
			}

			// Add the non tagged items.
			lex.Add(buildTagPattern(start, false), func(scan *lexmachine.Scanner, match *machines.Match) (v interface{}, e error) {
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
	}

	lex.Add([]byte("(\\s|,)+"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		return skipToken, nil
	})

	lex.Add([]byte(";[^\\n]*(\\n)?"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		return skipToken, nil
	})

	if err := lex.CompileNFA(); err != nil {
		return nil, err
	}

	lexer.lex = lex

	return lexer, nil
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
		return nil, MakeErrorWithFormat(ErrParserError, "Expected one result, got: %d = %+v", len(elems), elems)
	}
}

func (lexer *lexerImpl) GetFactory(elementType ElementType, tag string) (ElementTypeFactory, bool) {
	if lexer.factories == nil {
		return nil, false
	}

	if _, has := lexer.factories[elementType]; !has {
		return nil, false
	}

	factory, has := lexer.factories[elementType][tag]
	return factory, has
}

func (lexer *lexerImpl) RemoveFactory(elementType ElementType, tag string) {
	if lexer.factories == nil {
		return
	}

	if _, has := lexer.factories[elementType]; has {
		if _, has = lexer.factories[elementType][tag]; has {
			delete(lexer.factories[elementType], tag)
		}

		if len(lexer.factories[elementType]) == 0 {
			delete(lexer.factories, elementType)
		}
	}
}

type primitiveDef struct {
	processor PrimitiveProcessor
	patterns  []string
	tag       string
}

// AddPattern will add a pattern to the lexer
func (lexer *lexerImpl) AddPrimitiveFactory(priority PrimitiveType, elemType ElementType, tag string, elemFactory ElementTypeFactory, processor PrimitiveProcessor, patterns ...string) error {

	if lexer.primitives == nil {
		lexer.primitives = make(map[PrimitiveType][]*primitiveDef)
	}

	lexer.primitives[priority] = append(lexer.primitives[priority], &primitiveDef{
		processor: processor,
		patterns:  patterns,
		tag:       tag,
	})

	if lexer.factories == nil {
		lexer.factories = make(map[ElementType]map[string]ElementTypeFactory)
	}

	if _, has := lexer.factories[elemType]; !has {
		lexer.factories[elemType] = make(map[string]ElementTypeFactory)
	}

	if _, has := lexer.factories[elemType][tag]; has {
		return MakeErrorWithFormat(ErrInvalidFactory, "duplicate: `%s`:`%s`", elemType, tag)
	}

	lexer.factories[elemType][tag] = elemFactory

	return nil
}

// AddCollectionPattern will add the collection pattern to this one.
func (lexer *lexerImpl) AddCollectionPattern(start string, end string, processor CollectionProcessor) {

	if lexer.collectionPatterns == nil {
		lexer.collectionPatterns = make(map[string]*collProcDef)
	}

	pattern := start + end
	if _, has := lexer.collectionPatterns[pattern]; !has {
		lexer.collectionPatterns[pattern] = &collProcDef{
			start:     start,
			end:       end,
			processor: processor,
		}
	}
}
