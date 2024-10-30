package main

import (
	"SHDB/pkg/etherscan"
	"SHDB/pkg/flow"
	"SHDB/pkg/protocol"
	safeharbor "SHDB/pkg/safeHarbor"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	// Import necessary types
	"SHDB/pkg/web3"
)

const updateRegistryInterval = 1 * time.Hour
const updateProtocolsInterval = 1 * time.Hour

// RegistryConfig represents the structure of the registry data
type RegistryConfig struct {
	CommitHash string `json:"commitHash"`
	Registries []struct {
		Address string `json:"address"`
		ChainID int64  `json:"chainID"`
	} `json:"registries"`
	Version string `json:"version"`
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
		Name:   "safe-harbor-monitor",
		Usage:  "Monitors Safe Harbor registry addresses for new agreements",
		Flags:  []cli.Flag{},
		Action: Run,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}

func Run(c *cli.Context) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	chainConfigs, err := loadChainConfigs()
	if err != nil {
		return fmt.Errorf("failed to load chain configs: %w", err)
	}

	// Create Firestore client
	firestoreClient, err := newFirestoreClient()
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}
	defer firestoreClient.Close()

	registryConfig, err := loadRegistryConfigFromFirestore(firestoreClient)
	if err != nil {
		return fmt.Errorf("failed to load registry configs: %w", err)
	}

	var wg sync.WaitGroup

	for _, registryConfig := range registryConfig.Registries {
		wg.Add(1)
		go func(chainID int64, registryAddress string) {
			defer wg.Done()
			err := monitorRegistry(chainID, registryAddress, chainConfigs, firestoreClient)
			if err != nil {
				log.Printf("Error monitoring registry for ChainID %d: %v", chainID, err)
			}
		}(registryConfig.ChainID, registryConfig.Address)
	}

	go monitorSafeHarborProtocols(firestoreClient, chainConfigs)

	wg.Wait()
	return nil
}

func monitorRegistry(
	chainID int64,
	registryAddress string,
	chainConfigs map[int64]safeharbor.ChainConfig,
	firestoreClient *firestore.Client,
) error {
	apiKey := chainConfigs[chainID].APIKey
	address, err := web3.HexToAddress(registryAddress)
	if err != nil {
		return fmt.Errorf("invalid registry address: %w", err)
	}

	// Initialize last processed block number
	var lastProcessedBlock int64 = 0

	log.Printf("Monitoring ChainID: %d, RegistryAddress: %s", chainID, registryAddress)

	for {
		lastProcessedBlock, err = processEvents(chainID, chainConfigs, firestoreClient, apiKey, address, lastProcessedBlock)
		if err != nil {
			log.Printf("Error processing events for ChainID %d: %v", chainID, err)
		}

		time.Sleep(updateRegistryInterval)
	}
}

func processEvents(
	chainID int64,
	chainConfigs map[int64]safeharbor.ChainConfig,
	firestoreClient *firestore.Client,
	apiKey string,
	address *web3.Address,
	lastProcessedBlock int64,
) (int64, error) {
	// Fetch events from Etherscan starting from the last processed block
	events, err := fetchEventsFromEtherscan(chainID, apiKey, *address, lastProcessedBlock)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch events: %w", err)
	}

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

	// Create a map to store the latest event for each deployer (entity)
	latestEvents := make(map[string]EtherscanEvent)

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

		// Store the latest event for each deployer (entity)
		deployerKey := deployer.String()
		latestEvent, exists := latestEvents[deployerKey]

		// If this event's block number is higher than the currently stored event, update it
		if !exists || blockNumber.Int64() > latestEvent.BlockNumber {
			latestEvents[deployerKey] = event
		}
	}

	// Collect the latest events into a slice
	var events []EtherscanEvent
	for _, event := range latestEvents {
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
	eoa := event.Deployer.String() // Use the deployer address as protocolID
	blockNumber := big.NewInt(event.BlockNumber)

	// Check if a newer agreement exists in Firestore
	exists, existingBlockNumber, err := checkExistingAgreement(firestoreClient, eoa)
	if err != nil {
		return fmt.Errorf("failed to check existing agreement: %w", err)
	}

	if exists && existingBlockNumber >= event.TimeStamp {
		log.Printf("Skipping older or same agreement for deployer %s", eoa)
		return nil // Skip processing
	}

	// Process the new agreement
	protocolID := strings.ToLower(event.Deployer.String())
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

	err = sendEmail(event, chainID, protocolID)

	return err
}

