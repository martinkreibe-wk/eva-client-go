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

var _ = Describe("base source test", func() {

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

			_, err = source.Query(nil, nil)
			Ω(err).Should(BeNil())
		})
	})

	Context("with the default marshaller", func() {
		It("", func() {
			config, err := NewConfiguration("{\"category\": \"foo\"}")
			Ω(err).Should(BeNil())
			Ω(config).ShouldNot(BeNil())

			tenant, err := NewTenant("foo")
			Ω(err).Should(BeNil())

			_, err = NewBaseSource(config, tenant, &mockSource{}, nil, nil)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidConfiguration))
		})
	})

	Context("with the default marshaller", func() {
		It("", func() {
			config, err := NewConfiguration("{\"category\": \"\"}")
			Ω(err).Should(BeNil())
			Ω(config).ShouldNot(BeNil())

			tenant, err := NewTenant("foo")
			Ω(err).Should(BeNil())

			_, err = NewBaseSource(config, tenant, &mockSource{}, makeMockConnChannel, mockQuery)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidConfiguration))
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
			Ω(source.Tenant()).ShouldNot(BeNil())
			Ω(source.Tenant()).Should(BeEquivalentTo(tenant))
			Ω(source.Category()).Should(BeEquivalentTo("foo"))

			serialzier, err := source.Serializer()
			Ω(err).Should(BeNil())
			Ω(serialzier).ShouldNot(BeNil())
			Ω(serialzier).Should(BeEquivalentTo(edn.EvaEdnMimeType))
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

			var conn ConnectionChannel
			conn, err = source.Connection("label")
			Ω(err).Should(BeNil())
			Ω(conn).ShouldNot(BeNil())
			Ω(conn.Label()).Should(BeEquivalentTo("label"))

			conn, err = source.Connection(edn.NewStringElement("label2"))
			Ω(err).Should(BeNil())
			Ω(conn).ShouldNot(BeNil())
			Ω(conn.Label()).Should(BeEquivalentTo("label2"))
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
			Ω(snap.Label()).Should(BeEquivalentTo("label"))

			snap, err = source.LatestSnapshot(edn.NewStringElement("label2"))
			Ω(err).Should(BeNil())
			Ω(snap).ShouldNot(BeNil())
			Ω(snap.Label()).Should(BeEquivalentTo("label2"))

			snap, err = source.AsOfSnapshot("label3", 123)
			Ω(err).Should(BeNil())
			Ω(snap).ShouldNot(BeNil())
			Ω(snap.Label()).Should(BeEquivalentTo("label3"))
			Ω(snap.AsOf()).ShouldNot(BeNil())
			Ω(*snap.AsOf()).Should(BeEquivalentTo(123))

			snap, err = source.AsOfSnapshot(edn.NewStringElement("label4"), edn.NewIntegerElement(456))
			Ω(err).Should(BeNil())
			Ω(snap).ShouldNot(BeNil())
			Ω(snap.Label()).Should(BeEquivalentTo("label4"))
			Ω(snap.AsOf()).ShouldNot(BeNil())
			Ω(*snap.AsOf()).Should(BeEquivalentTo(456))
		})
	})

})
