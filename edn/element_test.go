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
	"strconv"

	"errors"
	"github.com/Workiva/eva-client-go/test"
	"github.com/mattrobenolt/gocql/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

const tag = NoTag

var _ = Describe("Elements in EDN", func() {
	Context("with the default marshaller", func() {
		It("should create an base element with no error", func() {

			t := ElementType(99)

			elem, err := baseFactory().make(nil, t, tag, func(serializer Serializer, tag string, i interface{}) (string, error) {
				return "", nil
			})
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeIdenticalTo(t))
		})

		It("should create an base element with no error", func() {

			t := ElementType(99)

			elem, err := baseFactory().make(nil, t, tag, nil)
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidElement))
		})

		It("should equal the same thing if they are actually equal", func() {

			value := "42"

			elem, err := baseFactory().make(value, StringType, tag, func(serializer Serializer, tag string, value interface{}) (out string, e error) {
				out = strconv.Quote(value.(string))
				return out, e
			})
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.Value()).Should(BeEquivalentTo(value))

			elem2, err := baseFactory().make(value, StringType, tag, func(serializer Serializer, tag string, value interface{}) (out string, e error) {
				out = strconv.Quote(value.(string))
				return out, e
			})
			Ω(err).Should(BeNil())
			Ω(elem2).ShouldNot(BeNil())
			Ω(elem2.Value()).Should(BeEquivalentTo(value))

			Ω(elem.Equals(elem2)).Should(BeTrue())
		})

		It("should set tags correctly", func() {

			value := "42"

			elem, err := baseFactory().make(value, StringType, tag, func(serializer Serializer, tag string, value interface{}) (out string, e error) {
				out = strconv.Quote(value.(string))
				return out, e
			})

			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			Ω(elem.Tag()).Should(BeEmpty())

			tag := "someValue"
			err = elem.SetTag(tag)
			Ω(err).Should(BeNil())
			Ω(elem.Tag()).Should(BeEquivalentTo(tag))

			tag = "anotherValue"
			err = elem.SetTag("#" + tag)
			Ω(err).Should(BeNil())
			Ω(elem.Tag()).Should(BeEquivalentTo(tag))

			orig := tag
			tag = "bad value"
			err = elem.SetTag(tag)
			Ω(err).Should(test.HaveMessage(ErrInvalidSymbol))
			Ω(elem.Tag()).Should(BeEquivalentTo(orig))
		})

		It("should create an base element with no error", func() {

			t := ElementType(99)

			elem, err := baseFactory().make(nil, t, tag, func(serializer Serializer, tag string, i interface{}) (string, error) {
				return "", errors.New("expected")

			})
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())

			Ω(func() { elem.String() }).Should(Panic())
		})

		It("should create an base element with no error", func() {
			elem, err := NewPrimitiveElement("nil")
			Ω(IsPrimitive("nil")).Should(BeTrue())
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(NilType))
		})

		It("should create an base element with no error", func() {
			elem, err := NewPrimitiveElement(nil)
			Ω(IsPrimitive(nil)).Should(BeFalse())
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(NilType))
		})

		It("should create an base element with no error", func() {
			elem, err := NewPrimitiveElement(true)
			Ω(IsPrimitive(true)).Should(BeTrue())
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(BooleanType))
			Ω(elem.Value()).Should(BeEquivalentTo(true))
		})

		It("should create an base element with no error", func() {
			elem, err := NewPrimitiveElement(int(1))
			Ω(IsPrimitive(int(1))).Should(BeTrue())
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(IntegerType))
			Ω(elem.Value()).Should(BeEquivalentTo(int(1)))
		})

		It("should create an base element with no error", func() {
			elem, err := NewPrimitiveElement(int32(1))
			Ω(IsPrimitive(int32(1))).Should(BeTrue())
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(IntegerType))
			Ω(elem.Value()).Should(BeEquivalentTo(int32(1)))
		})

		It("should create an base element with no error", func() {
			elem, err := NewPrimitiveElement(int64(1))
			Ω(IsPrimitive(int64(1))).Should(BeTrue())
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(IntegerType))
			Ω(elem.Value()).Should(BeEquivalentTo(int64(1)))
		})

		It("should create an base element with no error", func() {
			elem, err := NewPrimitiveElement(float32(1.2))
			Ω(IsPrimitive(float32(1.2))).Should(BeTrue())
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(FloatType))
			Ω(elem.Value()).Should(BeEquivalentTo(float32(1.2)))
		})

		It("should create an base element with no error", func() {
			elem, err := NewPrimitiveElement(float64(1.2))
			Ω(IsPrimitive(float64(1.2))).Should(BeTrue())
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(FloatType))
			Ω(elem.Value()).Should(BeEquivalentTo(float64(1.2)))
		})

		It("should create an base element with no error", func() {
			now := time.Now()
			Ω(IsPrimitive(now)).Should(BeTrue())
			elem, err := NewPrimitiveElement(now)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(InstantType))
			Ω(elem.Value()).Should(BeEquivalentTo(now))
		})

		It("should create an base element with no error", func() {
			id := uuid.RandomUUID()
			Ω(IsPrimitive(id)).Should(BeTrue())
			elem, err := NewPrimitiveElement(id)
			Ω(err).Should(BeNil())
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(UUIDType))
			Ω(elem.Value()).Should(BeEquivalentTo(id))
		})

		It("should create an base element with no error", func() {
			id := uuid.RandomUUID()
			Ω(IsPrimitive(id)).Should(BeTrue())

			lexer := &lexerImpl{}
			err := lexer.init()
			Ω(err).Should(BeNil())

			err = lexer.RemoveFactory(UUIDType, UUIDElementTag)
			Ω(err).Should(BeNil())

			elem, err := NewPrimitiveElementWithLexer(lexer, id)
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidElement))
		})

		It("should create an base element with no error", func() {
			id := struct{}{}
			Ω(IsPrimitive(id)).Should(BeFalse())
			elem, err := NewPrimitiveElement(id)
			Ω(err).ShouldNot(BeNil())
			Ω(elem).Should(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})
	})
})
