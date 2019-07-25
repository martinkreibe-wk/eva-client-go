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

			stream := NewStringStream()
			err = EvaEdnMimeType.SerializeTo(stream, group)
			edn := stream.String()
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("()"))
		})

		It("should error with a nil item", func() {
			group, err := NewList(nil)
			Ω(err).Should(test.HaveMessage(ErrInvalidElement))
			Ω(group).Should(BeNil())
		})

		It("should create a list element with the initial values", func() {
			elem, err := NewStringElement("foo")
			Ω(err).Should(BeNil())

			group, err := NewList(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(ListType))
			Ω(group.Len()).Should(BeEquivalentTo(1))
		})

		It("should be able to append", func() {
			elem, err := NewStringElement("foo")
			Ω(err).Should(BeNil())
			elem2, err := NewStringElement("bar")
			Ω(err).Should(BeNil())

			group, err := NewList(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(ListType))
			Ω(group.Len()).Should(BeEquivalentTo(1))

			err = group.Append(elem2)
			Ω(err).Should(BeNil())
			Ω(group.Len()).Should(BeEquivalentTo(2))

			e1, err := group.Get(0)
			Ω(err).Should(BeNil())
			Ω(e1.Value()).Should(BeEquivalentTo(elem.Value()))

			e2, err := group.Get(1)
			Ω(err).Should(BeNil())
			Ω(e2.Value()).Should(BeEquivalentTo(elem2.Value()))
		})

		It("should be able to prepend", func() {
			elem, err := NewStringElement("foo")
			Ω(err).Should(BeNil())
			elem2, err := NewStringElement("bar")
			Ω(err).Should(BeNil())

			group, err := NewList(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(ListType))
			Ω(group.Len()).Should(BeEquivalentTo(1))

			err = group.Prepend(elem2)
			Ω(err).Should(BeNil())
			Ω(group.Len()).Should(BeEquivalentTo(2))

			e1, err := group.Get(1)
			Ω(err).Should(BeNil())
			Ω(e1.Value()).Should(BeEquivalentTo(elem.Value()))

			e2, err := group.Get(0)
			Ω(err).Should(BeNil())
			Ω(e2.Value()).Should(BeEquivalentTo(elem2.Value()))
		})

		It("should serialize a single nil entry in a list correctly", func() {
			elem, err := NewNilElement()
			Ω(err).Should(BeNil())

			group, err := NewList(elem)
			Ω(err).Should(BeNil())

			stream := NewStringStream()
			err = EvaEdnMimeType.SerializeTo(stream, group)
			edn := stream.String()
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("(nil)"))
		})

		It("should serialize some nil entries in a list correctly", func() {
			var elem1, elem2, elem3 Element
			var group CollectionElement
			var err error

			elem1, err = NewStringElement("foo")
			Ω(err).Should(BeNil())
			elem2, err = NewStringElement("bar")
			Ω(err).Should(BeNil())
			elem3, err = NewStringElement("faz")
			Ω(err).Should(BeNil())

			group, err = NewList(elem1, elem2, elem3)
			Ω(err).Should(BeNil())

			stream := NewStringStream()
			err = EvaEdnMimeType.SerializeTo(stream, group)
			edn := stream.String()
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

		str, err := NewStringElement("()")
		if err != nil {
			panic(err)
		}

		a, err := NewStringElement("a")
		if err != nil {
			panic(err)
		}

		one, err := NewIntegerElement(1)
		if err != nil {
			panic(err)
		}

		two, err := NewIntegerElement(2)
		if err != nil {
			panic(err)
		}

		three, err := NewIntegerElement(3)
		if err != nil {
			panic(err)
		}

		runParserTests(ListType,
			&testDefinition{"()", func() (elements map[string]Element, err error) {
				return elements, err
			}},
			&testDefinition{"(\"()\")", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": str,
				}
				return elements, err
			}},
			&testDefinition{"(1)", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": one,
				}
				return elements, err
			}},
			&testDefinition{"(1 2 3)", func() (elements map[string]Element, err error) {
				elements = map[string]Element{
					"0": one,
					"1": two,
					"2": three,
				}
				return elements, err
			}},
			&testDefinition{"(#foo 1 2 #bar 3)", func() (elements map[string]Element, err error) {

				one, err := NewIntegerElement(1)
				Ω(err).Should(BeNil())
				three, err := NewIntegerElement(3)
				Ω(err).Should(BeNil())

				err = one.SetTag("foo")

				if err == nil {
					err = three.SetTag("bar")
				}

				if err == nil {
					elements = map[string]Element{
						"0": one,
						"1": two,
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
						"0": a,
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
						"1": a,
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
