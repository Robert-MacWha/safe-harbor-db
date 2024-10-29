package web3

// BlockBody represents an Ethereum block body.
type BlockBody struct {
	BlockHead

	Transactions []TxBody `json:"transactions"`
}
