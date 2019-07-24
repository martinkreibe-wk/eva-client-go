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

	// SetStartLiteral is the start of an EDN group element.
	SetStartLiteral = "#{"

	// SetEndLiteral is the end of an EDN group element.
	SetEndLiteral = "}"

	// GroupingSeparatorLiteral is the separator between item in a collection
	SetSeparatorLiteral = " "
)

// NewSet creates a new vector
func NewSet(elements ...Element) (elem CollectionElement, err error) {

	// check for errors
	for _, child := range elements {
		if child == nil {
			err = MakeError(ErrInvalidElement, "nil child")
			break
		}
	}

	if err == nil {
		coll := &collectionElemImpl{
			startSymbol:     SetStartLiteral,
			endSymbol:       SetEndLiteral,
			separatorSymbol: SetSeparatorLiteral,
			collection:      map[string][2]Element{},
		}

		var base *baseElemImpl
		if base, err = baseFactory().make(coll, SetType, collectionSerialization(false)); err == nil {
			coll.baseElemImpl = base
			if err = coll.Append(elements...); err == nil {
				elem = coll
			}
		}
	}

	return elem, err
}
