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
	"github.com/mattrobenolt/gocql/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UUID in EDN", func() {

	Context("", func() {

		It("should create elements from the factory", func() {
			uuidValue := "12345678-90ab-cdef-9876-0123456789ab"
			v, err := uuid.ParseUUID(uuidValue)
			Ω(err).Should(BeNil())

			fact, err := DefaultLexer.GetFactory(UUIDType, UUIDElementTag)
			Ω(err).Should(BeNil())
			elem, err := fact(v)
			Ω(err).Should(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(UUIDType))
			Ω(elem.Value()).Should(BeEquivalentTo(v))
		})

		It("should not create elements from the factory if the input is not a the right type", func() {
			v := "foo"

			fact, err := DefaultLexer.GetFactory(UUIDType, UUIDElementTag)
			Ω(err).Should(BeNil())
			elem, err := fact(v)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrParserError))
			Ω(elem).Should(BeNil())
		})
	})

	Context("with the default marshaller", func() {

		uuidValue := "12345678-90ab-cdef-9876-0123456789ab"
		testValue, err := uuid.ParseUUID(uuidValue)
		if err != nil {
			panic(err)
		}

		It("should create an uuid value with no error", func() {

			elem, err := NewUUIDElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(UUIDType))
			Ω(elem.Value()).Should(BeEquivalentTo(testValue))
		})

		It("should serialize the uuid without an issue", func() {
			elem, err := NewUUIDElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("#uuid " + uuidValue))
		})

		It("should serialize the uuid without an issue", func() {
			elem, err := NewUUIDElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			_, err = elem.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})
	})

	Context("Parsing", func() {
		runParserTests(UUIDType,
			&testDefinition{"#uuid \"6ba7b810-9dad-11d1-80b4-00c04fd430c8\"", func() (string, interface{}, error) {
				tag := "uuid"
				v, e := uuid.ParseUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
				return tag, v, e
			}},
		)
	})
})
