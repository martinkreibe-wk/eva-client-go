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

type ClientErrorKeyword string

//{
//	:message "",
//	:ex-info {
//		:explanation "Malformed transact request.",
//		:type "IncorrectTransactSyntax",
//		:code 3000},
//	:ex-data ""
//}

const (
	ErrSourceError = edn.ErrorMessage("source error")
)

// Get the error from just the code.
// -- What about the other data?

type ClientError interface {
	error
	Name() string
	Keyword() edn.SymbolElement
	Description() string
	Code() int
	//	Details() [] // TODO
}

// Error is the error type.
type clientErrorImpl struct {
	err  *edn.Error
	code int64
	//	details  // TODO
}

// Error returns the error message.
func (e *clientErrorImpl) Error() string {
	return e.err.Error()
}

func DecodeError(code edn.Element) (err error) {

	if code.ElementType() == edn.IntegerType {
		err = &clientErrorImpl{
			err:  edn.MakeError(ErrSourceError, nil),
			code: code.Value().(int64),
		}
	}

	return err
}
