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
	"time"

	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Instant in EDN", func() {
	Context("", func() {

		It("should create elements from the factory", func() {
			v := time.Date(2017, 12, 28, 22, 20, 30, 450, time.UTC)

			lexer, err := newLexer()
			Ω(err).Should(BeNil())
			fact, has := lexer.GetFactory(InstantType, InstantElementTag)
			Ω(has).Should(BeTrue())
			elem, err := fact(v)
			Ω(err).Should(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(InstantType))
			Ω(elem.Value()).Should(BeEquivalentTo(v))
		})

		It("should not create elements from the factory if the input is not a the right type", func() {
			v := "foo"

			lexer, err := newLexer()
			Ω(err).Should(BeNil())
			fact, has := lexer.GetFactory(InstantType, InstantElementTag)
			Ω(has).Should(BeTrue())
			elem, err := fact(v)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
			Ω(elem).Should(BeNil())
		})
	})

	Context("with the default marshaller", func() {

		testValue := time.Date(2017, 12, 28, 22, 20, 30, 450, time.UTC)

		It("should create an instant value with no error", func() {
			elem, err := NewInstantElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(InstantType))
			Ω(elem.Value()).Should(BeEquivalentTo(testValue))
		})

		It("should serialize the instant without an issue", func() {
			elem, err := NewInstantElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("#inst 2017-12-28T22:20:30Z"))
		})

		It("should serialize the instant without an issue", func() {
			elem, err := NewInstantElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			_, err = elem.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})
	})

	Context("Parsing", func() {
		runParserTests(InstantType,
			&testDefinition{"#inst \"1985-04-12T23:20:50.52Z\"", func() (string, interface{}, error) {
				tag := "inst"
				v, e := time.Parse(time.RFC3339, "1985-04-12T23:20:50.52Z")
				return tag, v, e
			}},
		)
	})
})
