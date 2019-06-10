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

// init will add the element factory to the collection of factories
func initKeyword(lexer Lexer) (err error) {
	if err = addElementTypeFactory(KeywordType, func(input interface{}) (elem Element, e error) {
		if v, ok := input.(string); ok {
			elem, e = NewKeywordElement(v)
		} else {
			e = MakeErrorWithFormat(ErrInvalidInput, "Value: %#v", input)
		}
		return elem, e
	}); err == nil {
		lexer.AddPattern(SymbolPrimitive, ":([*!?$%&=<>]|\\w)([-+*!?$%&=<>.#]|\\w)*(/([-+*!?$%&=<>.#]|\\w)*)?", func(tag string, tokenValue string) (el Element, e error) {
			tokenValue = strings.TrimSuffix(tokenValue, KeywordPrefix)
			if el, e = NewKeywordElement(tokenValue); e == nil {
				e = el.SetTag(tag)
			}

			return el, e
		})
	}

	return err
}

// NewKeywordElement creates a new character element or an error.
//
// Keywords are identifiers that typically designate themselves. They are semantically akin to enumeration values.
// Keywords follow the rules of symbols, except they can (and must) begin with :, e.g. :fred or :my/fred. If the target
// platform does not have a keyword type distinct from a symbol type, the same type can be used without conflict, since
// the mandatory leading : of keywords is disallowed for symbols. Per the symbol rules above, :/ and :/anything are not
// legal keywords. A keyword cannot begin with ::
func NewKeywordElement(parts ...string) (elem SymbolElement, err error) {

	// remove the : symbol if it is the first character.
	switch len(parts) {
	case 0:
		err = MakeError(ErrInvalidKeyword, "0 len")

	default:
		if strings.HasPrefix(parts[0], KeywordPrefix) {
			parts[0] = strings.TrimPrefix(parts[0], KeywordPrefix)
		}

		// Per the symbol rules above, :/ and :/anything are not legal keywords.
		if strings.HasPrefix(parts[0], SymbolSeparator) {
			err = MakeError(ErrInvalidKeyword, "found ':/'")
		}
	}

	if err == nil {
		var symbol SymbolElement
		if symbol, err = NewSymbolElement(parts...); err == nil {

			impl := symbol.(*symbolElemImpl)
			impl.baseElemImpl.elemType = KeywordType
			impl.modifier = KeywordPrefix

			elem = impl
		}
	}

	if ErrInvalidSymbol.IsEquivalent(err) {
		if myErr, is := err.(*Error); is {
			err = MakeErrorWithFormat(ErrInvalidKeyword, "msg: %s - %s", myErr.message, myErr.details)
		}
	}

	return elem, err
}
