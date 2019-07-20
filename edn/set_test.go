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

var _ = Describe("Set in EDN", func() {
	Context("with the default marshaller", func() {
		It("should create an empty set with no error", func() {
			group, err := NewSet()
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(SetType))
			Ω(group.Len()).Should(BeEquivalentTo(0))
		})

		It("should serialize an empty set correctly", func() {
			group, err := NewSet()
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("#{}"))
		})

		It("should serialize an empty set correctly", func() {
			group, err := NewSet()
			Ω(err).Should(BeNil())

			_, err = group.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})

		It("should error with a nil item", func() {
			group, err := NewSet(nil)
			Ω(err).Should(test.HaveMessage(ErrInvalidElement))
			Ω(group).Should(BeNil())
		})

		It("should create a set element with the initial values", func() {
			elem, err := NewStringElement("foo")
			Ω(err).Should(BeNil())

			group, err := NewSet(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(SetType))
			Ω(group.Len()).Should(BeEquivalentTo(1))
		})

		It("should serialize a single nil entry in a set correctly", func() {
			elem, err := NewNilElement()
			Ω(err).Should(BeNil())

			group, err := NewSet(elem)
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("#{nil}"))
		})

		It("should serialize some nil entries in a set correctly", func() {
			elem1, err := NewStringElement("foo")
			Ω(err).Should(BeNil())
			elem2, err := NewStringElement("bar")
			Ω(err).Should(BeNil())
			elem3, err := NewStringElement("faz")
			Ω(err).Should(BeNil())
			keys := []string{
				"foo",
				"bar",
				"faz",
			}

			group, err := NewSet(elem1, elem2, elem3)
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(HavePrefix("#{"))
			Ω(edn).Should(HaveSuffix("}"))

			for _, v := range keys {
				Ω(edn).Should(ContainSubstring("\"" + v + "\""))
			}
		})

		It("should error if two elements are the same", func() {

			elem1, err := NewStringElement("foo")
			Ω(err).Should(BeNil())
			elem2, err := NewStringElement("foo")
			Ω(err).Should(BeNil())

			group, err := NewSet(elem1, elem2)
			Ω(err).Should(test.HaveMessage(ErrDuplicateKey))
			Ω(group).Should(BeNil())
		})
	})

	Context("Parsing", func() {

		set, err := NewStringElement("#{}")
		if err != nil {
			panic(err)
		}

		a, err := NewStringElement("a")
		if err != nil {
			panic(err)
		}

		zero, err := NewIntegerElement(0)
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

		runParserTests(SetType,
			&testDefinition{"#{}", func() (elements map[string][2]Element, err error) {
				return elements, err
			}},

			&testDefinition{"#{\"#{}\"}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"0": {zero, set},
				}
				return elements, err
			}},
			&testDefinition{"#{1}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"0": {zero, one},
				}
				return elements, err
			}},
			&testDefinition{"#{1 2 3}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"0": {zero, one},
					"1": {one, two},
					"2": {two, three},
				}
				return elements, err
			}},
			&testDefinition{"#{#foo 1 2 #bar 3}", func() (elements map[string][2]Element, err error) {

				onei, err := NewIntegerElement(1)
				if err != nil {
					return nil, err
				}
				threei, err := NewIntegerElement(3)
				if err != nil {
					return nil, err
				}

				err = onei.SetTag("foo")

				if err == nil {
					err = threei.SetTag("bar")
				}

				if err == nil {
					elements = map[string][2]Element{
						"0": {zero, onei},
						"1": {one, two},
						"2": {two, threei},
					}
				}
				return elements, err
			}},

			&testDefinition{"#{#{}}", func() (elements map[string][2]Element, err error) {
				var subList CollectionElement
				if subList, err = NewSet(); err == nil {
					elements = map[string][2]Element{
						"0": {zero, subList},
					}
				}
				return elements, err
			}},
			&testDefinition{"#{\"a\" #{}}", func() (elements map[string][2]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewSet(); err == nil {
					elements = map[string][2]Element{
						"0": {zero, a},
						"1": {one, subList1},
					}
				}
				return elements, err
			}},
			&testDefinition{"#{#{} \"a\"}", func() (elements map[string][2]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewSet(); err == nil {
					elements = map[string][2]Element{
						"0": {zero, subList1},
						"1": {one, a},
					}
				}
				return elements, err
			}},
			&testDefinition{"#{#foo #{} #bar #{}}", func() (elements map[string][2]Element, err error) {
				var subList1 CollectionElement
				var subList2 CollectionElement
				if subList1, err = NewSet(); err == nil {
					if subList2, err = NewSet(); err == nil {
						if err = subList1.SetTag("foo"); err == nil {
							if err = subList2.SetTag("bar"); err == nil {
								elements = map[string][2]Element{
									"0": {zero, subList1},
									"1": {one, subList2},
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
