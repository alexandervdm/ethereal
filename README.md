# ethereal

[![Tag](https://img.shields.io/github/tag/wealdtech/ethereal.svg)](https://github.com/wealdtech/ethereal/releases/)
[![License](https://img.shields.io/github/license/wealdtech/ethereal.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/wealdtech/ethereal?status.svg)](https://godoc.org/github.com/wealdtech/ethereal)
[![Travis CI](https://img.shields.io/travis/wealdtech/ethereal.svg)](https://travis-ci.org/wealdtech/ethereal)

A command-line tool for managing common tasks in Ethereum.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install


`ethereal` is a standard Go program which can be installed with:

```sh
go get github.com/wealdtech/ethereal
```

## Usage

Ethereal contains a large number of features that are useful for day-to-day interactions with the Ethereum blockchain.

### Access to local wallets

Ethereal works with Geth, MIST and Parity wallets in the standard locations.  A simple way to check the addresses that can be seen by Ethereal is to run `ethereal account list` which will list all accounts that Ethereal can see.  If you expect an address to show up and it doesn't then please raise an issue with the relevant details.

If you use Parity and want to import a private key or a wallet from another system please see https://github.com/paritytech/parity/wiki/Backing-up-&-Restoring#restoring-options

If you use Geth and want to import a private key or a wallet from another system please see https://github.com/ethereum/go-ethereum/wiki/Managing-your-accounts

When accessing local wallets a `--passphrase` option is required to unlock the account.  Note that this is not shown in the examples

Alternatively you can use a private key directly with the `--privatekey` option, although be aware that this can leave your private key in command history.

### Access to Ethereum networks

Ethereal supports all main Ethereum networks  It auto-detects the network by querying the connected node for the network ID.  The connection should be geth-compatible, so either geth itself or parity with the `--geth` flag to enable geth compatibility mode.  The connection could be a local node or a network service such as Infura.

TODO explain default connection and how to connect locally.

### Configuration file
TODO

### Quiet, Verbose and Debug
TODO

### Transaction
TODO Options
TODO CLI output
TODO Logging

### ENS

Ethereal fully supports ENS.  Wherever an address is seen in the examples below an ENS name can be used instead.

### `account` commands

Account commands focus on information about local accounts, generally those used by Geth and Parity but also those from hardware devices.

#### `checksum`

`ethereal account checksum` generates or verifies the [EIP-55](https://eips.ethereum.org/EIPS/eip-55) checksum for a provided account address.  With the `--check` flag it checks if the supplied address is correctly checksummed, otherwise it generates a correctly checksummed version of the supplied address.  For example:

```sh
$ ethereal account checksum --address=0x7e5f4552091a69125d5dfcb7b8c2659029395bdf --check
Checksum is incorrect
$ ethereal account checksum --address=0x7e5f4552091a69125d5dfcb7b8c2659029395bdf
0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
$ ethereal account checksum --address=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf --check
Checksum is correct
```

#### `keys`

`ethereal account keys` shows the private key, public key and Ethereum address for a given account or private key.  For example:

```sh
$ ethereal account keys --privatekey=0x0000000000000000000000000000000000000000000000000000000000000001
Private key:            0x00000000000000000000000000000001
Public key:             0x0479be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8
Ethereum address:       0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
```

#### `list`

`ethereal account list` shows the Ethereum addresses of known accounts on the local computer.  For example:

```sh
$ ethereal account list
0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69
0x1efF47bc3a10a45D4B230B5d10E37751FE6AA718
0x003F53E95e293D08dc34C69ABcAbF5b577E50Cf5
```

With the `--verbose` flag this will provide the location of the keystore, current Ether funds and next nonce.  For example:

```sh
$ ethereal account list --verbose
Location:       keystore:///home/ethereum/.ethereum/keystore/UTC--2019-03-12T10-12-47.585144239Z--7e5f4552091a69125d5dfcb7b8c2659029395bdf
Address:        0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
Balance:        0
Next nonce:     243
...
```

#### `nonce`

`ethereal account nonce` shows the next nonce of an Ethereum address.  For example:

```sh
$ ethereal account nonce --address=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
243
```

### `block` commands

Block commands focus on information about specific blocks.

#### `info`

`ethereal block info` provides information about a block.  For example:

```sh
$ ethereal block info --block=5188504
Number:                 5188504
Hash:                   0x01262b8549472c95714993135f9aa1cb09685bd33076541522e3db0481f63fe7
Block time:             1552386092 (2019-03-12 10:21:32 +0000 GMT)
Gas limit:              8000000
Gas used:               7983831 (99.80%)
Uncles:                 0
Transactions:           66
```

With the `--verbose` flag this will provide information about the miner.  For example:

```sh
$ ethereal block info --block=5188504 --verbose
Number:                 5188504
Hash:                   0x01262b8549472c95714993135f9aa1cb09685bd33076541522e3db0481f63fe7
Block time:             1552386092 (2019-03-12 10:21:32 +0000 GMT)
Mined by:               0x6212Dd88f890FefE0Af24D1404d96aDF488e4E3B
Extra:                  ؃geth�go1.10.4�linux
Difficulty:             2232415661
Gas limit:              8000000
Gas used:               7983831 (99.80%)
Transactions:           66
```

Note that `--block=latest` will provide information on the latest mined block.

#### `overview`

`ethereal block overview` provides high-level statistics about the last few mined blocks.  For example:

```sh
$ ethereal block overview
Block    Gas used/Gas limit     Block time              Gap     Coinbase
5188514   7882080/  8000000     19/03/12 10:24:00               0xCd626bc764E1d553e0D75a42f5c4156B91a63F23
5188513   7994026/  8000000     19/03/12 10:23:58       2s      0xCd626bc764E1d553e0D75a42f5c4156B91a63F23
5188512   7964126/  8000029     19/03/12 10:23:48       10s     0x635B4764D1939DfAcD3a8014726159abC277BecC
5188511   7981224/  8000000     19/03/12 10:23:43       5s      0x6212Dd88f890FefE0Af24D1404d96aDF488e4E3B
5188510   7958922/  8000000     19/03/12 10:23:35       8s      0x6212Dd88f890FefE0Af24D1404d96aDF488e4E3B
```

With the `--verbose` flag this will provide column headers.  For example:

```sh
$ ethereal block overview --verbose
Block    Gas used/Gas limit     Block time              Gap     Coinbase
5188514   7882080/  8000000     19/03/12 10:24:00               0xCd626bc764E1d553e0D75a42f5c4156B91a63F23
5188513   7994026/  8000000     19/03/12 10:23:58       2s      0xCd626bc764E1d553e0D75a42f5c4156B91a63F23
5188512   7964126/  8000029     19/03/12 10:23:48       10s     0x635B4764D1939DfAcD3a8014726159abC277BecC
5188511   7981224/  8000000     19/03/12 10:23:43       5s      0x6212Dd88f890FefE0Af24D1404d96aDF488e4E3B
5188510   7958922/  8000000     19/03/12 10:23:35       8s      0x6212Dd88f890FefE0Af24D1404d96aDF488e4E3B
```

The number of blocks displayed in the overview can be altered using the `--blocks` parameter.

### `contract` commands

Contract commands focus on deploying and interacting with Ethereum smart contracts.

The examples of the commands below use the following contract at `SampleContract.sol`:

```solidity
pragma solidity ^0.5.0;
  
contract SampleContract {
    uint256 private value;

    constructor(uint256 _value) public {
        value = _value;
    }

    function getValue() public view returns (uint256) {
        return value;
    }

    function setValue(uint256 _value) public {
        value = _value;
    }
}
```

which is compiled using the command line:

```sh
$ solc --optimize --combined-json=bin,abi SampleContract.sol >SampleContract.json
```

Note that best results the names of the files should be the same as the name of the contract (ignoring the suffix), as per the example above.

#### `call`

`ethereal contract call` calls a contract function locally on the connected node.  For example:

```sh
$ ethereal contract call --contract=0x3c24F71e826D3762f5145f6a27d41545A7dfc8cF --json=SampleContract.json --call='getValue()' --from=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
5
```

#### `deploy`

`ethereal contract deploy` deploys a contract to the Ethereum blockchain.

The binary to deploy can be supplied in two different ways.  The simplest is to compile the contract using the `--combined-json=abi,json` option of `solc` to provide a JSON file containing both the binary data and the contract's ABI and deploy using that.  For example:

```sh
$ ethereal contract deploy --json=SampleContract.json --constructor='constructor(5)' --from=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
```

Alternatively the binary data and constructor arguments can be supplied directly on the command-line.  For example:

```sh
$ BIN=`solc --optimize --bin SampleContract.sol | egrep -A 2 SampleContract.sol:SampleContract | tail -1`
$ CONSTRUCTORARGS=`0000000000000000000000000000000000000000000000000000000000000005`
$ ethereal contract deploy --data="${BIN}${CONSTRUCTORARGS}" --from=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
```

#### `send`

`ethereal contract send` sends a contract transaction to the Ethereum blockchain.  For example:

```sh
$ ethereal contract send --contract=0x3c24F71e826D3762f5145f6a27d41545A7dfc8cF --json=SampleContract.json --call='setValue(6)' --from=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
```

#### `storage`

`ethereal contract storage` accesses contract storage directly.  Key values depend on the value stored; for more details see [this article](https://medium.com/aigang-network/how-to-read-ethereum-contract-storage-44252c8af925).

```sh
$ ethereal contract storage --contract=0x3c24F71e826D3762f5145f6a27d41545A7dfc8cF --key=0x00
0x0000000000000000000000000000000000000000000000000000000000000006
```

### `dns` commands

DNS commands focus on interacting with the [EthDNS](https://medium.com/@jgm.orinoco/ethdns-an-ethereum-backend-for-the-domain-name-system-d52dabd904b3) system to allow DNS records to be stored on Ethereum.

Getting and setting DNS records works on the basis of a DNS resource record set.  A resource record set is defined by the tuple (domain,name,resource record type) for example (ehdns.xyz,www.ethdns.xyz,A) would return all 'A' (address) records help for www.ethdns.xyz in the domain ethdns.xyz.

#### `clear`

`ethereal dns clear` clears all resource records for a DNS zone.

#### `get`

`ethereal dns get` obtains a single resource record set for the (domain,name,resource record type) tuple.  For example:

```sh
$ ethereal dns get --domain=ethdns.xyz --name=www --resource=CNAME
www.ethdns.xyz. 21600   IN      CNAME   ethdns.xyz.
```

Resource record sets on the root domain can be fetched by omitting the `name` argument. For example:

```sh
$ ethereal dns get --domain=ethdns.xyz --resource=NS
ethdns.xyz.     43200   IN      NS      ns1.ethdns.xyz.
ethdns.xyz.     43200   IN      NS      ns2.ethdns.xyz.
```

#### `set`

`ethereal dns set` sets a single resource record set for the (domain,name,resource record type) tuple.  For example:

```sh
$ ethereal dns set --domain=ethdns.xyz --name=www --resource=CNAME --value="ethdns.xyz."
```

Resource record sets with multiple values can be supplied by separating them with "&&".  For example:

```sh
$ ethereal dns set --domain=ethdns.xyz --resource=NS --value="ns1.ethdns.xyz&&ns2.ethdns.xyz"
```

### `ens` commands

ENS commands focus on interacting with the [Ethereum Name Service](https://ens.domains/) contracts that address resources using human-readable names.

#### `address clear`

`ethereal ens address clear` removes an address associated with an ENS domain.  For example:

```sh
$ ethereal ens address clear --domain=mydomain.eth
```

#### `address get`

`ethereal ens address get` gets the address associated with an ENS domain.  For example:

```sh
$ ethereal ens address get --domain=mydomain.eth
0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
```

#### `address set`

`ethereal ens address set` sets the address associated with an ENS domain.  For example:

```sh
$ ethereal ens address set --domain=mydomain.eth --address=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
```

#### `contenthash clear`

`ethereal ens contenthash clear`

#### `contenthash get`

`ethereal ens contenthash get`

#### `contenthash set`

`ethereal ens contenthash set`

#### `info`

`ethereal ens info` obtains various information about a domain.  For example:

```sh
$ ethereal ens info --domain=mydomain.eth
Domain owner is 0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
Owned since 2017-05-22 15:40:25 +0100 BST
Locked value is 0.01 Ether
Highest bid was 0.01 Ether
Deed owner is 0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
Domain owner is 0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
Previous owner was 0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
Resolver is 0x5FfC014343cd971B7eb70732021E26C35B744cc4
mydomain.eth resolves to 0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69
0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69 does not resolve to a domain
```

With the `--verbose` flag this will provide more information about the domain.  For example:

```sh
$ ethereal ens info --domain=mydomain.eth --verbose
Normalised domain is mydomain.eth
Top-level domain is eth
Name hash is 0x2ae48b85f1a7764f6b87d040180e633b395ca2471ceadf46a0cd8e6cb91d7ed7
Domain owner is 0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
Owned since 2017-05-22 15:40:25 +0100 BST
Locked value is 0.01 Ether
Highest bid was 0.01 Ether
Deed owner is 0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
Domain owner is 0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
Previous owner was 0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
Resolver is 0x5FfC014343cd971B7eb70732021E26C35B744cc4
mydomain.eth resolves to 0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69
0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69 does not resolve to a domain
```

#### `name clear`

`ethereal ens name clear` clears the ENS reverse resolution name of an address.  For example:

```sh
$ ethereal ens name clear --address=0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69
```

#### `name get`

`ethereal ens name get` gets the ENS reverse resolution name of an address.  For example:

```sh
$ ethereal ens name get --address=0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69
mydomain.eth
```

#### `name set`

`ethereal ens name set` sets the ENS reverse resolution name of an address.  For example:

```sh
$ ethereal ens name set --address=0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69 --name=mydomain.eth
```

#### `owner get`

`ethereal ens owner get` gets the owner of the domain.  For example:

```sh
$ ethereal ens owner get --domain=mydomain.eth
0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
```

#### `owner set`

`ethereal ens owner set` sets the owner of the domain.  For example:

```sh
$ ethereal ens owner set --domain=mydomain.eth --owner=0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
```

#### `resolver clear`

`ethereal ens resolver clear` clears the resolver contract for the domain.  For example:

```sh
$ ethereal ens resolver clear --domain=mydomain.eth
```

#### `resolver get`

`ethereal ens resolver get` gets the address of the resolver contract for the domain.  For example:

```sh
$ ethereal ens resolver get --domain=mydomain.eth
0x5FfC014343cd971B7eb70732021E26C35B744cc4
```

#### `resolver set`

`ethereal ens resolver set` sets the resolver contract for the domain.  If the standard public resolver (found at `resolver.eth`) is required then just the domain is required to set it.  For example:

```sh
$ ethereal ens resolver set --domain=mydomain.eth
```

If a non-standard resolver is required it can be supplied with the `--resolver` argument.  For example:

```sh
$ ethereal ens resolver set --domain=mydomain.eth --resolver=0x4d9b7D10e3a42E81659A90fDbaB51Bf19DD9bba7
```

#### `subdomain create`

`ethereal ens subdomain create` creates a subdomain of an existing ENS domain.  For example:

```sh
$ ethereal ens subdomain create --domain=mydomain.eth --subdomain=mysub
```

The subdomain will be owned by the domain owner.

#### `transfer`

`ethereal ens transfer` transfers ownership of a name and its associated deed to another address.  For example:

```sh
$ ethereal ens transfer --domain=mydomain.eth --newowner=0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
```

### `ether` commands

Ether commands focus on information about and movement of Ether.

#### `balance`

`ethereal ether balance` provides the Ether balance of an address.  For example:

```sh
$ ethereal ether balance --address=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf
5189.916425903288395771 Ether
```

If required the balance can be supplied in Wei with the `--wei` option.  For example:

```sh
$ ethereal ether balance --address=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf --wei
5189916425903288395771
```

#### `sweep`

`ethereal ether sweep` sweeps all Ether from one address to another, leaving 0 behind.  For example:

```sh
$ ethereal ether sweep --from=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf --to=0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF
```

#### `transfer`

`ethereal ether transfer` transfers a set amount of Ether from one address to another.  For example:

```sh
$ ethereal ether transfer --from=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf --to=0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF --amount="1.2 Ether"
```

### `gas` commands

#### `price`

`ethereal gas price` calaculates a gas price from historical information that should allow a transaction to be included within a certain number of blocks.  For example:

```sh
TODO
```

### `network` commands

#### `blocktime`

`ethereal network blocktime` calculates the average blocktime over a number of blocks.  For example:

```sh
$ ethereal network blocktime --blocks=30
11.76s
```

Instead of using `--blocks` it is possible to specify the time over which to calculate the blocktime.  For example:

```sh
$ ethereal network blocktime --time=24h
11.52s
```

With the `--verbose` flag this will provide more information about the start and end block for the calculation.  For example:

```sh
$ ethereal network blocktime --time=2h --verbose
Block 7370164 mined at 2019-03-14 23:48:53 +0000 GMT
Block 7369614 mined at 2019-03-14 21:49:28 +0000 GMT
13.02s
```
#### `gps`

`ethereal network gps` provides a gas-per-second metric for the Ethereum network over a number of blocks.  For example:

```sh
$ ethereal network gps --blocks=20
339950
```

With the `--verbose` flag this will provide more information about each block.  For example:

```sh
$ ethereal network gps --blocks=5 --verbose
Block 7370156 used 7407093 gas in 2 seconds
Block 7370155 used 6831267 gas in 9 seconds
Block 7370154 used 7787059 gas in 2 seconds
Block 7370153 used 219751 gas in 4 seconds
Block 7370152 used 1113156 gas in 6 seconds
1015579
```

#### `tps`

`ethereal network tps` provides a transactions-per-second metric for the Ethereum network over a number of blocks.  For example:

```sh
$ ethereal network tps --blocks=20
4.72
```
With the `--verbose` flag this will provide more information about each block.  For example:

```sh
$ ethereal network tps --blocks=5 --verbose
Block 7370160 processed 5 transactions in 3 seconds
Block 7370159 processed 30 transactions in 5 seconds
Block 7370158 processed 40 transactions in 8 seconds
Block 7370157 processed 48 transactions in 9 seconds
Block 7370156 processed 31 transactions in 2 seconds
```

5.70


### `registry` commands

Ether commands focus on use of the ERC-1820 registry.

### `signature` commands

Signature commands focus on generation and verification of signatures within Ethereum.

### `token` commands

Token commands focus on information and management of ERC-20 and ERC-777 tokens.

### `transaction` commands

Transaction commands focus on information and management of Ethereum transactions.

### `version`

`ethereal version` provides the current version of Ethereal.  For example:

```sh
$ ethereal version
2.0.889
```

## Examples

Note that for most of the examples below additional information can be obtained by adding the `--verbose` flag to the command line.

### Increase the gas price for transaction
You have submitted a transaction to the network but it's taking a long time to process because the gas price is too low.

```
ethereal transaction up --transaction=0x5219b09d629158c2759035c97b11b604f57d0c733515738aaae0d2dafb41ab98 --gasprice=20GWei
```
where `transaction` is the hash of the pending transactions, `gasprice` is the price you want to set for gas, and `passphrase` is the passphrase for the account that sent the transaction.

### Cancel a transaction
You have submitted a transaction to the network by mistake and want to cancel it.
```
ethereal transaction cancel --transaction=0x5219b09d629158c2759035c97b11b604f57d0c733515738aaae0d2dafb41ab98 --passphrase=secret
```
where `transaction` is the hash of the pending transactions and `passphrase` is the passphrase for the account that sent the transaction.

### Sweep Ether
You want to transfer all Ether in one account to another.
```
ethereal ether sweep --from=0x5FfC014343cd971B7eb70732021E26C35B744cc4 --to=0x52f1A3027d3aA514F17E454C93ae1F79b3B12d5d --passphrase=secret
```
where `from` is the address from which the Ether will be transferred, `to` is the address to which the Ether will be transferred, and `passphrase` is the passphrase for the `from` account.

### Transfer a token
You want to transfer an ERC-20 token to another account.
```
ethereal token transfer --token=omg --from=0x5FfC014343cd971B7eb70732021E26C35B744cc4 --to=0x52f1A3027d3aA514F17E454C93ae1F79b3B12d5d --amount=10.2 --passphrase=secret
```
where `token` is the token to transfer, `from` is the address from which the token will be transferred, `to` is the address to which the token will be transferred, `amount` is the amount of the token to transfer, `gasprice` is the price you want to set for gas, and `passphrase` is the passphrase for the `from` account.

*Please note that before using a token name such as 'omg' you should confirm that the contract address matches the expected contract address by using `ethereal info --token=omg` or similar.*

## Obtain information about a transaction

```
ethereal transaction info --transaction=0x5097a149236b675a5807ea78c657b64c71da48789476828fede68126769b24be
```

### Deploy a contract
Ethereal can deploy a contract in various ways, but the easiest is to use `solc` to create a file containing both the ABI and bytecode, and use that.  For example:

```
solc --optimize --combined-json=abi,bin MyContract.sol >MyContract.json
ethereal contract deploy --json=MyContract.json --from=0xdd8686E0Ea24bc74ea6a4688926b5397D167930E --passphrase=secret --constructor='constructor("hello")'
```

If the contract does not have a constructor the `--constructor` argument can be omitted.

### Call a contract
You want to obtain information directly from a contract using its ABI, for example call the `balanceOf()` call of an ERC-20 token.
```
ethereal contract call --abi='./erc20.abi' --contract=0xd26114cd6EE289AccF82350c8d8487fedB8A0C07 --from=0x5FfC014343cd971B7eb70732021E26C35B744cc4 --call='balanceOf(0x5FfC014343cd971B7eb70732021E26C35B744cc4)'
```
where `abi` is the path to the contract's ABI, `contract` is the address of the contract to call, and `call` is the ABI method to call.

You can also use the JSON file as referenced in 'Deploy a contract' above, for example:

```
ethereal contract call --json=MyContract.json --from=0xdd8686E0Ea24bc74ea6a4688926b5397D167930E --call="getString()"
```
## Maintainers

Jim McDonald: [@mcdee](https://github.com/mcdee).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/wealdtech/ethereal/issues).

## License

[Apache-2.0](LICENSE) © 2017-2019 Weald Technology Trading Ltd