func checkExistingAgreement(firestoreClient *firestore.Client, eoa string) (bool, int64, error) {
	// Query Firestore to find the agreement for the specified EOA
	ctx := context.Background()
	docRef := firestoreClient.Collection("safeHarborAgreements").Where("entity", "==", eoa).Limit(1)
	iter := docRef.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return false, 0, nil // No document found
	}
	if err != nil {
		return false, 0, fmt.Errorf("failed to query Firestore: %w", err)
	}

	// Extract and return the `createdAt` timestamp from the document
	data := doc.Data()
	createdAt, ok := data["createdAt"].(time.Time)
	if !ok {
		return false, 0, fmt.Errorf("failed to parse createdAt field")
	}

	return true, int64(createdAt.Unix()), nil
}

func loadChainConfigs() (map[int64]safeharbor.ChainConfig, error) {
	data := os.Getenv("CHAIN_CONFIG")
	if data == "" {
		return nil, fmt.Errorf("CHAIN_CONFIGS environment variable not set")
	}

	var chainConfigs map[int64]safeharbor.ChainConfig
	err := json.Unmarshal([]byte(data), &chainConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chain config JSON: %w", err)
	}

	return chainConfigs, nil
}

// loadRegistryConfigFromFirestore fetches a single registry config from Firestore
func loadRegistryConfigFromFirestore(client *firestore.Client) (RegistryConfig, error) {
	ctx := context.Background()

	// Get the document by its ID in the "registries" collection
	docRef := client.Collection("registries").Doc("977TQDEJbjykCRzqVnNZ")
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return RegistryConfig{}, fmt.Errorf("failed to get document: %w", err)
	}

	// Unmarshal the document data into a RegistryConfig struct
	var config RegistryConfig
	if err := docSnap.DataTo(&config); err != nil {
		return RegistryConfig{}, fmt.Errorf("failed to unmarshal document data: %w", err)
	}

	return config, nil
}

func newFirestoreClient() (*firestore.Client, error) {
	// Load from env
	creds := os.Getenv("FIREBASE_CREDENTIALS")
	if creds == "" {
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS environment variable not set")
	}

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "skylock-xyz", option.WithCredentialsJSON([]byte(creds)))
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}

	return client, nil
}

// monitorSafeHarborProtocols runs every updateProtocolsInterval time and processes all active Safe Harbor protocols.
func monitorSafeHarborProtocols(
	firestoreClient *firestore.Client,
	chainConfigs map[int64]safeharbor.ChainConfig,
) {
	for {
		log.Println("Fetching and processing current Safe Harbor protocols...")

		// Fetch all the current Safe Harbor protocols from Firestore
		err := fetchAndProcessSafeHarborAgreements(firestoreClient, chainConfigs)
		if err != nil {
			log.Printf("Error fetching or processing agreements: %v", err)
		}

		err = fetchAndProcessProtocols(firestoreClient)
		if err != nil {
			log.Printf("Error fetching or processing protocols: %v", err)
		}

		time.Sleep(updateProtocolsInterval)
	}
}

