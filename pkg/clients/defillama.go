// Package clients provides access to protocol information via external APIs.
package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// protocolDetail contains detailed information about a specific protocol including
// the name, URL, total value locked (TVL) data, Twitter handle, and logo.
type protocolDetail struct {
	Name    string    `json:"name"`
	URL     string    `json:"url"`
	TVL     []tvlData `json:"tvl"`
	Twitter string    `json:"twitter"`
	Logo    string    `json:"logo"`
}

// tvlData represents a time-specific entry of TVL (Total Value Locked)
// providing the USD value locked at a particular time.
type tvlData struct {
	TotalLiquidityUSD float64 `json:"totalLiquidityUSD"`
}

// protocolSummary is used for unmarshalling the list of all protocols
// retrieved from an API, storing each protocol's Twitter handle and category.
type protocolSummary struct {
	Twitter  string `json:"twitter"`
	Category string `json:"category"`
}

// GetProtocolInfo fetches and returns detailed information about a specific protocol
// including its name, URL, latest TVL, category based on Twitter handle, Twitter handle,
// and logo URL. It returns an error if any issues occur during the API calls or data processing.
func GetProtocolInfo(protocolName string) (name string, url string, tvl float64, category string, twitter string, logo string, err error) {
	protocolDetail, err := fetchProtocolDetail(protocolName)
	if err != nil {
		return "", "", 0, "", "", "", err
	}

	protocols, err := fetchAllProtocols()
	if err != nil {
		return "", "", 0, "", "", "", err
	}

	category = matchTwitter(protocols, protocolDetail.Twitter)

	return protocolDetail.Name, protocolDetail.URL, getLastTVL(protocolDetail.TVL), category, protocolDetail.Twitter, protocolDetail.Logo, nil
}

// fetchProtocolDetail retrieves the protocol details from the external API
// for a given protocol name. It unmarshals the JSON response into a ProtocolDetail struct.
func fetchProtocolDetail(protocolName string) (protocolDetail, error) {
	var detail protocolDetail
	resp, err := http.Get(fmt.Sprintf("https://api.llama.fi/protocol/%s", protocolName))
	if err != nil {
		return detail, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return detail, err
	}
	err = json.Unmarshal(body, &detail)
	return detail, err
}

// fetchAllProtocols retrieves a list of all protocols from the external API
// and unmarshals the JSON response into a slice of ProtocolSummary.
func fetchAllProtocols() ([]protocolSummary, error) {
	var summaries []protocolSummary
	resp, err := http.Get("https://api.llama.fi/protocols")
	if err != nil {
		return summaries, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return summaries, err
	}
	err = json.Unmarshal(body, &summaries)
	return summaries, err
}

// getLastTVL extracts the most recent TVL value from a slice of TVLData.
func getLastTVL(tvl []tvlData) float64 {
	if len(tvl) > 0 {
		return tvl[len(tvl)-1].TotalLiquidityUSD
	}
	return 0
}

// matchTwitter searches through a slice of ProtocolSummary for a matching Twitter handle
// to find the corresponding category of the protocol.
func matchTwitter(protocols []protocolSummary, twitter string) string {
	for _, p := range protocols {
		if p.Twitter == twitter {
			return p.Category
		}
	}
	return ""
}
