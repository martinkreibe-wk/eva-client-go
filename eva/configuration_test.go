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
	"fmt"
	"github.com/Workiva/eva-client-go/edn"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func makeConfig(source string, category string) string {
	return makeConfigWithMime(source, category, edn.EvaEdnMimeType.String())
}

func makeConfigWithMime(source string, category string, mime string) string {

	return fmt.Sprintf(`{
		"source": {
			"type": "%s",
			"mime": "%s",
			"thing": "value"
		},
		"category": "%s"
	}`, source, mime, category)
}

var _ = Describe("Configuration test", func() {

	category := "example"
	sourceType := "s"

	Context("normally", func() {

		It("should error with bad service name", func() {
			config, err := NewConfiguration(makeConfig(sourceType, category))
			Ω(err).Should(BeNil())
			Ω(config).ShouldNot(BeNil())

			Ω(config.Category()).Should(BeEquivalentTo(category))

			var src SourceConfiguration
			src, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(src).ShouldNot(BeNil())

			var ser edn.Serializer
			ser, err = src.Serializer()
			Ω(err).Should(BeNil())

			Ω(ser.MimeType()).Should(BeEquivalentTo(edn.EvaEdnMimeType))

			var srcType string
			srcType = src.Type()
			Ω(srcType).Should(BeEquivalentTo(sourceType))
			val, has := src.Setting("thing")
			Ω(has).Should(BeTrue())
			Ω(val).Should(BeEquivalentTo("value"))
		})
	})

	Context("normally", func() {

		It("should error with bad service name", func() {
			config, err := NewConfiguration(fmt.Sprintf(`{
		"source": {
			"type": "%s"
		},
		"category": "%s"
	}`, sourceType, category))

			Ω(err).Should(BeNil())
			Ω(config).ShouldNot(BeNil())

			Ω(config.Category()).Should(BeEquivalentTo(category))

			var src SourceConfiguration
			src, err = config.Source()
			Ω(err).Should(BeNil())
			Ω(src).ShouldNot(BeNil())

			var ser edn.Serializer
			ser, err = src.Serializer()
			Ω(err).Should(BeNil())

			Ω(ser.MimeType()).Should(BeEquivalentTo(edn.EvaEdnMimeType))
		})
	})

	Context("normally", func() {

		It("should error with bad service name", func() {
			config, err := NewConfiguration("{Nope = witha, a side ] of nope")

			Ω(err).ShouldNot(BeNil())
			Ω(config).Should(BeNil())
		})
	})
})
