package main

import (
	"SHDB/pkg/clients"
	"context"
	"fmt"
	"log"
	"time"

	firestoreClients "SHDB/pkg/firestore"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/api/option"
)

const (
	apiKey           = "PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3"
	rpcEndpoint      = "https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/"
	ethEndpoint      = "https://floral-wandering-telescope.quiknode.pro/3dab78094f801284adcf7ebf710a1f2bcf73bb42/"
	checkInterval    = 1 * time.Minute
	firebaseCredFile = "/home/dwu/SHDB/skylock-xyz-firebase-adminsdk-36s2d-bd6e795bf3.json"
)

func main() {
	ctx := context.Background()

	// Initialize RPC client
	rpcClient, err := rpc.Dial(rpcEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect to RPC client: %v", err)
	}
	defer rpcClient.Close()

	ethRPCClient, err := rpc.Dial(ethEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum RPC client: %v", err)
	}

	// Initialize Firestore client
	sa := option.WithCredentialsFile(firebaseCredFile)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	for {
		log.Println("Checking for new Safe Harbor adoptions...")
		err := checkAndProcessNewAdoptions(ctx, rpcClient, ethRPCClient, firestoreClient)
		if err != nil {
			log.Printf("Error processing adoptions: %v", err)
		}

		log.Printf("Sleeping for %v before next check...", checkInterval)
		time.Sleep(checkInterval)
	}
}

func checkAndProcessNewAdoptions(ctx context.Context, rpcClient *rpc.Client, ethRPCClient *rpc.Client, firestoreClient *firestore.Client) error {
	// Get all Safe Harbor adoptions
	adoptions, err := clients.GetSafeHarborAdoptions(apiKey, rpcClient)
	if err != nil {
		return fmt.Errorf("failed to get Safe Harbor adoptions: %w", err)
	}

	for _, adoption := range adoptions {
		// Check if this adoption already exists in Firestore
		fmt.Println(adoption)
		exists, err := checkAdoptionExists(ctx, firestoreClient, adoption.ProtocolName)
		if err != nil {
			log.Printf("Error checking adoption existence: %v", err)
			continue
		}

		fmt.Println("Exists: ", exists)

		if !exists {
			fmt.Println("Processing new adoption...")
			err = processNewAdoption(ctx, ethRPCClient, firestoreClient, adoption)
			if err != nil {
				log.Printf("Error processing new adoption for %s: %v", adoption.ProtocolName, err)
			}
		}
	}

	return nil
}

func checkAdoptionExists(ctx context.Context, client *firestore.Client, protocolName string) (bool, error) {
	docs, err := client.Collection("safeHarborAgreements").Where("protocol.name", "==", protocolName).Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}
	return len(docs) > 0, nil
}

func processNewAdoption(ctx context.Context, rpcClient *rpc.Client, firestoreClient *firestore.Client, adoption clients.AgreementDetailsV1) error {
	// Get protocol info
	protocolInfo, err := clients.GetProtocolInfo(adoption.ProtocolName)
	if err != nil {
		return fmt.Errorf("failed to get protocol info: %w", err)
	}

	// Generate Monitored data
	monitored, err := adoption.ToMonitored(rpcClient, apiKey)
	if err != nil {
		return fmt.Errorf("failed to generate monitored data: %w", err)
	}

	// Write protocol info to Firestore
	protocolRef, err := firestoreClients.WriteProtocolInfoToFirestore(firestoreClient, *protocolInfo)
	if err != nil {
		return fmt.Errorf("failed to write protocol info: %w", err)
	}

	// Write agreement details to Firestore
	agreementRef, err := firestoreClients.WriteAgreementDetailsToFirestore(firestoreClient, adoption, protocolRef)
	if err != nil {
		return fmt.Errorf("failed to write agreement details: %w", err)
	}

	// Write Monitored data to Firestore
	_, err = firestoreClients.WriteMonitoredToFirestore(firestoreClient, *monitored, protocolRef, agreementRef, nil)
	if err != nil {
		return fmt.Errorf("failed to write monitored data: %w", err)
	}

	log.Printf("Successfully processed new adoption for %s", adoption.ProtocolName)
	return nil
}
