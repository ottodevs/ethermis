package ethereum

import (
	"crypto/ecdsa"
	"math"
	"math/big"
	"strings"

	"golang.org/x/net/context"

	"github.com/alanchchen/ethermis/api"
	"github.com/alanchchen/ethermis/api/ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
)

func NewController() api.Controller {
	key, _ := crypto.GenerateKey()

	return &ethereumController{
		key: key,
		backend: backends.NewSimulatedBackend(
			core.GenesisAccount{
				Address: crypto.PubkeyToAddress(key.PublicKey),
				Balance: big.NewInt(math.MaxInt64),
			},
		),
	}
}

// ----------------------------------------------------------------------------

type ethereumController struct {
	key     *ecdsa.PrivateKey
	backend *backends.SimulatedBackend
}

func (c *ethereumController) Deploy(ctx context.Context, contract *ethereum.CompiledContract) (*ethereum.DeploymentInfo, error) {
	auth := bind.NewKeyedTransactor(c.key)
	auth.GasLimit = big.NewInt(4700000)
	auth.Value = big.NewInt(47000000)

	parsedABI, err := abi.JSON(strings.NewReader(contract.Abi))
	if err != nil {
		return nil, err
	}

	// Deploy a contract on the simulated blockchain
	address, tx, _, err := bind.DeployContract(
		auth,
		parsedABI,
		[]byte(contract.Code),
		c.backend,
		"Hello")
	if err != nil {
		glog.Errorf("Failed to deploy new token contract: %v", err)
		return nil, err
	}

	return &ethereum.DeploymentInfo{
		DeployedAddress: address.Hex(),
		TransactionId:   tx.Hash().Hex(),
	}, nil
}
