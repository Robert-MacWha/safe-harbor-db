package safeharbor

import (
	"SHDB/pkg/process"
	"strconv"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/log"
)

// ChainConfig holds information about each chain's configuration
type ChainConfig struct {
	APIKey string `json:"apiKey"`
	RPCURL string `json:"rpcURL"`
}

func ProcessSafeHarborAgreement(
	safeHarbor *SafeHarborAgreement,
	chainConfigs map[int64]ChainConfig,
	startBlock int,
) (*SafeHarborAgreement, error) {
	// Process the SafeHarborAgreement
	for _, chain := range safeHarbor.Chains {
		chainIdInt, err := strconv.Atoi(chain.ID)
		if err != nil {
			log.Error("Failed to parse chain ID", "error", err)
			continue
		}

		chainId := int64(chainIdInt)

		if _, exists := chainConfigs[chainId]; !exists {
			log.Error("Chain ID not found in chainConfigs", "chainID", chain.ID)
			continue
		}
		chainConfig := chainConfigs[chainId]

		rpcClient, err := rpc.Dial(chainConfig.RPCURL)
		if err != nil {
			log.Error("Failed to connect to the RPC client", "error", err)
			continue
		}

		var latestBlock int
		err = rpcClient.Call(&latestBlock, "eth_blockNumber")
		if err != nil {
			log.Error("Failed to get the latest block number", "error", err)
			continue
		}

		apiKey := chainConfig.APIKey

		chain.LastIndexedBlock = latestBlock

		for i := range chain.Accounts {
			web3Address, err := web3.HexToAddress(chain.Accounts[i].Address)
			if err != nil {
				log.Error("Failed to parse address", "error", err)
				return nil, err
			}

			chain.Accounts[i].Name = process.GetNameOrEmpty(*web3Address, chainId, apiKey)

			if chain.Accounts[i].ChildContractScope == ChildContractScopeAll || chain.Accounts[i].ChildContractScope == ChildContractScopeExistingOnly {
				children, err := process.GetAllSubContractAddresses(
					rpcClient,
					chainId,
					apiKey,
					*web3Address,
					startBlock,
				)
				if err != nil {
					log.Error("Failed to get child contract addresses", "error", err)
					return nil, err
				}

				childAccounts := []ChildAccount{}
				for _, child := range children {
					childAccount := ChildAccount{
						Name:    process.GetNameOrEmpty(child, chainId, apiKey),
						Address: child.String(),
					}
					childAccounts = append(childAccounts, childAccount)
				}
				chain.Accounts[i].Children = childAccounts
			}
		}
	}

	return safeHarbor, nil
}
