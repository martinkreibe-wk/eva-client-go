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
	"strings"

	"github.com/mattrobenolt/gocql/uuid"
)

const (

	// UUIDElementTag defines the uuid tag value.
	UUIDElementTag = "uuid"

	// 6ba7b810-9dad-11d1-80b4-00c04fd430c8
	UUIDPattern = "\"" +
		"[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]" + // 8
		"[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-" +
		"[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-" + // 4
		"[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-" + // 4
		"[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-" + // 4
		"[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]" + // 12
		"[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]" +
		"[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]\""
)

// uuidStringProcessor used the string processor but will accurately create the uuid.
func uuidStringProcessor(tokenValue string) (Element, error) {
	return uuidFactory(tokenValue)
}

func uuidFactory(input interface{}) (Element, error) {
	switch v := input.(type) {
	case string:
		if !strings.HasPrefix(v, "\"") || !strings.HasSuffix(v, "\"") {
			return nil, MakeError(ErrParserError, "Expected uuid to start and end with quotes.")
		}

		id, err := uuid.ParseUUID(v[1 : len(v)-1])
		if err != nil {
			return nil, err
		}

		return NewUUIDElement(id)
	case uuid.UUID:
		return NewUUIDElement(v)
	default:
		return nil, MakeError(ErrInvalidInput, input)
	}
}

// NewInstantElement creates a new instant element or an error.
func NewUUIDElement(value uuid.UUID) (Element, error) {
	return baseFactory().make(value, UUIDType, UUIDElementTag)
}
