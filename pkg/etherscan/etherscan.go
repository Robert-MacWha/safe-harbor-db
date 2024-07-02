package etherscan

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var chainIDBaseURLs = map[int64]string{
	1:        "https://api.etherscan.io/api",
	11155111: "https://api-sepolia.etherscan.io/api",
}

// baseConfig is an interface that's required for basic API calls
type basicConfig interface {
	toQueryParams() url.Values
	getBaseURL() string
}

// blockBasedConfig is an interface that's required for API calls that query based on block
type blockBasedConfig interface {
	toQueryParams() url.Values
	getStartBlock() int
	setStartBlock(block int)
	getBaseURL() string
}

type blockApiConfig struct {
	Module     string
	Action     string
	Address    string
	StartBlock int
	EndBlock   string
	Page       int
	Offset     int
	Sort       string
	APIKey     string
	BaseURL    string
}

// ResultProcessor is a type constraint that defines the behavior of processing functions
type ResultProcessor[T any] func([]byte) ([]T, int, error)

const maximumOffset = 10000

// General API configuration
func (c blockApiConfig) toQueryParams() url.Values {
	query := url.Values{}
	query.Set("module", c.Module)
	query.Set("action", c.Action)
	query.Set("address", c.Address)
	query.Set("startblock", strconv.Itoa(c.StartBlock))
	query.Set("endblock", c.EndBlock)
	query.Set("page", strconv.Itoa(c.Page))
	query.Set("offset", strconv.Itoa(c.Offset))
	query.Set("sort", c.Sort)
	query.Set("apikey", c.APIKey)
	return query
}

func (c blockApiConfig) getStartBlock() int {
	return c.StartBlock
}

func (c *blockApiConfig) setStartBlock(block int) {
	c.StartBlock = block
}

func (c blockApiConfig) getBaseURL() string {
	return c.BaseURL
}

// API calling functions

func callEtherscanAPI(config basicConfig) ([]byte, error) {
	baseURL := config.getBaseURL()
	query := config.toQueryParams()

	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, query.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return responseBytes, nil
}

// Result processing functions

// Main function to fetch all data

// fetchAllData uses generics to fetch and process data based on the given ResultProcessor
func fetchAllData[T any](
	config blockBasedConfig,
	processResults ResultProcessor[T],
) ([]T, error) {
	var allResults []T
	fromBlock := config.getStartBlock()

	for {

		config.setStartBlock(fromBlock)
		responseBytes, err := callEtherscanAPI(config)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		results, newLastBlock, err := processResults(responseBytes)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, results...)
		if newLastBlock == 0 || len(results) < maximumOffset {
			break
		}
		fromBlock = newLastBlock + 1

		// 200 miliseconds because etherscan max rate is 5 requests per second
		time.Sleep(200 * time.Millisecond) //nolint
	}

	return allResults, nil
}
