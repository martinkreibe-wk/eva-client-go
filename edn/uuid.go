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
func uuidStringProcessor(tokenValue string) (Element, error) {
	id, err := uuid.ParseUUID(tokenValue)
	if err != nil {
		return nil, err
	}

	return NewUUIDElement(id)
}

func uuidFactory(input interface{}) (Element, error) {
	v, ok := input.(uuid.UUID)
	if !ok {
		return nil, MakeError(ErrInvalidInput, input)
	}

	return NewUUIDElement(v)
}

func uuidSerializer(serializer Serializer, tag string, value interface{}) (string, error) {
	switch serializer.MimeType() {
	case EvaEdnMimeType:
		var out string
		if len(tag) > 0 {
			out = TagPrefix + tag + " "
		}
		return out + value.(uuid.UUID).String(), nil
	default:
		return "", MakeError(ErrUnknownMimeType, serializer.MimeType())
	}
}

// init will add the element factory to the collection of factories
func initUUID(_ Lexer) (err error) {
	return addElementTypeFactory(UUIDType, uuidFactory)
}

// NewInstantElement creates a new instant element or an error.
func NewUUIDElement(value uuid.UUID) (Element, error) {

	elem, err := baseFactory().make(value, UUIDType, uuidSerializer)
	if err != nil {
		return nil, err
	}

	err = elem.SetTag(UUIDElementTag)
	if err != nil {
		return nil, err
	}

	return elem, nil
}
