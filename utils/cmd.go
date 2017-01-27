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

package utils

import (
	"log"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/logger"
	"github.com/ethereum/go-ethereum/node"
	"github.com/golang/glog"
)

func StartNode(stack *node.Node) {
	if err := stack.Start(); err != nil {
		log.Fatalf("Error starting protocol stack: %v", err)
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)
		defer signal.Stop(sigc)
		// <-sigc
		glog.V(logger.Info).Infoln("Got interrupt, shutting down...")
		go stack.Stop()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				glog.V(logger.Info).Infof("Already shutting down, interrupt %d more times for panic.", i-1)
			}
		}
	}()
}
