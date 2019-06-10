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

import "github.com/Workiva/eva-client-go/edn"

type rawStringImpl string

func RawString(item string) edn.Serializable {
	return rawStringImpl(item)
}

// String creates a raw string.
func (item rawStringImpl) String() string {
	return string(item)
}

// Serialize will convert this structure to an edn string.
func (item rawStringImpl) Serialize(serialize edn.Serializer) (string, error) {
	return item.String(), nil
}
