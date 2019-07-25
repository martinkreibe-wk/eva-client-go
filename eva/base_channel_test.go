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

var _ = Describe("Base Channel test", func() {

	Context("normally", func() {

		It("should fail with invalid source.", func() {

			ct := ChannelType("test")

			var channel *BaseChannel
			var err error
			channel, err = NewBaseChannel(ct, nil, nil)
			Ω(channel).Should(BeNil())
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidSource))
		})

		It("should accept references with no fields", func() {

			ct := ChannelType("test")
			src := &mockSource{}

			var channel *BaseChannel
			var err error
			channel, err = NewBaseChannel(ct, src, nil)
			Ω(channel).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			Ω(channel.Type()).Should(BeEquivalentTo(ct))
			Ω(channel.Source()).Should(BeEquivalentTo(src))

			Ω(channel.Reference()).ShouldNot(BeNil())
			Ω(channel.Reference().Type()).Should(BeEquivalentTo(ct))
		})

		It("should accept references with fields", func() {

			ct := ChannelType("test")
			src := &mockSource{}

			elem, err := edn.NewStringElement("test")
			Ω(err).Should(BeNil())

			var channel *BaseChannel
			channel, err = NewBaseChannel(ct, src, map[string]edn.Element{
				"foo": elem,
			})
			Ω(channel).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			Ω(channel.Type()).Should(BeEquivalentTo(ct))
			Ω(channel.Source()).Should(BeEquivalentTo(src))

			Ω(channel.Reference()).ShouldNot(BeNil())
			Ω(channel.Reference().Type()).Should(BeEquivalentTo(ct))

			ref := channel.Reference()
			Ω(err).Should(BeNil())

			prop, err := ref.GetProperty("foo")
			Ω(err).Should(BeNil())

			Ω(prop.Value()).Should(BeEquivalentTo("test"))
		})

		It("should accept references with fields", func() {

			ct := ChannelType("test")
			src := &mockSource{}

			elem, err := edn.NewIntegerElement(123)
			Ω(err).Should(BeNil())

			var channel *BaseChannel
			channel, err = NewBaseChannel(ct, src, map[string]edn.Element{
				LabelReferenceProperty: elem,
			})
			Ω(channel).ShouldNot(BeNil())
			Ω(err).Should(BeNil())

			Ω(channel.Type()).Should(BeEquivalentTo(ct))
			Ω(channel.Source()).Should(BeEquivalentTo(src))

			label, err := channel.Label()
			Ω(err).Should(BeNil())
			Ω(label).Should(BeEquivalentTo("123"))
		})
	})
})
