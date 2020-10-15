package faucet

import (
	"encoding/json"
	"fmt"
)

func (f *Faucet) Send(recipient string) error {
	_, err := f.executeCli([]string{"tx", "send", f.keyName, recipient,
		fmt.Sprintf("%d%s", f.creditAmount, f.denom), "--yes", "--chain-id", f.chainID},
		f.keyringPassword, f.keyringPassword, f.keyringPassword)
	return err
}

func (f *Faucet) GetTotalSent(recipient string) (uint64, error) {
	args := []string{
		"query", "txs", "--events",
		fmt.Sprintf("transfer.sender=%s&transfer.recipient=%s", f.faucetAddress, recipient),
		"--page", "1",
		"--limit", "1000",
		"--trust-node",
	}

	output, err := f.executeCli(args)
	if err != nil {
		return 0, err
	}

	var txs txsQueryResult
	if err := json.Unmarshal([]byte(output), &txs); err != nil {
		return 0, err
	}

	var total uint64
	for _, tx := range txs.Txs {
		amount, err := getTxAmount(tx.Tx)
		if err != nil {
			return 0, err
		}
		total += amount
	}
	return total, nil
}
