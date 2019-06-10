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

	"github.com/mattrobenolt/gocql/uuid"
)

const (

	// InvalidElement defines an invalid element was encountered.
	ErrInvalidElement = ErrorMessage("Invalid Element")

	// TagPrefix defines the prefix for tags.
	TagPrefix = "#"
)

// Element defines the interface for EDN elements.
type Element interface {
	Serializable

	// ElementType returns the current type of this element.
	ElementType() ElementType

	// Value of the element
	Value() interface{}

	// HasTag returns true if the element has a tag prefix
	HasTag() bool

	// Tag returns the prefixed tag if it exists.
	Tag() string

	// SetTag sets the tag to the incoming value. If the value is an empty string then the tag is unset.
	SetTag(string) (err error)

	// Equals checks if the input element is equal to this element.
	Equals(e Element) (result bool)
}

// stereotypePrimitive returns the cleaned value and stereotype, or it returns an error.
func stereotypePrimitive(value interface{}) (_ interface{}, stereotype ElementType, err error) {

	stereotype = UnknownType
	switch v := value.(type) {
	case int:
		stereotype = IntegerType
		value = int64(v)
	case int32:
		stereotype = IntegerType
		value = int64(v)
	case bool:
		stereotype = BooleanType
		value = v
	case int64:
		stereotype = IntegerType
	case float32:
		stereotype = FloatType
		value = float64(v)
	case float64:
		stereotype = FloatType
	case string:
		if v == "nil" {
			stereotype = NilType
			value = nil
		} else {
			stereotype = StringType
			if len(v) > 0 && v[0] == ':' {
				stereotype = KeywordType
			}
		}
	case time.Time:
		stereotype = InstantType
	case uuid.UUID:
		stereotype = UUIDType
	default:
		err = MakeErrorWithFormat(ErrUnknownMimeType, "[%T]: %#v", v, v)
	}

	return value, stereotype, err
}

// NewPrimitiveElement creates a new primitive element from the inputs.
func NewPrimitiveElement(value interface{}) (elem Element, err error) {

	if value == nil {
		elem = NewNilElement()
	} else {
		var is bool
		if elem, is = value.(Element); !is {

			var stereotype ElementType
			var val interface{}

			if val, stereotype, err = stereotypePrimitive(value); err == nil {
				if factory, has := typeFactories[stereotype]; has {
					elem, err = factory(val)
				} else {
					err = MakeErrorWithFormat(ErrInvalidElement, "type: %s", stereotype.Name())
				}

			}
		}
	}

	return elem, err
}

// IsPrimitive checks to see if the input variable is
func IsPrimitive(value interface{}) bool {
	_, stereotype, _ := stereotypePrimitive(value)
	return stereotype != UnknownType
}
