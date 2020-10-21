package faucet

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Faucet struct {
	appCli          string
	chainID         string
	keyringPassword string
	keyName         string
	faucetAddress   string
	keyMnemonic     string
	denom           string
	creditAmount    uint64
	maxCredit       uint64
	cdc             *codec.Codec
}

func NewFaucet(opts ...Option) (*Faucet, error) {
	options := &defaultOptions
	for _, opt := range opts {
		opt(options)
	}

	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	cdc.RegisterConcrete(auth.StdTx{}, stdTxCodecType, nil)
	cdc.RegisterConcrete(bank.MsgSend{}, msgSendCodecType, nil)

	chainID, err := getChainID(cdc, options.AppCli)
	if err != nil {
		return nil, err
	}

	e := Faucet{
		appCli:          options.AppCli,
		keyringPassword: options.KeyringPassword,
		keyName:         options.KeyName,
		keyMnemonic:     options.KeyMnemonic,
		denom:           options.Denom,
		creditAmount:    options.CreditAmount,
		maxCredit:       options.MaxCredit,
		chainID:         chainID,
		cdc:             cdc,
	}
	return &e, e.loadKey()
}

func getChainID(cdc *codec.Codec, executable string) (string, error) {
	output, err := cmdexec(executable, []string{"status"})
	if err != nil {
		return "", err
	}

	var status ctypes.ResultStatus
	if err := cdc.UnmarshalJSON([]byte(output), &status); err != nil {
		return "", err
	}

	return status.NodeInfo.Network, nil
}
