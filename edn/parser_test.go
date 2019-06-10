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
	"errors"
	"fmt"
	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strconv"
	"strings"
)

type keywordValue struct {
	namespace string
	name      string
}

type testDefinition struct {
	expression string
	expected   interface{}
}

type testInstance struct {
	testDefinition
	elemType ElementType
	tag      string
}

const (
	// IMPORTANT!!!!
	// These flags are to help limit the huge number of tests...
	// but they should be:

	onlyTarget = ElementType("") // SetType
)

var tags = []string{
	"tagVal",
	"tag/Val",
	"another/tag-val",
}

var simpleFormats = []string{
	"%s %s",
}

var complexFormats = []string{
	"%s\n\r%s",
	"  %s %s",
	"  %s\n\r%s",
	"\n\r%s %s",
	"\n\r%s\n\r%s",
	"%s\n\r%s  ",
	"%s\n%s\n\r",
	"; comment\n\r%s %s",
	"; comment\n\r%s\n\r%s",
	"%s; comment\n%s",
	"%s; comment\n%s; comment\n",
	";comment\n%s; comment\n%s; comment\n",

	"%s %s; comment\n",
	"%s\n\r%s; comment\n",
}

func message(ser Serializer, label string, index int, test *testInstance, elem Element) string {
	var expected interface{}
	switch exp := test.expected.(type) {
	case *keywordValue:
		expected = exp.namespace + " / " + exp.name
	case func() (map[string]Element, error):
		v, _ := exp()
		if m, e := NewMap(); e == nil {
			for key, value := range v {
				m.Append(NewStringElement(key), value)
			}
			if expected, e = m.Serialize(ser); e != nil {
				panic(e)
			}
		} else {
			panic(e)
		}
	case func() (string, interface{}, error):
		_, expected, _ = exp()
	default:
		expected = test.expected
	}

	comp := "<nil>"
	if elem != nil {
		var ee error
		comp, ee = elem.Serialize(ser)
		if ee != nil {
			panic(ee)
		}
	}

	return fmt.Sprintf("[%d:%s] Test: '%s' -> [%v] (tag: %s) \n\tExpected: %#v\n\tActual: %s", index, label, test.expression, test.elemType, test.tag, expected, comp)
}

