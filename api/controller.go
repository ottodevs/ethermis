package api

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/alanchchen/ethermis/api/ethereum"
)

type Controller interface {
	Deploy(ctx context.Context, contract *ethereum.CompiledContract) (*ethereum.DeploymentInfo, error)
}

type controller struct{}

func (c *controller) Deploy(ctx context.Context, contract *ethereum.CompiledContract) (*ethereum.DeploymentInfo, error) {
	fmt.Println(contract)
	return &ethereum.DeploymentInfo{
		DeployedAddress: "0x1234567890",
		TransactionId:   "0xABCDEF",
	}, nil
}
