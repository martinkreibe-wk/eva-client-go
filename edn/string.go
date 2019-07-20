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
)

type stringProcessor func(string) (Element, error)

var stringProcessors = map[string]stringProcessor{
	UUIDElementTag:    uuidStringProcessor,
	InstantElementTag: instStringProcessor,
}

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

// normalStringProcessor defines the rule for normal string processing.
func normalStringProcessor(tokenValue string) (el Element, e error) {
	length := len(tokenValue)

	var out []rune
	for i := 0; length > i && e == nil; i++ {
		if current := rune(tokenValue[i]); current == '\\' {
			if next := i + 1; length > next {
				nextCh := rune(tokenValue[next])
				if ch, has := specialStrings[nextCh]; has {
					i++
					out = append(out, ch)
				} else if nextCh == 'u' && length > next+4 {
					i++ // remove the 'u'

					// Look for the next 4 characters
					unicode := tokenValue[next+1 : next+5]
					var v int64
					if v, e = strconv.ParseInt(unicode, 16, 16); e == nil {
						i = i + 4
						out = append(out, rune(v))
					}
				} else {
					e = MakeErrorWithFormat(ErrParserError, "Invalid escape character: %#U", ch)
				}
			} else {
				e = MakeError(ErrParserError, "Escape character found at end of string.")
			}
		} else {
			out = append(out, current)
		}
	}

	if e == nil {
		return NewStringElement(string(out))
	}

	return el, e
}

// init will add the element factory to the collection of factories
func initString(lexer Lexer) (err error) {
	if err = addElementTypeFactory(StringType, func(input interface{}) (elem Element, e error) {
		v, ok := input.(string)
		if !ok {
			return nil, MakeError(ErrInvalidInput, input)
		}
		return NewStringElement(v)
	}); err == nil {
		lexer.AddPattern(StringPrimitive, "\"(\\w|\\d| |[-+*!?$%&=<>.#:()\\[\\]@^;,/{}'|`~]|\\\\([tbnrf\"'\\\\]|u[0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f]))*\"", func(tag string, tokenValue string) (el Element, e error) {
			var proc stringProcessor
			var has bool

			if proc, has = stringProcessors[tag]; !has {
				proc = normalStringProcessor
			}

			el, e = proc(tokenValue[1 : len(tokenValue)-1])

			if e == nil && el != nil {
				e = el.SetTag(tag)
			}

			return el, e
		})
	}

	return err
}

// NewStringElement creates a new string element or an error.
func NewStringElement(value string) (Element, error) {

	return baseFactory().make(value, StringType, func(serializer Serializer, tag string, value interface{}) (out string, e error) {
		switch serializer.MimeType() {
		case EvaEdnMimeType:
			if len(tag) > 0 {
				out = TagPrefix + tag + " "
			}
			out += strconv.Quote(value.(string))
		default:
			e = MakeError(ErrUnknownMimeType, serializer.MimeType())
		}
		return out, e
	})
}
