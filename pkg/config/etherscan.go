package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type ChainCfg struct {
	RpcUrl  string `json:"rpcURL"`
	ScanUrl string
	ScanKey string `json:"scanKey"`
}

var scanUrls = map[int]string{
	1:        "https://api.etherscan.io/api?",
	137:      "https://api.polygonscan.com/api?",
	17000:    "https://api-holesky.etherscan.io/api?",
	11155111: "https://api-sepolia.etherscan.io/api?",
	42161:    "https://api.arbiscan.io/api?",
	8453:     "https://api.basescan.org/api?",
	10:       "https://api-optimistic.etherscan.io/api?",
	56:       "https://api.bscscan.com/api?",
	1101:     "https://api-zkevm.polygonscan.com/api?",
	43114:    "https://api.routescan.io/v2/network/mainnet/evm/43114/etherscan/api?",
	100:      "https://api.gnosisscan.io/api?",
}

// Loads and unmarshals the CHAIN_CONFIG environment variable
func LoadChainCfg() (map[int]ChainCfg, error) {
	chainCfgStr := os.Getenv("CHAIN_CONFIG")
	if chainCfgStr == "" {
		return nil, fmt.Errorf("missing CHAIN_CONFIG env")
	}

	chainCfgStr = strings.Trim(chainCfgStr, "'")

	var chainCfg map[int]ChainCfg
	err := json.Unmarshal([]byte(chainCfgStr), &chainCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chain config: %w", err)
	}

	for chain, cfg := range chainCfg {
		if scanUrl, exists := scanUrls[chain]; !exists {
			slog.Warn("Missing scan for chain, removing from config", "chain", chain)
			delete(chainCfg, chain)
		} else {
			cfg.ScanUrl = scanUrl
			chainCfg[chain] = cfg
		}
	}

	return chainCfg, nil
}
