// Copyright © 2017 Weald Technology Trading
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
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wealdtech/ethereal/cli"
	"github.com/wealdtech/ethereal/util"
	ens "github.com/wealdtech/go-ens"
)

var tokenSweepFromAddress string
var tokenSweepToAddress string
var tokenSweepData string

// tokenSweepCmd represents the token sweep command
var tokenSweepCmd = &cobra.Command{
	Use:   "sweep",
	Short: "Sweep tokens to a given address",
	Long: `Sweep token from one address to another.  For example:

    ethereal token sweep --token=omg --from=0x5FfC014343cd971B7eb70732021E26C35B744cc4 --to=0x52f1A3027d3aA514F17E454C93ae1F79b3B12d5d --passphrase=secret

In quiet mode this will return 0 if the transfer transaction is successfully sent, otherwise 1.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Assert(!offline, quiet, "Offline mode not supported at current with this command")

		cli.Assert(tokenSweepFromAddress != "", quiet, "--from is required")
		fromAddress, err := ens.Resolve(client, tokenSweepFromAddress)
		cli.ErrCheck(err, quiet, fmt.Sprintf("Failed to resolve from address %s", tokenSweepFromAddress))

		cli.Assert(tokenSweepToAddress != "", quiet, "--to is required")
		toAddress, err := ens.Resolve(client, tokenSweepToAddress)
		cli.ErrCheck(err, quiet, fmt.Sprintf("Failed to resolve to address %s", tokenSweepToAddress))

		cli.Assert(tokenStr != "", quiet, "--token is required")
		token, err := tokenContract(tokenStr)
		cli.ErrCheck(err, quiet, "Failed to obtain token contract")

		// Obtain the balance of the address
		balance, err := token.BalanceOf(nil, fromAddress)
		cli.ErrCheck(err, quiet, "Failed to obtain balance of address from which to send funds")
		cli.Assert(balance.Cmp(big.NewInt(0)) > 0, quiet, "No balance")

		if verbose {
			symbol, err := token.Symbol(nil)
			if err == nil {
				decimals, err := token.Decimals(nil)
				if err == nil {
					fmt.Printf("Sweeping %s %s\n", util.TokenValueToString(balance, decimals, false), symbol)
				}
			}
		}

		opts, err := generateTxOpts(fromAddress)
		cli.ErrCheck(err, quiet, "Failed to generate transaction options")

		signedTx, err := token.Transfer(opts, toAddress, balance)
		cli.ErrCheck(err, quiet, "Failed to create transaction")

		if offline {
			if !quiet {
				buf := new(bytes.Buffer)
				signedTx.EncodeRLP(buf)
				fmt.Printf("0x%s\n", hex.EncodeToString(buf.Bytes()))
			}
		} else {
			logTransaction(signedTx, log.Fields{
				"group":   "token",
				"command": "sweep",
				"token":   tokenStr,
				"from":    fromAddress.Hex(),
				"to":      toAddress.Hex(),
				"amount":  balance.String(),
			})

			if quiet {
				os.Exit(0)
			}

			fmt.Println(signedTx.Hash().Hex())
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenSweepCmd)
	tokenFlags(tokenSweepCmd)
	tokenSweepCmd.Flags().StringVar(&tokenSweepFromAddress, "from", "", "Address from which to sweep tokens")
	tokenSweepCmd.Flags().StringVar(&tokenSweepToAddress, "to", "", "Address to which to sweep tokens")
	addTransactionFlags(tokenSweepCmd, "the address from which to sweep tokens")
}
