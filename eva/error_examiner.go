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
	"github.com/Workiva/eva-client-go/edn"
)

func init() {
	PanicOnError(func() (err error) {
		if exInfoKeyword, err = edn.NewKeywordElement("ex-info"); err == nil {
			codeKeyword, err = edn.NewKeywordElement("code")
		}
		return err
	})
}

var exInfoKeyword edn.SymbolElement
var codeKeyword edn.SymbolElement

// ErrorExaminer will retrieve the error from the payload
type ErrorExaminer func([]byte) error

// GetErrorExaminer will return the error examiner requested, or an error
func GetErrorExaminer(serializer edn.Serializer) (examiner ErrorExaminer, err error) {

	if serializer != nil {
		switch serializer.MimeType() {
		case edn.EvaEdnMimeType:
			examiner = ednErrorExaminer
		default:
			err = edn.MakeError(ErrInvalidSerializer, serializer)
		}
	} else {
		examiner = ednErrorExaminer
	}

	return examiner, err
}

// ednErrorExaminer will examine the payload for an error.
func ednErrorExaminer(body []byte) (err error) {
	var elem edn.Element
	if elem, err = edn.Parse(string(body)); elem != nil {
		if elem.ElementType() == edn.MapType {
			coll := elem.(edn.CollectionElement)
			var exceptionElem edn.Element
			if exceptionElem, err = coll.Get(exInfoKeyword); err == nil {
				if elem.ElementType() == edn.MapType {
					innerMap := exceptionElem.(edn.CollectionElement)

					var code edn.Element
					if code, err = innerMap.Get(codeKeyword); err == nil {
						err = DecodeError(code)
					}
				}
			}
		}
	}
	return err
}
