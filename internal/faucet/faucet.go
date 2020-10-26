package faucet

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/std"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	stdTxCodecType   = "cosmos-sdk/StdTx"
	msgSendCodecType = "cosmos-sdk/MsgSend"
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
	cdc             *codec.ProtoCodec
}

func NewFaucet(opts ...Option) (*Faucet, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	chainID, err := getChainID(options.AppCli)
	if err != nil {
		return nil, err
	}

	reg := types.NewInterfaceRegistry()
	std.RegisterInterfaces(reg)
	bank.RegisterInterfaces(reg)

	e := Faucet{
		appCli:          options.AppCli,
		keyringPassword: options.KeyringPassword,
		keyName:         options.KeyName,
		keyMnemonic:     options.KeyMnemonic,
		denom:           options.Denom,
		creditAmount:    options.CreditAmount,
		maxCredit:       options.MaxCredit,
		chainID:         chainID,
		cdc:             codec.NewProtoCodec(reg),
	}
	return &e, e.loadKey()
}

func getChainID(executable string) (string, error) {
	output, err := cmdexec(executable, []string{"status"})
	if err != nil {
		return "", err
	}

	cdc := codec.NewLegacyAmino()
	codec.RegisterEvidences(cdc)
	cryptocodec.RegisterCrypto(cdc)

	var status ctypes.ResultStatus
	if err := cdc.UnmarshalJSON([]byte(output), &status); err != nil {
		return "", err
	}

	return status.NodeInfo.Network, nil
}
