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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Elements in EDN", func() {
	Context("with the default marshaller", func() {
		It("should handle complex embedding", func() {

			part, err := NewKeywordElement("db.part/db")
			Ω(err).Should(BeNil())

			var id CollectionElement
			id, err = NewVector(part)
			Ω(err).Should(BeNil())

			err = id.SetTag("db/id")
			Ω(err).Should(BeNil())

			var str string
			str, err = id.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())

			Ω(str).Should(BeEquivalentTo("#db/id [:db.part/db]"))

			key, err := NewKeywordElement("db/id")
			Ω(err).Should(BeNil())

			var pair Pair
			pair, err = NewPair(key, id)

			var attr CollectionElement
			attr, err = NewMap(pair)
			Ω(err).Should(BeNil())

			str, err = attr.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())

			Ω(str).Should(HavePrefix("{"))
			Ω(str).Should(HaveSuffix("}"))

			Ω(str).Should(ContainSubstring(":db/id #db/id [:db.part/db]"))
		})
	})
})
