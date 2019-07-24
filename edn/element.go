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
func stereotypePrimitive(value interface{}) (interface{}, ElementType, string, error) {

	switch v := value.(type) {
	case int:
		return int64(v), IntegerType, NoTag, nil
	case int32:
		return int64(v), IntegerType, NoTag, nil
	case bool:
		return v, BooleanType, NoTag, nil
	case int64:
		return int64(v), IntegerType, NoTag, nil
	case float32:
		return float64(v), FloatType, NoTag, nil
	case float64:
		return v, FloatType, NoTag, nil
	case string:
		if v == "nil" {
			return nil, NilType, NoTag, nil
		}

		if len(v) > 0 && v[0] == ':' {
			return v, KeywordType, NoTag, nil
		}

		return v, StringType, NoTag, nil
	case time.Time:
		return v, InstantType, InstantElementTag, nil
	case uuid.UUID:
		return v, UUIDType, UUIDElementTag, nil
	default:
		return nil, UnknownType, NoTag, MakeErrorWithFormat(ErrUnknownMimeType, "[%T]: %#v", v, v)
	}
}

// NewPrimitiveElement creates a new primitive element from the inputs.
func NewPrimitiveElement(value interface{}) (Element, error) {
	return NewPrimitiveElementWithLexer(DefaultLexer, value)
}

// NewPrimitiveElement creates a new primitive element from the inputs.
func NewPrimitiveElementWithLexer(lexer Lexer, value interface{}) (Element, error) {

	if value == nil {
		return NewNilElement()
	}

	if elem, is := value.(Element); is {
		return elem, nil
	}

	val, stereotype, tag, err := stereotypePrimitive(value)
	if err != nil {
		return nil, err
	}

	factory, has := lexer.GetFactory(stereotype, tag)
	if !has {
		return nil, MakeErrorWithFormat(ErrInvalidElement, "type: `%s`, tag: `%s`", stereotype.Name(), tag)
	}

	return factory(val)
}

// IsPrimitive checks to see if the input variable is
func IsPrimitive(value interface{}) bool {
	_, stereotype, _, _ := stereotypePrimitive(value)
	return stereotype != UnknownType
}
