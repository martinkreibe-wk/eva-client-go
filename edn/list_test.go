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

var _ = Describe("List in EDN", func() {
	Context("with the default marshaller", func() {
		It("should create an empty group with no error", func() {
			group, err := NewList()
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(ListType))
			Ω(group.Len()).Should(BeEquivalentTo(0))
		})

		It("should serialize an empty list correctly", func() {
			group, err := NewList()
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("()"))
		})

		It("should serialize an empty list correctly", func() {
			group, err := NewList()
			Ω(err).Should(BeNil())

			_, err = group.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})

		It("should error with a nil item", func() {
			group, err := NewList(nil)
			Ω(err).Should(test.HaveMessage(ErrInvalidElement))
			Ω(group).Should(BeNil())
		})

		It("should create a list element with the initial values", func() {
			elem := NewStringElement("foo")

			group, err := NewList(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(ListType))
			Ω(group.Len()).Should(BeEquivalentTo(1))
		})

		It("should be able to append", func() {
			elem := NewStringElement("foo")
			elem2 := NewStringElement("bar")

			group, err := NewList(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(ListType))
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

			group, err := NewList(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(ListType))
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

		It("should serialize a single nil entry in a list correctly", func() {
			elem := NewNilElement()

			group, err := NewList(elem)
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("(nil)"))
		})

		It("should serialize some nil entries in a list correctly", func() {
			var elem1, elem2, elem3 Element
			var group CollectionElement
			var err error

			elem1 = NewStringElement("foo")
			elem2 = NewStringElement("bar")
			elem3 = NewStringElement("faz")

			group, err = NewList(elem1, elem2, elem3)
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("(\"foo\" \"bar\" \"faz\")"))

			breakCount := 2
			templateError := ErrorMessage("This is the expected error")
			err = group.IterateChildren(func(key, value Element) (e error) {
				if breakCount--; breakCount == 0 {
					e = MakeError(templateError, 0)
				}
				return e
			})

			Ω(err).Should(test.HaveMessage(templateError))
		})
	})

	Context("Parsing", func() {
		runParserTests(ListType,
			&testDefinition{"()", func() (elements map[string]Element, err error) {
				return elements, err
			}},
			&testDefinition{"(\"()\")", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": NewStringElement("()"),
				}
				return elements, err
			}},
			&testDefinition{"(1)", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": NewIntegerElement(1),
				}
				return elements, err
			}},
			&testDefinition{"(1 2 3)", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": NewIntegerElement(1),
					"1": NewIntegerElement(2),
					"2": NewIntegerElement(3),
				}
				return elements, err
			}},
			&testDefinition{"(#foo 1 2 #bar 3)", func() (elements map[string]Element, err error) {

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
			&testDefinition{"(())", func() (elements map[string]Element, err error) {
				var subList CollectionElement
				if subList, err = NewList(); err == nil {
					elements = map[string]Element{
						"0": subList,
					}
				}
				return elements, err
			}},
			&testDefinition{"(\"a\" ())", func() (elements map[string]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewList(); err == nil {
					elements = map[string]Element{
						"0": NewStringElement("a"),
						"1": subList1,
					}
				}
				return elements, err
			}},
			&testDefinition{"(() \"a\")", func() (elements map[string]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewList(); err == nil {
					elements = map[string]Element{
						"0": subList1,
						"1": NewStringElement("a"),
					}
				}
				return elements, err
			}},
			&testDefinition{"(#foo () #bar ())", func() (elements map[string]Element, err error) {
				var subList1 CollectionElement
				var subList2 CollectionElement
				if subList1, err = NewList(); err == nil {
					if subList2, err = NewList(); err == nil {
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
