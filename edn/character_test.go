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

var _ = Describe("Character in EDN", func() {
	Context("", func() {

		It("should initialize without issue", func() {
			lexer, err := newLexer()
			Ω(err).Should(BeNil())

			lexer.RemoveFactory(CharacterType, NoTag)
			err = initCharacter(lexer)
			Ω(err).Should(BeNil())
			_, has := lexer.GetFactory(CharacterType, NoTag)
			Ω(has).Should(BeTrue())

			err = initCharacter(lexer)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidFactory))
		})

		It("should create elements from the factory", func() {
			v := 'g'

			lexer, err := newLexer()
			Ω(err).Should(BeNil())
			fact, has := lexer.GetFactory(CharacterType, NoTag)
			Ω(has).Should(BeTrue())
			elem, err := fact(v)
			Ω(err).Should(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(CharacterType))
			Ω(elem.Value()).Should(BeEquivalentTo(v))
		})

		It("should not create elements from the factory if the input is not a the right type", func() {
			v := "foo"

			lexer, err := newLexer()
			Ω(err).Should(BeNil())
			fact, has := lexer.GetFactory(CharacterType, NoTag)
			Ω(has).Should(BeTrue())
			elem, err := fact(v)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
			Ω(elem).Should(BeNil())
		})
	})

	Context("with the default marshaller", func() {

		c := 'c'
		runes := map[rune]string{
			c:    "\\c",
			'\n': "\\newline",
			'\r': "\\return",
			' ':  "\\space",
			'\t': "\\tab",
			'⌘':  "\\u2318",
		}

		It("should create an character value with no error", func() {
			elem, err := NewCharacterElement(c)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(CharacterType))
			Ω(elem.Value()).Should(BeEquivalentTo(c))
		})

		It("should serialize the character without an issue", func() {
			for r, ser := range runes {
				elem, err := NewCharacterElement(r)
				Ω(err).Should(BeNil())
				Ω(elem).ShouldNot(BeNil())
				Ω(elem.Value()).Should(BeEquivalentTo(r), fmt.Sprintf("For rune: %+q", r))

				edn, err := elem.Serialize(EvaEdnMimeType)
				Ω(err).Should(BeNil())
				Ω(edn).Should(BeEquivalentTo(ser), fmt.Sprintf("For rune: %+q", r))
			}
		})

		It("should serialize the character without an issue", func() {
			elem, err := NewCharacterElement('x')
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.Value()).Should(BeEquivalentTo('x'), fmt.Sprintf("For rune: %+q", 'x'))

			_, err = elem.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})

		It("should create an character value with no error", func() {
			elem, err := NewCharacterElement(c)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(CharacterType))
			Ω(elem.Value()).Should(BeEquivalentTo(c))
		})
	})

	Context("Parsing", func() {
		runParserTests(CharacterType,
			&testDefinition{"\\return", '\r'},
			&testDefinition{"\\newline", '\n'},
			&testDefinition{"\\space", ' '},
			&testDefinition{"\\tab", '\t'},
			&testDefinition{"\\r", 'r'},
			&testDefinition{"\\n", 'n'},
			&testDefinition{"\\s", 's'},
			&testDefinition{"\\u2318", '⌘'},
			&testDefinition{"\\u20AC", '€'},
		)
	})
})
