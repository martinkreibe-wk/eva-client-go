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

	// MapStartLiteral is the start of an EDN group element.
	MapStartLiteral = "{"

	// MapEndLiteral is the end of an EDN group element.
	MapEndLiteral = "}"

	// MapSeparatorLiteral is the separator between item in a collection
	MapSeparatorLiteral = ", "

	// MapKeyValueSeparatorLiteral is the separator for keys and values
	MapKeyValueSeparatorLiteral = " "

	// ErrDuplicateKey defines the duplicate key error
	ErrDuplicateKey = ErrorMessage("Duplicate key found")
)

// NewMap creates a new vector
func NewMap(pairs ...Pair) (elem CollectionElement, err error) {

	coll := &collectionElemImpl{
		startSymbol:             MapStartLiteral,
		endSymbol:               MapEndLiteral,
		separatorSymbol:         MapSeparatorLiteral,
		keyValueSeparatorSymbol: MapKeyValueSeparatorLiteral,
		collection:              map[string][2]Element{}, // { serialized_key, [key, value] }
	}

	var base *baseElemImpl
	if base, err = baseFactory().make(coll, MapType, collectionSerialization(true)); err == nil {
		coll.baseElemImpl = base

		// check for errors
		keys := make([]Element, 0)
		for _, pair := range pairs {
			if pair == nil || pair.Key() == nil {
				err = MakeError(ErrInvalidPair, "nil pair or nil key")
			} else {

				key := pair.Key()
				for _, k := range keys {

					if key.Equals(k) {
						err = MakeErrorWithFormat(ErrDuplicateKey, "%s = %s", k, key)
						break
					}
				}

				if err == nil {
					keys = append(keys, key)
					err = coll.Append(key, pair.Value())
				}
			}

			if err != nil {
				break
			}
		}

		if err == nil {
			elem = coll
		}
	}

	return elem, err
}
