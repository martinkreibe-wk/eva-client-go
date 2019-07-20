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

// fromFloat convert the float passed in (through the interface) to an Element.
func fromFloat(input interface{}) (elem Element, e error) {
	if v, ok := input.(float64); ok {
		return NewFloatElement(v)
	} else {
		e = MakeError(ErrInvalidInput, input)
	}
	return elem, e
}

// parseFloatElem parses the string into a float Element
func parseFloatElem(tag string, tokenValue string) (el Element, err error) {

	if strings.HasSuffix(tokenValue, "M") {
		tokenValue = strings.TrimSuffix(tokenValue, "M")
	}

	var v float64
	if v, err = strconv.ParseFloat(tokenValue, 64); err == nil {
		el, err = NewFloatElement(v)
		if err != nil {
			return nil, err
		}

		err = el.SetTag(tag)
	}

	return el, err
}

// floatSerialize takes the input value and serialize it.
func floatSerialize(serializer Serializer, tag string, value interface{}) (out string, e error) {
	switch serializer.MimeType() {
	case EvaEdnMimeType:
		if len(tag) > 0 {
			out = TagPrefix + tag + " "
		}
		out += strconv.FormatFloat(value.(float64), 'E', -1, 64)
	default:
		e = MakeError(ErrUnknownMimeType, serializer.MimeType())
	}
	return out, e
}

// init will add the element factory to the collection of factories
func initFloat(lexer Lexer) (err error) {
	if err = addElementTypeFactory(FloatType, fromFloat); err == nil {
		lexer.AddPattern(FloatPrimitive, "[-+]?(0|[1-9][0-9]*)(\\.[0-9]*)?([eE][-+]?[0-9]+)?M?", parseFloatElem)
	}

	return err
}

// NewFloatElement creates a new float point element or an error.
func NewFloatElement(value float64) (Element, error) {
	return baseFactory().make(value, FloatType, floatSerialize)
}
