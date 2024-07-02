package etherscan

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

// GetSourceCode represents a result from the GetSourceCode API call
type GetSourceCode struct {
	// ContractName string `json:"ContractName"`
	ABI                  string `json:"ABI"`
	ContractName         string `json:"ContractName"`
	CompilerVersion      string `json:"CompilerVersion"`
	OptimizationUsed     string `json:"OptimizationUsed"`
	Runs                 string `json:"Runs"`
	ConstructorArguments string `json:"ConstructorArguments"`
	EVMVersion           string `json:"EVMVersion"`
	Library              string `json:"Library"`
	LicenseType          string `json:"LicenseType"`
	Proxy                string `json:"Proxy"`
	Implementation       string `json:"Implementation"`
	SwarmSource          string `json:"SwarmSource"`
}

type apiGetSourceCodeResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  []GetSourceCode `json:"result"`
}

type apiGetSourceCodeConfig struct {
	Module  string
	Action  string
	Address string
	APIKey  string
	BaseURL string
}

func (c apiGetSourceCodeConfig) toQueryParams() url.Values {
	query := url.Values{}
	query.Set("module", c.Module)
	query.Set("action", c.Action)
	query.Set("address", c.Address)
	query.Set("apikey", c.APIKey)
	return query
}

func (c apiGetSourceCodeConfig) getBaseURL() string {
	return c.BaseURL
}

func processGetSourceCode(responseBytes []byte) (*GetSourceCode, error) {
	var apiResponse apiGetSourceCodeResponse
	if err := json.Unmarshal(responseBytes, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse.Result[0], nil
}

// FetchSourceCode fetches the source code of a contract from the Etherscan API
func FetchSourceCode(
	chainID int64,
	apiKey string,
	address web3.Address,
) (*GetSourceCode, error) {
	baseURL, ok := chainIDBaseURLs[chainID]
	if !ok {
		return nil, fmt.Errorf("chain ID %d not supported", chainID)
	}

	config := &apiGetSourceCodeConfig{
		Module:  "contract",
		Action:  "getsourcecode",
		Address: address.String(),
		APIKey:  apiKey,
		BaseURL: baseURL,
	}

	responseBytes, err := callEtherscanAPI(config)
	if err != nil {
		return nil, fmt.Errorf("error calling Etherscan API: %w", err)
	}

	return processGetSourceCode(responseBytes)
}
