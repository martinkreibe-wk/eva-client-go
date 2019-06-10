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

	// ErrInvalidPair defines the error for pair issues
	ErrInvalidPair = ErrorMessage("Invalid pair")
)

// Pair of elements
type Pair interface {

	// Key the key part of the pair
	Key() Element

	// Value the value part of the pair
	Value() Element
}

// pairImpl implements the pair interface
type pairImpl struct {
	key   Element
	value Element
}

// Key the key part of the pair
func (pair *pairImpl) Key() Element {
	return pair.key
}

// Value the value part of the pair
func (pair *pairImpl) Value() Element {
	return pair.value
}

// NewPair creates a new pair from the key and value supplied.
func NewPair(key, value interface{}) (pair Pair, err error) {

	if key != nil && value != nil {
		var k, v Element
		if k, err = NewPrimitiveElement(key); err == nil {
			if v, err = NewPrimitiveElement(value); err == nil {
				pair = &pairImpl{
					key:   k,
					value: v,
				}
			}
		}
	} else {
		err = MakeErrorWithFormat(ErrInvalidPair, "Key is nil [%t] or value is nil [%t]", key == nil, value == nil)
	}

	return pair, err
}

// Pairs defines a collection of pairs.
type Pairs struct {
	data []Pair
}

// Append will append the entities
func (pairs *Pairs) Append(key, value Element) (err error) {

	var pair Pair
	if pair, err = NewPair(key, value); err == nil {
		pairs.AppendPair(pair)
	}

	return err
}

// AppendPair will append the pair
func (pairs *Pairs) AppendPair(pair Pair) {
	pairs.data = append(pairs.data, pair)
}

// Raw returns the internal pair collection
func (pairs *Pairs) Raw() []Pair {
	return pairs.data
}

// RawElements returns the internal element collection
func (pairs *Pairs) RawElements() []Element {

	var elems []Element
	for _, val := range pairs.Raw() {
		elems = append(elems, val.Key(), val.Value())
	}

	return elems
}

// Len returns the pair collection length
func (pairs *Pairs) Len() int {
	return len(pairs.data)
}
