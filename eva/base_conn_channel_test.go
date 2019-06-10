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

var _ = Describe("base connection channel test", func() {

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

			var result Result
			result, err = conn.Transact("foo")
			Ω(err).Should(BeNil())
			Ω(result).ShouldNot(BeNil())

			var str string
			var has bool
			str, has = result.String()
			Ω(has).Should(BeTrue())
			Ω(str).Should(BeEquivalentTo("test"))

			result, err = conn.Transact(edn.NewStringElement("trx"))
			Ω(err).Should(BeNil())
			Ω(result).ShouldNot(BeNil())

			str, has = result.String()
			Ω(has).Should(BeTrue())
			Ω(str).Should(BeEquivalentTo("test"))

			result, err = conn.Transact()
			Ω(err).ShouldNot(BeNil())
			Ω(result).Should(BeNil())
			Ω(err).Should(test.HaveMessage(edn.ErrInvalidInput))
		})
	})
})