func runParserTests(elemType ElementType, definitions ...*testDefinition) {

	if len(onlyTarget) != 0 && elemType != onlyTarget {
		return
	}

	ser := EvaEdnMimeType

	formats := simpleFormats

	if full := os.Getenv("FULL_TESTS"); len(full) > 0 {
		if fullTests, e := strconv.ParseBool(full); e == nil && fullTests {
			formats = append(formats, complexFormats...)
		}
	}

	for _, inst := range definitions {
		var tests []*testInstance

		build := func(tagVal string) {
			for _, format := range formats {

				strTag := ""
				if len(tagVal) > 0 {
					strTag = "#" + tagVal
				}

				buildAndAppendInstance := func(ex string) {
					def := &testInstance{
						testDefinition: testDefinition{
							expression: ex,
							expected:   inst.expected,
						},
						elemType: elemType,
						tag:      tagVal,
					}

					tests = append(tests, def)
				}

				ex := fmt.Sprintf(format, strTag, inst.expression)
				buildAndAppendInstance(ex)

				if ex2 := strings.TrimSpace(ex); len(formats) > len(simpleFormats) && ex2 != ex {
					buildAndAppendInstance(ex2)
				}
			}
		}

		build("")

		if !strings.HasPrefix(inst.expression, "#") {
			for _, tag := range tags {
				build(tag)
			}
		}

		for index, testCase := range tests {
			It(fmt.Sprintf("should parse the expressions: `%s`", testCase.expression), func() {
				elem, err := Parse(testCase.expression)
				Ω(err).Should(BeNil(), message(ser, "err", index, testCase, elem))
				Ω(elem).ShouldNot(BeNil(), message(ser, "elem", index, testCase, elem))
				Ω(elem.ElementType()).Should(BeEquivalentTo(testCase.elemType), message(ser, "type", index, testCase, elem))

				if testCase.expected == nil {
					Ω(elem.Value()).Should(BeNil(), message(ser, "value", index, testCase, elem))
				} else {
					switch exp := testCase.expected.(type) {
					case *keywordValue:
						Ω(elem.Tag()).Should(BeEquivalentTo(testCase.tag), message(ser, "tag", index, testCase, elem))
						Ω(elem.HasTag()).Should(BeEquivalentTo(len(testCase.tag) > 0), message(ser, "hasTag", index, testCase, elem))

						e := elem.(SymbolElement)
						Ω(e.Name()).Should(BeEquivalentTo(exp.name), message(ser, "name", index, testCase, elem))
						Ω(e.Prefix()).Should(BeEquivalentTo(exp.namespace), message(ser, "namespace", index, testCase, elem))
					case func() (map[string]Element, error):
						Ω(elem.ElementType().IsCollection()).Should(BeTrue(), message(ser, "is collection", index, testCase, elem))
						coll, e := exp()

						Ω(elem.Tag()).Should(BeEquivalentTo(testCase.tag), message(ser, "tag", index, testCase, elem))
						Ω(elem.HasTag()).Should(BeEquivalentTo(len(testCase.tag) > 0), message(ser, "hasTag", index, testCase, elem))

						Ω(e).Should(BeNil(), message(ser, "func err", index, testCase, elem))

						if collElem, is := elem.(CollectionElement); is {
							Ω(collElem.Len()).Should(BeEquivalentTo(len(coll)), message(ser, "collection size", index, testCase, elem))

							err = collElem.IterateChildren(func(key Element, value Element) (iterErr error) {

								var comp string
								comp, iterErr = key.Serialize(ser)
								Ω(iterErr).Should(BeNil(), message(ser, "key comp err", index, testCase, elem))

								if toCompare, has := coll[comp]; has {
									Ω(value.Equals(toCompare)).Should(BeTrue(), message(ser, "child check", index, testCase, elem))
								} else {
									iterErr = errors.New("missing key: " + comp)
								}

								return iterErr
							})

							if elem.ElementType() == SetType && len(coll) != 0 {
								Fail(fmt.Sprintf("Expected set collection to be empty: %v", coll))
							}

							Ω(err).Should(BeNil(), message(ser, "iteration err", index, testCase, elem))
						} else {
							Fail(message(ser, "collection casting", index, testCase, elem))
						}

					case func() (map[string][2]Element, error):
						Ω(elem.ElementType().IsCollection()).Should(BeTrue(), message(ser, "is collection", index, testCase, elem))
						coll, e := exp()

						Ω(elem.Tag()).Should(BeEquivalentTo(testCase.tag), message(ser, "tag", index, testCase, elem))
						Ω(elem.HasTag()).Should(BeEquivalentTo(len(testCase.tag) > 0), message(ser, "hasTag", index, testCase, elem))

						Ω(e).Should(BeNil(), message(ser, "func err", index, testCase, elem))

						if collElem, is := elem.(CollectionElement); is {
							Ω(collElem.Len()).Should(BeEquivalentTo(len(coll)), message(ser, "collection size", index, testCase, elem))

							err = collElem.IterateChildren(func(key Element, value Element) (iterErr error) {

								if elem.ElementType() == SetType {

									for k, val := range coll {
										if value.Equals(val[1]) {
											delete(coll, k)
										}
									}

								} else {
									var comp string
									comp, iterErr = key.Serialize(ser)
									Ω(iterErr).Should(BeNil(), message(ser, "key comp err", index, testCase, elem))

									has := false
									keys := make([]string, 0)
									for key, val := range coll {
										keys = append(keys, key)

										if key == comp {
											has = true
											Ω(value.Equals(val[1])).Should(BeTrue(), message(ser, "child check", index, testCase, elem))
										}

										if has {
											break
										}
									}

									if !has {
										iterErr = errors.New(fmt.Sprintf("missing key: %s in %v\n: %#v", comp, keys, coll))
									}

								}

								return iterErr
							})

							if elem.ElementType() == SetType && len(coll) != 0 {
								Fail(fmt.Sprintf("Expected set collection to be empty: %v", coll))
							}

							Ω(err).Should(BeNil(), message(ser, "iteration err", index, testCase, elem))
						} else {
							Fail(message(ser, "collection casting", index, testCase, elem))
						}

					case func() (string, interface{}, error):
						t, v, e := exp()
						Ω(t).Should(BeEquivalentTo(t), message(ser, "tag", index, testCase, elem))
						Ω(elem.HasTag()).Should(BeEquivalentTo(len(t) > 0), message(ser, "hasTag", index, testCase, elem))
						Ω(e).Should(BeNil(), message(ser, "func err", index, testCase, elem))
						Ω(elem.Value()).Should(BeEquivalentTo(v), message(ser, "value", index, testCase, elem))
					default:
						Ω(elem.Tag()).Should(BeEquivalentTo(testCase.tag), message(ser, "tag", index, testCase, elem))
						Ω(elem.HasTag()).Should(BeEquivalentTo(len(testCase.tag) > 0), message(ser, "hasTag", index, testCase, elem))
						Ω(elem.Value()).Should(BeEquivalentTo(testCase.expected), message(ser, "value", index, testCase, elem))
					}
				}
			})
		}
	}
}

var _ = Describe("Collection Parser", func() {
	It("", func() {

		var coll CollectionElement
		var err error

		coll, err = ParseCollection("[]")
		Ω(err).Should(BeNil())
		Ω(coll).ShouldNot(BeNil())

		coll, err = ParseCollection("42")
		Ω(err).ShouldNot(BeNil())
		Ω(coll).Should(BeNil())
		Ω(err).Should(test.HaveMessage(ErrParserError))
	})
})
