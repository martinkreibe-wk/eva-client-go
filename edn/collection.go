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
	"strconv"
)

const (

	// ErrNoValue is returned when no value is found in the collection.
	ErrNoValue = ErrorMessage("No value found")
)

// ChildIterator is the iterator for children elements, if this is a list based item, the key will be an index of the
// current item, the value is mapped value to the key. To stop a loop mid iteration, set the error to non nil.
type ChildIterator func(key Element, value Element) (err error)

// GroupElement defines the element for the EDN grouping construct. A group is a sequence of values. Groups are
// represented by zero or more elements enclosed in parentheses (). Note that Group can be heterogeneous.
type CollectionElement interface {
	Element

	// Len return the quantity of items in this collection.
	Len() int

	// IterateChildren will iterate over all children.
	IterateChildren(iterator ChildIterator) (err error)

	// Append the elements into this collection.
	Append(children ...Element) (err error)

	// Prepend the elements into this collection.
	Prepend(children ...Element) (err error)

	// Get the key from the collection.
	Get(key interface{}) (Element, error)

	// Merge one collection into another.
	Merge(CollectionElement) error
}

// collectionElemImpl is the implementation to the GroupElement interface.
type collectionElemImpl struct {
	*baseElemImpl

	// startSymbol defines the start symbol
	startSymbol string

	// endSymbol defines the end symbol
	endSymbol string

	// separatorSymbol for the elements
	separatorSymbol string

	// keyValueSeparatorSymbol for the element
	keyValueSeparatorSymbol string

	// collection of elements
	collection interface{}
}

// Len return the quantity of items in this collection.
func (elem *collectionElemImpl) Len() (l int) {
	switch v := elem.collection.(type) {
	case []Element:
		l = len(v)
	case map[string][2]Element:
		l = len(v)
	}
	return l
}

// IterateChildren will iterate over the child elements within this collection.
func (elem *collectionElemImpl) IterateChildren(iterator ChildIterator) (err error) {
	switch v := elem.collection.(type) {
	case []Element:
		for i, c := range v {
			iElem, _ := NewIntegerElement(int64(i))
			if err = iterator(iElem, c); err != nil {
				break
			}
		}
	case map[string][2]Element:
		for _, c := range v {
			err = iterator(c[0], c[1])
			if err != nil {
				break
			}
		}
	}
	return err
}

// Equals checks if the input element is equal to this element.
func (elem *collectionElemImpl) Equals(e Element) (result bool) {
	if elem.ElementType() == e.ElementType() {
		if elem.Tag() == e.Tag() {
			if other := e.(*collectionElemImpl); elem.Len() == other.Len() {
				if elem.Len() > 0 {
					switch v := elem.collection.(type) {
					case []Element:
						otherChildren := other.collection.([]Element)

						for index, child := range v {
							// if the children are different then we don't need to look any more.
							if result = child.Equals(otherChildren[index]); !result {
								break
							}
						}

					case map[string][2]Element:
						otherChildren := other.collection.(map[string][2]Element)

						for key, child := range v {

							// if the children are different then we don't need to look any more.
							if otherChild, has := otherChildren[key]; has {
								if result = child[1].Equals(otherChild[1]); !result {
									break
								}
							} else {
								result = false
								break
							}
						}
					}
				} else {
					// both lengths are 0
					result = true
				}
			}
		}
	}
	return result
}

func (elem *collectionElemImpl) add(addFunc func([]Element, []Element) []Element, children []Element) (err error) {

	if len(children) != 0 {
		switch v := elem.collection.(type) {
		case []Element:
			elem.collection = addFunc(v, children)
		case map[string][2]Element:

			setSize := 2 // This is for maps...
			if len(elem.keyValueSeparatorSymbol) == 0 {
				setSize = 1 // This is for sets...
			}

			if len(children)%setSize == 0 {
				childOffset := setSize - 1

				for i := 0; i < len(children); i += setSize {
					var k string
					if str, is := children[i].(*baseElemImpl); is && str.elemType == StringType {
						k = str.value.(string)
					} else {
						stream := NewStringStream()
						err = EvaEdnMimeType.SerializeTo(stream, children[i])
						k = stream.String()
					}

					if err == nil {
						if _, has := v[k]; !has {
							v[k] = [2]Element{children[i], children[i+childOffset]}
						} else {
							err = MakeErrorWithFormat(ErrDuplicateKey, "Key: %s", k)
						}
					}

					if err != nil {
						break
					}
				}
			} else {
				err = MakeError(ErrInvalidInput, "must have an even number of inputs.")
			}

		default:
			err = MakeErrorWithFormat(ErrInvalidElement, "type: %T", v)
		}
	}

	return err
}

// Append will add the appropriate children. Note that a map must have 2 parameters.
func (elem *collectionElemImpl) Append(children ...Element) error {
	return elem.add(func(in []Element, c []Element) []Element {
		return append(in, c...)
	}, children)
}

// Prepend will add the appropriate children. Note that a map must have 2 parameters.
func (elem *collectionElemImpl) Prepend(children ...Element) error {
	return elem.add(func(in []Element, c []Element) []Element {
		return append(c, in...)
	}, children)
}

// Get the value from the collection.
func (elem *collectionElemImpl) Get(key interface{}) (value Element, err error) {

	var realKey string

	switch k := key.(type) {
	case int, int32, int64:
		realKey = fmt.Sprintf("%d", k)
	case string:
		realKey = k
	case Element:
		stream := &StringStream{}
		err = EvaEdnMimeType.SerializeTo(stream, k)
		if err != nil {
			return nil, err
		}
		realKey = stream.String()
	default:
		err = MakeErrorWithFormat(ErrInvalidInput, "key type: %T", k)
	}

	if err == nil {
		switch v := elem.collection.(type) {
		case []Element:
			var index int
			if index, err = strconv.Atoi(realKey); err == nil {
				if index >= 0 && index < len(v) {
					value = v[index]
				} else {
					err = MakeError(ErrNoValue, realKey)
				}
			}
		case map[string][2]Element:
			var has bool
			var pair [2]Element
			if pair, has = v[realKey]; has {
				value = pair[1]
			} else {
				err = MakeError(ErrNoValue, realKey)
			}
		default:
			err = MakeErrorWithFormat(ErrInvalidElement, "type: %T", v)
		}
	}

	return value, err
}

// Merge one collection into another.
func (elem *collectionElemImpl) Merge(child CollectionElement) (err error) {
	if child != nil {
		switch elem.collection.(type) {
		case []Element:
			err = child.IterateChildren(func(_ Element, child Element) error {
				return elem.Append(child)
			})
		case map[string][2]Element:
			err = child.IterateChildren(func(key Element, child Element) error {
				return elem.Append(key, child)
			})
		}
	}

	return err
}
