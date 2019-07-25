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

var _ = Describe("Keyword in EDN", func() {
	Context("", func() {

		It("should create elements from the factory", func() {
			v := "testKeyword"

			fact, err := DefaultLexer.GetFactory(KeywordType, NoTag)
			Ω(err).Should(BeNil())
			elem, err := fact(v)
			Ω(err).Should(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(KeywordType))

			symbol, is := elem.(SymbolElement)
			Ω(is).Should(BeTrue())
			Ω(symbol.Name()).Should(BeEquivalentTo(v))
		})

		It("should not create elements from the factory if the input is not a the right type", func() {
			v := 123

			fact, err := DefaultLexer.GetFactory(KeywordType, NoTag)
			Ω(err).Should(BeNil())
			elem, err := fact(v)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
			Ω(elem).Should(BeNil())
		})
	})

	Context("with the default marshaller", func() {

		badKeywords := []string{
			":/",
			"::",
			":",
			"/",
			":1",
			"1",
			"bad1/1worse",
			"/bad",
			"bad/",
			"bad/worse/wrong",
			"worse/+1bad",
			"+1bad",
			"-1bad",
			".1",
		}

		It("should create a keyword with just one parameter value with no error", func() {
			prefix := ""
			name := "foobar"

			elem, err := NewKeywordElement(name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(KeywordType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should create a keyword with two parameter value with no error", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewKeywordElement(prefix, name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(KeywordType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should create a keyword with one parameter (but with the separator) value with no error", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewKeywordElement(prefix + SymbolSeparator + name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(KeywordType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should create a keyword with just one parameter (first : prefixed) value with no error", func() {
			prefix := ""
			name := "foobar"

			elem, err := NewKeywordElement(KeywordPrefix + name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(KeywordType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should create a keyword with two parameter (first : prefixed) value with no error", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewKeywordElement(KeywordPrefix+prefix, name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(KeywordType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should create a keyword with one parameter (: prefixed and with the separator) value with no error", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewKeywordElement(KeywordPrefix + prefix + SymbolSeparator + name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(KeywordType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should create a keyword with zero parameter value with an error", func() {
			elem, err := NewKeywordElement()
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidKeyword))
		})

		It("should create a keyword with three parameter value with an error", func() {
			elem, err := NewKeywordElement("a", "b", "c")
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidKeyword))
		})

		It("should serialize the keyword with one parameter without an issue", func() {
			name := "foobar"

			elem, err := NewKeywordElement(name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(KeywordPrefix + name))
		})

		It("should serialize the keyword with one parameter without an issue", func() {
			name := "foobar"

			elem, err := NewKeywordElement(name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			_, err = elem.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})

		It("should serialize the keyword with two parameter without an issue", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewKeywordElement(prefix, name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(KeywordPrefix + prefix + SymbolSeparator + name))
		})

		It("should serialize the keyword with one (but with the separator) parameter without an issue", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewKeywordElement(prefix + SymbolSeparator + name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(KeywordPrefix + prefix + SymbolSeparator + name))
		})

		It("should serialize the keyword with one parameter (with : prefix) without an issue", func() {
			name := "foobar"

			elem, err := NewKeywordElement(KeywordPrefix + name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(KeywordPrefix + name))
		})

		It("should serialize the keyword with two parameter (with : prefix) without an issue", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewKeywordElement(KeywordPrefix+prefix, name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(KeywordPrefix + prefix + SymbolSeparator + name))
		})

		It("should serialize the keyword with one (with : prefix and with the separator) parameter without an issue", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewKeywordElement(KeywordPrefix + prefix + SymbolSeparator + name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(KeywordPrefix + prefix + SymbolSeparator + name))
		})

		It("should not process all odd invalid keywords", func() {

			for _, keyword := range badKeywords {
				elem, err := NewKeywordElement(keyword)
				Ω(elem).Should(BeNil())
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(test.HaveMessage(ErrInvalidKeyword))
			}
		})
	})

	Context("Parsing", func() {
		runParserTests(KeywordType,
			&testDefinition{":foo", &keywordValue{"", "foo"}},
			&testDefinition{":bar/foo", &keywordValue{"bar", "foo"}},

			&testDefinition{":*", &keywordValue{"", "*"}},
			&testDefinition{":!", &keywordValue{"", "!"}},
			&testDefinition{":?", &keywordValue{"", "?"}},
			&testDefinition{":$", &keywordValue{"", "$"}},
			&testDefinition{":%", &keywordValue{"", "%"}},
			&testDefinition{":&", &keywordValue{"", "&"}},
			&testDefinition{":=", &keywordValue{"", "="}},
			&testDefinition{":<", &keywordValue{"", "<"}},
			&testDefinition{":>", &keywordValue{"", ">"}},

			&testDefinition{":*-", &keywordValue{"", "*-"}},
			&testDefinition{":!-", &keywordValue{"", "!-"}},
			&testDefinition{":?-", &keywordValue{"", "?-"}},
			&testDefinition{":$-", &keywordValue{"", "$-"}},
			&testDefinition{":%-", &keywordValue{"", "%-"}},
			&testDefinition{":&-", &keywordValue{"", "&-"}},
			&testDefinition{":=-", &keywordValue{"", "=-"}},
			&testDefinition{":<-", &keywordValue{"", "<-"}},
			&testDefinition{":>-", &keywordValue{"", ">-"}},

			&testDefinition{":*+", &keywordValue{"", "*+"}},
			&testDefinition{":!+", &keywordValue{"", "!+"}},
			&testDefinition{":?+", &keywordValue{"", "?+"}},
			&testDefinition{":$+", &keywordValue{"", "$+"}},
			&testDefinition{":%+", &keywordValue{"", "%+"}},
			&testDefinition{":&+", &keywordValue{"", "&+"}},
			&testDefinition{":=+", &keywordValue{"", "=+"}},
			&testDefinition{":<+", &keywordValue{"", "<+"}},
			&testDefinition{":>+", &keywordValue{"", ">+"}},

			&testDefinition{":*.", &keywordValue{"", "*."}},
			&testDefinition{":!.", &keywordValue{"", "!."}},
			&testDefinition{":?.", &keywordValue{"", "?."}},
			&testDefinition{":$.", &keywordValue{"", "$."}},
			&testDefinition{":%.", &keywordValue{"", "%."}},
			&testDefinition{":&.", &keywordValue{"", "&."}},
			&testDefinition{":=.", &keywordValue{"", "=."}},
			&testDefinition{":<.", &keywordValue{"", "<."}},
			&testDefinition{":>.", &keywordValue{"", ">."}},

			&testDefinition{":*#", &keywordValue{"", "*#"}},
			&testDefinition{":!#", &keywordValue{"", "!#"}},
			&testDefinition{":?#", &keywordValue{"", "?#"}},
			&testDefinition{":$#", &keywordValue{"", "$#"}},
			&testDefinition{":%#", &keywordValue{"", "%#"}},
			&testDefinition{":&#", &keywordValue{"", "&#"}},
			&testDefinition{":=#", &keywordValue{"", "=#"}},
			&testDefinition{":<#", &keywordValue{"", "<#"}},
			&testDefinition{":>#", &keywordValue{"", ">#"}},

			&testDefinition{":bar/*", &keywordValue{"bar", "*"}},
			&testDefinition{":bar/!", &keywordValue{"bar", "!"}},
			&testDefinition{":bar/?", &keywordValue{"bar", "?"}},
			&testDefinition{":bar/$", &keywordValue{"bar", "$"}},
			&testDefinition{":bar/%", &keywordValue{"bar", "%"}},
			&testDefinition{":bar/&", &keywordValue{"bar", "&"}},
			&testDefinition{":bar/=", &keywordValue{"bar", "="}},
			&testDefinition{":bar/<", &keywordValue{"bar", "<"}},
			&testDefinition{":bar/>", &keywordValue{"bar", ">"}},

			&testDefinition{":bar/*-", &keywordValue{"bar", "*-"}},
			&testDefinition{":bar/!-", &keywordValue{"bar", "!-"}},
			&testDefinition{":bar/?-", &keywordValue{"bar", "?-"}},
			&testDefinition{":bar/$-", &keywordValue{"bar", "$-"}},
			&testDefinition{":bar/%-", &keywordValue{"bar", "%-"}},
			&testDefinition{":bar/&-", &keywordValue{"bar", "&-"}},
			&testDefinition{":bar/=-", &keywordValue{"bar", "=-"}},
			&testDefinition{":bar/<-", &keywordValue{"bar", "<-"}},
			&testDefinition{":bar/>-", &keywordValue{"bar", ">-"}},

			&testDefinition{":bar/*+", &keywordValue{"bar", "*+"}},
			&testDefinition{":bar/!+", &keywordValue{"bar", "!+"}},
			&testDefinition{":bar/?+", &keywordValue{"bar", "?+"}},
			&testDefinition{":bar/$+", &keywordValue{"bar", "$+"}},
			&testDefinition{":bar/%+", &keywordValue{"bar", "%+"}},
			&testDefinition{":bar/&+", &keywordValue{"bar", "&+"}},
			&testDefinition{":bar/=+", &keywordValue{"bar", "=+"}},
			&testDefinition{":bar/<+", &keywordValue{"bar", "<+"}},
			&testDefinition{":bar/>+", &keywordValue{"bar", ">+"}},

			&testDefinition{":bar/*.", &keywordValue{"bar", "*."}},
			&testDefinition{":bar/!.", &keywordValue{"bar", "!."}},
			&testDefinition{":bar/?.", &keywordValue{"bar", "?."}},
			&testDefinition{":bar/$.", &keywordValue{"bar", "$."}},
			&testDefinition{":bar/%.", &keywordValue{"bar", "%."}},
			&testDefinition{":bar/&.", &keywordValue{"bar", "&."}},
			&testDefinition{":bar/=.", &keywordValue{"bar", "=."}},
			&testDefinition{":bar/<.", &keywordValue{"bar", "<."}},
			&testDefinition{":bar/>.", &keywordValue{"bar", ">."}},

			&testDefinition{":bar/*#", &keywordValue{"bar", "*#"}},
			&testDefinition{":bar/!#", &keywordValue{"bar", "!#"}},
			&testDefinition{":bar/?#", &keywordValue{"bar", "?#"}},
			&testDefinition{":bar/$#", &keywordValue{"bar", "$#"}},
			&testDefinition{":bar/%#", &keywordValue{"bar", "%#"}},
			&testDefinition{":bar/&#", &keywordValue{"bar", "&#"}},
			&testDefinition{":bar/=#", &keywordValue{"bar", "=#"}},
			&testDefinition{":bar/<#", &keywordValue{"bar", "<#"}},
			&testDefinition{":bar/>#", &keywordValue{"bar", ">#"}},

			&testDefinition{":*/bar", &keywordValue{"*", "bar"}},
			&testDefinition{":!/bar", &keywordValue{"!", "bar"}},
			&testDefinition{":?/bar", &keywordValue{"?", "bar"}},
			&testDefinition{":$/bar", &keywordValue{"$", "bar"}},
			&testDefinition{":%/bar", &keywordValue{"%", "bar"}},
			&testDefinition{":&/bar", &keywordValue{"&", "bar"}},
			&testDefinition{":=/bar", &keywordValue{"=", "bar"}},
			&testDefinition{":</bar", &keywordValue{"<", "bar"}},
			&testDefinition{":>/bar", &keywordValue{">", "bar"}},

			&testDefinition{":*-/bar", &keywordValue{"*-", "bar"}},
			&testDefinition{":!-/bar", &keywordValue{"!-", "bar"}},
			&testDefinition{":?-/bar", &keywordValue{"?-", "bar"}},
			&testDefinition{":$-/bar", &keywordValue{"$-", "bar"}},
			&testDefinition{":%-/bar", &keywordValue{"%-", "bar"}},
			&testDefinition{":&-/bar", &keywordValue{"&-", "bar"}},
			&testDefinition{":=-/bar", &keywordValue{"=-", "bar"}},
			&testDefinition{":<-/bar", &keywordValue{"<-", "bar"}},
			&testDefinition{":>-/bar", &keywordValue{">-", "bar"}},

			&testDefinition{":*+/bar", &keywordValue{"*+", "bar"}},
			&testDefinition{":!+/bar", &keywordValue{"!+", "bar"}},
			&testDefinition{":?+/bar", &keywordValue{"?+", "bar"}},
			&testDefinition{":$+/bar", &keywordValue{"$+", "bar"}},
			&testDefinition{":%+/bar", &keywordValue{"%+", "bar"}},
			&testDefinition{":&+/bar", &keywordValue{"&+", "bar"}},
			&testDefinition{":=+/bar", &keywordValue{"=+", "bar"}},
			&testDefinition{":<+/bar", &keywordValue{"<+", "bar"}},
			&testDefinition{":>+/bar", &keywordValue{">+", "bar"}},

			&testDefinition{":*./bar", &keywordValue{"*.", "bar"}},
			&testDefinition{":!./bar", &keywordValue{"!.", "bar"}},
			&testDefinition{":?./bar", &keywordValue{"?.", "bar"}},
			&testDefinition{":$./bar", &keywordValue{"$.", "bar"}},
			&testDefinition{":%./bar", &keywordValue{"%.", "bar"}},
			&testDefinition{":&./bar", &keywordValue{"&.", "bar"}},
			&testDefinition{":=./bar", &keywordValue{"=.", "bar"}},
			&testDefinition{":<./bar", &keywordValue{"<.", "bar"}},
			&testDefinition{":>./bar", &keywordValue{">.", "bar"}},

			&testDefinition{":*#/bar", &keywordValue{"*#", "bar"}},
			&testDefinition{":!#/bar", &keywordValue{"!#", "bar"}},
			&testDefinition{":?#/bar", &keywordValue{"?#", "bar"}},
			&testDefinition{":$#/bar", &keywordValue{"$#", "bar"}},
			&testDefinition{":%#/bar", &keywordValue{"%#", "bar"}},
			&testDefinition{":&#/bar", &keywordValue{"&#", "bar"}},
			&testDefinition{":=#/bar", &keywordValue{"=#", "bar"}},
			&testDefinition{":<#/bar", &keywordValue{"<#", "bar"}},
			&testDefinition{":>#/bar", &keywordValue{">#", "bar"}},
		)
	})
})
