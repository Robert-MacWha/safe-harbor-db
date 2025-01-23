package defiliama

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Protocol struct {
	Name     string  `firestore:"name"`
	Slug     string  `firestore:"slug"`
	Website  string  `firestore:"website" json:"url"`
	Icon     string  `firestore:"icon" json:"logo"`
	TVL      float64 `firestore:"tvl"`
	Category string  `firestore:"category"`
}

type tvl struct {
	TotalLiquidityUSD float64 `json:"totalLiquidityUSD"`
}

type protocolDetail struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Logo    string `json:"logo"`
	Twitter string `json:"twitter"`
	TVL     []tvl  `json:"tvl"`
}

type protocolCategory struct {
	Twitter  string `json:"twitter"`
	Category string `json:"category"`
}

func GetTvl(slug string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.llama.fi/protocol/%s", slug))
	if err != nil {
		return 0, fmt.Errorf("http.Get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var details protocolDetail
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return 0, fmt.Errorf("json.Decode: %w", err)
	}

	return getLastTVL(details.TVL), nil
}

func GetProtocol(slug string) (Protocol, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.llama.fi/protocol/%s", slug))
	if err != nil {
		return Protocol{}, fmt.Errorf("http.Get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Protocol{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var details protocolDetail
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return Protocol{}, fmt.Errorf("json.Decode: %w", err)
	}

	return Protocol{
		Name:     details.Name,
		Slug:     slug,
		Website:  details.URL,
		Icon:     details.Logo,
		TVL:      getLastTVL(details.TVL),
		Category: getProtocolCategory(details.Twitter),
	}, nil
}

func getProtocolCategory(twitter string) string {
	resp, err := http.Get("https://api.llama.fi/protocols")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var protocols []protocolCategory
	err = json.NewDecoder(resp.Body).Decode(&protocols)
	if err != nil {
		return ""
	}

	category := ""
	for _, p := range protocols {
		if p.Twitter == twitter {
			if category != "" {
				return "Multiple"
			}

			category = p.Category
		}
	}

	return category
}

func getLastTVL(tvlData []tvl) float64 {
	if len(tvlData) > 0 {
		return tvlData[len(tvlData)-1].TotalLiquidityUSD
	}
	return 0
}
