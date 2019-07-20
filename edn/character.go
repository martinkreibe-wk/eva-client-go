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
	CharacterPrefix = `\`
)

var specialCharacters = map[string]rune{
	"return":  '\r',
	"newline": '\n',
	"space":   ' ',
	"tab":     '\t',
}

// fromChar convert the character passed in (through the interface) to an Element.
func fromChar(input interface{}) (Element, error) {
	v, ok := input.(rune)
	if !ok {
		return nil, MakeError(ErrInvalidInput, input)
	}

	return NewCharacterElement(v)
}

// parseSpecialCharElem parses the special string into a character Element
func parseSpecialCharElem(tag string, tokenValue string) (Element, error) {

	// strip the first '\' char.
	if !strings.HasPrefix(tokenValue, string(CharacterPrefix)) {
		return nil, MakeError(ErrParserError, "Missing character prefix.")
	}

	var r rune
	var has bool
	if r, has = specialCharacters[strings.TrimPrefix(tokenValue, string(CharacterPrefix))]; !has {
		return nil, MakeErrorWithFormat(ErrParserError, "Unknown character %s", tokenValue)
	}

	elem, err := NewCharacterElement(r)
	if err != nil {
		return nil, err
	}
	if err = elem.SetTag(tag); err != nil {
		return nil, err
	}

	return elem, nil
}

// parseSpecialCharElem parses the unicode string into a character Element
func parseUnicodeCharElem(tag string, tokenValue string) (Element, error) {
	tokenValue = strings.TrimPrefix(tokenValue, CharacterPrefix+"u")

	// It isn't possible to get anything other then 4 characters, so checking isn't needed.
	v, err := strconv.ParseInt(tokenValue, 16, 16)
	if err != nil {
		return nil, err
	}

	elem, err := NewCharacterElement(rune(v))
	if err != nil {
		return nil, err
	}

	err = elem.SetTag(tag)
	if err != nil {
		return nil, err
	}

	return elem, nil
}

// parseSpecialCharElem parses the standard string into a character Element
func parseCharElem(tag string, tokenValue string) (Element, error) {

	tokenValue = strings.TrimPrefix(tokenValue, CharacterPrefix)
	runes := []rune(tokenValue)

	// It isn't possible to get anything other then a single character, so checking isn't needed.
	elem, err := NewCharacterElement(runes[0])
	if err != nil {
		return nil, err
	}

	err = elem.SetTag(tag)
	if err != nil {
		return nil, err
	}

	return elem, nil
}

// initCharacter will add the element factory to the collection of factories
func initCharacter(lexer Lexer) error {
	if err := addElementTypeFactory(CharacterType, fromChar); err != nil {
		return err
	}

	for v := range specialCharacters {
		lexer.AddPattern(CharacterPrimitive, `\\`+v, parseSpecialCharElem)
	}

	lexer.AddPattern(CharacterPrimitive, `\\u[0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f]`, parseUnicodeCharElem)
	lexer.AddPattern(CharacterPrimitive, `\\\w`, parseCharElem)

	return nil
}

// charSerializer takes the input value and serialize it.
func charSerializer(serializer Serializer, tag string, value interface{}) (string, error) {

	switch serializer.MimeType() {
	case EvaEdnMimeType:
		var out string
		if len(tag) > 0 {
			out = TagPrefix + tag + " "
		}

		val := value.(rune)

		// look at the special characters first.
		for v, r := range specialCharacters {
			if val == r {
				return out + CharacterPrefix + v, nil
			}
		}

		// if there is no special character, then quote the rune, remove the single quotes around this, then
		// if it is an ASCII then make sure to prefix is intact.
		char := strings.Trim(fmt.Sprintf("%+q", val), "'")
		if strings.HasPrefix(char, CharacterPrefix) {
			return out + char, nil
		}

		return out + CharacterPrefix + char, nil

	default:
		return "", MakeError(ErrUnknownMimeType, serializer.MimeType())
	}
}

// NewCharacterElement creates a new character element or an error.
func NewCharacterElement(value rune) (Element, error) {
	return baseFactory().make(value, CharacterType, charSerializer)
}
