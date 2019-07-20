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

var _ = Describe("Float in EDN", func() {
	Context("", func() {

		It("should initialize without issue", func() {
			lexer, err := newLexer()
			Ω(err).Should(BeNil())

			delete(typeFactories, FloatType)
			err = initFloat(lexer)
			Ω(err).Should(BeNil())
			_, has := typeFactories[FloatType]
			Ω(has).Should(BeTrue())

			err = initFloat(lexer)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidFactory))
		})

		It("should create elements from the factory", func() {
			v := float64(1.234)

			elem, err := typeFactories[FloatType](v)
			Ω(err).Should(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(FloatType))
			Ω(elem.Value()).Should(BeEquivalentTo(v))
		})

		It("should not create elements from the factory if the input is not a the right type", func() {
			v := "foo"

			elem, err := typeFactories[FloatType](v)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
			Ω(elem).Should(BeNil())
		})
	})

	Context("with the default marshaller", func() {

		testValue := float64(12345.67)

		It("should create an float value with no error", func() {
			elem, err := NewFloatElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(FloatType))
			Ω(elem.Value()).Should(BeEquivalentTo(testValue))
		})

		It("should serialize the float without an issue", func() {
			elem, err := NewFloatElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("1.234567E+04"))
		})

		It("should serialize the float without an issue", func() {
			elem, err := NewFloatElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			_, err = elem.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})
	})

	Context("Parsing", func() {
		runParserTests(FloatType,
			&testDefinition{"0.0", 0.0},
			&testDefinition{"+0.0", 0.0},
			&testDefinition{"-0.0", 0.0},
			&testDefinition{"1.0", 1.0},
			&testDefinition{"-1.0", -1.0},
			&testDefinition{"1234.0", 1234.0},
			&testDefinition{"12.340", 12.34},
			&testDefinition{"12.34", 12.34},
			&testDefinition{"0M", 0.0},
			&testDefinition{"+0M", 0.0},
			&testDefinition{"-0M", 0.0},
			&testDefinition{"1M", 1.0},
			&testDefinition{"-1M", -1.0},
			&testDefinition{"1234E-2", 12.34},
			&testDefinition{"1.234E1", 12.34},
			&testDefinition{"12.34E0", 12.34},
		)
	})
})
