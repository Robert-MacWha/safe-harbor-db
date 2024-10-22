package main

import (
	"SHDB/pkg/etherscan"
	"SHDB/pkg/flow"
	safeharbor "SHDB/pkg/safeHarbor"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/option"

	// Import necessary types
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

type RegistryConfig struct {
	ChainID         int64  `json:"chainId"`
	RegistryAddress string `json:"registryAddress"`
}

type EtherscanEvent struct {
	TransactionHash   web3.Hash
	BlockNumber       int64
	TimeStamp         int64
	Deployer          web3.Address
	SafeHarborAddress web3.Address
}

func main() {
	app := &cli.App{
		Name:  "safe-harbor-monitor",
		Usage: "Monitors Safe Harbor registry addresses for new agreements",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to the registry config JSON file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "chainConfigs",
				Aliases:  []string{"cc"},
				Usage:    "Path to the chain configurations JSON file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "creds",
				Aliases:  []string{"f"},
				Usage:    "Path to Firestore credentials file",
				Required: true,
			},
		},
		Action: Run,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}

func Run(c *cli.Context) error {
	registryConfigs, err := loadRegistryConfigs(c.String("config"))
	if err != nil {
		return fmt.Errorf("failed to load registry configs: %w", err)
	}

	chainConfigs, err := loadChainConfigs(c.String("chainConfigs"))
	if err != nil {
		return fmt.Errorf("failed to load chain configs: %w", err)
	}

	credsPath := c.String("creds")

	// Create Firestore client
	firestoreClient, err := newFirestoreClient(credsPath)
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}
	defer firestoreClient.Close()

	var wg sync.WaitGroup

	for _, registryConfig := range registryConfigs {
		wg.Add(1)
		go func(rc RegistryConfig) {
			defer wg.Done()
			err := monitorRegistry(rc, chainConfigs, firestoreClient)
			if err != nil {
				log.Printf("Error monitoring registry for ChainID %d: %v", rc.ChainID, err)
			}
		}(registryConfig)
	}

	wg.Wait()
	return nil
}

func monitorRegistry(
	registryConfig RegistryConfig,
	chainConfigs map[int64]safeharbor.ChainConfig,
	firestoreClient *firestore.Client,
) error {
	chainID := registryConfig.ChainID
	apiKey := chainConfigs[chainID].APIKey
	address, err := web3.HexToAddress(registryConfig.RegistryAddress)
	if err != nil {
		return fmt.Errorf("invalid registry address: %w", err)
	}

	// Initialize last processed block number
	var lastProcessedBlock int64 = 0

	for {
		log.Printf("Monitoring ChainID: %d, RegistryAddress: %s", chainID, registryConfig.RegistryAddress)
		lastProcessedBlock, err = processEvents(registryConfig, chainConfigs, firestoreClient, apiKey, address, lastProcessedBlock)
		if err != nil {
			log.Printf("Error processing events for ChainID %d: %v", chainID, err)
		}

		time.Sleep(1 * time.Hour)
	}
}

func processEvents(
	registryConfig RegistryConfig,
	chainConfigs map[int64]safeharbor.ChainConfig,
	firestoreClient *firestore.Client,
	apiKey string,
	address *web3.Address,
	lastProcessedBlock int64,
) (int64, error) {
	chainID := registryConfig.ChainID

	// Fetch events from Etherscan starting from the last processed block
	events, err := fetchEventsFromEtherscan(chainID, apiKey, *address, lastProcessedBlock)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch events: %w", err)
	}

	// Update the last processed block number
	for _, event := range events {
		// Process each event
		err := processAgreementEvent(event, chainID, chainConfigs, firestoreClient)
		if err != nil {
			log.Printf("Error processing agreement event: %v", err)
			continue
		}
	}

	client, err := ethclient.Dial(chainConfigs[chainID].RPCURL)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to Ethereum RPC client: %w", err)
	}
	defer client.Close()

	blockNum, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block number: %w", err)
	}

	return int64(blockNum), nil
}

