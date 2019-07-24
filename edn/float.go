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
func fromFloat(input interface{}) (Element, error) {
	v, ok := input.(float64)
	if !ok {
		return nil, MakeError(ErrInvalidInput, input)
	}
	return NewFloatElement(v)
}

// parseFloatElem parses the string into a float Element
func parseFloatElem(tokenValue string) (Element, error) {

	if strings.HasSuffix(tokenValue, "M") {
		tokenValue = strings.TrimSuffix(tokenValue, "M")
	}

	v, err := strconv.ParseFloat(tokenValue, 64)
	if err != nil {
		return nil, err
	}

	return NewFloatElement(v)
}

// floatSerialize takes the input value and serialize it.
func floatSerialize(serializer Serializer, tag string, value interface{}) (string, error) {
	switch serializer.MimeType() {
	case EvaEdnMimeType:
		var out string
		if len(tag) > 0 {
			out = TagPrefix + tag + " "
		}
		return out + strconv.FormatFloat(value.(float64), 'E', -1, 64), nil
	default:
		return "", MakeError(ErrUnknownMimeType, serializer.MimeType())
	}
}

// NewFloatElement creates a new float point element or an error.
func NewFloatElement(value float64) (Element, error) {
	return baseFactory().make(value, FloatType, floatSerialize)
}
