package deduab

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Deployed struct {
	Address     string `json:"address"`
	BlockNumber int    `json:"block_number"`
	Name        string `json:"name"`
}

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c Client) GetDeployed(address string, blockNumber int, txIndex int, vmStepStart int, limit int) ([]Deployed, error) {
	url := fmt.Sprintf(
		"https://api.dedaub.com/api/account/ethereum/%s/deployed?block_number=%d&tx_index=%d&vm_step_start=%d&limit=%d",
		address, blockNumber, txIndex, vmStepStart, limit,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deployed []Deployed
	if err := json.NewDecoder(resp.Body).Decode(&deployed); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	return deployed, nil
}
