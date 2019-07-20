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
	"time"
)

const (

	// InstantElementTag defines the instant tag value.
	InstantElementTag = "inst"
)

// instStringProcessor used the string processor but will accurately create the instances.
func instStringProcessor(tokenValue string) (Element, error) {
	timeVal, err := time.Parse(time.RFC3339, tokenValue)
	if err != nil {
		return nil, err
	}

	return NewInstantElement(timeVal)
}

func instantFactory(input interface{}) (Element, error) {
	v, ok := input.(time.Time)
	if !ok {
		return nil, MakeError(ErrInvalidInput, input)
	}
	return NewInstantElement(v)
}

// init will add the element factory to the collection of factories
func initInstant(_ Lexer) error {
	return addElementTypeFactory(InstantType, instantFactory)
}

func instantSerializer(serializer Serializer, tag string, value interface{}) (string, error) {
	switch serializer.MimeType() {
	case EvaEdnMimeType:
		var out string
		if len(tag) > 0 {
			out = TagPrefix + tag + " "
		}
		return out + value.(time.Time).Format(time.RFC3339), nil
	default:
		return "", MakeError(ErrUnknownMimeType, serializer.MimeType())
	}
}

// NewInstantElement creates a new instant element or an error.
func NewInstantElement(value time.Time) (Element, error) {

	elem, err := baseFactory().make(value, InstantType, instantSerializer)
	if err != nil {
		return nil, err
	}

	if err = elem.SetTag(InstantElementTag); err != nil {
		return nil, err
	}

	return elem, err
}
