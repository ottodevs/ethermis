package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// EventRegistrationRequest event registration request
// swagger:model EventRegistrationRequest
type EventRegistrationRequest struct {

	// contract
	Contract *EventRegistrationRequestContract `json:"contract,omitempty"`

	// event
	Event *EventRegistrationRequestEvent `json:"event,omitempty"`
}

// Validate validates this event registration request
func (m *EventRegistrationRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContract(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateEvent(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EventRegistrationRequest) validateContract(formats strfmt.Registry) error {

	if swag.IsZero(m.Contract) { // not required
		return nil
	}

	if m.Contract != nil {

		if err := m.Contract.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("contract")
			}
			return err
		}
	}

	return nil
}

func (m *EventRegistrationRequest) validateEvent(formats strfmt.Registry) error {

	if swag.IsZero(m.Event) { // not required
		return nil
	}

	if m.Event != nil {

		if err := m.Event.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("event")
			}
			return err
		}
	}

	return nil
}

// EventRegistrationRequestContract event registration request contract
// swagger:model EventRegistrationRequestContract
type EventRegistrationRequestContract struct {

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this event registration request contract
func (m *EventRegistrationRequestContract) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// EventRegistrationRequestEvent event registration request event
// swagger:model EventRegistrationRequestEvent
type EventRegistrationRequestEvent struct {

	// abi
	Abi interface{} `json:"abi,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this event registration request event
func (m *EventRegistrationRequestEvent) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}