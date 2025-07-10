package config

import (
	"fmt"
	"log/slog"
	"os"
)

type ChainCfg struct {
	RpcUrl  string `json:"rpcURL"`
	ScanUrl string
	ScanKey string `json:"scanKey"`
}

var scanUrls = map[int]string{
	1:        "https://api.etherscan.io/api?",
	10:       "https://api-optimistic.etherscan.io/api?",
	56:       "https://api.bscscan.com/api?",
	100:      "https://api.gnosisscan.io/api?",
	137:      "https://api.polygonscan.com/api?",
	1101:     "https://api-zkevm.polygonscan.com/api?",
	8453:     "https://api.basescan.org/api?",
	17000:    "https://api-holesky.etherscan.io/api?",
	42161:    "https://api.arbiscan.io/api?",
	43114:    "https://api.routescan.io/v2/network/mainnet/evm/43114/etherscan/api?",
	11155111: "https://api-sepolia.etherscan.io/api?",
}

// Loads and unmarshals the CHAIN_CONFIG environment variable
func LoadChainCfg() (map[int]ChainCfg, error) {
	chainCfg := make(map[int]ChainCfg)
	for chain, scanUrl := range scanUrls {
		//? Stored as environment secrets
		rpcUrlEnv := fmt.Sprintf("RPC_URL_%d", chain)
		scanKeyEnv := fmt.Sprintf("SCAN_KEY_%d", chain)

		rpcUrl := os.Getenv(rpcUrlEnv)
		if rpcUrl == "" {
			slog.Warn("Missing RPC URL for chain, removing from config", "chain", chain)
			continue
		}

		scanKey := os.Getenv(scanKeyEnv)
		if scanKey == "" {
			slog.Warn("Missing scan key for chain, removing from config", "chain", chain)
			continue
		}

		chainCfg[chain] = ChainCfg{
			RpcUrl:  rpcUrl,
			ScanUrl: scanUrl,
			ScanKey: scanKey,
		}
	}

	return chainCfg, nil
}
