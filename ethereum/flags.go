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
	"math/big"
	"runtime"

	flag "github.com/spf13/pflag"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
)

var (
	EthereumFlags = flag.NewFlagSet("ethereum", flag.ExitOnError)

	dataDir                 string
	keyStoreDir             string
	networkID               int
	identity                string
	natspecEnabled          bool
	fastSync                bool
	lightMode               bool
	lightServ               int
	lightPeers              int
	lightKDF                bool
	cacheSize               int
	trieCacheNum            int
	miningEnabled           bool
	minerThreads            int
	targetGasLimit          string
	autoDAG                 bool
	etherbase               string
	gasPrice                string
	extraData               string
	unlockedAccount         string
	passwordFile            string
	ethstatsURL             string
	metricsEnabled          bool
	fakePOW                 bool
	maxPeers                int
	maxPendingPeers         int
	p2pPort                 int
	bootNodes               string
	nodeKeyFile             string
	nodeKeyHex              string
	natSetting              string
	noDiscover              bool
	discoveryV5             bool
	netRestrict             string
	solcPath                string
	gpoMinGasPrice          string
	gpoMaxGasPrice          string
	gpoFullBlockRatio       int
	gpoBaseStepDown         int
	gpoBaseStepUp           int
	gpoBaseCorrectionFactor int
)

// These are all the command line flags we support.
// If you add to this list, please remember to include the
// flag in the appropriate command definition.
//
// The flags are defined here so their names and help texts
// are the same for all commands.

