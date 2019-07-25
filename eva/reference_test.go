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

type mockSerializable struct {
}

func (mock *mockSerializable) String() string {
	return "foo"
}

// Serialize will convert this structure to an edn string.
func (mock *mockSerializable) Serialize(serialize edn.Serializer) (string, error) {
	return "", nil
}

var _ = Describe("reference test", func() {

	Context("normally", func() {

		It("create conn ref with a label", func() {
			ref, err := NewConnectionReference("label")
			Ω(err).Should(BeNil())
			Ω(ref).ShouldNot(BeNil())
			Ω(ref.Type()).Should(BeEquivalentTo(ConnectionReferenceType))

			stream := edn.NewStringStream()
			err = edn.EvaEdnMimeType.SerializeTo(stream, ref)
			Ω(err).Should(BeNil())
			v := stream.String()
			Ω(v).Should(BeEquivalentTo("#eva.client.service/connection-ref {:label \"label\"}"))
		})

		It("create snap ref with a label", func() {
			ref, err := NewSnapshotReference("label")
			Ω(err).Should(BeNil())
			Ω(ref).ShouldNot(BeNil())
			Ω(ref.Type()).Should(BeEquivalentTo(SnapshotReferenceType))

			stream := edn.NewStringStream()
			err = edn.EvaEdnMimeType.SerializeTo(stream, ref)
			Ω(err).Should(BeNil())
			v := stream.String()
			Ω(v).Should(BeEquivalentTo("#eva.client.service/snapshot-ref {:label \"label\"}"))
		})

		It("create snap ref with a label and asOf", func() {
			elem, err := edn.NewIntegerElement(123)
			Ω(err).Should(BeNil())

			ref, err := NewSnapshotAsOfReference("label", elem)
			Ω(err).Should(BeNil())
			Ω(ref).ShouldNot(BeNil())
			Ω(ref.Type()).Should(BeEquivalentTo(SnapshotReferenceType))

			stream := edn.NewStringStream()
			err = edn.EvaEdnMimeType.SerializeTo(stream, ref)
			Ω(err).Should(BeNil())
			v := stream.String()
			Ω(v).Should(HavePrefix("#eva.client.service/snapshot-ref {"))
			Ω(v).Should(ContainSubstring(":as-of 123"))
			Ω(v).Should(ContainSubstring(":label \"label\""))
			Ω(v).Should(HaveSuffix("}"))
		})
	})

	Context("normally", func() {

		It("should not panic if there is no error", func() {
			chanType := ChannelType("taco")

			ref, err := newReference(chanType, nil)
			Ω(err).Should(BeNil())
			Ω(ref).ShouldNot(BeNil())
			Ω(ref.Type()).Should(BeEquivalentTo(chanType))

			stream := edn.NewStringStream()
			err = edn.EvaEdnMimeType.SerializeTo(stream, ref)
			Ω(err).Should(BeNil())
			v := stream.String()
			Ω(v).Should(BeEquivalentTo(fmt.Sprintf("#%s {}", chanType)))
		})

		It("handle properties correctly", func() {
			chanType := ChannelType("taco")

			prop1, err := edn.NewStringElement("value1")
			Ω(err).Should(BeNil())
			prop2, err := edn.NewStringElement("value2")
			Ω(err).Should(BeNil())
			prop3, err := edn.NewStringElement("str")
			Ω(err).Should(BeNil())
			prop4, err := edn.NewIntegerElement(123)
			Ω(err).Should(BeNil())

			ref, err := newReference(chanType, map[string]edn.Element{
				"prop1": prop1,
				"prop2": prop2,
				"prop3": prop3,
				"prop4": prop4,
			})
			Ω(err).Should(BeNil())
			Ω(ref).ShouldNot(BeNil())
			Ω(ref.Type()).Should(BeEquivalentTo(chanType))

			stream := edn.NewStringStream()
			err = edn.EvaEdnMimeType.SerializeTo(stream, ref)
			Ω(err).Should(BeNil())
			v := stream.String()

			Ω(v).Should(ContainSubstring(fmt.Sprintf(":prop1 \"value1\"")))
			Ω(v).Should(ContainSubstring(fmt.Sprintf(":prop2 \"value2\"")))
			Ω(v).Should(ContainSubstring(fmt.Sprintf(":prop3 \"str\"")))
			Ω(v).Should(ContainSubstring(fmt.Sprintf(":prop4 123")))
			Ω(v).Should(HavePrefix(fmt.Sprintf("#%s {", chanType)))
			Ω(v).Should(HaveSuffix("}"))

			prop5, err := edn.NewStringElement("new-str")
			Ω(err).Should(BeNil())

			prop6, err := edn.NewStringElement("value2.2")
			Ω(err).Should(BeNil())

			err = ref.AddProperty("prop5", prop5)
			Ω(err).Should(BeNil())
			err = ref.AddProperty("prop3", nil)
			Ω(err).Should(BeNil())
			err = ref.AddProperty("prop2", prop6)
			Ω(err).Should(BeNil())

			stream = edn.NewStringStream()
			err = edn.EvaEdnMimeType.SerializeTo(stream, ref)
			Ω(err).Should(BeNil())
			v = stream.String()
			Ω(v).Should(HavePrefix("#" + string(chanType)))
			Ω(v).Should(ContainSubstring(":prop1 \"value1\""))
			Ω(v).Should(ContainSubstring(":prop2 \"value2.2\""))
			Ω(v).Should(ContainSubstring(":prop5 \"new-str\""))
			Ω(v).Should(ContainSubstring(":prop4 123"))
			Ω(v).Should(HaveSuffix("}"))
		})
	})
})
