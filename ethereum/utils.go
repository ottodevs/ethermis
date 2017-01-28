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

// Package utils contains internal helper functions for go-ethereum commands.
package ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethstats"
	"github.com/ethereum/go-ethereum/les"
	"github.com/ethereum/go-ethereum/logger/glog"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"github.com/ethereum/go-ethereum/params"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv2"
)

func init() {
	// 	cli.AppHelpTemplate = `{{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

	// VERSION:
	//    {{.Version}}

	// COMMANDS:
	//    {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
	//    {{end}}{{if .Flags}}
	// GLOBAL OPTIONS:
	//    {{range .Flags}}{{.}}
	//    {{end}}{{end}}
	// `

	// 	cli.CommandHelpTemplate = `{{.Name}}{{if .Subcommands}} command{{end}}{{if .Flags}} [command options]{{end}} [arguments...]
	// {{if .Description}}{{.Description}}
	// {{end}}{{if .Subcommands}}
	// SUBCOMMANDS:
	// 	{{range .Subcommands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
	// 	{{end}}{{end}}{{if .Flags}}
	// OPTIONS:
	// 	{{range .Flags}}{{.}}
	// 	{{end}}{{end}}
	// `
}

func StartNode(stack *node.Node) {
	if err := stack.Start(); err != nil {
		glog.Fatalf("Error starting protocol stack: %v", err)
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)
		defer signal.Stop(sigc)
		// <-sigc
		glog.Infoln("Got interrupt, shutting down...")
		go stack.Stop()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				glog.Infof("Already shutting down, interrupt %d more times for panic.", i-1)
			}
		}
	}()
}

// MakeDataDir retrieves the currently requested data directory, terminating
// if none (or the empty string) is specified. If the node is starting a testnet,
// the a subdirectory of the specified datadir will be used.
func MakeDataDir() string {
	if path := dataDir; path != "" {
		return path
	}
	glog.Fatalf("Cannot determine default data directory, please set manually (--datadir)")
	return ""
}

// MakeNodeKey creates a node key from set command line flags, either loading it
// from a file or as a specified hex value. If neither flags were provided, this
// method returns nil and an emphemeral key is to be generated.
func MakeNodeKey() *ecdsa.PrivateKey {
	var (
		hex  = nodeKeyHex
		file = nodeKeyFile

		key *ecdsa.PrivateKey
		err error
	)
	switch {
	case file != "" && hex != "":
		glog.Fatalf("Options nodekeyhex and nodekey are mutually exclusive")

	case file != "":
		if key, err = crypto.LoadECDSA(file); err != nil {
			glog.Fatalf("Option nodekey: %v", err)
		}

	case hex != "":
		if key, err = crypto.HexToECDSA(hex); err != nil {
			glog.Fatalf("Option nodekeyhex: %v", err)
		}
	}
	return key
}

// makeNodeUserIdent creates the user identifier from CLI flags.
func makeNodeUserIdent() string {
	var comps []string
	if identity := identity; len(identity) > 0 {
		comps = append(comps, identity)
	}

	return strings.Join(comps, "/")
}

// MakeBootstrapNodes creates a list of bootstrap nodes from the command line
// flags, reverting to pre-configured ones if none have been specified.
func MakeBootstrapNodes() []*discover.Node {
	urls := []string{}
	if bootNodes != "" {
		urls = strings.Split(bootNodes, ",")
	}

	bootnodes := make([]*discover.Node, 0, len(urls))
	for _, url := range urls {
		node, err := discover.ParseNode(url)
		if err != nil {
			glog.Infof("Bootstrap URL %s: %v\n", url, err)
			continue
		}
		bootnodes = append(bootnodes, node)
	}
	return bootnodes
}

// MakeBootstrapNodesV5 creates a list of bootstrap nodes from the command line
// flags, reverting to pre-configured ones if none have been specified.
func MakeBootstrapNodesV5() []*discv5.Node {
	urls := params.DiscoveryV5Bootnodes
	if bootNodes != "" {
		urls = strings.Split(bootNodes, ",")
	}

	bootnodes := make([]*discv5.Node, 0, len(urls))
	for _, url := range urls {
		node, err := discv5.ParseNode(url)
		if err != nil {
			glog.Errorf("Bootstrap URL %s: %v\n", url, err)
			continue
		}
		bootnodes = append(bootnodes, node)
	}
	return bootnodes
}

// MakeListenAddress creates a TCP listening address string from set command
// line flags.
func MakeListenAddress() string {
	return fmt.Sprintf(":%d", p2pPort)
}