// fetchAndProcessSafeHarborAgreements retrieves all active Safe Harbor protocols from Firestore and processes them.
func fetchAndProcessSafeHarborAgreements(
	firestoreClient *firestore.Client,
	chainConfigs map[int64]safeharbor.ChainConfig,
) error {
	ctx := context.Background()
	collection := firestoreClient.Collection("safeHarborAgreements")
	iter := collection.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error iterating Firestore documents: %w", err)
		}

		// Extract SafeHarborAgreement from Firestore document
		var agreement safeharbor.SafeHarborAgreement
		err = doc.DataTo(&agreement)
		if err != nil {
			log.Printf("Failed to parse Firestore data for agreement: %v", err)
			continue
		}

		// Process the Safe Harbor agreement for each chain
		for _, chain := range agreement.Chains {
			chainID, err := strconv.Atoi(chain.ID)
			if err != nil {
				log.Printf("Invalid chain ID for agreement: %s", err)
				continue
			}

			// Get the latest indexed block for this chain
			startBlock := chain.LastIndexedBlock

			// Call ProcessSafeHarborAgreement function
			processedAgreement, err := safeharbor.ProcessSafeHarborAgreement(&agreement, chainConfigs, startBlock)
			if err != nil {
				log.Printf("Failed to process agreement for chain %d: %v", chainID, err)
				continue
			}

			err = processedAgreement.Upload(firestoreClient, doc.Ref.ID)
			if err != nil {
				log.Printf("Failed to upload processed agreement to Firestore: %v", err)
				continue
			}

			// Update Firestore with the latest Safe Harbor agreement details
			err = safeharbor.SetProtocol(firestoreClient, doc.Ref.ID)
			if err != nil {
				log.Printf("Failed to update agreement in Firestore: %v", err)
				continue
			}

			log.Printf("Successfully processed and updated agreement for chain %d", chainID)
		}
	}

	return nil
}

// fetchAndProcessProtocols retrieves all protocols from Firestore, processes them, and uploads them back.
func fetchAndProcessProtocols(
	firestoreClient *firestore.Client,
) error {
	ctx := context.Background()
	collection := firestoreClient.Collection("protocols")
	iter := collection.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error iterating Firestore documents: %w", err)
		}

		// Get the protocol name (assuming it's in a field called "name")
		protocolName := doc.Data()["slug"].(string)
		protocolSlug := strings.ToLower(protocolName)

		// Fetch the protocol details
		resultProtocol, err := protocol.GetProtocol(protocolSlug)
		if err != nil {
			log.Printf("Failed to get protocol '%s': %v", protocolName, err)
			continue
		}

		// Upload the processed protocol back to Firestore
		err = resultProtocol.Upload(firestoreClient, protocolSlug)
		if err != nil {
			log.Printf("Failed to upload protocol '%s' to Firestore: %v", protocolSlug, err)
			continue
		}

		err = protocol.SetSafeHarborAgreement(firestoreClient, protocolSlug)
		if err != nil {
			log.Printf("Failed to set Safe Harbor Agreement for protocol '%s': %v", protocolSlug, err)
			continue
		}

		// Set the protocol reference in Safe Harbor
		err = safeharbor.SetProtocol(firestoreClient, protocolSlug)
		if err != nil {
			log.Printf("Failed to set protocol reference for Safe Harbor: %v", err)
			continue
		}

		log.Printf("Successfully processed and uploaded protocol '%s'", protocolSlug)
	}

	return nil
}

func sendEmail(event EtherscanEvent, chainID int64, protocolID string) error {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("SENDGRID_API_KEY environment variable not set")
	}

	mg := mailgun.NewMailgun("sandbox8e9c6c34bc304b27a1f1c840e59a11b3.mailgun.org", apiKey)

	sender := "mailgun@sandbox8e9c6c34bc304b27a1f1c840e59a11b3.mailgun.org"
	subject := "Yo New Protocol"
	body := fmt.Sprintf(
		"New Protocol plz check\nChainID: %d \nTxHash: %s \nDeployer: %s \nProtocol ID: %s",
		chainID,
		event.TransactionHash.String(),
		event.Deployer.String(),
		protocolID,
	)
	recipient := "dickson@skylock.xyz"

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return fmt.Errorf("mailgun error: %v", err)
	}

	fmt.Println(resp, id)

	return nil
}
