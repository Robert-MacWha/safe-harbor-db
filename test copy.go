package main

// import (
// 	"SHDB/pkg/protocol"
// 	safeharbor "SHDB/pkg/safeHarbor"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"strings"

// 	"cloud.google.com/go/firestore"
// 	"google.golang.org/api/option"
// )

// func main() {
// 	// // Example usage

// 	// // Connect to an Ethereum node
// 	// client, err := ethclient.Dial("https://wandering-warmhearted-glade.matic.quiknode.pro/052abdbd5b5edfa400f9bed3ff2d14501d31759f")
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to connect to the Ethereum client: %v", err)
// 	// }

// 	// txHash, err := web3.HexToHash("0x62a554a7a8f8a7ab49f41b4df5b72eea6ca30680adc2f61f608d3c2d47296685")
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to parse txHash: %v", err)
// 	// }

// 	// // List of safeHarborAddresses to query
// 	// safeHarborAddress, err := web3.HexToAddress("0x2f6748580b200b9b2ace5774edc2657ff7ccc56b")
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to parse safeHarborAddress: %v", err)
// 	// }

// 	// deployer, err := web3.HexToAddress("0x566345a70d70ce724cc1a441dca748b6b6c31265")
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to parse deployer: %v", err)
// 	// }

// 	// blockNumber := *big.NewInt(63210447)

// 	// commonSafeHarborAddress := safeHarborAddress.ToCommon()

// 	// // Fetch agreement details
// 	// agreement, _, err := safeharbor.FetchAgreementDetails(blockNumber, deployer.String(), txHash.String(), commonSafeHarborAddress, client)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to fetch agreement details: %v", err)
// 	// }

// 	// log.Printf("Agreement: %+v", agreement)

// 	// chainConfigs := map[int64]safeharbor.ChainConfig{
// 	// 	1: {
// 	// 		APIKey: "PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3",
// 	// 		RPCURL: "https://quick-distinguished-haze.quiknode.pro/766c1d076c4d839418033b62f97582feaf0f6c97",
// 	// 	},
// 	// 	137: {
// 	// 		APIKey: "2UYYF5VKP189UKN33F57UAFWPW7Q658GNB",
// 	// 		RPCURL: "https://wandering-warmhearted-glade.matic.quiknode.pro/052abdbd5b5edfa400f9bed3ff2d14501d31759f",
// 	// 	},
// 	// }

// 	// // Process agreement details
// 	// agreement, err = safeharbor.ProcessSafeHarborAgreement(agreement, chainConfigs)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to process agreement details: %v", err)
// 	// }

// 	// data, err := json.Marshal(agreement)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to marshal agreement: %v", err)
// 	// }

// 	// // save agreement details to a file
// 	// err = os.WriteFile("agreementSave.json", data, 0644)

// 	// Unmarshal agreement details from a file
// 	data, err := os.ReadFile("agreementSave.json")
// 	if err != nil {
// 		log.Fatalf("Failed to read agreementSave.json: %v", err)
// 	}

// 	agreement := &safeharbor.SafeHarborAgreement{}
// 	err = json.Unmarshal(data, agreement)
// 	if err != nil {
// 		log.Fatalf("Failed to unmarshal agreement: %v", err)
// 	}

// 	// log.Printf("Agreement: %+v", agreement)

// 	protocolName := "polymarket"

// 	resultProtocol, err := protocol.GetProtocol(protocolName)
// 	if err != nil {
// 		log.Fatalf("Failed to get protocol: %v", err)
// 	}

// 	credsFile, err := os.Open("test_data/skylock-xyz-firebase-adminsdk-36s2d-bd6e795bf3.json")
// 	if err != nil {
// 		log.Fatalf("Failed to open creds file: %v", err)
// 	}

// 	firestoreClient, err := newClient(credsFile)
// 	if err != nil {
// 		log.Fatalf("Failed to create Firestore client: %v", err)
// 	}

// 	protocolID := strings.ToLower(protocolName)

// 	// Upload protocol to Firestore
// 	resultProtocol.Upload(firestoreClient, protocolID)

// 	// Upload safe harbor agreement to Firestore
// 	agreement.Upload(firestoreClient, protocolID)

// 	// Set safe harbor agreement reference in protocol
// 	err = protocol.SetSafeHarborAgreement(firestoreClient, protocolID)

// }

// func newClient(creds io.Reader) (*firestore.Client, error) {
// 	ctx := context.Background()

// 	p, err := io.ReadAll(creds)
// 	if err != nil {
// 		return nil, fmt.Errorf("io.ReadAll: %w", err)
// 	}

// 	client, err := firestore.NewClient(ctx, "skylock-xyz", option.WithCredentialsJSON(p))
// 	if err != nil {
// 		return nil, fmt.Errorf("firestore.NewClient: %w", err)
// 	}

// 	return client, nil
// }