// MakeDiscoveryV5Address creates a UDP listening address string from set command
// line flags for the V5 discovery protocol.
func MakeDiscoveryV5Address() string {
	return fmt.Sprintf(":%d", p2pPort+1)
}

// MakeNAT creates a port mapper from set command line flags.
func MakeNAT() nat.Interface {
	natif, err := nat.Parse(natSetting)
	if err != nil {
		glog.Fatalf("Invalid NAT '%s': %v", natSetting, err)
	}
	return natif
}

// MakeDatabaseHandles raises out the number of allowed file handles per process
// for Geth and returns half of the allowance to assign to the database.
func MakeDatabaseHandles() int {
	if err := raiseFdLimit(2048); err != nil {
		glog.Fatalf("Failed to raise file descriptor allowance: %v", err)
	}
	limit, err := getFdLimit()
	if err != nil {
		glog.Fatalf("Failed to retrieve file descriptor allowance: %v", err)
	}
	if limit > 2048 { // cap database file descriptors even if more is available
		limit = 2048
	}
	return limit / 2 // Leave half for networking and other stuff
}

// MakeAddress converts an account specified directly as a hex encoded string or
// a key index in the key store to an internal account representation.
func MakeAddress(accman *accounts.Manager, account string) (accounts.Account, error) {
	// If the specified account is a valid address, return it
	if common.IsHexAddress(account) {
		return accounts.Account{Address: common.HexToAddress(account)}, nil
	}
	// Otherwise try to interpret the account as a keystore index
	index, err := strconv.Atoi(account)
	if err != nil {
		return accounts.Account{}, fmt.Errorf("invalid account address or index %q", account)
	}
	return accman.AccountByIndex(index)
}

// MakeEtherbase retrieves the etherbase either from the directly specified
// command line flags or from the keystore if CLI indexed.
func MakeEtherbase(accman *accounts.Manager) common.Address {
	accounts := accman.Accounts()
	if etherbase == "" && len(accounts) == 0 {
		glog.Infoln("WARNING: No etherbase set and no accounts found as default")
		return common.Address{}
	}

	if etherbase == "" {
		return common.Address{}
	}
	// If the specified etherbase is a valid address, return it
	account, err := MakeAddress(accman, etherbase)
	if err != nil {
		glog.Fatalf("Invalid etherbase %q: %v", etherbase, err)
	}
	return account.Address
}

// MakeMinerExtra resolves extradata for the miner from the set command line flags
// or returns a default one composed on the client, runtime and OS metadata.
func MakeMinerExtra(extra []byte) []byte {
	if extraData != "" {
		return []byte(extraData)
	}
	return extra
}

// MakePasswordList reads password lines from the file specified by --password.
func MakePasswordList() []string {
	path := passwordFile
	if path == "" {
		return nil
	}
	text, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Fatalf("Failed to read password file: %v", err)
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines
}

// MakeNode configures a node with no services from command line flags.
func MakeNode(name, gitCommit string) *node.Node {
	vsn := params.Version
	if gitCommit != "" {
		vsn += "-" + gitCommit[:8]
	}

	config := &node.Config{
		DataDir:           MakeDataDir(),
		KeyStoreDir:       keyStoreDir,
		UseLightweightKDF: lightKDF,
		PrivateKey:        MakeNodeKey(),
		Name:              name,
		Version:           vsn,
		UserIdent:         makeNodeUserIdent(),
		NoDiscovery:       noDiscover || lightMode,
		DiscoveryV5:       discoveryV5 || lightMode || lightServ > 0,
		DiscoveryV5Addr:   MakeDiscoveryV5Address(),
		BootstrapNodes:    MakeBootstrapNodes(),
		BootstrapNodesV5:  MakeBootstrapNodesV5(),
		ListenAddr:        MakeListenAddress(),
		NAT:               MakeNAT(),
		MaxPeers:          maxPeers,
		MaxPendingPeers:   maxPendingPeers,
	}

	if netrestrict := netRestrict; netrestrict != "" {
		list, err := netutil.ParseNetlist(netrestrict)
		if err != nil {
			glog.Fatalf("Option netrestrict: %v", err)
		}
		config.NetRestrict = list
	}

	stack, err := node.New(config)
	if err != nil {
		glog.Fatalf("Failed to create the protocol stack: %v", err)
	}
	return stack
}

