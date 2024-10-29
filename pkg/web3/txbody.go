package web3

// TxBody represents a transaction body.
type TxBody struct {
	From     Address  `json:"from"`
	To       *Address `json:"to,omitempty"`
	Hash     string   `json:"hash"`
	Gas      BigInt   `json:"gas"`
	GasPrice BigInt   `json:"gasPrice"`
	Input    Bytes    `json:"input"`
	Value    BigInt   `json:"value"`
	// If Block is nil, assume the "latest" block.
	Block   *BigInt `json:"blockNumber,omitempty"`
	Nonce   BigInt  `json:"nonce"`
	ChainID BigInt  `json:"chainId"`
	// a.k.a. GasFeeCap - EIP-1559
	MaxFeePerGas *BigInt `json:"maxFeePerGas,omitempty"`
	// a.k.a. GasTipCap - EIP-1559
	MaxPriorityFeePerGas *BigInt `json:"maxPriorityFeePerGas,omitempty"`
}
