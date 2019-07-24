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

type lexerImpl struct {
	primitives         map[PrimitiveType][]*primitiveDef
	collectionPatterns map[string]*collProcDef
	factories          map[ElementType]map[string]ElementTypeFactory
	lex                *lexmachine.Lexer
}

type collProcDef struct {
	start     string
	end       string
	processor CollectionProcessor
}

// newLexer will create a new lexer.
func newLexer() (Lexer, error) {
	lexer := &lexerImpl{
		primitives:         make(map[PrimitiveType][]*primitiveDef),
		collectionPatterns: make(map[string]*collProcDef),
		factories:          make(map[ElementType]map[string]ElementTypeFactory),
	}

	err := lexer.AddPrimitiveFactory(LiteralPrimitive, NilType, NoTag, fromNil, parseNil, "nil")
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(LiteralPrimitive, BooleanType, NoTag, fromBool, parseBoolElem, "true", "false")
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(StringPrimitive, StringType, NoTag, stringFactory, parseString, StringPattern)
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(CharacterPrimitive, CharacterType, NoTag, fromChar, parseCharElem, `\\([A-Z0-9a-mo-qv-z]|n(ewline)?|r(eturn)?|s(pace)?|t(ab)?|u([0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f])?)`)
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(SymbolPrimitive, SymbolType, NoTag, fromSymbol, parseSymbol, "[*!?$%&=<>_a-zA-Z.]([-+*!?$%&=<>_.#]|\\w)*(/([-+*!?$%&=<>_.#]|\\w)*)?")
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(SymbolPrimitive, KeywordType, NoTag, fromKeyword, parseKeyword, ":([*!?$%&=<>]|\\w)([-+*!?$%&=<>.#]|\\w)*(/([-+*!?$%&=<>.#]|\\w)*)?")
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(IntegerPrimitive, IntegerType, NoTag, fromInt64, parseInt64Elem, int64Regex)
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(FloatPrimitive, FloatType, NoTag, fromFloat, parseFloatElem, "[-+]?(0|[1-9][0-9]*)(\\.[0-9]*)?([eE][-+]?[0-9]+)?M?")
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(StringPrimitive, InstantType, InstantElementTag, instantFactory, instStringProcessor, InstPattern)
	if err != nil {
		return nil, err
	}

	err = lexer.AddPrimitiveFactory(StringPrimitive, UUIDType, UUIDElementTag, uuidFactory, uuidStringProcessor, UUIDPattern)
	if err != nil {
		return nil, err
	}

	/////

	/*
		URIType
		BytesType
		BigIntType
		BigDecType
		DoubleType
		RefType
	*/

	/////

	lexer.AddCollectionPattern(ListStartLiteral, ListEndLiteral, func(elements []Element) (el CollectionElement, e error) {
		return NewList(elements...)
	})

	lexer.AddCollectionPattern(VectorStartLiteral, VectorEndLiteral, func(elements []Element) (el CollectionElement, e error) {
		return NewVector(elements...)
	})

	makePairs := func(elements []Element) (pairs Pairs, err error) {

		if l := len(elements); l%2 == 0 {
			for i := 0; i < l; i = i + 2 {
				if err = pairs.Append(elements[i], elements[i+1]); err != nil {
					return Pairs{}, err
				}
			}
		} else {
			err = MakeError(ErrInvalidPair, "Map input are not paired up.")
		}

		return pairs, err
	}

	lexer.AddCollectionPattern(MapStartLiteral, MapEndLiteral, func(elements []Element) (el CollectionElement, e error) {
		var pairs Pairs
		if pairs, e = makePairs(elements); e != nil {
			return nil, e
		}
		return NewMap(pairs.Raw()...)
	})

	lexer.AddCollectionPattern(SetStartLiteral, SetEndLiteral, func(elements []Element) (el CollectionElement, e error) {
		return NewSet(elements...)
	})

	lex := lexmachine.NewLexer()

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

				var el CollectionElement
				if e == nil {
					el, e = processor(children)
				}

				if e == nil {
					e = el.SetTag(tag)
				}

				return el, e
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

	lexer.primitives[priority] = append(lexer.primitives[priority], &primitiveDef{
		processor: processor,
		patterns:  patterns,
		tag:       tag,
	})

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

	pattern := start + end
	if _, has := lexer.collectionPatterns[pattern]; !has {
		lexer.collectionPatterns[pattern] = &collProcDef{
			start:     start,
			end:       end,
			processor: processor,
		}
	}
}
