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

package http

import (
	"github.com/Workiva/eva-client-go/edn"
	"github.com/Workiva/eva-client-go/eva"
	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection Channel Test", func() {

	Context("with the default marshaller", func() {

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			label := edn.NewStringElement("test")

			channel, err := newHttpConnChannel(label, &mockSource{})
			Ω(err).Should(BeNil())
			Ω(channel).ShouldNot(BeNil())
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())
				label := edn.NewStringElement("test")

				channel, err := newHttpConnChannel(label, source)
				Ω(err).Should(BeNil())
				Ω(channel).ShouldNot(BeNil())

				result, err := channel.Transact("")
				Ω(err).Should(BeNil())
				Ω(result).ShouldNot(BeNil())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})

		It("compile the wildcard pattern correctly", func() {

			var err error
			var config eva.Configuration
			config, err = eva.NewConfiguration(`{
				"source": {
					"type":   "http",
					"server": "localhost"
				},
				"category": "test"
			}`)
			Ω(err).Should(BeNil())

			var tenant eva.Tenant
			tenant, err = eva.NewTenant("tenant")

			source, err := initHttpSource(config, tenant)
			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			if httpSource, is := source.(*httpSourceImpl); is {
				httpSource.callClient = fakeGoodCaller(edn.EvaEdnMimeType.String())
				label := edn.NewStringElement("test")

				channel, err := newHttpConnChannel(label, &mockSource{})
				Ω(err).Should(BeNil())
				Ω(channel).ShouldNot(BeNil())

				result, err := channel.Transact("")
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(test.HaveMessage(ErrUnsupportedType))
				Ω(result).Should(BeNil())
			} else {
				Fail("Expected the binding to be a *httpSourceImpl")
			}
		})
	})
})
