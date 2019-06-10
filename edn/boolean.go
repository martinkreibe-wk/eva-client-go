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

import "strconv"

// initBoolean will add the element factory to the collection of factories
func initBoolean(lexer Lexer) (err error) {
	if err = addElementTypeFactory(BooleanType, func(input interface{}) (elem Element, e error) {
		if v, ok := input.(bool); ok {
			elem = NewBooleanElement(v)
		} else {
			e = MakeError(ErrInvalidInput, input)
		}
		return elem, e
	}); err == nil {
		lexer.AddPattern(LiteralPrimitive, "true", func(tag string, tokenValue string) (Element, error) {
			elem := NewBooleanElement(true)
			return elem, elem.SetTag(tag)
		})
		lexer.AddPattern(LiteralPrimitive, "false", func(tag string, tokenValue string) (Element, error) {
			elem := NewBooleanElement(false)
			return elem, elem.SetTag(tag)
		})
	}

	return err
}

// NewBooleanElement creates a new boolean element or an error.
func NewBooleanElement(value bool) (elem Element) {

	var err error
	if elem, err = baseFactory().make(value, BooleanType, func(serializer Serializer, tag string, value interface{}) (out string, e error) {

		switch serializer.MimeType() {
		case EvaEdnMimeType:
			if len(tag) > 0 {
				out = TagPrefix + tag + " "
			}
			out += strconv.FormatBool(value.(bool))
		default:
			e = MakeError(ErrUnknownMimeType, serializer.MimeType())
		}

		return out, e
	}); err != nil {
		panic(err)
	}

	return elem
}
