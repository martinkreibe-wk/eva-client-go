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

const (

	// EvaEdnMimeType defines the mime type for the eva edn data.
	EvaEdnMimeType SerializerMimeType = "application/vnd.eva+edn"
)

type SerializerMimeType string

// String this Mime Type.
func (smt SerializerMimeType) String() string {
	return string(smt)
}

// MimeType of serializer this is.
func (smt SerializerMimeType) MimeType() SerializerMimeType {
	t, _ := scrapeOptionFromMime(smt, nil)
	return t
}

// Options returns the option value or false.
func (smt SerializerMimeType) Options(name string) (option string, has bool) {

	scrapeOptionFromMime(smt, func(o string, value string) bool {
		if has = name == o; has {
			option = value
		}

		return !has
	})

	return option, has
}
