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
	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vector in EDN", func() {

	Context("with the default marshaller", func() {
		It("should create an empty vector with no error", func() {
			group, err := NewVector()
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(VectorType))
			Ω(group.Len()).Should(BeEquivalentTo(0))
		})

		It("should serialize an empty vector correctly", func() {
			group, err := NewVector()
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("[]"))
		})

		It("should serialize an empty vector correctly", func() {
			group, err := NewVector()
			Ω(err).Should(BeNil())

			_, err = group.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})

		It("should error with a nil item", func() {
			group, err := NewVector(nil)
			Ω(err).Should(test.HaveMessage(ErrInvalidElement))
			Ω(group).Should(BeNil())
		})

		It("should create a vector element with the initial values", func() {
			elem := NewStringElement("foo")

			group, err := NewVector(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(VectorType))
			Ω(group.Len()).Should(BeEquivalentTo(1))

			var data interface{}
			data, err = group.Get(42)
			Ω(data).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrNoValue))
		})

		It("should serialize a single nil entry in a vector correctly", func() {
			elem := NewNilElement()

			group, err := NewVector(elem)
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("[nil]"))
		})

		It("should serialize some nil entries in a vector correctly", func() {
			elem1 := NewStringElement("foo")
			elem2 := NewStringElement("bar")
			elem3 := NewStringElement("faz")

			group, err := NewVector(elem1, elem2, elem3)
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("[\"foo\" \"bar\" \"faz\"]"))

			index := NewIntegerElement(0)
			v, err := group.Get(index)
			Ω(err).Should(BeNil())
			Ω(v).ShouldNot(BeNil())

			v, err = group.Get(0)
			Ω(err).Should(BeNil())
			Ω(v).ShouldNot(BeNil())

			v, err = group.Get(&struct{}{})
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
			Ω(v).Should(BeNil())
		})

		It("should be able to append", func() {
			elem := NewStringElement("foo")
			elem2 := NewStringElement("bar")

			group, err := NewVector(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(VectorType))
			Ω(group.Len()).Should(BeEquivalentTo(1))

			group.Append(elem2)
			Ω(group.Len()).Should(BeEquivalentTo(2))

			e1, err := group.Get(0)
			Ω(err).Should(BeNil())
			Ω(e1.String()).Should(BeEquivalentTo(elem.String()))

			e2, err := group.Get(1)
			Ω(err).Should(BeNil())
			Ω(e2.String()).Should(BeEquivalentTo(elem2.String()))
		})

		It("should be able to prepend", func() {
			elem := NewStringElement("foo")
			elem2 := NewStringElement("bar")

			group, err := NewVector(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(VectorType))
			Ω(group.Len()).Should(BeEquivalentTo(1))

			group.Prepend(elem2)
			Ω(group.Len()).Should(BeEquivalentTo(2))

			e1, err := group.Get(1)
			Ω(err).Should(BeNil())
			Ω(e1.String()).Should(BeEquivalentTo(elem.String()))

			e2, err := group.Get(0)
			Ω(err).Should(BeNil())
			Ω(e2.String()).Should(BeEquivalentTo(elem2.String()))
		})

		It("should serialize some nil entries in a vector correctly", func() {
			elem1 := NewStringElement("foo")
			elem2 := NewStringElement("bar")
			elem3 := NewStringElement("faz")
			elem4 := NewStringElement("baz")

			group, err := NewVector(elem1, elem2)
			Ω(err).Should(BeNil())

			group2, err := NewVector(elem3, elem4)
			Ω(err).Should(BeNil())

			err = group.Merge(group2)
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("[\"foo\" \"bar\" \"faz\" \"baz\"]"))
		})

		It("should Handle Equality right with tags", func() {
			// #db/id [-9223363240760753529 -1000074] = #db/id [-9223372036853775739 -1000075]
			elem1a := NewIntegerElement(-9223363240760753529)
			elem1b := NewIntegerElement(-1000074)

			elem2a := NewIntegerElement(-9223372036853775739)
			elem2b := NewIntegerElement(-1000075)

			group1, err := NewVector(elem1a, elem1b)
			group1.SetTag("db/id")
			Ω(err).Should(BeNil())

			group2, err := NewVector(elem2a, elem2b)
			group2.SetTag("db/id")
			Ω(err).Should(BeNil())

			Ω(group1.Equals(group2)).Should(BeFalse())
		})

		It("should Handle Equality right without tags", func() {
			// #db/id [-9223363240760753529 -1000074] = #db/id [-9223372036853775739 -1000075]
			elem1a := NewIntegerElement(-9223363240760753529)
			elem1b := NewIntegerElement(-1000074)

			elem2a := NewIntegerElement(-9223372036853775739)
			elem2b := NewIntegerElement(-1000075)

			group1, err := NewVector(elem1a, elem1b)
			Ω(err).Should(BeNil())

			group2, err := NewVector(elem2a, elem2b)
			Ω(err).Should(BeNil())

			Ω(group1.Equals(group2)).Should(BeFalse())
		})
	})

	Context("Parsing", func() {
		runParserTests(VectorType,
			&testDefinition{"[]", func() (elements map[string]Element, err error) {
				return elements, err
			}},
			&testDefinition{"[\"[]\"]", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": NewStringElement("[]"),
				}
				return elements, err
			}},
			&testDefinition{"[*]", func() (elements map[string]Element, err error) {
				var sym Element
				sym, err = NewSymbolElement("*")
				elements = map[string]Element{
					"0": sym,
				}
				return elements, err
			}},
			&testDefinition{"[1]", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": NewIntegerElement(1),
				}
				return elements, err
			}},
			&testDefinition{"[1 2 3]", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": NewIntegerElement(1),
					"1": NewIntegerElement(2),
					"2": NewIntegerElement(3),
				}
				return elements, err
			}},
			&testDefinition{"[#foo 1 2 #bar 3]", func() (elements map[string]Element, err error) {

				one := NewIntegerElement(1)
				three := NewIntegerElement(3)

				err = one.SetTag("foo")

				if err == nil {
					err = three.SetTag("bar")
				}

				if err == nil {
					elements = map[string]Element{
						"0": one,
						"1": NewIntegerElement(2),
						"2": three,
					}
				}
				return elements, err
			}},

			&testDefinition{"[[\"I just rewrote a docstring.\"]]", func() (elements map[string]Element, err error) {
				str := NewStringElement("I just rewrote a docstring.")
				var subList CollectionElement
				if subList, err = NewVector(str); err == nil {
					elements = map[string]Element{
						"0": subList,
					}
				}
				return elements, err
			}},

			&testDefinition{"[[]]", func() (elements map[string]Element, err error) {
				var subList CollectionElement
				if subList, err = NewVector(); err == nil {
					elements = map[string]Element{
						"0": subList,
					}
				}
				return elements, err
			}},
			&testDefinition{"[\"a\" []]", func() (elements map[string]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewVector(); err == nil {
					elements = map[string]Element{
						"0": NewStringElement("a"),
						"1": subList1,
					}
				}
				return elements, err
			}},
			&testDefinition{"[[] \"a\"]", func() (elements map[string]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewVector(); err == nil {
					elements = map[string]Element{
						"0": subList1,
						"1": NewStringElement("a"),
					}
				}
				return elements, err
			}},
			&testDefinition{"[#foo [] #bar []]", func() (elements map[string]Element, err error) {
				var subList1 CollectionElement
				var subList2 CollectionElement
				if subList1, err = NewVector(); err == nil {
					if subList2, err = NewVector(); err == nil {
						if err = subList1.SetTag("foo"); err == nil {
							if err = subList2.SetTag("bar"); err == nil {
								elements = map[string]Element{
									"0": subList1,
									"1": subList2,
								}
							}
						}
					}
				}
				return elements, err
			}},
		)
	})
})
