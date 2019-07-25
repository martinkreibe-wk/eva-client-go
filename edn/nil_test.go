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

	It("should create elements from the factory", func() {
		var v interface{}

		fact, err := DefaultLexer.GetFactory(NilType, NoTag)
		Ω(err).Should(BeNil())
		elem, err := fact(v)
		Ω(err).Should(BeNil())
		Ω(elem.ElementType()).Should(BeEquivalentTo(NilType))
		Ω(elem.Value()).Should(BeNil())
	})

	It("should not create elements from the factory if the input is not a the right type", func() {
		v := "foo"

		fact, err := DefaultLexer.GetFactory(NilType, NoTag)
		Ω(err).Should(BeNil())
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

			stream := NewStringStream()
			err = EvaEdnMimeType.SerializeTo(stream, elem)
			edn := stream.String()
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("nil"))
		})
	})

	Context("Parsing", func() {
		runParserTests(NilType,
			&testDefinition{"nil", nil},
		)
	})
})
