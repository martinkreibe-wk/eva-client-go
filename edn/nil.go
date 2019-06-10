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

// initNil will add the element factory to the collection of factories
func initNil(lexer Lexer) (err error) {
	if err = addElementTypeFactory(NilType, func(input interface{}) (elem Element, e error) {
		if input == nil {
			elem = NewNilElement()
		} else {
			e = MakeError(ErrInvalidInput, input)
		}
		return elem, e
	}); err == nil {
		lexer.AddPattern(LiteralPrimitive, "nil", func(tag string, tokenValue string) (Element, error) {
			elem := NewNilElement()
			return elem, elem.SetTag(tag)
		})
	}

	return err
}

// NewNilElement returns the nil element or an error.
func NewNilElement() (elem Element) {

	var err error
	if elem, err = baseFactory().make(nil, NilType, func(serializer Serializer, tag string, value interface{}) (out string, e error) {
		switch serializer.MimeType() {
		case EvaEdnMimeType:
			if len(tag) > 0 {
				out = TagPrefix + tag + " "
			}
			out += "nil"
		default:
			e = MakeError(ErrUnknownMimeType, serializer.MimeType())
		}
		return out, e
	}); err != nil {
		panic(err)
	}
	return elem
}
