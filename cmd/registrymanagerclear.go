// Copyright © 2019 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wealdtech/ethereal/cli"
	ens "github.com/wealdtech/go-ens"
	erc1820 "github.com/wealdtech/go-erc1820"
)

// registryManagerClearCmd represents the registry manager set command
var registryManagerClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the address of an ERC-1820 address manager",
	Long: `Clear the manager of an address in the ERC-1820 registry.  For example:

    ethereal registry manager clear --address=0x1234...5678

In quiet mode this will return 0 if the transaction to clear the manager is sent successfully, otherwise 1.`,

	Run: func(cmd *cobra.Command, args []string) {
		address, err := ens.Resolve(client, registryManagerAddressStr)
		cli.ErrCheck(err, quiet, "failed to resolve address")

		registry, err := erc1820.NewRegistry(client)
		cli.ErrCheck(err, quiet, "failed to obtain ERC-1820 registry")

		existingManager, err := registry.Manager(&address)
		cli.ErrCheck(err, quiet, "failed to obtain existing manager")

		opts, err := generateTxOpts(*existingManager)
		cli.ErrCheck(err, quiet, "failed to generate transaction options")
		signedTx, err := registry.SetManager(opts, &address, &ens.UnknownAddress)
		cli.ErrCheck(err, quiet, "failed to send transaction")

		logTransaction(signedTx, log.Fields{
			"group":   "registry/manager",
			"command": "clear",
			"address": address.Hex(),
		})

		if quiet {
			os.Exit(0)
		}

		fmt.Println(signedTx.Hash().Hex())
		os.Exit(0)
	},
}

func init() {
	registryManagerFlags(registryManagerClearCmd)
	registryManagerCmd.AddCommand(registryManagerClearCmd)
	addTransactionFlags(registryManagerClearCmd, "passphrase for the address")
}