/*
 * Copyright 2019, authproxy authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package errors

import (
	"net/http"
	"fmt"
)

// APIStatus is exposed by errors that can be converted to an api.Status object
// for finer grained details.
type APIStatus interface {
	Status() int
}

type StatusError struct {
	HTTPStatus int
	Message    string
}

func (e *StatusError) Error() string {
	return e.Message
}

// Status allows access to e's http status without having to know the detailed workings
// of StatusError.
func (e *StatusError) Status() int {
	return e.HTTPStatus
}

// ReasonForError returns the HTTP status for a particular error.
func ReasonForError(err error) int {
	switch t := err.(type) {
	case APIStatus:
		return t.Status()
	}
	return 0
}

// NewUnauthorized returns an error indicating the client is not authorized to perform the requested
// action.
func NewUnauthorized(reason string) *StatusError {
	msg := reason
	if len(msg) == 0 {
		msg = "not authorized"
	}
	return &StatusError{
		HTTPStatus: http.StatusUnauthorized,
		Message:    msg,
	}
}

// NewInternalError returns an error indicating the item is invalid and cannot be processed.
func NewInternalError(err error) *StatusError {
	return &StatusError{
		HTTPStatus: http.StatusInternalServerError,
		Message:    fmt.Sprintf("Internal error occurred: %v", err),
	}
}

// IsUnauthorized determines if err is an error which indicates that the request is unauthorized and
// requires authentication by the user.
func IsUnauthorized(err error) bool {
	return ReasonForError(err) == http.StatusUnauthorized
}

// IsInternalError determines if err is an error which indicates an internal server error.
func IsInternalError(err error) bool {
	return ReasonForError(err) == http.StatusInternalServerError
}
