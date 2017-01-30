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

package ethereum

import (
	"runtime"

	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/les"
	"github.com/ethereum/go-ethereum/logger/glog"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

// Backend handles the chain database and VM
type Backend struct {
	ethereum *eth.Ethereum
	config   *eth.Config
}

// NewBackend creates a new Backend
func NewBackend(ctx *node.ServiceContext, config *eth.Config) (*Backend, error) {
	ethereum, err := eth.New(ctx, config)
	if err != nil {
		return nil, err
	}

	if ethereum != nil && config.LightServ > 0 {
		ls, _ := les.NewLesServer(ethereum, config)
		ethereum.AddLesServer(ls)
	}

	ethBackend := &Backend{
		ethereum: ethereum,
		config:   config,
	}

	return ethBackend, nil
}

// APIs returns the collection of RPC services the ethereum package offers.
func (s *Backend) APIs() []rpc.API {
	return s.Ethereum().APIs()
}

// Start implements node.Service, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *Backend) Start(srvr *p2p.Server) error {
	return s.Ethereum().Start(srvr)
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *Backend) Stop() error {
	return s.Ethereum().Stop()
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Backend) Protocols() []p2p.Protocol {
	return s.Ethereum().Protocols()
}

// Ethereum returns the underlying the ethereum object
func (s *Backend) Ethereum() *eth.Ethereum {
	return s.ethereum
}

// Config returns the eth.Config
func (s *Backend) Config() *eth.Config {
	return s.config
}

func MakeFullNode(version uint, identifier string, gitCommit string) *node.Node {
	// Create the default extradata and construct the base node
	var clientInfo = struct {
		Version   uint
		Name      string
		GoVersion string
		Os        string
	}{version, identifier, runtime.Version(), runtime.GOOS}
	extra, err := rlp.EncodeToBytes(clientInfo)
	if err != nil {
		glog.Warning("error setting canonical miner information:", err)
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize.Uint64() {
		glog.Warning("error setting canonical miner information: extra exceeds ", params.MaximumExtraDataSize)
		glog.Warningf("extra: %x\n", extra)
		extra = nil
	}
	stack := MakeNode(identifier, gitCommit)
	RegisterEthService(stack, extra)
	RegisterEthStatsService(stack, ethstatsURL)

	// // Add the API service
	// if err := stack.Register(func(serviceContext *node.ServiceContext) (node.Service, error) {
	// 	return api.New(
	// 		api.UseServiceContext(serviceContext),
	// 		api.UseController(
	// 			NewController(stack),
	// 		),
	// 	)
	// }); err != nil {
	// 	glog.Fatalf("Failed to register the API service: %v", err)
	// }

	return stack
}
