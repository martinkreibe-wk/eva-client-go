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

var _ = Describe("Serializer tests", func() {
	Context("", func() {
		It("the string version of the getter should be easily gotten", func() {
			ser, err := GetSerializer(string(EvaEdnMimeType))
			Ω(err).Should(BeNil())
			Ω(ser).ShouldNot(BeNil())
			Ω(ser).Should(BeEquivalentTo(EvaEdnMimeType))

			op, has := ser.Options("anything")
			Ω(op).Should(BeEquivalentTo(""))
			Ω(has).Should(BeFalse())
		})

		It("the string version of the getter should be easily gotten with options", func() {
			ser, err := GetSerializer(string(EvaEdnMimeType) + ";option=foo")
			Ω(err).Should(BeNil())
			Ω(ser).ShouldNot(BeNil())
			Ω(ser).ShouldNot(BeEquivalentTo(EvaEdnMimeType))

			op, has := ser.Options("option")
			Ω(has).Should(BeTrue())
			Ω(op).Should(BeEquivalentTo("foo"))
		})

		It("the string version of the getter should be easily gotten", func() {
			ser, err := GetSerializer("trash")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
			Ω(ser).Should(BeNil())
		})
	})
})
