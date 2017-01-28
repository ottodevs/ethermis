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

package api

import (
	"github.com/ethereum/go-ethereum/logger/glog"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-openapi/loads"

	"github.com/alanchchen/ethermis/api/restapi"
	"github.com/alanchchen/ethermis/api/restapi/operations"
)

type service struct {
	context *node.ServiceContext
	server  *restapi.Server
}

func New(serviceContext *node.ServiceContext) (node.Service, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		glog.Errorf("Failed to load swagger spec, err=%s", err)
		return nil, err
	}

	api := operations.NewBlockchainEventAPI(swaggerSpec)

	server := restapi.NewServer(api)
	server.EnabledListeners = []string{"http"}

	server.ConfigureFlags()
	server.ConfigureAPI()

	// validate the API descriptor, to ensure we don't have any unhandled operations
	if err := api.Validate(); err != nil {
		glog.Errorf("Invalid API handler, err=%s", err)
		return nil, err
	}

	service := &service{
		context: serviceContext,
		server:  server,
	}

	return service, nil
}

func (s *service) APIs() []rpc.API {
	return nil
}

// Start implements node.Service, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *service) Start(srvr *p2p.Server) error {
	return s.server.Serve()
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *service) Stop() error {
	return s.server.Shutdown()
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *service) Protocols() []p2p.Protocol {
	return nil
}
