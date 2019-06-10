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

const (

	// VectorStartLiteral is the start of an EDN group element.
	VectorStartLiteral = "["

	// VectorEndLiteral is the end of an EDN group element.
	VectorEndLiteral = "]"

	// GroupingSeparatorLiteral is the separator between item in a collection
	VectorSeparatorLiteral = " "
)

func initVector(lexer Lexer) (err error) {

	lexer.AddCollectionPattern(VectorStartLiteral, VectorEndLiteral, func(tag string, elements []Element) (el Element, e error) {
		if el, e = NewVector(elements...); e == nil {
			e = el.SetTag(tag)
		}
		return el, e
	})

	return err
}

// NewVector creates a new vector
func NewVector(elements ...Element) (elem CollectionElement, err error) {

	// check for errors
	for _, child := range elements {
		if child == nil {
			err = MakeError(ErrInvalidElement, "nil child")
			break
		}
	}

	if err == nil {
		coll := &collectionElemImpl{
			startSymbol:     VectorStartLiteral,
			endSymbol:       VectorEndLiteral,
			separatorSymbol: VectorSeparatorLiteral,
			collection:      []Element{},
		}

		var base *baseElemImpl
		if base, err = baseFactory().make(coll, VectorType, collectionSerialization(false)); err == nil {
			coll.baseElemImpl = base
			elem = coll
			err = elem.Append(elements...)
		}
	}

	return elem, err
}
