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

var _ = Describe("Base Channel", func() {

	ct := ChannelType("ch-type")
	label := edn.NewStringElement("label")

	Context("with the default marshaller", func() {
		It("", func() {
			src := &mockSource{}

			bc, err := NewBaseChannel(ct, src, map[string]edn.Serializable{})
			Ω(err).Should(BeNil())
			Ω(bc.Type()).Should(BeEquivalentTo(ct))
			Ω(bc.Reference()).ShouldNot(BeNil())
		})
	})

	tenant := "org"

	asOfSnapshot := func(edn.Serializable) (channel SnapshotChannel, err error) {

		channel = &BaseSnapshotChannel{}

		return channel, err
	}

	goodTransact := func(data edn.Serializable) (result Result, err error) {
		return nil, nil
	}

	Context("with the default marshaller", func() {
		It("", func() {
			src := &mockSource{}

			_, err := NewBaseConnectionChannel(label, src, goodTransact, nil)

			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(edn.ErrInvalidInput))

		})
	})

	Context("with the default marshaller", func() {
		It("", func() {
			src := &mockSource{}

			bcc, err := NewBaseConnectionChannel(label, src, goodTransact, asOfSnapshot)

			Ω(err).Should(BeNil())
			Ω(bcc.Type()).Should(BeEquivalentTo(ConnectionReferenceType))
			Ω(bcc.Reference()).ShouldNot(BeNil())

			tenant, err := NewTenant(tenant)
			Ω(err).Should(BeNil())

			_, err = bcc.Transact(tenant, nil)
			Ω(err).ShouldNot(BeNil())
		})
	})
})
