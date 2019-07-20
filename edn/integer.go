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
	int64Regex = "[-+]?(0|[1-9][0-9]*)N?"
)

// fromInt64 convert the integer64 passed in (through the interface) to an Element.
func fromInt64(input interface{}) (elem Element, e error) {

	var v int64
	var ok bool
	if v, ok = input.(int64); !ok {
		return nil, MakeError(ErrInvalidInput, input)
	}

	return NewIntegerElement(v)
}

// parseInt64Elem parses the string into a int64 Element
func parseInt64Elem(tag string, tokenValue string) (elem Element, err error) {

	if strings.HasSuffix(tokenValue, "N") {
		tokenValue = strings.TrimSuffix(tokenValue, "N")
	}

	var v int64
	if v, err = strconv.ParseInt(tokenValue, 10, 64); err != nil {
		return nil, err
	}
	elem, err = NewIntegerElement(v)
	if err != nil {
		return nil, err
	}

	if err = elem.SetTag(tag); err != nil {
		return nil, err
	}

	return elem, nil
}

// int64Serializer takes the input value and serialize it.
func int64Serializer(serializer Serializer, tag string, value interface{}) (out string, e error) {
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
}

// initInteger will add the element factory to the collection of factories
func initInteger(lexer Lexer) (err error) {
	if err = addElementTypeFactory(IntegerType, fromInt64); err == nil {
		lexer.AddPattern(IntegerPrimitive, int64Regex, parseInt64Elem)
	}

	return err
}

// NewIntegerElement creates a new integer element or an error.
func NewIntegerElement(value int64) (Element, error) {
	return baseFactory().make(value, IntegerType, int64Serializer)
}
