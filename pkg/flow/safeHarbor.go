package flow

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"SHDB/pkg/protocol"
	safeharbor "SHDB/pkg/safeHarbor"

	"cloud.google.com/go/firestore"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ChainConfig represents the configuration for a blockchain chain
type ChainConfig struct {
	APIKey string `json:"apiKey"`
	RPCURL string `json:"rpcURL"`
}

// ProcessSafeHarborAgreement processes the Safe Harbor agreement, uploads it to Firestore, and optionally sets the protocol reference.
func ProcessSafeHarborAgreement(
	chainConfigs map[int64]safeharbor.ChainConfig, // Chain configurations
	txHash web3.Hash, // Ethereum transaction hash
	safeHarborAddress web3.Address, // Safe Harbor contract address
	deployer web3.Address, // Deployer address
	chainId int, // Chain ID
	blockNumber *big.Int, // Block number
	protocolID string, // Firestore Protocol ID
	firestoreClient *firestore.Client, // Firestore client
	setProtocol bool, // Flag to set Safe Harbor agreement in protocol
) error {
	// Step 1: Connect to Ethereum node using the RPC URL from the chain config
	client, err := ethclient.Dial(chainConfigs[int64(chainId)].RPCURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Step 2: Fetch agreement details
	commonSafeHarborAddress := safeHarborAddress.ToCommon()
	agreement, _, err := safeharbor.FetchAgreementDetails(*blockNumber, txHash.String(), fmt.Sprintf("%d", chainId), deployer.String(), commonSafeHarborAddress, client)
	if err != nil {
		return fmt.Errorf("failed to fetch agreement details: %w", err)
	}

	log.Printf("Fetched agreement: %+v", agreement)

	// Step 3: Process the agreement details (e.g., with additional business logic)
	agreement, err = safeharbor.ProcessSafeHarborAgreement(agreement, chainConfigs)
	if err != nil {
		return fmt.Errorf("failed to process agreement details: %w", err)
	}

	// Step 4: Upload the agreement to Firestore
	protocolIDLower := strings.ToLower(protocolID)
	err = agreement.Upload(firestoreClient, protocolIDLower)
	if err != nil {
		return fmt.Errorf("failed to upload agreement to Firestore: %w", err)
	}

	// Step 5 (optional): Set Safe Harbor Agreement reference in protocol
	if setProtocol {
		err = safeharbor.SetProtocol(firestoreClient, protocolIDLower)
		if err != nil {
			return fmt.Errorf("failed to set Safe Harbor Agreement reference: %w", err)
		}
		err = protocol.SetSafeHarborAgreement(firestoreClient, protocolIDLower)
		if err != nil {
			return fmt.Errorf("failed to set protocol reference: %w", err)
		}
	}

	log.Printf("Successfully processed and uploaded protocol: %s", protocolID)
	return nil
}
