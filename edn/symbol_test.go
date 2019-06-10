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

var _ = Describe("Symbol in EDN", func() {
	Context("with the default marshaller", func() {

		goodSymbols := map[string]struct {
			prefix string
			name   string
		}{
			"fine/good":    {"fine", "good"},
			"fine:good":    {"", "fine:good"},
			"fine#good":    {"", "fine#good"},
			"fine-good":    {"", "fine-good"},
			"ok/fine-good": {"ok", "fine-good"},
			"a":            {"", "a"},
			"good":         {"", "good"},
			"good1":        {"", "good1"},
			"fine/+good":   {"fine", "+good"},
			"fine/>good":   {"fine", ">good"},
			"*good":        {"", "*good"},
			"*1good":       {"", "*1good"},
			"!good":        {"", "!good"},
			"!1good":       {"", "!1good"},
			"_good":        {"", "_good"},
			"_1good":       {"", "_1good"},
			"?good":        {"", "?good"},
			"?1good":       {"", "?1good"},
			"$good":        {"", "$good"},
			"$1good":       {"", "$1good"},
			"%good":        {"", "%good"},
			"%1good":       {"", "%1good"},
			"&good":        {"", "&good"},
			"&1good":       {"", "&1good"},
			"=good":        {"", "=good"},
			"=1good":       {"", "=1good"},
			"<good":        {"", "<good"},
			"<1good":       {"", "<1good"},
			">good":        {"", ">good"},
			">1good":       {"", ">1good"},
			"+good":        {"", "+good"},
			"-good":        {"", "-good"},
			".good":        {"", ".good"},
			"*":            {"", "*"},
			"!":            {"", "!"},
			"_":            {"", "_"},
			"?":            {"", "?"},
			"$":            {"", "$"},
			"%":            {"", "%"},
			"&":            {"", "&"},
			"=":            {"", "="},
			"<":            {"", "<"},
			">":            {"", ">"},
			"/":            {"", "/"},
			".":            {"", "."},
			"-":            {"", "-"},
			"+":            {"", "+"},
		}

		badSymbols := []string{
			":/",
			"::",
			":",
			":bad",
			"#bad",
			"1bad",
			"bad1/1worse",
			"/bad",
			"bad/",
			"bad/worse/wrong",
			"worse/+1bad",
			"+1bad",
			"-1bad",
			".1bad",
			".1",
		}

		It("should create a symbol with just one parameter value with no error", func() {
			prefix := ""
			name := "foobar"

			elem, err := NewSymbolElement(name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(SymbolType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should create a symbol with two parameter value with no error", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewSymbolElement(prefix, name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(SymbolType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should create a symbol with one parameter (but with the separator) value with no error", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewSymbolElement(prefix + SymbolSeparator + name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(SymbolType))
			Ω(elem.Prefix()).Should(BeEquivalentTo(prefix))
			Ω(elem.Name()).Should(BeEquivalentTo(name))
		})

		It("should be equal to each other it they are the same.", func() {
			prefix := "namespace"
			name := "foobar"

			elem1, err1 := NewSymbolElement(prefix, name)
			Ω(err1).Should(BeNil())
			Ω(elem1).ShouldNot(BeNil())
			elem2, err2 := NewSymbolElement(prefix + SymbolSeparator + name)
			Ω(err2).Should(BeNil())
			Ω(elem2).ShouldNot(BeNil())
			Ω(elem1.Equals(elem2)).Should(BeTrue())
		})

		It("should not equal to each other it they are the not the same.", func() {
			prefix := "namespace"
			name := "foobar"

			elem1, err1 := NewSymbolElement(prefix + SymbolSeparator + name + "2")
			Ω(err1).Should(BeNil())
			Ω(elem1).ShouldNot(BeNil())
			elem2, err2 := NewSymbolElement(prefix + SymbolSeparator + name)
			Ω(err2).Should(BeNil())
			Ω(elem2).ShouldNot(BeNil())
			Ω(elem1.Equals(elem2)).Should(BeFalse())
		})

		It("should not equal to each other it they are the not the same.", func() {
			namespace := "namespace"
			prefix := "prefix"
			name := "foobar"

			elem1, err1 := NewSymbolElement(namespace + NamespaceSeparator + prefix + SymbolSeparator + name + "2")
			Ω(err1).Should(BeNil())
			Ω(elem1).ShouldNot(BeNil())
			elem2, err2 := NewSymbolElement(namespace + NamespaceSeparator + prefix + SymbolSeparator + name)
			Ω(err2).Should(BeNil())
			Ω(elem2).ShouldNot(BeNil())
			Ω(elem1.Equals(elem2)).Should(BeFalse())
		})

		It("should create a symbol with zero parameter value with an error", func() {
			elem, err := NewSymbolElement()
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidSymbol))
		})

		It("should create a symbol with three parameter value with an error", func() {
			elem, err := NewSymbolElement("a", "b", "c")
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidSymbol))
		})

		It("should serialize the symbol with one parameter without an issue", func() {
			name := "foobar"

			elem, err := NewSymbolElement(name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(name))
		})

		It("should serialize the symbol with one parameter without an issue", func() {
			name := "foobar"

			elem, err := NewSymbolElement(name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			_, err = elem.Serialize(SerializerMimeType("InvalidSerializer"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})

		It("should serialize the symbol with one parameter without an issue", func() {
			name := "foobar"

			elem, err := NewSymbolElement(name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(name))
		})

		It("validate the namespace", func() {
			Ω(IsValidNamespace("namespace")).Should(BeTrue())
		})

		It("should serialize the symbol with one (but with the separator) parameter without an issue", func() {
			prefix := "namespace"
			name := "foobar"

			elem, err := NewSymbolElement(prefix + SymbolSeparator + name)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo(prefix + SymbolSeparator + name))
		})

		It("should not create an element with a bad namespace", func() {
			prefix := "1bad"
			name := "foobar"

			elem, err := NewSymbolElement(prefix, name)
			Ω(elem).Should(BeNil())
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidSymbol))
		})

		It("should not create an element with a bad name", func() {
			prefix := "namespace"
			name := "1bad"

			elem, err := NewSymbolElement(prefix, name)
			Ω(elem).Should(BeNil())
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidSymbol))
		})

		It("should process all odd, but good symbols", func() {
			for symbol, result := range goodSymbols {

				message := fmt.Sprintf("Expecting good from: %s", symbol)

				elem, err := NewSymbolElement(symbol)
				Ω(err).Should(BeNil(), message)
				Ω(elem).ShouldNot(BeNil(), message)

				edn, err := elem.Serialize(EvaEdnMimeType)
				Ω(err).Should(BeNil(), message)
				Ω(edn).Should(BeEquivalentTo(symbol), message)
				Ω(elem.Prefix()).Should(BeEquivalentTo(result.prefix), message)
				Ω(elem.Name()).Should(BeEquivalentTo(result.name), message)
			}
		})

		It("should not process all odd invalid symbols", func() {

			for _, symbol := range badSymbols {
				elem, err := NewSymbolElement(symbol)
				Ω(elem).Should(BeNil())
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(test.HaveMessage(ErrInvalidSymbol))
			}
		})
	})

	Context("Parsing", func() {
		runParserTests(SymbolType,
			&testDefinition{"foo", &keywordValue{"", "foo"}},
			&testDefinition{"bar/foo", &keywordValue{"bar", "foo"}},

			// TODO: These are not working...
			// &testDefinition{".", &keywordValue{"", "."}},
			// &testDefinition{"...", &keywordValue{"", "..."}},

			&testDefinition{"*", &keywordValue{"", "*"}},
			&testDefinition{"!", &keywordValue{"", "!"}},
			&testDefinition{"?", &keywordValue{"", "?"}},
			&testDefinition{"$", &keywordValue{"", "$"}},
			&testDefinition{"%", &keywordValue{"", "%"}},
			&testDefinition{"&", &keywordValue{"", "&"}},
			&testDefinition{"=", &keywordValue{"", "="}},
			&testDefinition{"<", &keywordValue{"", "<"}},
			&testDefinition{">", &keywordValue{"", ">"}},

			&testDefinition{"*-", &keywordValue{"", "*-"}},
			&testDefinition{"!-", &keywordValue{"", "!-"}},
			&testDefinition{"?-", &keywordValue{"", "?-"}},
			&testDefinition{"$-", &keywordValue{"", "$-"}},
			&testDefinition{"%-", &keywordValue{"", "%-"}},
			&testDefinition{"&-", &keywordValue{"", "&-"}},
			&testDefinition{"=-", &keywordValue{"", "=-"}},
			&testDefinition{"<-", &keywordValue{"", "<-"}},
			&testDefinition{">-", &keywordValue{"", ">-"}},

			&testDefinition{"*+", &keywordValue{"", "*+"}},
			&testDefinition{"!+", &keywordValue{"", "!+"}},
			&testDefinition{"?+", &keywordValue{"", "?+"}},
			&testDefinition{"$+", &keywordValue{"", "$+"}},
			&testDefinition{"%+", &keywordValue{"", "%+"}},
			&testDefinition{"&+", &keywordValue{"", "&+"}},
			&testDefinition{"=+", &keywordValue{"", "=+"}},
			&testDefinition{"<+", &keywordValue{"", "<+"}},
			&testDefinition{">+", &keywordValue{"", ">+"}},

			&testDefinition{"*.", &keywordValue{"", "*."}},
			&testDefinition{"!.", &keywordValue{"", "!."}},
			&testDefinition{"?.", &keywordValue{"", "?."}},
			&testDefinition{"$.", &keywordValue{"", "$."}},
			&testDefinition{"%.", &keywordValue{"", "%."}},
			&testDefinition{"&.", &keywordValue{"", "&."}},
			&testDefinition{"=.", &keywordValue{"", "=."}},
			&testDefinition{"<.", &keywordValue{"", "<."}},
			&testDefinition{">.", &keywordValue{"", ">."}},

			&testDefinition{"*#", &keywordValue{"", "*#"}},
			&testDefinition{"!#", &keywordValue{"", "!#"}},
			&testDefinition{"?#", &keywordValue{"", "?#"}},
			&testDefinition{"$#", &keywordValue{"", "$#"}},
			&testDefinition{"%#", &keywordValue{"", "%#"}},
			&testDefinition{"&#", &keywordValue{"", "&#"}},
			&testDefinition{"=#", &keywordValue{"", "=#"}},
			&testDefinition{"<#", &keywordValue{"", "<#"}},
			&testDefinition{">#", &keywordValue{"", ">#"}},

			&testDefinition{"bar/*", &keywordValue{"bar", "*"}},
			&testDefinition{"bar/!", &keywordValue{"bar", "!"}},
			&testDefinition{"bar/?", &keywordValue{"bar", "?"}},
			&testDefinition{"bar/$", &keywordValue{"bar", "$"}},
			&testDefinition{"bar/%", &keywordValue{"bar", "%"}},
			&testDefinition{"bar/&", &keywordValue{"bar", "&"}},
			&testDefinition{"bar/=", &keywordValue{"bar", "="}},
			&testDefinition{"bar/<", &keywordValue{"bar", "<"}},
			&testDefinition{"bar/>", &keywordValue{"bar", ">"}},

			&testDefinition{"bar/*-", &keywordValue{"bar", "*-"}},
			&testDefinition{"bar/!-", &keywordValue{"bar", "!-"}},
			&testDefinition{"bar/?-", &keywordValue{"bar", "?-"}},
			&testDefinition{"bar/$-", &keywordValue{"bar", "$-"}},
			&testDefinition{"bar/%-", &keywordValue{"bar", "%-"}},
			&testDefinition{"bar/&-", &keywordValue{"bar", "&-"}},
			&testDefinition{"bar/=-", &keywordValue{"bar", "=-"}},
			&testDefinition{"bar/<-", &keywordValue{"bar", "<-"}},
			&testDefinition{"bar/>-", &keywordValue{"bar", ">-"}},

			&testDefinition{"bar/*+", &keywordValue{"bar", "*+"}},
			&testDefinition{"bar/!+", &keywordValue{"bar", "!+"}},
			&testDefinition{"bar/?+", &keywordValue{"bar", "?+"}},
			&testDefinition{"bar/$+", &keywordValue{"bar", "$+"}},
			&testDefinition{"bar/%+", &keywordValue{"bar", "%+"}},
			&testDefinition{"bar/&+", &keywordValue{"bar", "&+"}},
			&testDefinition{"bar/=+", &keywordValue{"bar", "=+"}},
			&testDefinition{"bar/<+", &keywordValue{"bar", "<+"}},
			&testDefinition{"bar/>+", &keywordValue{"bar", ">+"}},

			&testDefinition{"bar/*.", &keywordValue{"bar", "*."}},
			&testDefinition{"bar/!.", &keywordValue{"bar", "!."}},
			&testDefinition{"bar/?.", &keywordValue{"bar", "?."}},
			&testDefinition{"bar/$.", &keywordValue{"bar", "$."}},
			&testDefinition{"bar/%.", &keywordValue{"bar", "%."}},
			&testDefinition{"bar/&.", &keywordValue{"bar", "&."}},
			&testDefinition{"bar/=.", &keywordValue{"bar", "=."}},
			&testDefinition{"bar/<.", &keywordValue{"bar", "<."}},
			&testDefinition{"bar/>.", &keywordValue{"bar", ">."}},

			&testDefinition{"bar/*#", &keywordValue{"bar", "*#"}},
			&testDefinition{"bar/!#", &keywordValue{"bar", "!#"}},
			&testDefinition{"bar/?#", &keywordValue{"bar", "?#"}},
			&testDefinition{"bar/$#", &keywordValue{"bar", "$#"}},
			&testDefinition{"bar/%#", &keywordValue{"bar", "%#"}},
			&testDefinition{"bar/&#", &keywordValue{"bar", "&#"}},
			&testDefinition{"bar/=#", &keywordValue{"bar", "=#"}},
			&testDefinition{"bar/<#", &keywordValue{"bar", "<#"}},
			&testDefinition{"bar/>#", &keywordValue{"bar", ">#"}},

			&testDefinition{"*/bar", &keywordValue{"*", "bar"}},
			&testDefinition{"!/bar", &keywordValue{"!", "bar"}},
			&testDefinition{"?/bar", &keywordValue{"?", "bar"}},
			&testDefinition{"$/bar", &keywordValue{"$", "bar"}},
			&testDefinition{"%/bar", &keywordValue{"%", "bar"}},
			&testDefinition{"&/bar", &keywordValue{"&", "bar"}},
			&testDefinition{"=/bar", &keywordValue{"=", "bar"}},
			&testDefinition{"</bar", &keywordValue{"<", "bar"}},
			&testDefinition{">/bar", &keywordValue{">", "bar"}},

			&testDefinition{"*-/bar", &keywordValue{"*-", "bar"}},
			&testDefinition{"!-/bar", &keywordValue{"!-", "bar"}},
			&testDefinition{"?-/bar", &keywordValue{"?-", "bar"}},
			&testDefinition{"$-/bar", &keywordValue{"$-", "bar"}},
			&testDefinition{"%-/bar", &keywordValue{"%-", "bar"}},
			&testDefinition{"&-/bar", &keywordValue{"&-", "bar"}},
			&testDefinition{"=-/bar", &keywordValue{"=-", "bar"}},
			&testDefinition{"<-/bar", &keywordValue{"<-", "bar"}},
			&testDefinition{">-/bar", &keywordValue{">-", "bar"}},

			&testDefinition{"*+/bar", &keywordValue{"*+", "bar"}},
			&testDefinition{"!+/bar", &keywordValue{"!+", "bar"}},
			&testDefinition{"?+/bar", &keywordValue{"?+", "bar"}},
			&testDefinition{"$+/bar", &keywordValue{"$+", "bar"}},
			&testDefinition{"%+/bar", &keywordValue{"%+", "bar"}},
			&testDefinition{"&+/bar", &keywordValue{"&+", "bar"}},
			&testDefinition{"=+/bar", &keywordValue{"=+", "bar"}},
			&testDefinition{"<+/bar", &keywordValue{"<+", "bar"}},
			&testDefinition{">+/bar", &keywordValue{">+", "bar"}},

			&testDefinition{"*./bar", &keywordValue{"*.", "bar"}},
			&testDefinition{"!./bar", &keywordValue{"!.", "bar"}},
			&testDefinition{"?./bar", &keywordValue{"?.", "bar"}},
			&testDefinition{"$./bar", &keywordValue{"$.", "bar"}},
			&testDefinition{"%./bar", &keywordValue{"%.", "bar"}},
			&testDefinition{"&./bar", &keywordValue{"&.", "bar"}},
			&testDefinition{"=./bar", &keywordValue{"=.", "bar"}},
			&testDefinition{"<./bar", &keywordValue{"<.", "bar"}},
			&testDefinition{">./bar", &keywordValue{">.", "bar"}},

			&testDefinition{"*#/bar", &keywordValue{"*#", "bar"}},
			&testDefinition{"!#/bar", &keywordValue{"!#", "bar"}},
			&testDefinition{"?#/bar", &keywordValue{"?#", "bar"}},
			&testDefinition{"$#/bar", &keywordValue{"$#", "bar"}},
			&testDefinition{"%#/bar", &keywordValue{"%#", "bar"}},
			&testDefinition{"&#/bar", &keywordValue{"&#", "bar"}},
			&testDefinition{"=#/bar", &keywordValue{"=#", "bar"}},
			&testDefinition{"<#/bar", &keywordValue{"<#", "bar"}},
			&testDefinition{">#/bar", &keywordValue{">#", "bar"}},
		)
	})
})
