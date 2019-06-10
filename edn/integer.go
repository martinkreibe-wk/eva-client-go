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

// init will add the element factory to the collection of factories
func initInteger(lexer Lexer) (err error) {
	if err = addElementTypeFactory(IntegerType, func(input interface{}) (elem Element, e error) {
		if v, ok := input.(int64); ok {
			elem = NewIntegerElement(v)
		} else {
			e = MakeError(ErrInvalidInput, input)
		}
		return elem, e
	}); err == nil {
		lexer.AddPattern(IntegerPrimitive, "[-+]?(0|[1-9][0-9]*)N?", func(tag string, tokenValue string) (el Element, e error) {

			if strings.HasSuffix(tokenValue, "N") {
				tokenValue = strings.TrimSuffix(tokenValue, "N")
			}

			var v int64
			if v, e = strconv.ParseInt(tokenValue, 10, 64); e == nil {
				el = NewIntegerElement(v)
				e = el.SetTag(tag)
			}

			return el, e
		})
	}

	return err
}

// NewIntegerElement creates a new integer element or an error.
func NewIntegerElement(value int64) (elem Element) {

	var err error
	if elem, err = baseFactory().make(value, IntegerType, func(serializer Serializer, tag string, value interface{}) (out string, e error) {
		switch serializer.MimeType() {
		case EvaEdnMimeType:
			if len(tag) > 0 {
				out = TagPrefix + tag + " "
			}
			out += strconv.FormatInt(value.(int64), 10)
		default:
			e = MakeError(ErrUnknownMimeType, serializer.MimeType())
		}

		return out, e
	}); err != nil {
		panic(err)
	}

	return elem
}
