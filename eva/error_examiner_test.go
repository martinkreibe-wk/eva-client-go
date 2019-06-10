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

package eva

import (
	"github.com/Workiva/eva-client-go/edn"
	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error examiner", func() {

	Context("GetErrorExaminer", func() {

		It("with nil", func() {
			examiner, err := GetErrorExaminer(nil)
			Ω(examiner).ShouldNot(BeNil())
			Ω(err).Should(BeNil())
		})

		It("with bad mime type", func() {
			examiner, err := GetErrorExaminer(edn.SerializerMimeType("nothing"))
			Ω(examiner).Should(BeNil())
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidSerializer))
		})

		It("with good mime type", func() {
			examiner, err := GetErrorExaminer(edn.EvaEdnMimeType)
			Ω(examiner).ShouldNot(BeNil())
			Ω(err).Should(BeNil())
		})
	})

	Context("ednErrorExaminer", func() {
		It("with nil", func() {
			err := ednErrorExaminer(nil)
			Ω(err).ShouldNot(BeNil())
		})

		It("with empty", func() {
			err := ednErrorExaminer([]byte(""))
			Ω(err).ShouldNot(BeNil())
		})

		It("with empty", func() {

			example := []byte(`{
	:message "",
	:ex-info {
		:explanation "Malformed transact request.",
		:type "IncorrectTransactSyntax",
		:code 3000}, 
	:ex-data ""
}`)

			err := ednErrorExaminer(example)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(BeAssignableToTypeOf(&clientErrorImpl{}))

			clientErr := err.(*clientErrorImpl)
			Ω(clientErr).ShouldNot(BeNil())

			Ω(clientErr.Error()).Should(ContainSubstring(ErrSourceError.Message()))
		})
	})
})
