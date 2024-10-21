package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetProtocol fetches and returns detailed information about a specific protocol.
// It returns a Protocol struct populated with data from the API.
func GetProtocol(protocolSlug string) (*Protocol, error) {
	// Fetch detailed information about the protocol.
	protocol, twitterHandle, err := fetchProtocolDetail(protocolSlug)
	if err != nil {
		return nil, err
	}

	// Fetch the list of all protocols to match the category using Twitter handle.
	category, err := fetchProtocolCategory(twitterHandle)
	if err != nil {
		return nil, err
	}

	// Assign the matched category to the protocol.
	protocol.Category = category

	return protocol, nil
}

// fetchProtocolDetail retrieves the protocol details from the external API
// for a given protocol slug. It returns a Protocol struct and Twitter handle.
func fetchProtocolDetail(protocolSlug string) (*Protocol, string, error) {
	var protocol Protocol

	// API endpoint for fetching protocol details.
	url := fmt.Sprintf("https://api.llama.fi/protocol/%s", protocolSlug)
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Temporary struct for unmarshalling the JSON response.
	var detail struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		Logo    string `json:"logo"`
		Twitter string `json:"twitter"`
		TVL     []struct {
			TotalLiquidityUSD float64 `json:"totalLiquidityUSD"`
		} `json:"tvl"`
	}

	// Unmarshal the JSON response into the temporary struct.
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, "", err
	}

	// Map the data from the temporary struct to the Protocol struct.
	protocol.Name = detail.Name
	protocol.Slug = protocolSlug
	protocol.Website = detail.URL
	protocol.Icon = detail.Logo
	protocol.TVL = getLastTVL(detail.TVL)
	protocol.ContactDetails = "" // Set appropriately if you have contact details.

	// Return the protocol and Twitter handle for category matching.
	return &protocol, detail.Twitter, nil
}

// fetchProtocolCategory retrieves the category of the protocol using the Twitter handle.
func fetchProtocolCategory(twitterHandle string) (string, error) {
	// API endpoint for fetching all protocols.
	resp, err := http.Get("https://api.llama.fi/protocols")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Temporary struct for unmarshalling the JSON response.
	var protocols []struct {
		Twitter  string `json:"twitter"`
		Category string `json:"category"`
	}

	// Unmarshal the JSON response into the temporary struct.
	if err := json.Unmarshal(body, &protocols); err != nil {
		return "", err
	}

	// Match the Twitter handle to find the category.
	for _, p := range protocols {
		if p.Twitter == twitterHandle {
			return p.Category, nil
		}
	}

	// Return an empty string if no category is found.
	return "", nil
}

// getLastTVL extracts the most recent TVL value from a slice of TVL data.
func getLastTVL(tvlData []struct {
	TotalLiquidityUSD float64 `json:"totalLiquidityUSD"`
}) float64 {
	if len(tvlData) > 0 {
		return tvlData[len(tvlData)-1].TotalLiquidityUSD
	}
	return 0
}
