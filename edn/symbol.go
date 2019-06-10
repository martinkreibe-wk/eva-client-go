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
	"fmt"
	"regexp"
	"strings"
)

const (

	// SymbolSeparator defines the symbol for separating the prefix with the name. If there is no separator, then the
	// symbol is just a name value.
	SymbolSeparator = "/"

	// NamespaceSeparator defines the symbol for separating the namespaces from each other. Note that this is not the
	// same as the namespace/name separator (SymbolSeparator).
	NamespaceSeparator = "."

	// ErrInvalidSymbol defines an invalid symbol
	ErrInvalidSymbol = ErrorMessage("Invalid Symbol")

	// symbols that can modify a numeric
	numericModifierSymbols = `\.|\+|-`

	// symbols other then alphanumeric and numeric modifiers that are legal
	legalFirstSymbols = `\*|!|_|\?|\$|%|&|=|<|>`

	// symbols that are marked as not being allowed to be first characters other then numeric
	specialSymbols = KeywordPrefix + `|` + TagPrefix

	// symbolRegex defines the valid symbols.
	// Symbols begin with a non-numeric character and can contain alphanumeric characters and . * + ! - _ ? $ % & = < >.
	// If -, + or . are the first character, the second character (if any) must be non-numeric. Additionally, : # are
	// allowed as constituent characters in symbols other than as the first character.
	symbolRegex = `^((` + numericModifierSymbols + `)|((((` + numericModifierSymbols + `)(` + legalFirstSymbols + `|[[:alpha:]]))|(` + legalFirstSymbols + `|[[:alpha:]]))+(` + numericModifierSymbols + `|` + legalFirstSymbols + `|` + specialSymbols + `|[[:alnum:]])*))$`
)

// init will add the element factory to the collection of factories
func initSymbol(lexer Lexer) (err error) {
	lexer.AddPattern(SymbolPrimitive, "[*!?$%&=<>_a-zA-Z.]([-+*!?$%&=<>_.#]|\\w)*(/([-+*!?$%&=<>_.#]|\\w)*)?", func(tag string, tokenValue string) (el Element, e error) {
		if el, e = NewSymbolElement(tokenValue); e == nil {
			e = el.SetTag(tag)
		}
		return el, e
	})

	return err
}

// symbolMatcher is the matching mechanism for symbols
var symbolMatcher = regexp.MustCompile(symbolRegex).MatchString

// IsValidNamespace checks if the namespace is valid.
func IsValidNamespace(namespace string) bool {
	return symbolMatcher(namespace)
}

// Symbols are used to represent identifiers, and should map to something other than strings, if possible.
type SymbolElement interface {
	Element

	// Modifier for this symbol
	Modifier() string

	// Prefix to this symbol
	Prefix() string

	// Name to this symbol
	Name() string

	// AppendNameOntoNamespace will append the input name onto the namespace.
	AppendNameOntoNamespace(string) string
}

// symbolElemImpl implements the symbolElemImpl
type symbolElemImpl struct {
	*baseElemImpl
	prefix   string
	name     string
	modifier string
}

func encodeSymbol(prefix string, name string) string {
	if len(prefix) > 0 {
		name = fmt.Sprintf("%s%s%s", prefix, SymbolSeparator, name)
	}
	return name
}

func decodeSymbol(parts ...string) (prefix string, name string, err error) {
	switch len(parts) {
	case 1:

		switch name = parts[0]; {

		// handle the case where the name was really sent in with the separator
		case name == SymbolSeparator:
			// Fine, break

		case strings.Contains(name, SymbolSeparator):
			if parts = strings.Split(name, SymbolSeparator); len(parts) == 2 {
				if prefix = parts[0]; len(prefix) != 0 && symbolMatcher(prefix) {
					if name = parts[1]; len(name) == 0 || !symbolMatcher(name) {
						err = MakeErrorWithFormat(ErrInvalidSymbol, "Name[0]: %#v", parts)
					}
				} else {
					err = MakeErrorWithFormat(ErrInvalidSymbol, "Prefix[0]: %#v", parts)
				}
			} else {
				err = MakeErrorWithFormat(ErrInvalidSymbol, "Name[1]: %#v", parts)
			}
		default:
			if !symbolMatcher(name) {
				err = MakeErrorWithFormat(ErrInvalidSymbol, "Invalid Name: %#v", parts)
			}
		}

	case 2:
		if prefix = parts[0]; len(prefix) != 0 && symbolMatcher(prefix) {
			if name = parts[1]; !symbolMatcher(name) {
				err = MakeErrorWithFormat(ErrInvalidSymbol, "Prefix[1]: %#v", parts)
			}
		} else {
			err = MakeErrorWithFormat(ErrInvalidSymbol, "Prefix[2]: %#v", parts)
		}
	default:
		err = MakeError(ErrInvalidSymbol, parts)
	}

	return prefix, name, err
}

// NewSymbolElement creates a new character element or an error.
func NewSymbolElement(parts ...string) (elem SymbolElement, err error) {

	var prefix string
	var name string
	if prefix, name, err = decodeSymbol(parts...); err == nil {

		symElem := &symbolElemImpl{
			prefix: prefix,
			name:   name,
		}

		var base *baseElemImpl
		if base, err = baseFactory().make(symElem, SymbolType, func(serializer Serializer, tag string, value interface{}) (out string, err error) {
			switch serializer.MimeType() {
			case EvaEdnMimeType:
				if len(tag) > 0 {
					out = TagPrefix + tag + " "
				}
				if elem, ok := value.(SymbolElement); ok {
					out += elem.AppendNameOntoNamespace(elem.Name())
				}
			default:
				err = MakeError(ErrUnknownMimeType, serializer.MimeType())
			}

			return out, err
		}); err == nil {

			symElem.baseElemImpl = base

			// equality for symbols are different then the normal path.
			symElem.baseElemImpl.equality = func(left, right Element) (result bool) {
				if leftSym, has := left.(SymbolElement); has {
					if rightSym, has := right.(SymbolElement); has {
						if leftSym.Name() == rightSym.Name() && leftSym.Prefix() == rightSym.Prefix() && leftSym.Modifier() == rightSym.Modifier() {
							result = true
						}
					}
				}

				return result
			}

			elem = symElem
		}
	}

	return elem, err
}

// AppendNameOntoNamespace will append the input name onto the namespace.
func (elem *symbolElemImpl) AppendNameOntoNamespace(name string) (out string) {
	return elem.Modifier() + encodeSymbol(elem.Prefix(), name)
}

// Equals checks if the input element is equal to this element.
func (elem *symbolElemImpl) Equals(e Element) (result bool) {
	if elem.ElementType() == e.ElementType() {
		if elem.Tag() == e.Tag() {
			result = elem.baseElemImpl.equality(elem, e)
		}
	}
	return result
}

// Prefix to this symbol
func (elem *symbolElemImpl) Prefix() string {
	return elem.prefix
}

// Name to this symbol
func (elem *symbolElemImpl) Name() string {
	return elem.name
}

// Modifier for this symbol
func (elem *symbolElemImpl) Modifier() string {
	return elem.modifier
}
