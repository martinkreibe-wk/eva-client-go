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
	"io"
	"strings"
)

const (
	ErrParserError = ErrorMessage("Parser error")
)

// globalLexer holds the global lexer.
var globalLexer Lexer

// getLexer returns the global lexer or an error.
func getLexer() (lexer Lexer, err error) {

	if globalLexer == nil {
		globalLexer, err = newLexer()
	}

	return globalLexer, err
}

// Parse the string into an edn element.
func ParseString(data string) (elem Element, err error) {
	return Parse(strings.NewReader(data))
}

// ParseCollection will parse a collection.
func ParseCollectionString(data string) (elem CollectionElement, err error) {
	return ParseCollection(strings.NewReader(data))
}

// Parse the string into an edn element.
func Parse(data io.Reader) (elem Element, err error) {

	var lex Lexer
	if lex, err = getLexer(); err == nil {
		elem, err = lex.Parse(data)
	}
	return elem, err
}

// ParseCollection will parse a collection.
func ParseCollection(data io.Reader) (elem CollectionElement, err error) {

	var rawElem Element
	if rawElem, err = Parse(data); err == nil {
		if rawElem.ElementType().IsCollection() {
			elem = rawElem.(CollectionElement)
		} else {
			err = MakeErrorWithFormat(ErrParserError, "Parsed an element, but was not a collection")
		}
	}

	return elem, err
}
