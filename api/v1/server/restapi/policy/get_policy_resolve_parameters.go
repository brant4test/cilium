// Code generated by go-swagger; DO NOT EDIT.

package policy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/cilium/cilium/api/v1/models"
)

// NewGetPolicyResolveParams creates a new GetPolicyResolveParams object
// with the default values initialized.
func NewGetPolicyResolveParams() GetPolicyResolveParams {
	var ()
	return GetPolicyResolveParams{}
}

// GetPolicyResolveParams contains all the bound params for the get policy resolve operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetPolicyResolve
type GetPolicyResolveParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*Context to provide policy evaluation on
	  In: body
	*/
	IdentityContext *models.IdentityContext
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetPolicyResolveParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.IdentityContext
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("identityContext", "body", "", err))
		} else {
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.IdentityContext = &body
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
