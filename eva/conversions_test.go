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

var _ = Describe("Decoding items", func() {

	Context("all types", func() {

		It("nil", func() {
			serializable, err := decodeSerializable(nil)
			Ω(serializable).Should(BeNil())
			Ω(err).Should(BeNil())
		})

		It("edn.Element", func() {
			serializable, err := decodeSerializable(edn.NewStringElement("foo"))
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("\"foo\""))
			Ω(err).Should(BeNil())
		})

		It("string", func() {
			serializable, err := decodeSerializable("foo")
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("foo"))
			Ω(err).Should(BeNil())
		})

		It("raw string", func() {
			serializable, err := decodeSerializable(RawString("foo"))
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("foo"))
			Ω(err).Should(BeNil())
		})

		It("int", func() {
			serializable, err := decodeSerializable(int(123))
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("123"))
			Ω(err).Should(BeNil())
		})

		It("int8", func() {
			serializable, err := decodeSerializable(int8(123))
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("123"))
			Ω(err).Should(BeNil())
		})

		It("int16", func() {
			serializable, err := decodeSerializable(int16(123))
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("123"))
			Ω(err).Should(BeNil())
		})

		It("int32", func() {
			serializable, err := decodeSerializable(int32(123))
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("123"))
			Ω(err).Should(BeNil())
		})

		It("int64", func() {
			serializable, err := decodeSerializable(int64(123))
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("123"))
			Ω(err).Should(BeNil())
		})

		It("raw int", func() {
			serializable, err := decodeSerializable(RawInt(123))
			Ω(serializable).ShouldNot(BeNil())
			Ω(serializable.String()).Should(BeEquivalentTo("123"))
			Ω(err).Should(BeNil())
		})

		It("should error on unknown type", func() {
			serializable, err := decodeSerializable(&mockSource{})
			Ω(serializable).Should(BeNil())
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(edn.ErrInvalidInput))
		})

	})
})
