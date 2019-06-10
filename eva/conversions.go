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

func decodeSerializable(item interface{}) (ser edn.Serializable, err error) {

	if item != nil {
		bad := true
		switch val := item.(type) {
		case edn.Element:
			ser = val
			bad = false
		case string:
			ser = RawString(val)
			bad = false
		case int:
			ser = RawInt(int64(val))
			bad = false
		case int8:
			ser = RawInt(int64(val))
			bad = false
		case int16:
			ser = RawInt(int64(val))
			bad = false
		case int32:
			ser = RawInt(int64(val))
			bad = false
		case int64:
			ser = RawInt(val)
			bad = false
		case rawStringImpl:
			ser = val
			bad = false
		case rawIntImpl:
			ser = val
			bad = false
		}

		if bad {
			err = edn.MakeErrorWithFormat(edn.ErrInvalidInput, "Unsupported type: %T", item)
		}
	}

	return ser, err
}
