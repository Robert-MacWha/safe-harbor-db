package safeharbor

import (
	"SHDB/pkg/process"
	"context"
	"log"
	"strconv"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
			log.Println("Failed to parse chain ID", "error", err)
			continue
		}

		chainId := int64(chainIdInt)

		if _, exists := chainConfigs[chainId]; !exists {
			log.Println("Chain ID not found in chainConfigs", "chainID", chain.ID)
			continue
		}
		chainConfig := chainConfigs[chainId]

		rpcClient, err := rpc.Dial(chainConfig.RPCURL)
		if err != nil {
			log.Println("Failed to connect to the RPC client", "error", err)
			continue
		}

		ethClient, err := ethclient.Dial(chainConfig.RPCURL)
		if err != nil {
			log.Println("Failed to connect to the Ethereum client", "error", err)
			continue
		}

		latestBlock, err := ethClient.BlockNumber(context.Background())
		if err != nil {
			log.Println("Failed to get the latest block number", "error", err)
			continue
		}

		apiKey := chainConfig.APIKey
		chain.LastIndexedBlock = int(latestBlock)

		for i := range chain.Accounts {
			web3Address, err := web3.HexToAddress(chain.Accounts[i].Address)
			if err != nil {
				log.Println("Failed to parse address", "error", err)
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
					log.Println("Failed to get child contract addresses", "error", err)
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