func init() {

	// General settings
	EthereumFlags.StringVar(&dataDir,
		"datadir",
		node.DefaultDataDir(),
		"Data directory for the databases",
	)

	EthereumFlags.StringVar(&keyStoreDir,
		"keystore",
		"",
		"Directory for the keystore (default = inside the datadir)",
	)

	EthereumFlags.IntVar(&networkID,
		"networkid",
		eth.NetworkId,
		"Network identifier (integer, 0=Olympic (disused), 1=Frontier, 2=Morden (disused), 3=Ropsten)",
	)

	EthereumFlags.StringVar(&identity,
		"identity",
		"",
		"Custom node name",
	)

	EthereumFlags.BoolVar(&natspecEnabled,
		"natspec",
		false,
		"Enable NatSpec confirmation notice",
	)

	EthereumFlags.BoolVar(&fastSync,
		"fast",
		false,
		"Enable fast syncing through state downloads",
	)

	EthereumFlags.BoolVar(&lightMode,
		"light",
		false,
		"Enable light client mode",
	)

	EthereumFlags.IntVar(&lightServ,
		"lightserv",
		0,
		"Maximum percentage of time allowed for serving LES requests (0-90)",
	)

	EthereumFlags.IntVar(&lightPeers,
		"lightpeers",
		20,
		"Maximum number of LES client peers",
	)

	EthereumFlags.BoolVar(&lightKDF,
		"lightkdf",
		false,
		"Reduce key-derivation RAM & CPU usage at some expense of KDF strength",
	)

	// Performance tuning settings
	EthereumFlags.IntVar(&cacheSize,
		"cache",
		128,
		"Megabytes of memory allocated to internal caching (min 16MB / database forced)",
	)

	EthereumFlags.IntVar(&trieCacheNum,
		"trie-cache-gens",
		int(state.MaxTrieCacheGen),
		"Number of trie node generations to keep in memory",
	)

	// Miner settings
	EthereumFlags.BoolVar(&miningEnabled,
		"mine",
		false,
		"Enable mining",
	)

	EthereumFlags.IntVar(&minerThreads,
		"minerthreads",
		runtime.NumCPU(),
		"Number of CPU threads to use for mining",
	)

	EthereumFlags.StringVar(&targetGasLimit,
		"targetgaslimit",
		params.GenesisGasLimit.String(),
		"Target gas limit sets the artificial target gas floor for the blocks to mine",
	)

	EthereumFlags.BoolVar(&autoDAG,
		"autodag",
		false,
		"Enable automatic DAG pregeneration",
	)

	EthereumFlags.StringVar(&etherbase,
		"etherbase",
		"",
		"Public address for block mining rewards (default = first account created)",
	)

	EthereumFlags.StringVar(&gasPrice,
		"gasprice",
		new(big.Int).Mul(big.NewInt(20), common.Shannon).String(),
		"Minimal gas price to accept for mining a transactions",
	)

	EthereumFlags.StringVar(&extraData,
		"extradata",
		"",
		"Block extra data set by the miner (default = client version)",
	)

	// Account settings
	EthereumFlags.StringVar(&unlockedAccount,
		"unlock",
		"",
		"Comma separated list of accounts to unlock",
	)

	EthereumFlags.StringVar(&passwordFile,
		"password",
		"",
		"Password file to use for non-inteactive password input",
	)

	// Logging and debug settings
	EthereumFlags.StringVar(&ethstatsURL,
		"ethstats",
		"",
		"Reporting URL of a ethstats service (nodename:secret@host:port)",
	)

	EthereumFlags.BoolVar(&metricsEnabled,
		metrics.MetricsEnabledFlag,
		false,
		"Enable metrics collection and reporting",
	)

	EthereumFlags.BoolVar(&fakePOW,
		"fakepow",
		false,
		"Disables proof-of-work verification",
	)

	// Network Settings
	EthereumFlags.IntVar(&maxPeers,
		"maxpeers",
		25,
		"Maximum number of network peers (network disabled if set to 0)",
	)

	EthereumFlags.IntVar(&maxPendingPeers,
		"maxpendpeers",
		0,
		"Maximum number of pending connection attempts (defaults used if set to 0)",
	)

	EthereumFlags.IntVar(&p2pPort,
		"p2pport",
		30303,
		"P2P network listening port",
	)

	EthereumFlags.StringVar(&bootNodes,
		"bootnodes",
		"",
		"Comma separated enode URLs for P2P discovery bootstrap",
	)

	EthereumFlags.StringVar(&nodeKeyFile,
		"nodekey",
		"",
		"P2P node key file",
	)

	EthereumFlags.StringVar(&nodeKeyHex,
		"nodekeyhex",
		"",
		"P2P node key as hex (for testing)",
	)

	EthereumFlags.StringVar(&natSetting,
		"nat",
		"any",
		"NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>)",
	)

	EthereumFlags.BoolVar(&noDiscover,
		"nodiscover",
		false,
		"Disables the peer discovery mechanism (manual peer addition)",
	)

	EthereumFlags.BoolVar(&discoveryV5,
		"v5disc",
		false,
		"Enables the experimental RLPx V5 (Topic Discovery) mechanism",
	)

	EthereumFlags.StringVar(&netRestrict,
		"netrestrict",
		"",
		"Restricts network communication to the given IP networks (CIDR masks)",
	)

	EthereumFlags.StringVar(&solcPath,
		"solc",
		"solc",
		"Solidity compiler command to be used",
	)

	// Gas price oracle settings
	EthereumFlags.StringVar(&gpoMinGasPrice,
		"gpomin",
		new(big.Int).Mul(big.NewInt(20), common.Shannon).String(),
		"Minimum suggested gas price",
	)

	EthereumFlags.StringVar(&gpoMaxGasPrice,
		"gpomax",
		new(big.Int).Mul(big.NewInt(500), common.Shannon).String(),
		"Maximum suggested gas price",
	)

	EthereumFlags.IntVar(&gpoFullBlockRatio,
		"gpofull",
		80,
		"Full block threshold for gas price calculation (%)",
	)

	EthereumFlags.IntVar(&gpoBaseStepDown,
		"gpobasedown",
		10,
		"Suggested gas price base step down ratio (1/1000)",
	)

	EthereumFlags.IntVar(&gpoBaseStepUp,
		"gpobaseup",
		100,
		"Suggested gas price base step up ratio (1/1000)",
	)

	EthereumFlags.IntVar(&gpoBaseCorrectionFactor,
		"gpobasecf",
		110,
		"Suggested gas price base correction factor (%)",
	)
}
