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
	"strings"
)

const (

	// KeywordPrefix defines the prefix for keywords
	KeywordPrefix = ":"

	// ErrInvalidKeyword defines the error for invalid keywords
	ErrInvalidKeyword = ErrorMessage("Invalid keyword")
)

func fromKeyword(input interface{}) (Element, error) {
	v, ok := input.(string)
	if !ok {
		return nil, MakeErrorWithFormat(ErrInvalidInput, "Value: %#v", input)
	}
	return NewKeywordElement(v)
}

func parseKeyword(tag string, tokenValue string) (Element, error) {
	tokenValue = strings.TrimSuffix(tokenValue, KeywordPrefix)

	elem, err := NewKeywordElement(tokenValue)
	if err != nil {
		return elem, err
	}

	if err = elem.SetTag(tag); err != nil {
		return nil, err
	}

	return elem, err
}

// init will add the element factory to the collection of factories
func initKeyword(lexer Lexer) error {
	if err := addElementTypeFactory(KeywordType, fromKeyword); err != nil {
		return err
	}

	lexer.AddPattern(SymbolPrimitive, ":([*!?$%&=<>]|\\w)([-+*!?$%&=<>.#]|\\w)*(/([-+*!?$%&=<>.#]|\\w)*)?", parseKeyword)
	return nil
}

// NewKeywordElement creates a new character element or an error.
//
// Keywords are identifiers that typically designate themselves. They are semantically akin to enumeration values.
// Keywords follow the rules of symbols, except they can (and must) begin with :, e.g. :fred or :my/fred. If the target
// platform does not have a keyword type distinct from a symbol type, the same type can be used without conflict, since
// the mandatory leading : of keywords is disallowed for symbols. Per the symbol rules above, :/ and :/anything are not
// legal keywords. A keyword cannot begin with ::
func NewKeywordElement(parts ...string) (SymbolElement, error) {

	// remove the : symbol if it is the first character.
	switch len(parts) {
	case 0:
		return nil, MakeError(ErrInvalidKeyword, "0 len")

	default:
		if strings.HasPrefix(parts[0], KeywordPrefix) {
			parts[0] = strings.TrimPrefix(parts[0], KeywordPrefix)
		}

		// Per the symbol rules above, :/ and :/anything are not legal keywords.
		if strings.HasPrefix(parts[0], SymbolSeparator) {
			return nil, MakeError(ErrInvalidKeyword, "found ':/'")
		}
	}

	symbol, err := NewSymbolElement(parts...)
	if err != nil {
		if ErrInvalidSymbol.IsEquivalent(err) {
			if myErr, is := err.(*Error); is {
				return nil, MakeErrorWithFormat(ErrInvalidKeyword, "msg: %s - %s", myErr.message, myErr.details)
			}
		}
		return nil, err
	}

	impl := symbol.(*symbolElemImpl)
	impl.baseElemImpl.elemType = KeywordType
	impl.modifier = KeywordPrefix

	return impl, nil
}
