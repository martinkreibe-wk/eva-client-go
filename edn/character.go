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

// parseSpecialCharElem parses the standard string into a character Element
func parseCharElem(value string) (Element, error) {

	// strip the first '\' char.
	if strings.HasPrefix(value, string(CharacterPrefix)) {
		value = strings.TrimPrefix(value, CharacterPrefix)
	}

	if strings.HasPrefix(value, "u") && len(value) > 1 {
		value = strings.TrimPrefix(value, "u")

		// It isn't possible to get anything other then 4 characters, so checking isn't needed.
		v, err := strconv.ParseInt(value, 16, 16)
		if err != nil {
			return nil, MakeErrorWithFormat(ErrParserError, "`%s` could not be converted to a hex value: %v", value, err)
		}

		return NewCharacterElement(rune(v))
	}

	if r, has := specialCharacters[value]; has {
		return NewCharacterElement(r)
	}

	runes := []rune(value)

	// It isn't possible to get anything other then a single character, so checking isn't needed.
	return NewCharacterElement(runes[0])
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
	return baseFactory().make(value, CharacterType, NoTag, charSerializer)
}
