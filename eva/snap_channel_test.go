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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("snapshot channel test", func() {

	Context("with the default marshaller", func() {
		It("", func() {
			config, err := NewConfiguration("{\"category\": \"foo\"}")
			Ω(err).Should(BeNil())
			Ω(config).ShouldNot(BeNil())

			tenant, err := NewTenant("foo")
			Ω(err).Should(BeNil())

			source, err := NewBaseSource(config, tenant, &mockSource{}, makeMockConnChannel, mockQuery)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			var snap SnapshotChannel
			snap, err = source.LatestSnapshot("label")
			Ω(err).Should(BeNil())
			Ω(snap).ShouldNot(BeNil())
			label, err := snap.Label()
			Ω(err).Should(BeNil())
			Ω(label).Should(BeEquivalentTo("label"))

			var result Result
			result, err = snap.Pull("[*]", "123")
			Ω(err).Should(BeNil())

			var test string
			var has bool
			test, has = result.String()
			Ω(has).Should(BeTrue())
			Ω(test).Should(BeEquivalentTo("test"))
		})
	})

	Context("with the default marshaller", func() {
		It("", func() {
			config, err := NewConfiguration("{\"category\": \"foo\"}")
			Ω(err).Should(BeNil())
			Ω(config).ShouldNot(BeNil())

			tenant, err := NewTenant("foo")
			Ω(err).Should(BeNil())

			source, err := NewBaseSource(config, tenant, &mockSource{}, makeMockConnChannel, mockQuery)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			var snap SnapshotChannel
			snap, err = source.LatestSnapshot("label")
			Ω(err).Should(BeNil())
			Ω(snap).ShouldNot(BeNil())
			label, err := snap.Label()
			Ω(err).Should(BeNil())
			Ω(label).Should(BeEquivalentTo("label"))

			var result Result
			result, err = snap.Invoke("func")
			Ω(err).Should(BeNil())

			var test string
			var has bool
			test, has = result.String()
			Ω(has).Should(BeTrue())
			Ω(test).Should(BeEquivalentTo("test"))
		})
	})
})
