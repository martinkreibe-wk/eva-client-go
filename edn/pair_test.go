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

var _ = Describe("Pair for maps in EDN", func() {
	Context("with the default usage", func() {
		It("should create a pair with no error", func() {
			key, err := NewStringElement("key")
			Ω(err).Should(BeNil())
			value, err := NewStringElement("value")
			Ω(err).Should(BeNil())

			pair, err := NewPair(key, value)
			Ω(err).Should(BeNil())
			Ω(pair).ShouldNot(BeNil())
			Ω(pair.Key()).Should(BeEquivalentTo(key))
			Ω(pair.Value()).Should(BeEquivalentTo(value))
		})

		It("should create an error with nil key", func() {
			value, err := NewStringElement("value")
			Ω(err).Should(BeNil())

			pair, err := NewPair(nil, value)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidPair))
			Ω(pair).Should(BeNil())
		})

		It("should create an error with nil value", func() {
			key, err := NewStringElement("key")
			Ω(err).Should(BeNil())

			pair, err := NewPair(key, nil)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidPair))
			Ω(pair).Should(BeNil())
		})

		It("should append a pair to the pair collection", func() {
			key, err := NewStringElement("key")
			Ω(err).Should(BeNil())

			value, err := NewStringElement("value")
			Ω(err).Should(BeNil())

			pairs := &Pairs{}
			err = pairs.Append(key, value)
			Ω(err).Should(BeNil())
			Ω(pairs.Len()).Should(BeEquivalentTo(1))
			Ω(pairs.Raw()).Should(HaveLen(1))

		})

		It("should append a pair to the pair collection", func() {
			key, err := NewStringElement("key")
			Ω(err).Should(BeNil())

			value, err := NewStringElement("value")
			Ω(err).Should(BeNil())

			pairs := &Pairs{}
			err = pairs.Append(key, value)
			Ω(err).Should(BeNil())

			elems := pairs.RawElements()
			Ω(elems).Should(HaveLen(2))
		})
	})
})
