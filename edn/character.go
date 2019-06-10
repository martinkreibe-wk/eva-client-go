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
	"strconv"
	"strings"
)

const (

	// CharacterPrefix defines the prefix for characters
	CharacterPrefix = "\\"
)

var specialCharacters = map[rune]string{
	'\r': "return",
	'\n': "newline",
	' ':  "space",
	'\t': "tab",
}

// init will add the element factory to the collection of factories
func initCharacter(lexer Lexer) (err error) {
	if err = addElementTypeFactory(CharacterType, func(input interface{}) (elem Element, e error) {
		if v, ok := input.(rune); ok {
			elem = NewCharacterElement(v)
		} else {
			e = MakeError(ErrInvalidInput, input)
		}
		return elem, e
	}); err == nil {
		for r, v := range specialCharacters {
			c := r // for some reason I need to use a local variable or things get mixed up.
			lexer.AddPattern(CharacterPrimitive, "\\\\"+v, func(tag string, tokenValue string) (Element, error) {
				el := NewCharacterElement(c)
				return el, el.SetTag(tag)
			})
		}

		lexer.AddPattern(CharacterPrimitive, "\\\\u[0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f]", func(tag string, tokenValue string) (el Element, e error) {
			tokenValue = strings.TrimPrefix(tokenValue, CharacterPrefix+"u")
			var v int64

			// It isn't possible to get anything other then 4 characters, so checking isn't needed.
			if v, e = strconv.ParseInt(tokenValue, 16, 16); e == nil {
				el = NewCharacterElement(rune(v))
				e = el.SetTag(tag)
			}

			return el, e
		})

		lexer.AddPattern(CharacterPrimitive, "\\\\\\w", func(tag string, tokenValue string) (el Element, e error) {

			tokenValue = strings.TrimPrefix(tokenValue, CharacterPrefix)
			runes := []rune(tokenValue)

			// It isn't possible to get anything other then a single character, so checking isn't needed.
			el = NewCharacterElement(runes[0])
			e = el.SetTag(tag)

			return el, e
		})
	}

	return err
}

// NewCharacterElement creates a new character element or an error.
func NewCharacterElement(value rune) (elem Element) {

	var err error
	if elem, err = baseFactory().make(value, CharacterType, func(serializer Serializer, tag string, value interface{}) (out string, e error) {

		switch serializer.MimeType() {
		case EvaEdnMimeType:
			if len(tag) > 0 {
				out = TagPrefix + tag + " "
			}

			r := value.(rune)
			if char, has := specialCharacters[r]; has {
				out += CharacterPrefix + char
			} else {

				// if there is no special character, then quote the rune, remove the single quotes around this, then
				// if it is an ASCII then make sure to prefix is intact.
				if char = strings.Trim(fmt.Sprintf("%+q", r), "'"); strings.HasPrefix(char, CharacterPrefix) {
					out += char
				} else {
					out += CharacterPrefix + char
				}
			}
		default:
			e = MakeError(ErrUnknownMimeType, serializer.MimeType())
		}

		return out, e
	}); err != nil {
		panic(err)
	}

	return elem
}
