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
	"strings"

	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type badReader struct {
}

var (
	testError = fmt.Errorf("test error")
)

func (r *badReader) Read(p []byte) (n int, err error) {
	return 0, testError
}

var _ = Describe("Lexer tests", func() {
	It("should error with a nil reader with a parse error", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse(nil)
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should error with a bad reader with a parse error", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse(&badReader{})
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should error with a empty string with a parse error", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse(strings.NewReader(""))
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should error with a comment only string with a parse error", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse(strings.NewReader(";this comment."))
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should error with a open collection string with a parse error", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse(strings.NewReader("[ Foo"))
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})

	It("should see that a dual blank string is an element token type", func() {
		tt := tokenType("  ")
		Ω(tt.String()).Should(BeEquivalentTo("[Element]"))
	})

	It("should be able to split a tag from the element", func() {
		tag, value := splitTag([]byte("#my/taco foobar"), "taco")
		Ω(tag).Should(BeEquivalentTo("my/"))
		Ω(value).Should(BeEquivalentTo("foobar"))
	})

	It("should error with a open collection, in a collection string with a parse error", func() {
		lexer, err := newLexer()
		Ω(err).Should(BeNil())

		var elem Element
		elem, err = lexer.Parse(strings.NewReader("[ { :Foo  :bar ]"))
		Ω(err).ShouldNot(BeNil())
		Ω(elem).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})
})
