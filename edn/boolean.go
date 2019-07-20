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

import "strconv"

// fromBool convert the boolean passed in (through the interface) to an Element.
func fromBool(input interface{}) (elem Element, e error) {
	if v, ok := input.(bool); ok {
		return NewBooleanElement(v)
	} else {
		e = MakeError(ErrInvalidInput, input)
	}
	return elem, e
}

// parseBoolElem parses the string into a boolean Element
func parseBoolElem(tag string, tokenValue string) (elem Element, err error) {

	var val bool
	switch tokenValue {
	case "true":
		val = true
	case "false":
		val = false
	default:
		return nil, MakeErrorWithFormat(ErrParserError, "Unknown bool: `%s`", tokenValue)
	}

	elem, err = NewBooleanElement(val)
	if err != nil {
		return nil, err
	}

	if err = elem.SetTag(tag); err != nil {
		return nil, err
	}

	return elem, nil
}

// boolSerializer takes the input value and serialize it.
func boolSerializer(serializer Serializer, tag string, value interface{}) (out string, e error) {

	switch serializer.MimeType() {
	case EvaEdnMimeType:
		if len(tag) > 0 {
			out = TagPrefix + tag + " "
		}
		out += strconv.FormatBool(value.(bool))
	default:
		e = MakeError(ErrUnknownMimeType, serializer.MimeType())
	}

	return out, e
}

// initBoolean will add the element factory to the collection of factories
func initBoolean(lexer Lexer) (err error) {
	if err = addElementTypeFactory(BooleanType, fromBool); err == nil {

		// For some reason, I can not do `true|false` or any other variation without
		// the symbol parser picking the bool up.
		lexer.AddPattern(LiteralPrimitive, "true", parseBoolElem)
		lexer.AddPattern(LiteralPrimitive, "false", parseBoolElem)
	}

	return err
}

// NewBooleanElement creates a new boolean element or an error.
func NewBooleanElement(value bool) (Element, error) {
	return baseFactory().make(value, BooleanType, boolSerializer)
}
