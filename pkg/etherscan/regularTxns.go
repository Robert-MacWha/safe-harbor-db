package etherscan

import (
	"encoding/json"
	"fmt"

	"SHDB/pkg/web3"
)

// RegularTransaction represents a regular transaction from the Etherscan API
type RegularTransaction struct {
	BlockNumber       web3.BigInt   `json:"blockNumber"`
	TimeStamp         web3.BigInt   `json:"timeStamp"`
	Hash              web3.Hash     `json:"hash"`
	Nonce             web3.BigInt   `json:"nonce"`
	BlockHash         web3.Hash     `json:"blockHash"`
	TransactionIndex  web3.BigInt   `json:"transactionIndex"`
	From              *web3.Address `json:"from,omitempty"`
	To                *web3.Address `json:"to,omitempty"`
	Value             web3.BigInt   `json:"value"`
	Gas               web3.BigInt   `json:"gas"`
	GasPrice          web3.BigInt   `json:"gasPrice"`
	IsError           string        `json:"isError"`
	TxreceiptStatus   string        `json:"txreceipt_status"`
	Input             web3.Bytes    `json:"input"`
	ContractAddress   *web3.Address `json:"contractAddress,omitempty"`
	CumulativeGasUsed web3.BigInt   `json:"cumulativeGasUsed"`
	GasUsed           web3.BigInt   `json:"gasUsed"`
	Confirmations     web3.BigInt   `json:"confirmations"`
	MethodID          string        `json:"methodId"`
	FunctionName      string        `json:"functionName"`
}

type apiRegularTransactionResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Result  []RegularTransaction `json:"result"`
}

func processRegularTransactions(responseBytes []byte) ([]RegularTransaction, int, error) {
	var apiResponse apiRegularTransactionResponse
	if err := json.Unmarshal(responseBytes, &apiResponse); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Result) == 0 {
		return nil, 0, nil
	}

	lastTxn := apiResponse.Result[len(apiResponse.Result)-1]

	return apiResponse.Result, int(lastTxn.BlockNumber.Int.Int64()), nil
}

// FetchRegularTransactions fetches regular transactions from the Etherscan API
func FetchRegularTransactions(
	chainID int64,
	apiKey string,
	address web3.Address,
	startBlock int,
) ([]RegularTransaction, error) {
	baseURL, ok := chainIDBaseURLs[chainID]
	if !ok {
		return nil, fmt.Errorf("chain ID %d not supported", chainID)
	}

	regularConfig := &blockApiConfig{
		Module:     "account",
		Action:     "txlist",
		Address:    address.String(),
		StartBlock: startBlock,
		EndBlock:   "latest",
		Page:       1,
		Offset:     maximumOffset,
		Sort:       "asc",
		APIKey:     apiKey,
		BaseURL:    baseURL,
	}
	return fetchAllData[RegularTransaction](regularConfig, processRegularTransactions)
}
