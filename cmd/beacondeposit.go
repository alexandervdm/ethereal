// Copyright © 2017-2020 Weald Technology Trading
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wealdtech/ethereal/cli"
	"github.com/wealdtech/ethereal/util"
	"github.com/wealdtech/ethereal/util/contracts"
	ens "github.com/wealdtech/go-ens/v3"
	string2eth "github.com/wealdtech/go-string2eth"
)

var beaconDepositData string
var beaconDepositFrom string
var beaconDepositForce bool
var beaconDepositContractAddress string

type ethdoDepositData struct {
	Account               string `json:"account"`
	PublicKey             string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Signature             string `json:"signature"`
	DepositDataRoot       string `json:"deposit_data_root"`
	Value                 uint64 `json:"value"`
	Version               uint64 `json:"version"`
}

type beaconDepositContract struct {
	name       string
	chainID    *big.Int
	address    []byte
	minVersion uint64
	maxVersion uint64
}

var beaconDepositContractWhitelists = []*beaconDepositContract{
	{
		name:       "Prysm topaz",
		chainID:    big.NewInt(5),
		address:    util.MustDecodeHexString("0x5ca1e00004366ac85f492887aaab12d0e6418876"),
		minVersion: 1,
		maxVersion: 1,
	},
	{
		name:       "Prysm onyx",
		chainID:    big.NewInt(5),
		address:    util.MustDecodeHexString("0x0f0f0fc0530007361933eab5db97d09acdd6c1c8"),
		minVersion: 1,
		maxVersion: 1,
	},
}

