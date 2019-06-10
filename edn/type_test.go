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
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types in EDN", func() {
	Context("with the default usage", func() {
		It("initialization should not panic the first iteration, but should on the second.", func() {

			for key := range typeFactories {
				delete(typeFactories, key)
			}
			Ω(initAll).ShouldNot(Panic())
			Ω(initAll).Should(Panic())
		})

		typeCollectionMap := map[ElementType]struct {
			value bool
			name  string
		}{
			NilType:       {false, ":db.type/nil"},
			BooleanType:   {false, ":db.type/boolean"},
			StringType:    {false, ":db.type/string"},
			CharacterType: {false, ":db.type/character"},
			SymbolType:    {false, ":db.type/symbol"},
			KeywordType:   {false, ":db.type/keyword"},
			IntegerType:   {false, ":db.type/long"},
			FloatType:     {false, ":db.type/float"},
			InstantType:   {false, ":db.type/instant"},
			UUIDType:      {false, ":db.type/uuid"},
			ListType:      {true, ":db.type/group"},
			VectorType:    {true, ":db.type/vector"},
			MapType:       {true, ":db.type/map"},
			SetType:       {true, ":db.type/set"},
		}

		It("distinguish collections from non collections", func() {

			for t, data := range typeCollectionMap {
				Ω(t.IsCollection()).Should(BeEquivalentTo(data.value), fmt.Sprintf("Expected %s to be: %v", data.name, data.value))
			}
		})

		It("should have the right name", func() {
			for t, data := range typeCollectionMap {
				Ω(t.Name()).Should(BeEquivalentTo(data.name), fmt.Sprintf("Expected %s to be: %v", data.name, data.value))
			}
		})

		It("should have the right name for unknown types", func() {
			Ω(UnknownType.Name()).Should(BeEquivalentTo(""))

			testType := ElementType("foo")
			Ω(testType.Name()).Should(BeEquivalentTo("foo"))
			Ω(testType.IsCollection()).Should(BeEquivalentTo(false))

		})
	})
})
