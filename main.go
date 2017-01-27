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

package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"gopkg.in/urfave/cli.v1"

	"github.com/alanchchen/ethermis/api"
	"github.com/alanchchen/ethermis/log"
	"github.com/alanchchen/ethermis/utils"
)

const (
	clientIdentifier = "Ethermis" // Client identifier to advertise over the network
	versionMajor     = 0          // Major version component of the current release
	versionMinor     = 1          // Minor version component of the current release
	versionPatch     = 0          // Patch version component of the current release
	versionMeta      = "unstable" // Version metadata to append to the version string
)

var (
	gitCommit string
	verString string // Combined textual representation of all the version components
	cliApp    *cli.App
	stack     *node.Node
)

func init() {
	verString = fmt.Sprintf("%d.%d.%d", versionMajor, versionMinor, versionPatch)
	if versionMeta != "" {
		verString += "-" + versionMeta
	}
	cliApp = newCliApp(verString, "the Ethermis command line interface")
	cliApp.Action = start
	cliApp.Commands = []cli.Command{
		{
			Action:      initCommand,
			Name:        "init",
			Usage:       "init genesis.json",
			Description: "",
		},
	}
	cliApp.HideVersion = true // we have a command to print the version

	cliApp.After = func(ctx *cli.Context) error {
		log.Flush()
		return nil
	}
}

func main() {
	log.Info("Starting Ethermis")

	if err := cliApp.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initCommand(ctx *cli.Context) error {
	// ethereum genesis.json
	genesisPath := ctx.Args().First()
	if len(genesisPath) == 0 {
		log.Fatalf("must supply path to genesis JSON file")
	}

	chainDb, err := ethdb.NewLDBDatabase(filepath.Join(utils.MakeDataDir(ctx), clientIdentifier, "chaindata"), 0, 0)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}

	genesisFile, err := os.Open(genesisPath)
	if err != nil {
		log.Fatalf("failed to read genesis file: %v", err)
	}

	block, err := core.WriteGenesisBlock(chainDb, genesisFile)
	if err != nil {
		log.Fatalf("failed to write genesis block: %v", err)
	}
	log.Infof("successfully wrote genesis block and/or chain rule set: %x", block.Hash())
	return nil
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		if user, err := user.Current(); err == nil {
			p = user.HomeDir + p[1:]
		}
	}
	return path.Clean(os.ExpandEnv(p))
}

type DirectoryString struct {
	Value string
}

func (self *DirectoryString) String() string {
	return self.Value
}

func (self *DirectoryString) Set(value string) error {
	self.Value = expandPath(value)
	return nil
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Ethermint")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Ethermint")
		} else {
			return filepath.Join(home, ".ethermint")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func newCliApp(version, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = "Alan Chen"
	app.Email = "alan@am.is"
	app.Version = version
	app.Usage = usage
	app.Flags = []cli.Flag{
		utils.IdentityFlag,
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.BootnodesFlag,
		utils.DataDirFlag,
		utils.KeyStoreDirFlag,
		utils.CacheFlag,
		utils.LightKDFFlag,
		utils.JSpathFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.EtherbaseFlag,
		utils.TargetGasLimitFlag,
		utils.GasPriceFlag,
		utils.NATFlag,
		utils.NatspecEnabledFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.IPCDisabledFlag,
		utils.IPCApiFlag,
		utils.IPCPathFlag,
		utils.ExecFlag,
		utils.PreloadJSFlag,
		utils.TestNetFlag,
		utils.VMForceJitFlag,
		utils.VMJitCacheFlag,
		utils.VMEnableJitFlag,
		utils.NetworkIdFlag,
		utils.RPCCORSDomainFlag,
		utils.MetricsEnabledFlag,
		utils.SolcPathFlag,
		utils.GpoMinGasPriceFlag,
		utils.GpoMaxGasPriceFlag,
		utils.GpoFullBlockRatioFlag,
		utils.GpobaseStepDownFlag,
		utils.GpobaseStepUpFlag,
		utils.GpobaseCorrectionFactorFlag,
		utils.EthStatsURLFlag,
		utils.HostFlag,
		utils.PortFlag,
	}
	return app
}

func start(ctx *cli.Context) error {
	stack = makeFullNode(ctx)

	utils.StartNode(stack)
	stack.Wait()

	return nil
}

func makeFullNode(ctx *cli.Context) *node.Node {
	// Create the default extradata and construct the base node
	var clientInfo = struct {
		Version   uint
		Name      string
		GoVersion string
		Os        string
	}{uint(versionMajor<<16 | versionMinor<<8 | versionPatch), clientIdentifier, runtime.Version(), runtime.GOOS}
	extra, err := rlp.EncodeToBytes(clientInfo)
	if err != nil {
		log.Warning("error setting canonical miner information:", err)
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize.Uint64() {
		log.Warning("error setting canonical miner information: extra exceeds", params.MaximumExtraDataSize)
		log.Warningf("extra: %x\n", extra)
		extra = nil
	}
	stack := utils.MakeNode(ctx, clientIdentifier, gitCommit)
	utils.RegisterEthService(ctx, stack, extra)

	// Add the Ethereum Stats daemon if requested
	if url := ctx.GlobalString(utils.EthStatsURLFlag.Name); url != "" {
		utils.RegisterEthStatsService(stack, url)
	}

	// Add the API service
	if err := stack.Register(func(serviceContext *node.ServiceContext) (node.Service, error) {
		return api.New(ctx, serviceContext)
	}); err != nil {
		log.Fatalf("Failed to register the API service: %v", err)
	}

	return stack
}
