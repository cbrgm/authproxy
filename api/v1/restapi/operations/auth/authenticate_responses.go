// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/cbrgm/authproxy/api/v1/models"
)

// AuthenticateOKCode is the HTTP code returned for type AuthenticateOK
const AuthenticateOKCode int = 200

/*AuthenticateOK OK

swagger:response authenticateOK
*/
type AuthenticateOK struct {

	/*
	  In: Body
	*/
	Payload *models.TokenReviewRequest `json:"body,omitempty"`
}

// NewAuthenticateOK creates AuthenticateOK with default headers values
func NewAuthenticateOK() *AuthenticateOK {

	return &AuthenticateOK{}
}

// WithPayload adds the payload to the authenticate o k response
func (o *AuthenticateOK) WithPayload(payload *models.TokenReviewRequest) *AuthenticateOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the authenticate o k response
func (o *AuthenticateOK) SetPayload(payload *models.TokenReviewRequest) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AuthenticateOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AuthenticateUnauthorizedCode is the HTTP code returned for type AuthenticateUnauthorized
const AuthenticateUnauthorizedCode int = 401

/*AuthenticateUnauthorized unauthorized

swagger:response authenticateUnauthorized
*/
type AuthenticateUnauthorized struct {
}

// NewAuthenticateUnauthorized creates AuthenticateUnauthorized with default headers values
func NewAuthenticateUnauthorized() *AuthenticateUnauthorized {

	return &AuthenticateUnauthorized{}
}

// WriteResponse to the client
func (o *AuthenticateUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// AuthenticateInternalServerErrorCode is the HTTP code returned for type AuthenticateInternalServerError
const AuthenticateInternalServerErrorCode int = 500

/*AuthenticateInternalServerError internal server error

swagger:response authenticateInternalServerError
*/
type AuthenticateInternalServerError struct {
}

// NewAuthenticateInternalServerError creates AuthenticateInternalServerError with default headers values
func NewAuthenticateInternalServerError() *AuthenticateInternalServerError {

	return &AuthenticateInternalServerError{}
}

// WriteResponse to the client
func (o *AuthenticateInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
