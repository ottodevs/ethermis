// Copyright 2017 The Ethermis Authors
// This file is part of Ethermis.
//
// Ethermis is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Ethermis is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Ethermis. If not, see <http://www.gnu.org/licenses/>.

package restapi

import (
	"encoding/json"

	"github.com/alanchchen/ethermis/api/models"
	"github.com/alanchchen/ethermis/api/restapi/operations"
	"github.com/alanchchen/ethermis/log"
	"github.com/go-openapi/runtime/middleware"
)

func Register(params operations.RegistrationParams) middleware.Responder {
	raw, _ := json.MarshalIndent(params.Body, "", "\t")
	log.Info(string(raw))
	response := operations.NewRegistrationOK()
	return response.WithPayload(&models.EventRegistrationResponse{
		Contract: &params.Body.Contract.Name,
	})
}
