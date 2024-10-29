package web3

// BlockHead represents an Ethereum block head.
type BlockHead struct {
	Difficulty       BigInt  `json:"difficulty"`
	ExtraData        Bytes   `json:"extraData"`
	GasLimit         BigInt  `json:"gasLimit"`
	GasUsed          BigInt  `json:"gasUsed"`
	LogsBloom        string  `json:"logsBloom"`
	Miner            Address `json:"miner"`
	Nonce            BigInt  `json:"nonce"`
	Number           BigInt  `json:"number"`
	ParentHash       Hash    `json:"parentHash"`
	ReceiptRoot      Hash    `json:"receiptRoot"`
	Sha3Uncles       Hash    `json:"sha3Uncles"`
	StateRoot        Hash    `json:"stateRoot"`
	Timestamp        BigInt  `json:"timestamp"`
	TransactionsRoot Hash    `json:"transactionsRoot"`
}
