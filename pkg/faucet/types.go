package faucet

type resultStatus struct {
	NodeInfo nodeInfo `json:"node_info"`
}

type nodeInfo struct {
	Network string `json:"network"`
}

type txsQueryResult struct {
	Txs []tx `json:"txs"`
}

type tx struct {
	Tx interface{} `json:"tx"`
}

type httpRequest struct {
	Address string `json:"address"`
}

type httpResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
