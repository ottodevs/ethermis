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

package cmd

import (
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/spf13/cobra"

	"github.com/alanchchen/ethermis/constant"
	"github.com/alanchchen/ethermis/ethereum"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap and initialize a new genesis blockd",
	Long:  `Bootstrap and initialize a new genesis block`,
	Run: func(cmd *cobra.Command, args []string) {
		// ethereum genesis.json
		genesisPath := args[0]
		if len(genesisPath) == 0 {
			cmd.Println("must supply path to genesis JSON file")
			return
		}

		chainDb, err := ethdb.NewLDBDatabase(filepath.Join(ethereum.MakeDataDir(), constant.ClientIdentifier, "chaindata"), 0, 0)
		if err != nil {
			cmd.Printf("could not open database: %v\n", err)
			return
		}

		genesisFile, err := os.Open(genesisPath)
		if err != nil {
			cmd.Printf("failed to read genesis file: %v\n", err)
			return
		}

		block, err := core.WriteGenesisBlock(chainDb, genesisFile)
		if err != nil {
			cmd.Printf("failed to write genesis block: %v\n", err)
			return
		}
		cmd.Printf("successfully wrote genesis block and/or chain rule set: %x\n", block.Hash())
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
