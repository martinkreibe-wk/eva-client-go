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

	// ListStartLiteral is the start of an EDN group element.
	ListStartLiteral = "("

	// ListEndLiteral is the end of an EDN group element.
	ListEndLiteral = ")"

	// ListSeparatorLiteral is the separator between item in a collection
	ListSeparatorLiteral = " "
)

// NewList creates a new group
func NewList(elements ...Element) (elem CollectionElement, err error) {

	// check for errors
	for _, child := range elements {
		if child == nil {
			err = MakeError(ErrInvalidElement, "nil child")
			break
		}
	}

	if err == nil {
		coll := &collectionElemImpl{
			startSymbol:     ListStartLiteral,
			endSymbol:       ListEndLiteral,
			separatorSymbol: ListSeparatorLiteral,
			collection:      []Element{},
		}

		var base *baseElemImpl
		if base, err = baseFactory().make(coll, ListType, collectionSerialization(false)); err == nil {
			coll.baseElemImpl = base
			elem = coll
			err = elem.Append(elements...)
		}
	}

	return elem, err
}
