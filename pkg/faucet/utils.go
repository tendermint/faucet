package faucet

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (f *Faucet) getChainID() (string, error) {
	output, err := f.executeCli([]string{"status"})
	if err != nil {
		return "", err
	}

	var status resultStatus
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &status); err != nil {
		return "", err
	}

	return status.NodeInfo.Network, nil
}

func getTxAmount(tx interface{}) (uint64, error) {
	if v, ok := tx.(map[string]interface{})["value"]; ok {
		if v, ok := v.(map[string]interface{})["msg"]; ok {
			if v := v.([]interface{}); len(v) > 0 {
				if v, ok := v[0].(map[string]interface{})["value"]; ok {
					if v, ok := v.(map[string]interface{})["amount"]; ok {
						if v := v.([]interface{}); len(v) > 0 {
							if v, ok := v[0].(map[string]interface{})["amount"]; ok {
								return strconv.ParseUint(v.(string), 10, 64)
							}
						}
					}
				}
			}
		}
	}
	return 0, fmt.Errorf("could not get transaction amount")
}
