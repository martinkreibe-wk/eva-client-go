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
func instStringProcessor(tokenValue string) (el Element, e error) {
	var t time.Time
	if t, e = time.Parse(time.RFC3339, tokenValue); e == nil {
		el = NewInstantElement(t)
	}

	return el, e
}

// init will add the element factory to the collection of factories
func initInstant(_ Lexer) error {
	return addElementTypeFactory(InstantType, func(input interface{}) (elem Element, err error) {
		if v, ok := input.(time.Time); ok {
			elem = NewInstantElement(v)
		} else {
			err = MakeError(ErrInvalidInput, input)
		}
		return elem, err
	})
}

// NewInstantElement creates a new instant element or an error.
func NewInstantElement(value time.Time) (elem Element) {

	var err error
	if elem, err = baseFactory().make(value, InstantType, func(serializer Serializer, tag string, value interface{}) (out string, e error) {
		switch serializer.MimeType() {
		case EvaEdnMimeType:
			if len(tag) > 0 {
				out = TagPrefix + tag + " "
			}
			out += value.(time.Time).Format(time.RFC3339)
		default:
			e = MakeError(ErrUnknownMimeType, serializer.MimeType())
		}

		return out, e
	}); err == nil {
		elem.SetTag(InstantElementTag)
	} else {
		panic(err)
	}

	return elem
}