// RegisterEthService configures eth.Ethereum from command line flags and adds it to the
// given node.
func RegisterEthService(stack *node.Node, extra []byte) {
	ethConf := &eth.Config{
		Etherbase:               MakeEtherbase(stack.AccountManager()),
		ChainConfig:             MakeChainConfig(stack),
		FastSync:                fastSync,
		LightMode:               lightMode,
		LightServ:               lightServ,
		LightPeers:              lightPeers,
		MaxPeers:                maxPeers,
		DatabaseCache:           cacheSize,
		DatabaseHandles:         MakeDatabaseHandles(),
		NetworkId:               networkID,
		MinerThreads:            minerThreads,
		ExtraData:               MakeMinerExtra(extra),
		NatSpec:                 natspecEnabled,
		GasPrice:                common.String2Big(gasPrice),
		GpoMinGasPrice:          common.String2Big(gpoMinGasPrice),
		GpoMaxGasPrice:          common.String2Big(gpoMaxGasPrice),
		GpoFullBlockRatio:       gpoFullBlockRatio,
		GpobaseStepDown:         gpoBaseStepDown,
		GpobaseStepUp:           gpoBaseStepUp,
		GpobaseCorrectionFactor: gpoBaseCorrectionFactor,
		SolcPath:                solcPath,
		AutoDAG:                 autoDAG || miningEnabled,
		// DocRoot:                 ""
	}

	// Override any global options pertaining to the Ethereum protocol
	if gen := trieCacheNum; gen > 0 {
		state.MaxTrieCacheGen = uint16(gen)
	}

	if ethConf.LightMode {
		if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			return les.New(ctx, ethConf)
		}); err != nil {
			glog.Fatalf("Failed to register the Ethereum light node service: %v", err)
		}
	} else {
		if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			return NewBackend(ctx, ethConf)
		}); err != nil {
			glog.Fatalf("Failed to register the Ethereum full node service: %v", err)
		}
	}
}

// RegisterShhService configures Whisper and adds it to the given node.
func RegisterShhService(stack *node.Node) {
	if err := stack.Register(func(*node.ServiceContext) (node.Service, error) { return whisper.New(), nil }); err != nil {
		glog.Fatalf("Failed to register the Whisper service: %v", err)
	}
}

// RegisterEthStatsService configures the Ethereum Stats daemon and adds it to
// th egiven node.
func RegisterEthStatsService(stack *node.Node, url string) {
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		// Retrieve both eth and les services
		var ethServ *Backend
		ctx.Service(&ethServ)

		var lesServ *les.LightEthereum
		ctx.Service(&lesServ)

		return ethstats.New(url, ethServ.Ethereum(), lesServ)
	}); err != nil {
		glog.Fatalf("Failed to register the Ethereum Stats service: %v", err)
	}
}

// SetupNetwork configures the system for either the main net or some test network.
func SetupNetwork() {
	params.TargetGasLimit = common.String2Big(targetGasLimit)
}

// MakeChainConfig reads the chain configuration from the database in ctx.Datadir.
func MakeChainConfig(stack *node.Node) *params.ChainConfig {
	db := MakeChainDatabase(stack)
	defer db.Close()

	return MakeChainConfigFromDb(db)
}

// MakeChainConfigFromDb reads the chain configuration from the given database.
func MakeChainConfigFromDb(db ethdb.Database) *params.ChainConfig {
	// If the chain is already initialized, use any existing chain configs
	config := new(params.ChainConfig)

	genesis := core.GetBlock(db, core.GetCanonicalHash(db, 0), 0)
	if genesis != nil {
		storedConfig, err := core.GetChainConfig(db, genesis.Hash())
		switch err {
		case nil:
			config = storedConfig
		case core.ChainConfigNotFoundErr:
			// No configs found, use empty, will populate below
		default:
			glog.Fatalf("Could not make chain configuration: %v", err)
		}
	}
	// set chain id in case it's zero.
	if config.ChainId == nil {
		config.ChainId = new(big.Int)
	}

	return config
}

func ChainDbName() string {
	if lightMode {
		return "lightchaindata"
	} else {
		return "chaindata"
	}
}

// MakeChainDatabase open an LevelDB using the flags passed to the client and will hard crash if it fails.
func MakeChainDatabase(stack *node.Node) ethdb.Database {
	var (
		cache   = cacheSize
		handles = MakeDatabaseHandles()
		name    = ChainDbName()
	)

	chainDb, err := stack.OpenDatabase(name, cache, handles)
	if err != nil {
		glog.Fatalf("Could not open database: %v", err)
	}
	return chainDb
}
