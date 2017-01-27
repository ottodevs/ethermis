package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/alanchchen/ethermis/api/models"
)

/*RegistrationOK registration o k

swagger:response registrationOK
*/
type RegistrationOK struct {

	/*
	  In: Body
	*/
	Payload *models.EventRegistrationResponse `json:"body,omitempty"`
}

// NewRegistrationOK creates RegistrationOK with default headers values
func NewRegistrationOK() *RegistrationOK {
	return &RegistrationOK{}
}

// WithPayload adds the payload to the registration o k response
func (o *RegistrationOK) WithPayload(payload *models.EventRegistrationResponse) *RegistrationOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the registration o k response
func (o *RegistrationOK) SetPayload(payload *models.EventRegistrationResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RegistrationOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
