// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// CreateUserOKCode is the HTTP code returned for type CreateUserOK
const CreateUserOKCode int = 200

/*
CreateUserOK OK

swagger:response createUserOK
*/
type CreateUserOK struct {

	/*
	  In: Body
	*/
	Payload *CreateUserOKBody `json:"body,omitempty"`
}

// NewCreateUserOK creates CreateUserOK with default headers values
func NewCreateUserOK() *CreateUserOK {

	return &CreateUserOK{}
}

// WithPayload adds the payload to the create user o k response
func (o *CreateUserOK) WithPayload(payload *CreateUserOKBody) *CreateUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create user o k response
func (o *CreateUserOK) SetPayload(payload *CreateUserOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateUserInternalServerErrorCode is the HTTP code returned for type CreateUserInternalServerError
const CreateUserInternalServerErrorCode int = 500

/*
CreateUserInternalServerError Ошибка на стороне сервера

swagger:response createUserInternalServerError
*/
type CreateUserInternalServerError struct {
}

// NewCreateUserInternalServerError creates CreateUserInternalServerError with default headers values
func NewCreateUserInternalServerError() *CreateUserInternalServerError {

	return &CreateUserInternalServerError{}
}

// WriteResponse to the client
func (o *CreateUserInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
