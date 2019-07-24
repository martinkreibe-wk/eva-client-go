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

var _ = Describe("Nil in EDN", func() {
	It("should initialize without issue", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		lexer.RemoveFactory(NilType, NoTag)
		err = initNil(lexer)
		Ω(err).Should(BeNil())
		_, has := lexer.GetFactory(NilType, NoTag)
		Ω(has).Should(BeTrue())

		err = initNil(lexer)
		Ω(err).ShouldNot(BeNil())
	})

	It("should create elements from the factory", func() {
		var v interface{}

		lexer, err := newLexer()
		Ω(err).Should(BeNil())
		fact, has := lexer.GetFactory(NilType, NoTag)
		Ω(has).Should(BeTrue())
		elem, err := fact(v)
		Ω(err).Should(BeNil())
		Ω(elem.ElementType()).Should(BeEquivalentTo(NilType))
		Ω(elem.Value()).Should(BeNil())
	})

	It("should not create elements from the factory if the input is not a the right type", func() {
		v := "foo"

		lexer, err := newLexer()
		Ω(err).Should(BeNil())
		fact, has := lexer.GetFactory(NilType, NoTag)
		Ω(has).Should(BeTrue())
		elem, err := fact(v)
		Ω(err).ShouldNot(BeNil())
		Ω(err).Should(test.HaveMessage(ErrInvalidInput))
		Ω(elem).Should(BeNil())
	})

	Context("with the default marshaller", func() {

		It("should create an nil with no error", func() {
			elem, err := NewNilElement()
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(NilType))
		})

		It("should serialize without an issue", func() {
			elem, err := NewNilElement()
			Ω(err).Should(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("nil"))
		})

		It("should serialize without an issue", func() {
			elem, err := NewNilElement()
			Ω(err).Should(BeNil())

			_, err = elem.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})
	})

	Context("Parsing", func() {
		runParserTests(NilType,
			&testDefinition{"nil", nil},
		)
	})
})
