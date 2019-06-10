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

package test

import (
	"fmt"
	"github.com/onsi/gomega/types"
)

type Messager interface {
	Message() string
}

func HaveMessage(message Messager) types.GomegaMatcher {
	return &errMessageMatcher{
		message: message,
	}
}

type errMessageMatcher struct {
	message Messager
}

func (matcher *errMessageMatcher) Match(actual interface{}) (success bool, err error) {
	if response, ok := actual.(Messager); ok {
		success = response.Message() == matcher.message.Message()
	} else {
		err = fmt.Errorf("Expected actual to have a Message() string method, instead found:  %T", actual)
	}

	return success, err
}

func (matcher *errMessageMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nhave message\n\t%#v", actual, matcher.message)
}

func (matcher *errMessageMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to have message\n\t%#v", actual, matcher.message)
}
