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

func parseSymbol(tokenValue string) (Element, error) {
	return NewSymbolElement(tokenValue)
}

func fromSymbol(input interface{}) (Element, error) {
	v, ok := input.(string)
	if !ok {
		return nil, MakeErrorWithFormat(ErrInvalidInput, "Value: %#v", input)
	}
	return NewSymbolElement(v)
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
	if len(prefix) == 0 {
		return name
	}

	return fmt.Sprintf("%s%s%s", prefix, SymbolSeparator, name)
}

func decodeSymbol(parts ...string) (string, string, error) {
	switch len(parts) {
	case 1:

		name := parts[0]
		switch {

		// handle the case where the name was really sent in with the separator
		case name == SymbolSeparator:
			return "", name, nil

		case strings.Contains(name, SymbolSeparator):
			if parts = strings.Split(name, SymbolSeparator); len(parts) != 2 {
				return "", "", MakeErrorWithFormat(ErrInvalidSymbol, "Name[1]: %#v", parts)
			}

			prefix := parts[0]
			if len(prefix) == 0 || !symbolMatcher(prefix) {
				return "", "", MakeErrorWithFormat(ErrInvalidSymbol, "Prefix[0]: %#v", parts)
			}

			if name = parts[1]; len(name) == 0 || !symbolMatcher(name) {
				return "", "", MakeErrorWithFormat(ErrInvalidSymbol, "Name[0]: %#v", parts)
			}

			return prefix, name, nil
		default:
			if !symbolMatcher(name) {
				return "", "", MakeErrorWithFormat(ErrInvalidSymbol, "Invalid Name: %#v", parts)
			}
			return "", name, nil
		}

	case 2:

		prefix := parts[0]
		if len(prefix) != 0 && !symbolMatcher(prefix) {
			return "", "", MakeErrorWithFormat(ErrInvalidSymbol, "Prefix[2]: %#v", parts)
		}

		name := parts[1]
		if !symbolMatcher(name) {
			return "", "", MakeErrorWithFormat(ErrInvalidSymbol, "Prefix[1]: %#v", parts)
		}

		return prefix, name, nil
	default:
		return "", "", MakeError(ErrInvalidSymbol, parts)
	}
}

func symbolSerializer(serializer Serializer, tag string, value interface{}) (string, error) {
	switch serializer.MimeType() {
	case EvaEdnMimeType:
		var out string
		if len(tag) > 0 {
			out = TagPrefix + tag + " "
		}

		if elem, is := value.(SymbolElement); is {
			return out + elem.AppendNameOntoNamespace(elem.Name()), nil
		}

		return out, nil
	default:
		return "", MakeError(ErrUnknownMimeType, serializer.MimeType())
	}
}

func symbolEqual(left, right Element) bool {
	leftSym, is := left.(SymbolElement)
	if !is {
		return false
	}

	rightSym, is := right.(SymbolElement)
	if !is {
		return false
	}

	return leftSym.Name() == rightSym.Name() &&
		leftSym.Prefix() == rightSym.Prefix() &&
		leftSym.Modifier() == rightSym.Modifier()
}

// NewSymbolElement creates a new character element or an error.
func NewSymbolElement(parts ...string) (SymbolElement, error) {

	prefix, name, err := decodeSymbol(parts...)
	if err != nil {
		return nil, err
	}

	symElem := &symbolElemImpl{
		prefix: prefix,
		name:   name,
	}

	base, err := baseFactory().make(symElem, SymbolType, NoTag, symbolSerializer)
	if err != nil {
		return nil, err
	}

	symElem.baseElemImpl = base

	// equality for symbols are different then the normal path.
	symElem.baseElemImpl.equality = symbolEqual

	return symElem, nil
}

// AppendNameOntoNamespace will append the input name onto the namespace.
func (elem *symbolElemImpl) AppendNameOntoNamespace(name string) string {
	return elem.Modifier() + encodeSymbol(elem.Prefix(), name)
}

// Equals checks if the input element is equal to this element.
func (elem *symbolElemImpl) Equals(e Element) bool {
	if elem.ElementType() != e.ElementType() {
		return false
	}

	if elem.Tag() != e.Tag() {
		return false
	}

	return elem.baseElemImpl.equality(elem, e)
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
