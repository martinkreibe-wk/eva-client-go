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

// fromNil convert the integer64 passed in (through the interface) to an Element.
func fromNil(input interface{}) (Element, error) {
	if input != nil {
		return nil, MakeError(ErrInvalidInput, input)
	}
	return NewNilElement()
}

// parseNil parses the string into a nil Element
func parseNil(tag string, _ string) (Element, error) {
	elem, err := NewNilElement()
	if err != nil {
		return nil, err
	}
	if err = elem.SetTag(tag); err != nil {
		return nil, err
	}
	return elem, nil
}

// nilSerializer takes the input value and serialize it.
func nilSerializer(serializer Serializer, tag string, _ interface{}) (out string, e error) {
	switch serializer.MimeType() {
	case EvaEdnMimeType:
		if len(tag) > 0 {
			out = TagPrefix + tag + " "
		}
		return out + "nil", nil
	default:
		return "", MakeError(ErrUnknownMimeType, serializer.MimeType())
	}
}

// initNil will add the element factory to the collection of factories
func initNil(lexer Lexer) error {
	if err := addElementTypeFactory(NilType, fromNil); err != nil {
		return err
	}

	lexer.AddPattern(LiteralPrimitive, "nil", parseNil)
	return nil
}

// NewNilElement returns the nil element or an error.
func NewNilElement() (Element, error) {
	return baseFactory().make(nil, NilType, nilSerializer)
}
