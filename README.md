# cosmos-faucet

A faucet that uses [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) executable binaries only.

The main purpose of this `faucet` is to avoid using RPC or API endpoints, and use the CLI binary instead, more
specifically, the commands:

```bash
$ {app}d tx bank send
```

and:

```bash
$ {app}d query txs
```

Since the faucet only uses the CLI binary, it is compatible with practically any blockchain built with
[cosmos-sdk](https://github.com/cosmos/cosmos-sdk) even if different types of keys are used (such as in
[ethermint](https://github.com/cosmos/ethermint) for example).


## Installation

You can build the faucet with:

```bash
$ make build
```

The executable binary will be avaialable in the `./build/` directory. To install it to `$GOPATH/bin`, use instead:

```bash
$ make install
```

## Usage

### Configuration

You can configure the faucet either using command line flags or environment variables. The following table
shows the available configuration options and respective defaults:

| flag             	| env              	| description                                                      	| default                      	|
|------------------	|------------------	|------------------------------------------------------------------	|------------------------------	|
| log-level        	| LOG_LEVEL        	| the log level to be used (trace, debug, info, warn or error)     	| info                         	|
| port             	| PORT             	| the port in which the server will be listening for HTTP requests 	| 8000                         	|
| key-name         	| KEY_NAME         	| the name of the key to be used by the faucet                     	| faucet                       	|
| mnemonic         	| MNEMONIC         	| a mnemonic to restore an existing key (this is optional)         	|                              	|
| keyring-password 	| KEYRING_PASSWORD 	| the password for accessing the keys keyring                      	|                              	|
| cli-name         	| CLI_NAME         	| the name of the cli executable                                   	| gaiad 	                    |
| denom            	| DENOM            	| the denomination of the coin to be distributed by the faucet     	| uatom                        	|
| credit-amount    	| CREDIT_AMOUNT    	| the amount to credit in each request                             	| 10000000                     	|
| max-credit       	| MAX_CREDIT       	| the maximum credit per account                                   	| 100000000                    	|
| sdk-version      	| SDK_VERSION      	| version of sdk (launchpad or stargate)                            | stargate                    	|
| node            	| NODE            	| the address of the node that will handle the requests             |                    	        |
| keyring-backend   | KEYRING_BACKEND   | keyring backend                                                   |                               |

### [gaia](https://github.com/cosmos/gaia) example

This faucet options default to work with [gaia](https://github.com/cosmos/gaia). So you can start the faucet with just:

```bash
$ faucet --keyring-password 12345678
INFO[0000] listening on :8000
```

or, with environment variables:

```bash
$ export KEYRING_PASSWORD=12345678
$ faucet
INFO[0000] listening on :8000
```

### [ethermint](https://github.com/cosmos/ethermint) example

Start the faucet with:

```bash
$ faucet --cli-name ethermintcli --denom ueth --keyring-password 12345678 --sdk-version launchpad
INFO[0000] listening on :8000
```

or, with environment variables:

```bash
$ export CLI_NAME=ethermintcli
$ export SDK_VERSION=launchpad
$ export DENOM=ueth
$ export KEYRING_PASSWORD=12345678
$ faucet
INFO[0000] listening on :8000
```

### [wasmd](https://github.com/CosmWasm/wasmd) example

Start the faucet with:

```bash
$ faucet --cli-name wasmcli --denom ucosm --keyring-password 12345678
INFO[0000] listening on :8000
```

or, with environment variables:

```bash
$ export CLI_NAME=wasmcli
$ export DENOM=ucosm
$ export KEYRING_PASSWORD=12345678
$ faucet
INFO[0000] listening on :8000
```

### Request tokens

You can request tokens by sending a `POST` request to any path on the server, with a key address in a `JSON`:

```bash
$ curl -X POST -d '{"address": "cosmos1kd63kkhtswlh5vcx5nd26fjmr9av74yd4sf8ve"}' http://localhost:8000
{"transfers":[{"coin":"10000000uatom","status":"ok"}]}
```