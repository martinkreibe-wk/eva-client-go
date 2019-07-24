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

var _ = Describe("String in EDN", func() {
	Context("", func() {
		It("should create elements from the factory", func() {
			v := "Hello world"

			lexer, err := newLexer()
			Ω(err).Should(BeNil())
			fact, has := lexer.GetFactory(StringType, NoTag)
			Ω(has).Should(BeTrue())
			elem, err := fact(v)
			Ω(err).Should(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(StringType))
			Ω(elem.Value()).Should(BeEquivalentTo(v))
		})

		It("should not create elements from the factory if the input is not a the right type", func() {
			v := 123

			lexer, err := newLexer()
			Ω(err).Should(BeNil())
			fact, has := lexer.GetFactory(StringType, NoTag)
			Ω(has).Should(BeTrue())
			elem, err := fact(v)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
			Ω(elem).Should(BeNil())
		})

		It("fail parsing invalid escapes.", func() {
			val := "\"this is not valued! \\u5\""
			elem, err := parseString(val)
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrParserError))
		})

		It("fail parsing invalid escapes.", func() {
			val := "this is not valued! \\"
			elem, err := parseString(val)
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrParserError))
		})
	})

	Context("with the default marshaller", func() {

		testValue := "This is my test value."

		It("should create a string value with no error", func() {
			elem, err := NewStringElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(StringType))
			Ω(elem.Value()).Should(BeEquivalentTo(testValue))
		})

		It("should serialize the string without an issue", func() {
			elem, err := NewStringElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("\"" + testValue + "\""))
		})

		It("should serialize the string without an issue", func() {
			elem, err := NewStringElement(testValue)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			_, err = elem.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})
	})

	Context("Parsing", func() {

		var tests []*testDefinition

		tests = append(tests,
			&testDefinition{"\"\"", ""},
			&testDefinition{"\"value\"", "value"},
			&testDefinition{"\"value's\"", "value's"},
			&testDefinition{"\" value\"", " value"},
			&testDefinition{"\"value \"", "value "},
			&testDefinition{"\"\\tvalue\"", "\tvalue"},
			&testDefinition{"\"value\\t\"", "value\t"},
			&testDefinition{"\"\\bvalue\"", "\bvalue"},
			&testDefinition{"\"value\\b\"", "value\b"},
			&testDefinition{"\"\\nvalue\"", "\nvalue"},
			&testDefinition{"\"value\\n\"", "value\n"},
			&testDefinition{"\"\\rvalue\"", "\rvalue"},
			&testDefinition{"\"value\\r\"", "value\r"},
			&testDefinition{"\"\\fvalue\"", "\fvalue"},
			&testDefinition{"\"value\\f\"", "value\f"},
			&testDefinition{"\"\\\"value\"", "\"value"},
			&testDefinition{"\"value\\\"\"", "value\""},
			&testDefinition{"\"\\'value\"", "'value"},
			&testDefinition{"\"value\\'\"", "value'"},
			&testDefinition{"\"\\\\value\"", "\\value"},
			&testDefinition{"\"value\\\\\"", "value\\"},
			&testDefinition{"\"\\t\"", "\t"},
			&testDefinition{"\"\\\\t\"", "\\t"},
			&testDefinition{"\"\\u2318\"", "⌘"},
			&testDefinition{"\"\\u20AC\"", "€"},
			&testDefinition{"\"value value\"", "value value"},
			&testDefinition{"\"()\"", "()"},
			&testDefinition{"\"[]\"", "[]"},
			&testDefinition{"\"6ba7b810-9dad-11d1-80b4-00c04fd430c8\"", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
		)

		charTests := []string{
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"`", "~",
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=",
			"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+",
			"[", "]", ";", ",", ".", "/",
			"{", "}", "|", ":", "<", ">", "?",
		}

		for _, v := range charTests {
			tests = append(tests,
				&testDefinition{"\"" + v + "\"", v},
			)
		}

		runParserTests(StringType,
			tests...,
		)
	})
})
