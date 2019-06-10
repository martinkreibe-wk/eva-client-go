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

var _ = Describe("Tenant Test", func() {
	It("Create the reference", func() {
		org := "org"

		tenant, err := NewTenant(org)
		Ω(err).Should(BeNil())

		Ω(tenant.Name()).Should(BeEquivalentTo(org))

		resCorr, has := tenant.CorrelationId()
		Ω(has).Should(BeFalse())
		Ω(resCorr).Should(BeEmpty())
	})

	It("Create the reference", func() {
		org := "org"
		corr := "corr-id"

		tenant, err := NewCorrelationTenant(org, corr)
		Ω(err).Should(BeNil())

		Ω(tenant.Name()).Should(BeEquivalentTo(org))

		resCorr, has := tenant.CorrelationId()
		Ω(has).Should(BeTrue())
		Ω(resCorr).Should(BeEquivalentTo(corr))
	})
})
