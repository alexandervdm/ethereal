// Copyright © 2017-2019 Weald Technology Trading
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
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wealdtech/ethereal/cli"
	ens "github.com/wealdtech/go-ens"
)

var ensOwnerSetOwnerStr string

// ensOwnerSetCmd represents the ens owner set command
var ensOwnerSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the owner of an ENS domain",
	Long: `Set the owner of a name registered with the Ethereum Name Service (ENS).  For example:

    ethereal ens owner set --domain=enstest.eth --owner=0x1234...5678 --passphrase="my secret passphrase"

The keystore for the account that owns the name must be local (i.e. listed with 'get accounts list') and unlockable with the supplied passphrase.

This will return an exit status of 0 if the transaction is successfully submitted (and mined if --wait is supplied), 1 if the transaction is not successfully submitted, and 2 if the transaction is successfully submitted but not mined within the supplied time limit.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Assert(!offline, quiet, "Offline mode not supported at current with this command")
		cli.Assert(ensDomain != "", quiet, "--domain is required")

		registry, err := ens.NewRegistry(client)
		cli.ErrCheck(err, quiet, "Cannot obtain ENS registry contract")

		// Fetch the owner of the name
		owner, err := registry.Owner(ensDomain)
		cli.ErrCheck(err, quiet, "Cannot obtain owner")
		cli.Assert(bytes.Compare(owner.Bytes(), ens.UnknownAddress.Bytes()) != 0, quiet, fmt.Sprintf("Owner of %s is not set", ensDomain))

		cli.Assert(ensOwnerSetOwnerStr != "", quiet, "--owner is required")
		newOwnerAddress, err := ens.Resolve(client, ensOwnerSetOwnerStr)
		cli.Assert(bytes.Compare(newOwnerAddress.Bytes(), ens.UnknownAddress.Bytes()) != 0, quiet, "Attempt to set owner to 0x00 disallowed")
		cli.ErrCheck(err, quiet, "Failed to obtain new owner address")

		opts, err := generateTxOpts(owner)
		cli.ErrCheck(err, quiet, "Failed to generate transaction options")
		signedTx, err := registry.SetOwner(opts, ensDomain, newOwnerAddress)
		cli.ErrCheck(err, quiet, "Failed to send transaction")

		handleSubmittedTransaction(signedTx, log.Fields{
			"group":     "ens/owner",
			"command":   "set",
			"ensdomain": ensDomain,
			"ensowner":  newOwnerAddress.Hex(),
		})
	},
}

func init() {
	ensOwnerCmd.AddCommand(ensOwnerSetCmd)
	ensOwnerFlags(ensOwnerSetCmd)
	ensOwnerSetCmd.Flags().StringVar(&ensOwnerSetOwnerStr, "owner", "", "The new owner's name or address")
	addTransactionFlags(ensOwnerSetCmd, "passphrase for the account that owns the domain")
}
