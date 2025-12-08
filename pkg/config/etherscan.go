package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type ChainCfg struct {
	RpcUrl  string `json:"rpcURL"`
	ScanUrl string
	ScanKey string `json:"scanKey"`
}

var ScanUrls = map[int]string{
	1:     "https://api.etherscan.io/api?",
	10:    "https://api-optimistic.etherscan.io/api?",
	56:    "https://api.bscscan.com/api?",
	100:   "https://api.gnosisscan.io/api?",
	137:   "https://api.polygonscan.com/api?",
	324:   "https://block-explorer-api.mainnet.zksync.io/api?",
	1101:  "https://api-zkevm.polygonscan.com/api?",
	8453:  "https://api.basescan.org/api?",
	17000: "https://api-holesky.etherscan.io/api?",
	42161: "https://api.arbiscan.io/api?",
	43114: "https://api.routescan.io/v2/network/mainnet/evm/43114/etherscan/api?",
}

var ScanSites = map[string]int{
	"etherscan.io":  1,
	"sonicscan.org": 146,
	"basescan.org":  8453,
}

// Loads and unmarshals the CHAIN_CONFIG environment variable
func LoadChainCfg() (map[int]ChainCfg, error) {
	chainCfg := make(map[int]ChainCfg)
	for chain, scanUrl := range ScanUrls {
		//? Stored as environment secrets
		rpcUrlEnv := fmt.Sprintf("RPC_URL_%d", chain)
		scanKeyEnv := "SCAN_KEY"

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

var addressRegex = regexp.MustCompile(`0x[a-fA-F0-9]{40}`)

func ParseAddressFromURL(raw string) (address string, chainId int, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", 0, fmt.Errorf("url.Parse: %w", err)
	}

	host := strings.ToLower(u.Hostname())
	address = addressRegex.FindString(u.Path)
	if address == "" {
		return "", 0, fmt.Errorf("no address found in URL: %s", raw)
	}

	chainID, ok := ScanSites[host]
	if !ok {
		return "", 0, fmt.Errorf("no chain ID found for host: %s", host)
	}

	return address, chainID, nil
}
