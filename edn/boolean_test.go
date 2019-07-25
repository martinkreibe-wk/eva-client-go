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

var _ = Describe("Boolean in EDN", func() {

	Context("", func() {

		It("should create elements from the factory", func() {
			fact, err := DefaultLexer.GetFactory(BooleanType, NoTag)
			Ω(err).Should(BeNil())
			elem, err := fact(true)
			Ω(err).Should(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(BooleanType))
			Ω(elem.Value()).Should(BeTrue())
		})

		It("should not create elements from the factory if the input is not a the right type", func() {
			fact, err := DefaultLexer.GetFactory(BooleanType, NoTag)
			Ω(err).Should(BeNil())
			elem, err := fact("true")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
			Ω(elem).Should(BeNil())
		})
	})

	Context("with the default marshaller", func() {

		It("should create an true value with no error", func() {
			elem, err := NewBooleanElement(true)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(BooleanType))
			Ω(elem.Value()).Should(BeTrue())
		})

		It("should create an false value with no error", func() {
			elem, err := NewBooleanElement(false)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(BooleanType))
			Ω(elem.Value()).Should(BeFalse())
		})

		It("should serialize true without an issue", func() {
			elem, err := NewBooleanElement(true)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			var edn string
			stream := NewStringStream()
			err = EvaEdnMimeType.SerializeTo(stream, elem)
			edn = stream.String()
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("true"))
		})
	})

	Context("Parsing", func() {
		runParserTests(BooleanType,
			&testDefinition{"true", true},
			&testDefinition{"false", false},
		)
	})
})
