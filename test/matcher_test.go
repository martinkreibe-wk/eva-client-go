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

package test

import (
	"github.com/Workiva/eva-client-go/edn"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Matcher test", func() {

	Context("in a test environment", func() {

		It("not panic if there is no error", func() {

			matcher := HaveMessage(edn.ErrInvalidInput)
			Ω(matcher).ShouldNot(BeNil())

			Ω(matcher.FailureMessage("foo")).Should(BeEquivalentTo("Expected\n\t\"foo\"\nhave message\n\t\"" + edn.ErrInvalidInput.Message() + "\""))
			Ω(matcher.NegatedFailureMessage("foo")).Should(BeEquivalentTo("Expected\n\t\"foo\"\nnot to have message\n\t\"" + edn.ErrInvalidInput.Message() + "\""))

			success, err := matcher.Match("foo")
			Ω(success).Should(BeFalse())
			Ω(err).ShouldNot(BeNil())
			Ω(err.Error()).Should(BeEquivalentTo("Expected actual to have a Message() string method, instead found:  string"))

			success, err = matcher.Match(edn.ErrInvalidInput)
			Ω(success).Should(BeTrue())
			Ω(err).Should(BeNil())
		})
	})
})
