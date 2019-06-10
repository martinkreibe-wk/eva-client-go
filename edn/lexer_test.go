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

var _ = Describe("Lexer tests", func() {
	It("should initialize without issue", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse("")
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should initialize without issue", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse(";this comment.")
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should initialize without issue", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse("[ Foo")
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should initialize without issue", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse("[ Foo")
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should break the creation if there is an error", func() {
		tt := tokenType("  ")
		Ω(tt.String()).Should(BeEquivalentTo("[Element]"))
	})

	It("should break the creation if there is an error", func() {
		tag, value := splitTag([]byte("#my/taco foobar"), "taco")
		Ω(tag).Should(BeEquivalentTo("my/"))
		Ω(value).Should(BeEquivalentTo("foobar"))
	})

	It("should initialize without issue", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse("[ { :Foo  :bar ]")
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})
})
