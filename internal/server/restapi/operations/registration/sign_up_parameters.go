// Code generated by go-swagger; DO NOT EDIT.

package registration

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewSignUpParams creates a new SignUpParams object
//
// There are no default values defined in the spec.
func NewSignUpParams() SignUpParams {

	return SignUpParams{}
}

// SignUpParams contains all the bound params for the sign up operation
// typically these are obtained from a http.Request
//
// swagger:parameters signUp
type SignUpParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Домен
	  Required: true
	  Max Length: 128
	  Min Length: 1
	  In: query
	*/
	Domain string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSignUpParams() beforehand.
func (o *SignUpParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qDomain, qhkDomain, _ := qs.GetOK("domain")
	if err := o.bindDomain(qDomain, qhkDomain, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindDomain binds and validates parameter Domain from query.
func (o *SignUpParams) bindDomain(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("domain", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("domain", "query", raw); err != nil {
		return err
	}
	o.Domain = raw

	if err := o.validateDomain(formats); err != nil {
		return err
	}

	return nil
}

// validateDomain carries on validations for parameter Domain
func (o *SignUpParams) validateDomain(formats strfmt.Registry) error {

	if err := validate.MinLength("domain", "query", o.Domain, 1); err != nil {
		return err
	}

	if err := validate.MaxLength("domain", "query", o.Domain, 128); err != nil {
		return err
	}

	return nil
}
