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

var _ = Describe("Raw String test", func() {

	Context("normally", func() {

		It("should wrap a string", func() {
			serializable := RawString("item")
			Ω(serializable).Should(BeEquivalentTo("item"))

			Ω(serializable.String()).Should(BeEquivalentTo("item"))

			str, err := serializable.Serialize(nil)
			Ω(str).Should(BeEquivalentTo("item"))
			Ω(err).Should(BeNil())
		})
	})
})
