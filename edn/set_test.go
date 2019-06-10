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
			elem := NewStringElement("foo")

			group, err := NewSet(elem)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(SetType))
			Ω(group.Len()).Should(BeEquivalentTo(1))
		})

		It("should serialize a single nil entry in a set correctly", func() {
			elem := NewNilElement()

			group, err := NewSet(elem)
			Ω(err).Should(BeNil())

			var edn string
			edn, err = group.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("#{nil}"))
		})

		It("should serialize some nil entries in a set correctly", func() {
			elem1 := NewStringElement("foo")
			elem2 := NewStringElement("bar")
			elem3 := NewStringElement("faz")
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

			elem1 := NewStringElement("foo")
			elem2 := NewStringElement("foo")

			group, err := NewSet(elem1, elem2)
			Ω(err).Should(test.HaveMessage(ErrDuplicateKey))
			Ω(group).Should(BeNil())
		})
	})

	Context("Parsing", func() {
		runParserTests(SetType,
			&testDefinition{"#{}", func() (elements map[string][2]Element, err error) {
				return elements, err
			}},

			&testDefinition{"#{\"#{}\"}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"0": {NewIntegerElement(0), NewStringElement("#{}")},
				}
				return elements, err
			}},
			&testDefinition{"#{1}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"0": {NewIntegerElement(0), NewIntegerElement(1)},
				}
				return elements, err
			}},
			&testDefinition{"#{1 2 3}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"0": {NewIntegerElement(0), NewIntegerElement(1)},
					"1": {NewIntegerElement(1), NewIntegerElement(2)},
					"2": {NewIntegerElement(2), NewIntegerElement(3)},
				}
				return elements, err
			}},
			&testDefinition{"#{#foo 1 2 #bar 3}", func() (elements map[string][2]Element, err error) {

				one := NewIntegerElement(1)
				three := NewIntegerElement(3)

				err = one.SetTag("foo")

				if err == nil {
					err = three.SetTag("bar")
				}

				if err == nil {
					elements = map[string][2]Element{
						"0": {NewIntegerElement(0), one},
						"1": {NewIntegerElement(1), NewIntegerElement(2)},
						"2": {NewIntegerElement(2), three},
					}
				}
				return elements, err
			}},

			&testDefinition{"#{#{}}", func() (elements map[string][2]Element, err error) {
				var subList CollectionElement
				if subList, err = NewSet(); err == nil {
					elements = map[string][2]Element{
						"0": {NewIntegerElement(0), subList},
					}
				}
				return elements, err
			}},
			&testDefinition{"#{\"a\" #{}}", func() (elements map[string][2]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewSet(); err == nil {
					elements = map[string][2]Element{
						"0": {NewIntegerElement(0), NewStringElement("a")},
						"1": {NewIntegerElement(1), subList1},
					}
				}
				return elements, err
			}},
			&testDefinition{"#{#{} \"a\"}", func() (elements map[string][2]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewSet(); err == nil {
					elements = map[string][2]Element{
						"0": {NewIntegerElement(0), subList1},
						"1": {NewIntegerElement(1), NewStringElement("a")},
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
									"0": {NewIntegerElement(0), subList1},
									"1": {NewIntegerElement(1), subList2},
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
