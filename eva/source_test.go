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

var _ = Describe("Source test", func() {

	category := "example"
	sourceType := "test"

	sourceFunc := func(config Configuration, tenant Tenant) (Source, error) {
		return &mockSource{}, nil
	}

	clean := func() {
		if _, has := sourceFactories[sourceType]; has {
			delete(sourceFactories, sourceType)
		}
	}

	It("sources", func() {
		clean()

		Ω(HasSource(sourceType)).Should(BeFalse())
		Ω(Sources()).ShouldNot(ContainElement(sourceType))

		e := AddSourceFactory(sourceType, sourceFunc)
		Ω(e).Should(BeNil())

		Ω(HasSource(sourceType)).Should(BeTrue())
		Ω(Sources()).Should(ContainElement(sourceType))

		e = AddSourceFactory(sourceType, sourceFunc)
		Ω(e).ShouldNot(BeNil())
		Ω(e).Should(test.HaveMessage(ErrDuplicateSourceType))

		clean()
	})

	It("sources", func() {
		_, err := NewSource(nil, nil)
		Ω(err).ShouldNot(BeNil())
		Ω(err).Should(test.HaveMessage(ErrInvalidConfiguration))
	})

	It("sources", func() {
		tenant, err := NewTenant("foo")
		Ω(err).Should(BeNil())

		_, err = NewSource(nil, tenant)
		Ω(err).ShouldNot(BeNil())
		Ω(err).Should(test.HaveMessage(ErrInvalidConfiguration))
	})

	It("sources", func() {
		clean()

		e := AddSourceFactory(sourceType, sourceFunc)
		Ω(e).Should(BeNil())

		config, err := NewConfiguration(makeConfig(sourceType, category))
		Ω(err).Should(BeNil())
		Ω(config).ShouldNot(BeNil())

		var tenant Tenant
		tenant, err = NewTenant("foo")
		Ω(err).Should(BeNil())

		source, err := NewSource(config, tenant)
		Ω(err).Should(BeNil())
		Ω(source).ShouldNot(BeNil())

		clean()
	})

	It("sources", func() {
		clean()

		config, err := NewConfiguration(makeConfig("ninja", category))
		Ω(err).Should(BeNil())
		Ω(config).ShouldNot(BeNil())

		var tenant Tenant
		tenant, err = NewTenant("foo")
		Ω(err).Should(BeNil())

		source, err := NewSource(config, tenant)
		Ω(err).ShouldNot(BeNil())
		Ω(err).Should(test.HaveMessage(ErrUnknownSourceType))
		Ω(source).Should(BeNil())

		clean()
	})

	It("sources", func() {
		clean()

		config, err := NewConfiguration("{\"source\":{}}")
		Ω(err).Should(BeNil())
		Ω(config).ShouldNot(BeNil())

		var tenant Tenant
		tenant, err = NewTenant("foo")
		Ω(err).Should(BeNil())

		source, err := NewSource(config, tenant)
		Ω(err).ShouldNot(BeNil())
		Ω(err).Should(test.HaveMessage(ErrUnknownSourceType))
		Ω(source).Should(BeNil())

		clean()
	})

	It("sources", func() {
		clean()

		config, err := NewConfiguration("{}")
		Ω(err).Should(BeNil())
		Ω(config).ShouldNot(BeNil())

		var tenant Tenant
		tenant, err = NewTenant("foo")
		Ω(err).Should(BeNil())

		source, err := NewSource(config, tenant)
		Ω(err).ShouldNot(BeNil())
		Ω(err).Should(test.HaveMessage(ErrUnknownSourceType))
		Ω(source).Should(BeNil())

		clean()
	})

	It("sources", func() {
		clean()

		tenant, err := NewTenant("foo")
		Ω(err).Should(BeNil())

		source, err := NewSource(nil, tenant)
		Ω(err).ShouldNot(BeNil())
		Ω(err).Should(test.HaveMessage(ErrInvalidConfiguration))
		Ω(source).Should(BeNil())

		clean()
	})

	queryImpl := func(interface{}, ...interface{}) (Result, error) {
		return nil, nil
	}

	Context("with the default marshaller", func() {
		It("", func() {
			config, err := NewConfiguration("{\"category\": \"foo\"}")
			Ω(err).Should(BeNil())
			Ω(config).ShouldNot(BeNil())

			tenant, err := NewTenant("foo")
			Ω(err).Should(BeNil())

			source, err := NewBaseSource(config, tenant, &mockSource{}, func(label edn.Element, source Source) (c ConnectionChannel, e error) {
				return c, e
			}, queryImpl)

			Ω(err).Should(BeNil())
			Ω(source).ShouldNot(BeNil())

			_, err = source.Query(nil, nil)
			Ω(err).Should(BeNil())
		})

		It("", func() {
			source, err := NewBaseSource(nil, nil, &mockSource{}, func(label edn.Element, source Source) (c ConnectionChannel, e error) {
				return c, e
			}, queryImpl)

			Ω(err).ShouldNot(BeNil())
			Ω(source).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidConfiguration))
		})
	})
})
