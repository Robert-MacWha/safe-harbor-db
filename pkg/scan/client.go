package scan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	ContractName(address string) (name string, err error)
}

type RateLimitedClient struct {
	apiKey  string
	chainID int
	client  *http.Client
}

const rateLimit = 300 * time.Millisecond

type EtherscanV2Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"` // Delay decoding
}

type SourceCodeResult struct {
	ContractName string `json:"ContractName"`
}

func NewRateLimitedClient(apiKey string, chainID int) *RateLimitedClient {
	return &RateLimitedClient{
		apiKey:  apiKey,
		chainID: chainID,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *RateLimitedClient) ContractName(address string) (string, error) {
	defer time.Sleep(rateLimit)

	url := fmt.Sprintf(
		"https://api.etherscan.io/v2/api?chainid=%d&module=contract&action=getsourcecode&address=%s&apikey=%s",
		c.chainID, address, c.apiKey,
	)

	// println(url)

	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var apiResp EtherscanV2Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}

	// If Status is "0", Result is a string (the error message)
	if apiResp.Status != "1" {
		var errMsg string
		_ = json.Unmarshal(apiResp.Result, &errMsg)
		return "", fmt.Errorf("etherscan error (%s): %s", apiResp.Message, errMsg)
	}

	// If Status is "1", Result is the expected array
	var results []SourceCodeResult
	if err := json.Unmarshal(apiResp.Result, &results); err != nil {
		return "", fmt.Errorf("failed to decode results: %w", err)
	}

	if len(results) == 0 || results[0].ContractName == "" {
		return "", fmt.Errorf("contract name not found/unverified for %s", address)
	}

	return results[0].ContractName, nil
}
