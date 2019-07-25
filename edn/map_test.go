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

	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map in EDN", func() {
	Context("with the default marshaller", func() {

		makePair := func(key string, value string) Pair {
			elem, err := NewStringElement(value)
			Ω(err).Should(BeNil())

			pair, err := NewPair(key, elem)
			Ω(err).Should(BeNil())
			Ω(pair.Key()).ShouldNot(BeNil())
			Ω(pair.Key().ElementType()).Should(BeEquivalentTo(KeywordType))
			Ω(pair.Key().Value().(SymbolElement).Name()).Should(BeEquivalentTo(key[1:]))

			Ω(pair.Value()).ShouldNot(BeNil())
			Ω(pair.Value().ElementType()).Should(BeEquivalentTo(StringType))
			Ω(pair.Value().Value()).Should(BeEquivalentTo(value))

			return pair
		}

		It("should create an empty map with no error", func() {
			group, err := NewMap()
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(MapType))
			Ω(group.Len()).Should(BeEquivalentTo(0))
		})

		It("should serialize an empty map correctly", func() {
			group, err := NewMap()
			Ω(err).Should(BeNil())

			stream := NewStringStream()
			err = EvaEdnMimeType.SerializeTo(stream, group)
			edn := stream.String()
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("{}"))
		})

		It("should error with a nil item", func() {
			group, err := NewMap(nil)
			Ω(err).Should(test.HaveMessage(ErrInvalidPair))
			Ω(group).Should(BeNil())
		})

		It("should create a map element with the initial values", func() {
			elem, err := NewStringElement("foo")
			Ω(err).Should(BeNil())

			pair, err := NewPair(elem, elem)
			Ω(err).Should(BeNil())

			group, err := NewMap(pair)
			Ω(err).Should(BeNil())
			Ω(group).ShouldNot(BeNil())
			Ω(group.ElementType()).Should(BeEquivalentTo(MapType))
			Ω(group.Len()).Should(BeEquivalentTo(1))

			var v Element
			v, err = group.Get("foo")
			Ω(err).Should(BeNil())
			Ω(v).ShouldNot(BeNil())

			v, err = group.Get("not-here")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrNoValue))
			Ω(v).Should(BeNil())
		})

		It("should serialize a single nil entry in a map correctly", func() {
			elem, err := NewNilElement()
			Ω(err).Should(BeNil())

			pair, err := NewPair(elem, elem)
			Ω(err).Should(BeNil())

			group, err := NewMap(pair)
			Ω(err).Should(BeNil())

			stream := NewStringStream()
			err = EvaEdnMimeType.SerializeTo(stream, group)
			edn := stream.String()
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("{nil nil}"))
		})

		It("should serialize some pairs entries in a map correctly", func() {
			keys := map[string]string{
				":key1": "val1",
				":key2": "val2",
				":key3": "val3",
				":key4": "val3", // same values are ok
			}

			var pairs []Pair
			for k, v := range keys {
				pairs = append(pairs, makePair(k, v))
			}

			group, err := NewMap(pairs...)
			Ω(err).Should(BeNil())

			stream := NewStringStream()
			err = EvaEdnMimeType.SerializeTo(stream, group)
			edn := stream.String()
			Ω(err).Should(BeNil())
			Ω(edn).Should(HavePrefix("{"))
			Ω(edn).Should(HaveSuffix("}"))

			for k, v := range keys {
				Ω(edn).Should(ContainSubstring(k + " " + "\"" + v + "\""))
			}
		})

		It("should not accept duplicate keys", func() {
			p1 := makePair(":key1", "val1")
			p2 := makePair(":key1", "val2")
			Ω(p1.Key().Equals(p2.Key())).Should(BeTrue())

			group, err := NewMap(
				p1,
				p2,
			)
			Ω(err).ShouldNot(BeNil())
			Ω(group).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrDuplicateKey))
		})

		It("should break the iteration and return the error", func() {

			keys := map[string]string{
				":key1": "val1",
				":key2": "val2",
				":key3": "val3",
				":key4": "val3", // same values are ok
			}

			var pairs []Pair
			for k, v := range keys {
				pairs = append(pairs, makePair(k, v))
			}

			group, err := NewMap(pairs...)
			Ω(err).Should(BeNil())

			breakCount := len(keys) / 2
			Ω(len(keys) > breakCount).Should(BeTrue())

			templateError := ErrorMessage("This is the expected error")
			err = group.IterateChildren(func(key, value Element) (e error) {
				if breakCount--; breakCount == 0 {
					e = MakeError(templateError, "")
				}
				return e
			})

			Ω(err).Should(test.HaveMessage(templateError))
		})

		It("should merge correctly", func() {

			keys := map[string]string{
				":key1": "val1",
				":key2": "val2",
			}

			keys2 := map[string]string{
				":key3": "val3",
				":key4": "val3", // same values are ok
			}

			var pairs []Pair
			for k, v := range keys {
				pairs = append(pairs, makePair(k, v))
			}

			group, err := NewMap(pairs...)
			Ω(group).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			pairs = []Pair{}
			for k, v := range keys2 {
				pairs = append(pairs, makePair(k, v))
			}

			group2, err := NewMap(pairs...)
			Ω(group2).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			err = group.Merge(group2)
			Ω(err).Should(BeNil())
			Ω(group.Len()).Should(BeEquivalentTo(4))
		})

		It("should prepend correctly", func() {

			keys := map[string]string{
				":key1": "val1",
				":key2": "val2",
			}

			keys2 := map[string]string{
				":key3": "val3",
				":key4": "val3", // same values are ok
			}

			var pairs []Pair
			for k, v := range keys {
				pairs = append(pairs, makePair(k, v))
			}

			group, err := NewMap(pairs...)
			Ω(group).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			for k, v := range keys2 {
				key, err := NewPrimitiveElement(k)
				Ω(err).Should(BeNil())
				value, err := NewPrimitiveElement(v)
				Ω(err).Should(BeNil())
				err = group.Prepend(key, value)
				Ω(err).Should(BeNil())
			}

			Ω(group.Len()).Should(BeEquivalentTo(4))
		})

		It("should append correctly", func() {

			keys := map[string]string{
				":key1": "val1",
				":key2": "val2",
			}

			keys2 := map[string]string{
				":key3": "val3",
				":key4": "val3", // same values are ok
			}

			var pairs []Pair
			for k, v := range keys {
				pairs = append(pairs, makePair(k, v))
			}

			group, err := NewMap(pairs...)
			Ω(group).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			for k, v := range keys2 {
				key, err := NewPrimitiveElement(k)
				Ω(err).Should(BeNil())
				value, err := NewPrimitiveElement(v)
				Ω(err).Should(BeNil())
				err = group.Append(key, value)
				Ω(err).Should(BeNil())
			}

			Ω(group.Len()).Should(BeEquivalentTo(4))
		})

		It("should equal correctly", func() {

			keys := map[string]string{
				":key1": "val1",
				":key2": "val2",
			}

			keys2 := map[string]string{
				":key1": "val1",
				":key2": "val2",
			}

			var pairs []Pair
			for k, v := range keys {
				pairs = append(pairs, makePair(k, v))
			}

			group, err := NewMap(pairs...)
			Ω(group).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			pairs = []Pair{}
			for k, v := range keys2 {
				pairs = append(pairs, makePair(k, v))
			}

			group2, err := NewMap(pairs...)
			Ω(group2).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			Ω(group.Equals(group2)).Should(BeTrue())
		})

		It("should equal correctly", func() {

			keys := map[string]string{
				":key1": "val1",
				":key2": "val2",
			}

			keys2 := map[string]string{
				":key1": "val1",
				":key2": "val3",
			}

			var pairs []Pair
			for k, v := range keys {
				pairs = append(pairs, makePair(k, v))
			}

			group, err := NewMap(pairs...)
			Ω(group).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			pairs = []Pair{}
			for k, v := range keys2 {
				pairs = append(pairs, makePair(k, v))
			}

			group2, err := NewMap(pairs...)
			Ω(group2).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			Ω(group.Equals(group2)).Should(BeFalse())
		})

		It("should equal correctly", func() {

			keys := map[string]string{
				":key1": "val1",
				":key2": "val2",
			}

			keys2 := map[string]string{
				":key1": "val1",
				":key3": "val2",
			}

			var pairs []Pair
			for k, v := range keys {
				pairs = append(pairs, makePair(k, v))
			}

			group, err := NewMap(pairs...)
			Ω(group).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			pairs = []Pair{}
			for k, v := range keys2 {
				pairs = append(pairs, makePair(k, v))
			}

			group2, err := NewMap(pairs...)
			Ω(group2).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			Ω(group.Equals(group2)).Should(BeFalse())
		})

		It("should break the iteration and return the error", func() {

			m, err := NewMap()
			Ω(err).Should(BeNil())

			elem, err := NewNilElement()
			Ω(err).Should(BeNil())

			err = m.Append(elem)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
		})

		It("should error if somehow the collection was not the type we were expecting.", func() {

			m, err := NewMap()
			Ω(err).Should(BeNil())

			raw := m.(*collectionElemImpl)
			raw.collection = &struct{}{} // overwrite the actual data.

			elem, err := NewNilElement()
			Ω(err).Should(BeNil())

			err = m.Append(elem)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidElement))

			_, err = m.Get("foo")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidElement))
		})

		It("should break the creation if there is an error", func() {

			p, err := NewPair(":key1", nil)
			Ω(err).ShouldNot(BeNil())

			_, err = NewMap(p)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidPair))
		})

		It("should break the creation if there is an error", func() {
			_, err := ParseString("{ :foo }")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidPair))
		})

		It("should break the creation if there is an error", func() {
			_, err := ParseString("{ :foo 2 ]")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrParserError))
		})

		It("should break the creation if there is an error", func() {
			_, err := ParseString("{ :foo 2 :taco")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrParserError))
		})

		It("should break the creation if there is an error", func() {
			_, err := ParseString("{ :foo 2 \n\t")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrParserError))
		})
	})

	Context("Parsing", func() {

		a, err := NewStringElement("a")
		if err != nil {
			panic(err)
		}

		zero, err := NewStringElement("0")
		if err != nil {
			panic(err)
		}
		one, err := NewStringElement("1")
		if err != nil {
			panic(err)
		}
		two, err := NewStringElement("2")
		if err != nil {
			panic(err)
		}
		braces, err := NewStringElement("[]")
		if err != nil {
			panic(err)
		}

		onei, err := NewIntegerElement(1)
		if err != nil {
			panic(err)
		}

		twoi, err := NewIntegerElement(2)
		if err != nil {
			panic(err)
		}

		threei, err := NewIntegerElement(3)
		if err != nil {
			panic(err)
		}

		runParserTests(MapType,
			&testDefinition{"{}", func() (elements map[string][2]Element, err error) {
				return elements, err
			}},
			&testDefinition{"{\"0\" \"[]\"}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"\"0\"": {zero, braces},
				}
				return elements, err
			}},
			&testDefinition{"{\"0\" 1}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"\"0\"": {zero, onei},
				}
				return elements, err
			}},
			&testDefinition{"{\"0\" 1 \"1\" 2 \"2\" 3}", func() (elements map[string][2]Element, err error) {
				elements = map[string][2]Element{
					"\"0\"": {zero, onei},
					"\"1\"": {one, twoi},
					"\"2\"": {two, threei},
				}
				return elements, err
			}},
			&testDefinition{"{#foo 1 2}", func() (elements map[string][2]Element, err error) {

				one, err := NewIntegerElement(1)
				Ω(err).Should(BeNil())

				err = one.SetTag("foo")

				if err == nil {
					t, err := NewIntegerElement(1)
					Ω(err).Should(BeNil())
					if e := t.SetTag("foo"); e != nil {
						Fail(fmt.Sprintf("should not have gotten this: %s", e))
					}

					tt, err := NewIntegerElement(2)
					Ω(err).Should(BeNil())

					elements = map[string][2]Element{
						"#foo 1": {t, tt},
					}
				}
				return elements, err
			}},
			&testDefinition{"{\"0\" {}}", func() (elements map[string][2]Element, err error) {
				var subList CollectionElement
				if subList, err = NewMap(); err == nil {
					elements = map[string][2]Element{
						"\"0\"": {zero, subList},
					}
				}
				return elements, err
			}},
			&testDefinition{"{\"a\" []}", func() (elements map[string][2]Element, err error) {
				var subList1 CollectionElement
				if subList1, err = NewVector(); err == nil {
					elements = map[string][2]Element{
						"\"a\"": {a, subList1},
					}
				}
				return elements, err
			}},
		)
	})
})
