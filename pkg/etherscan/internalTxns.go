package etherscan

import (
	"encoding/json"
	"fmt"

	"SHDB/pkg/web3"
)

// InternalTransaction represents an internal transaction from the Etherscan API
type InternalTransaction struct {
	BlockNumber     web3.BigInt   `json:"blockNumber"`
	TimeStamp       web3.BigInt   `json:"timeStamp"`
	Hash            web3.Hash     `json:"hash"`
	From            *web3.Address `json:"from"`
	To              *web3.Address `json:"to,omitempty"`
	Value           web3.BigInt   `json:"value"`
	ContractAddress *web3.Address `json:"contractAddress"`
	Input           web3.Bytes    `json:"input"`
	Type            string        `json:"type"`
	Gas             web3.BigInt   `json:"gas"`
	GasUsed         web3.BigInt   `json:"gasUsed"`
	TraceID         string        `json:"traceId"`
	IsError         string        `json:"isError"`
	ErrCode         string        `json:"errCode"`
}

type apiInternalTransactionResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Result  []InternalTransaction `json:"result"`
}

func processInternalTransactions(
	responseBytes []byte,
) ([]InternalTransaction, int, error) {
	var apiResponse apiInternalTransactionResponse
	if err := json.Unmarshal(responseBytes, &apiResponse); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Result) == 0 {
		return nil, 0, nil
	}

	lastInternalTxn := apiResponse.Result[len(apiResponse.Result)-1]

	return apiResponse.Result, int(lastInternalTxn.BlockNumber.Int.Int64()), nil
}

// FetchInternalTransactions fetches internal transactions from the Etherscan API
func FetchInternalTransactions(
	chainID int64,
	apiKey string,
	address web3.Address,
	startBlock int,
) ([]InternalTransaction, error) {
	baseURL, ok := chainIDBaseURLs[chainID]
	if !ok {
		return nil, fmt.Errorf("chain ID %d not supported", chainID)
	}

	internalConfig := &blockApiConfig{
		Module:     "account",
		Action:     "txlistinternal",
		Address:    address.String(),
		StartBlock: startBlock,
		EndBlock:   "latest",
		Page:       1,
		Offset:     maximumOffset,
		Sort:       "asc",
		APIKey:     apiKey,
		BaseURL:    baseURL,
	}
	return fetchAllData[InternalTransaction](internalConfig, processInternalTransactions)
}