func fetchEventsFromEtherscan(
	chainID int64,
	apiKey string,
	address web3.Address,
	startBlock int64,
) ([]EtherscanEvent, error) {
	topic0, err := web3.HexToHash("0xfb9c334c719c97ecac9e4d31dec8572d1e2cf193a6af229da967437a30dc7010") // AgreementCreated event
	if err != nil {
		return nil, fmt.Errorf("failed to parse topic0: %w", err)
	}

	logs, err := etherscan.FetchLogs(chainID, apiKey, address, *topic0, int(startBlock))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logs: %w", err)
	}

	var events []EtherscanEvent
	for _, logEntry := range logs {
		if len(logEntry.Topics) < 2 {
			continue
		}

		// Get deployer (entity) from Topics[1]
		deployerBytes := logEntry.Topics[1][12:] // Last 20 bytes
		deployer, err := web3.BytesToAddress(deployerBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse deployer address: %w", err)
		}

		// Get safeHarborAddress from data
		if len(logEntry.Data) < 64 {
			continue
		}
		safeHarborBytes := logEntry.Data[44:64]
		safeHarborAddress, err := web3.BytesToAddress(safeHarborBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse safeHarborAddress: %w", err)
		}

		blockNumber := logEntry.BlockNumber
		txHash := logEntry.TransactionHash
		timeStamp := logEntry.TimeStamp.Int64()

		event := EtherscanEvent{
			TransactionHash:   txHash,
			BlockNumber:       blockNumber.Int64(),
			TimeStamp:         timeStamp,
			Deployer:          *deployer,
			SafeHarborAddress: *safeHarborAddress,
		}

		events = append(events, event)
	}

	return events, nil
}

func processAgreementEvent(
	event EtherscanEvent,
	chainID int64,
	chainConfigs map[int64]safeharbor.ChainConfig,
	firestoreClient *firestore.Client,
) error {
	protocolID := event.Deployer.String() // Use the deployer address as protocolID
	blockNumber := big.NewInt(event.BlockNumber)

	// Check if a newer agreement exists in Firestore
	exists, existingBlockNumber, err := checkExistingAgreement(firestoreClient, protocolID)
	if err != nil {
		return fmt.Errorf("failed to check existing agreement: %w", err)
	}

	if exists && existingBlockNumber >= event.BlockNumber {
		log.Printf("Skipping older or same agreement for deployer %s", protocolID)
		return nil // Skip processing
	}

	// Process the new agreement
	err = flow.ProcessSafeHarborAgreement(
		chainConfigs,
		event.TransactionHash,
		event.SafeHarborAddress,
		event.Deployer,
		int(chainID),
		blockNumber,
		protocolID,
		firestoreClient,
		false, // setProtocol is false
	)
	if err != nil {
		return fmt.Errorf("failed to process Safe Harbor agreement: %w", err)
	}

	// TODO: Send an email notification

	return nil
}

func checkExistingAgreement(firestoreClient *firestore.Client, eoa string) (bool, int64, error) {
	// Query Firestore to find if there's an existing agreement from the same EOA
	// Return whether it exists and its block time
	ctx := context.Background()
	iter := firestoreClient.Collection("safeHarborAgreements").Where("deployer", "==", eoa).Documents(ctx)
	defer iter.Stop()

	docs, err := iter.GetAll()
	if err != nil {
		return false, 0, fmt.Errorf("failed to query Firestore: %w", err)
	}

	if len(docs) == 0 {
		return false, 0, nil
	}

	// Assume we're interested in the latest agreement
	var latestBlockTime int64
	for _, doc := range docs {
		data := doc.Data()
		blockTime, ok := data["blockTime"].(int64)
		if ok && blockTime > latestBlockTime {
			latestBlockTime = blockTime
		}
	}

	return true, latestBlockTime, nil
}

func loadRegistryConfigs(filePath string) ([]RegistryConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry config file: %w", err)
	}

	var registryConfigs []RegistryConfig
	err = json.Unmarshal(data, &registryConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal registry config JSON: %w", err)
	}

	return registryConfigs, nil
}

func loadChainConfigs(filePath string) (map[int64]safeharbor.ChainConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read chain config file: %w", err)
	}

	var chainConfigs map[int64]safeharbor.ChainConfig
	err = json.Unmarshal(data, &chainConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chain config JSON: %w", err)
	}

	return chainConfigs, nil
}

func newFirestoreClient(credsPath string) (*firestore.Client, error) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "your-project-id", option.WithCredentialsFile(credsPath))
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}

	return client, nil
}
