# cosmos-faucet

This faucet was developed to use [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) executable binaries only. Hence
it is compatible with pratically any blockchain built with [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) even if
different types of keys are used (as [ethermint](https://github.com/cosmos/ethermint) for example).

The `master` version supports `launchpad` only. Far `stargate` please use
[stargate](https://github.com/allinbits/cosmos-faucet/tree/stargate) branch.

## Installation

You can build the faucet with:

```bash
$ make build
```

The executable binary will be avaialable in the `./build/` directory. To install it to `$GOPATH/bin`, use:

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
| cli-name         	| CLI_NAME         	| the name of the cli executable                                   	| gaiacli (gaiad for stargate) 	|
| denom            	| DENOM            	| the denomination of the coin to be distributed by the faucet     	| uatom                        	|
| credit-amount    	| CREDIT_AMOUNT    	| the amount to credit in each request                             	| 10000000                     	|
| max-credit       	| MAX_CREDIT       	| the maximum credit per account                                   	| 100000000                    	|

### Example

Start the faucet with:

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

### Request tokens

You can request tokens by sending a `POST` request to any path on the server, with a key address in a `JSON`:

```bash
$ curl -X POST -d '{"address": "cosmos1rlumjuvfjss4hq0vykrk0pwt7ws62vt3dj7cj2"}'
{"status": "ok"}
```

---
**NOTE**

In order to make the API of this faucet compatible with others, it is possible to include other fields in the `JSON`
sent to the server - it will only read `address` field though. Additionally, the post can be made to any path in the
server. This is compatible for example with [cosmjs faucet](https://github.com/cosmos/cosmjs/tree/master/packages/faucet).

---