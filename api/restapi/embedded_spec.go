package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import "encoding/json"

// SwaggerJSON embedded version of the swagger document used at generation time
var SwaggerJSON json.RawMessage

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Blockchain Event API",
    "version": "1.0"
  },
  "paths": {
    "/v1/event/register": {
      "post": {
        "operationId": "Registration",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/EventRegistrationRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "schema": {
              "$ref": "#/definitions/EventRegistrationResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "EventRegistrationRequest": {
      "type": "object",
      "properties": {
        "contract": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            }
          }
        },
        "event": {
          "type": "object",
          "properties": {
            "abi": {
              "type": "object"
            },
            "name": {
              "type": "string"
            }
          }
        }
      }
    },
    "EventRegistrationResponse": {
      "type": "object",
      "required": [
        "contract"
      ],
      "properties": {
        "contract": {
          "type": "string"
        }
      }
    }
  }
}`))
}