// beaconDepositCmd represents the beacon deposit command
var beaconDepositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit Ether to the beacon contract.",
	Long: `Deposit Ether to the Ethereum 2 beacon contract, either creating or supplementing a validator.  For example:

    ethereal becon deposit --data=/home/me/depositdata.json --from=0x.... --passphrase="my secret passphrase"

Note that at current this deposits Ether to the Prysm test deposit contract on Goerli.  Other networks and deposit contracts are not supported.

The depositdata.json file can be generated by ethdo.  The data can be an array of deposits, in which case they will be processed sequentially.

The keystore for the account that owns the name must be local (i.e. listed with 'get accounts list') and unlockable with the supplied passphrase.

This will return an exit status of 0 if the transaction is successfully submitted (and mined if --wait is supplied), 1 if the transaction is not successfully submitted, and 2 if the transaction is successfully submitted but not mined within the supplied time limit.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Assert(!offline, quiet, "Offline mode not supported at current with this command")

		cli.Assert(chainID.Cmp(big.NewInt(5)) == 0, quiet, "This command is only supported on the Goerli network")

		cli.Assert(beaconDepositData != "", quiet, "--data is required")
		var err error
		var data []byte
		// Data could be JSON or a path to JSON
		if strings.HasPrefix(beaconDepositData, "{") {
			// Looks like JSON
			data = []byte("[" + beaconDepositData + "]")
		} else if strings.HasPrefix(beaconDepositData, "[") {
			// Looks like JSON array
			data = []byte(beaconDepositData)
		} else {
			// Assume it's a path to JSON
			data, err = ioutil.ReadFile(beaconDepositData)
			cli.ErrCheck(err, quiet, "Failed to find deposit data file")
			if data[0] == '{' {
				data = []byte("[" + string(data) + "]")
			}
		}
		var depositData []ethdoDepositData
		err = json.Unmarshal(data, &depositData)
		cli.ErrCheck(err, quiet, "Data is not valid JSON")
		cli.Assert(len(depositData) > 0, quiet, "No deposit data supplied")
		minVersion := depositData[0].Version
		maxVersion := depositData[0].Version
		for i := range depositData {
			cli.Assert(depositData[i].PublicKey != "", quiet, fmt.Sprintf("No public key for deposit %d", i))
			cli.Assert(depositData[i].DepositDataRoot != "", quiet, fmt.Sprintf("No data root for deposit %d", i))
			cli.Assert(depositData[i].Signature != "", quiet, fmt.Sprintf("No signature for deposit %d", i))
			cli.Assert(depositData[i].WithdrawalCredentials != "", quiet, fmt.Sprintf("No withdrawal credentials for deposit %d", i))
			cli.Assert(depositData[i].Value >= 1000000000, quiet, fmt.Sprintf("Value too small for deposit %d", i))
			if depositData[i].Version > maxVersion {
				maxVersion = depositData[i].Version
			}
			if depositData[i].Version < minVersion {
				minVersion = depositData[i].Version
			}
		}

		cli.Assert(beaconDepositFrom != "", quiet, "--from is required")
		fromAddress, err := ens.Resolve(client, beaconDepositFrom)
		cli.ErrCheck(err, quiet, "Failed to obtain address for --from")

		// Fetch the address of the contract.
		cli.Assert(beaconDepositContractAddress != "", quiet, "--address is required")
		depositContractAddress, err := ens.Resolve(client, beaconDepositContractAddress)
		cli.ErrCheck(err, quiet, "Failed to obtain address of deposit contract")
		// Ensure this contract is whitelisted.
		contractName := ""
		for _, whitelistEntry := range beaconDepositContractWhitelists {
			if chainID.Cmp(whitelistEntry.chainID) == 0 &&
				bytes.Equal(depositContractAddress.Bytes(), whitelistEntry.address) {
				cli.Assert(beaconDepositForce || minVersion >= whitelistEntry.minVersion, quiet, `Data generated by ethdo is old and possibly inaccurate.  This means you need to upgrade your version of ethdo (or you are sending your deposit to the wrong contract or network); please do so by visiting https://github.com/wealdtech/ethdo and following the installation instructions there.  Once you have done this please regenerate your deposit data and try again.`)
				cli.Assert(beaconDepositForce || maxVersion <= whitelistEntry.maxVersion, quiet, `Data generated by ethdo is newer than supported.  This means you need to upgrade your version of ethereal (or you are sending your deposit to the wrong contract or network); please do so by visiting https://github.com/wealdtech/ethereal and following the installation instructions there.  Once you have done this please try again.`)
				contractName = whitelistEntry.name
				break
			}
		}
		cli.Assert(beaconDepositForce || contractName != "", quiet, `Deposit contract address is unknown.  This means you are either running an old version of ethereal, or are attempting to send to the wrong network or a custom contract.  You should confirm that you are on the latest version of Ethereal by comparing the output of running "ethereal version" with the release information at https://github.com/wealdtech/ethereal/releases and upgrading where appropriate.

If you are *completely sure* you know what you are doing, you can use the --force option to carry out this transaction.  Otherwise, please seek support to ensure you do not lose your Ether.`)
		outputIf(verbose, fmt.Sprintf("Deposit contract is %s", contractName))

		contract, err := contracts.NewEth2Deposit(depositContractAddress, client)
		cli.ErrCheck(err, quiet, "Failed to obtain deposit contract")

		for _, deposit := range depositData {
			opts, err := generateTxOpts(fromAddress)
			cli.ErrCheck(err, quiet, "Failed to generate deposit options")
			// Need to override the value with the info from the JSON
			opts.Value = new(big.Int).Mul(new(big.Int).SetUint64(deposit.Value), big.NewInt(1000000000))

			// Need to set gas limit because it moves around a fair bit with the merkle tree calculations.
			opts.GasLimit = 600000

			pubKey, err := hex.DecodeString(deposit.PublicKey)
			cli.ErrCheck(err, quiet, "Failed to parse deposit public key")
			withdrawalCredentials, err := hex.DecodeString(deposit.WithdrawalCredentials)
			cli.ErrCheck(err, quiet, "Failed to parse deposit withdrawal credentials")
			signature, err := hex.DecodeString(deposit.Signature)
			cli.ErrCheck(err, quiet, "Failed to parse deposit signature")
			dataRootTmp, err := hex.DecodeString(deposit.DepositDataRoot)
			cli.ErrCheck(err, quiet, "Failed to parse deposit data root")
			var dataRoot [32]byte
			copy(dataRoot[:], dataRootTmp)

			// TODO recalculate signature to ensure correcteness (needs a pure Go BLS implementation).

			// TODO check Ethereum 2 node to see if there is already a deposit for this validator public key (needs an Ethereum 2 node).

			outputIf(verbose, fmt.Sprintf("Creating %s deposit for %s", string2eth.WeiToString(big.NewInt(int64(deposit.Value)), true), deposit.Account))

			nextNonce(fromAddress)
			signedTx, err := contract.Deposit(opts, pubKey, withdrawalCredentials, signature, dataRoot)
			cli.ErrCheck(err, quiet, "Failed to send deposit")

			handleSubmittedTransaction(signedTx, log.Fields{
				"group":                        "beacon",
				"command":                      "deposit",
				"depositPublicKey":             pubKey,
				"depositWithdrawalCredentials": withdrawalCredentials,
				"depositSignature":             signature,
				"depositDataRoot":              dataRoot,
			}, false)
		}
	},
}

func init() {
	beaconCmd.AddCommand(beaconDepositCmd)
	beaconFlags(beaconDepositCmd)
	beaconDepositCmd.Flags().StringVar(&beaconDepositData, "data", "", "The data for the deposit, provided by ethdo or a similar command")
	beaconDepositCmd.Flags().StringVar(&beaconDepositFrom, "from", "", "The account from which to send the deposit")
	beaconDepositCmd.Flags().BoolVar(&beaconDepositForce, "force", false, "Force send data to non-whitelisted contracts (not recommended)")
	beaconDepositCmd.Flags().StringVar(&beaconDepositContractAddress, "address", "eth2bridge.eth", "The address to which to send the deposit")
	addTransactionFlags(beaconDepositCmd, "passphrase for the account that owns the account")
}
