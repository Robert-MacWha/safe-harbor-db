package etherscan

import (
	"encoding/json"
	"fmt"
	"net/url"

	"SHDB/pkg/web3"
)

// Log represents a log (event) from the Etherscan API
type Log struct {
	Address          web3.Address `json:"address"`
	Topics           []web3.Hash  `json:"topics"`
	Data             web3.Bytes   `json:"data"`
	BlockNumber      web3.BigInt  `json:"blockNumber"`
	TimeStamp        web3.BigInt  `json:"timeStamp"`
	GasPrice         web3.BigInt  `json:"gasPrice"`
	GasUsed          web3.BigInt  `json:"gasUsed"`
	LogIndex         web3.BigInt  `json:"logIndex"`
	TransactionHash  web3.Hash    `json:"transactionHash"`
	TransactionIndex web3.BigInt  `json:"transactionIndex"`
}

type apiLogResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []Log  `json:"result"`
}

// Logs-specific API configuration
type logsAPIConfig struct {
	blockApiConfig
	Topic0 string
}

func (c logsAPIConfig) toQueryParams() url.Values {
	query := c.blockApiConfig.toQueryParams()

	// Since Logs API uses from and to block instead of start and end block
	query.Del("startblock")
	query.Set("fromBlock", fmt.Sprintf("%d", c.StartBlock))
	query.Del("endblock")
	query.Set("toBlock", "latest")

	query.Set("topic0", c.Topic0)
	return query
}

func (c logsAPIConfig) getStartBlock() int {
	return c.StartBlock
}

func (c *logsAPIConfig) setStartBlock(block int) {
	c.StartBlock = block
}

func (c logsAPIConfig) getBaseURL() string {
	return c.BaseURL
}

func processLogsResponse(responseBytes []byte) ([]Log, int, error) {
	var apiResponse apiLogResponse
	if err := json.Unmarshal(responseBytes, &apiResponse); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Result) == 0 {
		return nil, 0, nil // Return early if no logs are found
	}

	// Assuming the logs are sorted by block number
	lastLog := apiResponse.Result[len(apiResponse.Result)-1]

	return apiResponse.Result, int(lastLog.BlockNumber.Int.Int64()), nil
}

// FetchLogs fetches logs (events) from the Etherscan API
func FetchLogs(
	chainID int64,
	apiKey string,
	address web3.Address,
	topic0 web3.Hash,
	startBlock int,
) ([]Log, error) {
	baseURL, ok := chainIDBaseURLs[chainID]
	if !ok {
		return nil, fmt.Errorf("chain ID %d not supported", chainID)
	}

	logsConfig := &logsAPIConfig{
		blockApiConfig: blockApiConfig{
			Module:     "logs",
			Action:     "getLogs",
			Address:    address.String(),
			StartBlock: startBlock,
			EndBlock:   "latest",
			Page:       1,
			Offset:     maximumOffset,
			Sort:       "asc",
			APIKey:     apiKey,
			BaseURL:    baseURL,
		},
		Topic0: topic0.String(),
	}
	return fetchAllData[Log](logsConfig, processLogsResponse)
}
