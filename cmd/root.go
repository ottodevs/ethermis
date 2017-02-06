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
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/alanchchen/ethermis/api"
	"github.com/alanchchen/ethermis/ethereum"
)

var cfgFile string
var stack *node.Node

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ethermis",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// stack = ethereum.MakeFullNode(
		// 	uint(constant.VersionMajor<<16|constant.VersionMinor<<8|constant.VersionPatch),
		// 	constant.ClientIdentifier,
		// 	constant.GitCommit,
		// )

		//utils.StartNode(stack)

		// Add the API service
		apiService := api.New(
			// api.UseController(
			ethereum.NewController(),
		// ),
		)
		if apiService == nil {
			cmd.Println("failed to initialize API service")
			return
		}
		apiService.Start()
		//stack.Wait()
		apiService.Stop()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ethermis.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().AddFlagSet(ethereum.EthereumFlags)
	RootCmd.PersistentFlags().AddFlagSet(api.APIServiceFlags)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ethermis") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
