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
)

// ErrorMessage defines a message used in an error
type ErrorMessage string

// Error is the error type.
type Error struct {
	message ErrorMessage
	details string
}

// Message will get the message part.
func (em ErrorMessage) Message() string {
	return string(em)
}

// Message will get the message part.
func (e *Error) Message() string {
	return e.message.Message()
}

// Error returns the error message.
func (e *Error) Error() string {
	return fmt.Sprintf("[%s]: %s", e.message, e.details)
}

// FormatError is an error that creates a unique message from the state at the time of creation.
type FormatError struct {
	message string
	items   []interface{}
}

// Error returns the error message.
func (e FormatError) Error() string {
	return fmt.Sprintf(e.message, e.items...)
}

// CumulativeError defines a collection of errors.
type CumulativeError struct {
	items []error
}

// Append the error to this error.
func (cumErr *CumulativeError) Append(err ...error) {
	for _, e := range err {
		if e != nil {
			switch v := e.(type) {
			case *CumulativeError:
				cumErr.Append(v.items...)
			default:
				cumErr.items = append(cumErr.items, e)
			}
		}
	}
}

// Error returns the error message.
func (cumErr *CumulativeError) Error() string {
	message := ""
	for index, part := range cumErr.items {
		message += fmt.Sprintf("%d: %s\n", index, part)
	}

	return message
}

// ErrorList returns the error collection.
func (cumErr *CumulativeError) ErrorList() []error {
	return cumErr.items
}

// MakeErrorWithFormat will create the error with a formatted string.
func MakeErrorWithFormat(message ErrorMessage, format string, details ...interface{}) (err *Error) {
	return MakeError(message, fmt.Sprintf(format, details...))
}

// MakeError will create the error
func MakeError(message ErrorMessage, details interface{}) (err *Error) {

	err = &Error{
		message: message,
	}

	if details != nil {
		err.details = fmt.Sprintf("%+v", details)
	}

	return err
}

// IsEquivalent checks if the errors are equivalent.
func (em ErrorMessage) IsEquivalent(err error) (eq bool) {
	if err != nil {
		if myErr, is := err.(*Error); is {
			eq = em == myErr.message
		}
	}

	return eq
}

// NewError creates a new error.
func NewError(message string, contents ...interface{}) (err error) {

	if len(contents) == 0 {
		err = &Error{
			message: ErrorMessage(message),
		}
	} else {
		err = &FormatError{
			message: message,
			items:   contents,
		}
	}

	return err
}

// AppendError will take the original error and append the new one.
func AppendError(errors ...error) (err error) {
	switch len(errors) {

	// handle the simple cases where one or both are nil.
	case 0:
		break
	case 1:
		err = errors[0]

	default:
		for i := 0; i < len(errors); i++ {
			if errors[i] != nil {
				switch v := errors[i].(type) {
				case *CumulativeError:
					v.Append(errors[i+1:]...)
					err = v
				default:
					if len(errors)-i > 1 {
						cumErr := &CumulativeError{}
						cumErr.Append(errors[i:]...)
						err = cumErr
					} else {
						err = errors[i]
					}
				}
				break
			}
		}
	}

	return err
}
