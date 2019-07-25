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
	"strconv"
	"strings"
)

const (
	StringPattern = "\"(\\w|\\d| |[-+*!?$%&=<>.#:()\\[\\]@^;,/{}'|`~]|\\\\([tbnrf\"'\\\\]|u[0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f]))*\""
)

var specialStrings = map[rune]rune{
	't':  '\t',
	'b':  '\b',
	'n':  '\n',
	'r':  '\r',
	'f':  '\f',
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
}

func stringFactory(input interface{}) (Element, error) {
	v, ok := input.(string)
	if !ok {
		return nil, MakeError(ErrInvalidInput, input)
	}
	return NewStringElement(v)
}

func parseString(tokenValue string) (Element, error) {

	if !strings.HasPrefix(tokenValue, "\"") || !strings.HasSuffix(tokenValue, "\"") {
		return nil, MakeError(ErrParserError, "Expected string to start and end with quotes.")
	}

	tokenValue = tokenValue[1 : len(tokenValue)-1]
	length := len(tokenValue)

	var out []rune
	for i := 0; length > i; i++ {
		current := rune(tokenValue[i])
		switch current {
		case '\\':

			next := i + 1
			if length <= next {
				return nil, MakeError(ErrParserError, "Escape character found at end of string.")
			}

			nextCh := rune(tokenValue[next])
			switch ch, has := specialStrings[nextCh]; {
			case has:
				i++
				out = append(out, ch)
			case nextCh == 'u' && length > next+4:
				i++ // remove the 'u'

				// Look for the next 4 characters
				unicode := tokenValue[next+1 : next+5]
				v, err := strconv.ParseInt(unicode, 16, 16)
				if err != nil {
					return nil, err
				}
				i = i + 4
				out = append(out, rune(v))
			default:
				return nil, MakeErrorWithFormat(ErrParserError, "Invalid escape character: %#U", ch)
			}
		default:
			out = append(out, current)
		}
	}

	return NewStringElement(string(out))
}

// NewStringElement creates a new string element or an error.
func NewStringElement(value string) (Element, error) {
	return baseFactory().make(value, StringType, NoTag)
}
