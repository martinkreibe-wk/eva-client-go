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

import "strings"

const (
	// ErrUnknownMimeType defines an unknown serialization type.
	ErrUnknownMimeType = ErrorMessage("unknown serialization mime type")

	// DefaultMimeType defines the default serializer type.
	DefaultMimeType = EvaEdnMimeType
)

var serializerFactories = map[SerializerMimeType]func(string) (Serializer, error){
	EvaEdnMimeType.MimeType(): func(s string) (Serializer, error) {
		return SerializerMimeType(s), nil
	},
}

// Serializer defines the interface for converting the entity into a serialized edn value.
type Serializer interface {

	// MimeType of serializer this is.
	MimeType() SerializerMimeType

	// Options returns the option value or false.
	Options(string) (string, bool)
}

// GetSerializer will return the serializer requested or an error
func GetSerializer(serializerType string) (serializer Serializer, err error) {
	return GetSerializerByType(SerializerMimeType(serializerType))
}

type optionsProcessor func(string, string) bool

func scrapeOptionFromMime(serializerType SerializerMimeType, processor optionsProcessor) (SerializerMimeType, string) {
	mimeType := serializerType
	strType := string(serializerType)
	if index := strings.Index(strType, ";"); index != -1 {
		mimeType = SerializerMimeType(strType[:index])

		if processor != nil {
			optionsStr := strType[index+1:]

			options := strings.Split(optionsStr, ",")
			for _, option := range options {
				index = strings.Index(option, "=")

				if !processor(option[:index], option[index+1:]) {
					break
				}
			}
		}
	}

	return mimeType, strType
}

// GetSerializerByType will return the serializer requested or an error
func GetSerializerByType(serializerType SerializerMimeType) (serializer Serializer, err error) {

	mimeType, strType := scrapeOptionFromMime(serializerType, nil)
	if factory, has := serializerFactories[mimeType]; has {
		serializer, err = factory(strType)
	} else {
		err = MakeError(ErrUnknownMimeType, serializerType)
	}

	return serializer, err
}
