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
	"io"
)

const (
	// ErrUnknownMimeType defines an unknown serialization type.
	ErrUnknownMimeType = ErrorMessage("unknown serialization mime type")
)

type Stream interface {
	io.Writer
	fmt.Stringer
}

type StringStream struct {
	value string
}

func NewStringStream() Stream {
	return &StringStream{}
}

func (stream *StringStream) Write(p []byte) (int, error) {
	stream.value += string(p)
	return 0, nil
}

func (stream *StringStream) String() string {
	return stream.value
}

// Serializer defines the interface for converting the entity into a serialized edn value.
type Serializer interface {

	// Serialize the tag and value
	SerializeTo(Stream, Element) error
}
