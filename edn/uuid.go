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
	"github.com/mattrobenolt/gocql/uuid"
)

const (

	// UUIDElementTag defines the uuid tag value.
	UUIDElementTag = "uuid"
)

// uuidStringProcessor used the string processor but will accurately create the uuid.
func uuidStringProcessor(tokenValue string) (el Element, e error) {
	var id uuid.UUID
	if id, e = uuid.ParseUUID(tokenValue); e == nil {
		el = NewUUIDElement(id)
	}

	return el, e
}

// init will add the element factory to the collection of factories
func initUUID(_ Lexer) (err error) {
	err = addElementTypeFactory(UUIDType, func(input interface{}) (elem Element, e error) {
		if v, ok := input.(uuid.UUID); ok {
			elem = NewUUIDElement(v)
		} else {
			e = MakeError(ErrInvalidInput, input)
		}
		return elem, e
	})

	return err
}

// NewInstantElement creates a new instant element or an error.
func NewUUIDElement(value uuid.UUID) (elem Element) {

	var err error
	if elem, err = baseFactory().make(value, UUIDType, func(serializer Serializer, tag string, value interface{}) (out string, e error) {
		switch serializer.MimeType() {
		case EvaEdnMimeType:
			if len(tag) > 0 {
				out = TagPrefix + tag + " "
			}
			out += value.(uuid.UUID).String()
		default:
			e = MakeError(ErrUnknownMimeType, serializer.MimeType())
		}
		return out, e
	}); err == nil {
		elem.SetTag(UUIDElementTag)
	} else {
		panic(err)
	}

	return elem
}
