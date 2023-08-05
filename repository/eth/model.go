package eth

type request struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      string `json:"id"`
}

type trans struct {
	From        string `json:"from"`
	To          string `json:"to"`
	BlockNumber string `json:"blockNumber"`
	Ind         string `json:"transactionIndex"`
	Value       string `json:"value"`
}

type block struct {
	Transactions []trans `json:"transactions"`
}
