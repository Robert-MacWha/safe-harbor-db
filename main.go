package main

import (
	"SHDB/pkg/clients"
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/rpc"

	firestoreClients "SHDB/pkg/firestore"

	firebase "firebase.google.com/go"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"google.golang.org/api/option"
)

const (
	apiKey           = "PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3"
	rpcEndpoint      = "https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/"
	checkInterval    = 1 * time.Hour
	firebaseCredFile = "/home/dwu/SHDB/skylock-xyz-firebase-adminsdk-36s2d-bd6e795bf3.json"
)

func addressOrPanic(address string) web3.Address {
	addr, err := web3.HexToAddress(address)
	if err != nil {
		panic(err)
	}
	return *addr
}

func hashOrPanic(hash string) web3.Hash {
	h, err := web3.HexToHash(hash)
	if err != nil {
		panic(err)
	}
	return *h
}

func main() {
	ctx := context.Background()

	rpcClient, err := rpc.Dial("https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/")
	if err != nil {
		panic(err)
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

	// firewallClient := clients.Firewall{
	// 	ProtocolName: "Native",
	// 	Accounts: []clients.AccountFW{
	// 		{
	// 			Address: addressOrPanic("0x85b0f66e83515ff4e825dfcaa58e040e08278ef9"),
	// 			RegisteredEvents: []clients.EventFW{
	// 				{
	// 					Topic0:              hashOrPanic("0xcf4eef09472ed02a0130e7d806d47fbd0db559056897eb16edb78506bdf45eed"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     12,
	// 				},
	// 				{
	// 					Topic0:              hashOrPanic("0x82cc5194ed1e660e8a6a4b2c99f7305283b16d8db2a1b2c9bb440096a4a07435"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     12,
	// 				},
	// 				{
	// 					Topic0:              hashOrPanic("0x583aa3202641e70abf6c4e526dd0bf713aa272b67f557ed30f2e40a37d985a8e"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     108,
	// 				},
	// 			},
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 	},
	// }

	firewallClient := clients.Firewall{
		ProtocolName: "Native",
		Accounts: []clients.AccountFW{
			{
				Address: addressOrPanic("0x85b0f66e83515ff4e825dfcaa58e040e08278ef9"),
				RegisteredEvents: []clients.EventFW{
					{
						Topic0:              hashOrPanic("0xcf4eef09472ed02a0130e7d806d47fbd0db559056897eb16edb78506bdf45eed"),
						AddressLocationType: clients.AddressLocationTypeData,
						AddressLocation:     12,
					},
					{
						Topic0:              hashOrPanic("0x82cc5194ed1e660e8a6a4b2c99f7305283b16d8db2a1b2c9bb440096a4a07435"),
						AddressLocationType: clients.AddressLocationTypeData,
						AddressLocation:     12,
					},
					{
						Topic0:              hashOrPanic("0x583aa3202641e70abf6c4e526dd0bf713aa272b67f557ed30f2e40a37d985a8e"),
						AddressLocationType: clients.AddressLocationTypeData,
						AddressLocation:     108,
					},
				},
				ChildContractScope: clients.ChildContractScopeNone,
				ChainIDs:           []int64{1},
			},
		},
	}

	protocolInfo, err := clients.GetProtocolInfo(firewallClient.ProtocolName)
	if err != nil {
		panic(err)
	}

	// Generate Monitored data
	monitored, err := firewallClient.ToMonitored(rpcClient, apiKey)
	if err != nil {
		panic(err)
	}

	// Write protocol info to Firestore
	protocolRef, err := firestoreClients.WriteProtocolInfoToFirestore(firestoreClient, *protocolInfo)
	if err != nil {
		panic(err)
	}

	// Write agreement details to Firestore
	firewallRef, err := firestoreClients.WriteFirewallToFirestore(firestoreClient, firewallClient, protocolRef)
	if err != nil {
		panic(err)
	}

	// Write Monitored data to Firestore
	_, err = firestoreClients.WriteMonitoredToFirestore(firestoreClient, *monitored, protocolRef, firewallRef, nil)
	if err != nil {
		panic(err)
	}

	// chainID := int64(1)
	// apiKey := "PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3"
	// address, err := web3.HexToAddress("0xdf3601014686674e53d1fa52f7602525483f9122")
	// if err != nil {
	// 	panic(err)
	// }
	// startBlock := 19017136
	// _, err = etherscan.FetchRegularTransactions(chainID, apiKey, *address, startBlock)
	// if err != nil {
	// 	panic(err)
	// }
}

// import (
// 	"SHDB/pkg/clients"
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"cloud.google.com/go/firestore"
// 	firebase "firebase.google.com/go"
// 	"github.com/ethereum/go-ethereum/rpc"
// 	"google.golang.org/api/option"
// )

// const (
// 	apiKey           = "PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3"
// 	rpcEndpoint      = "https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/"
// 	checkInterval    = 1 * time.Hour
// 	firebaseCredFile = "/home/dwu/SHDB/skylock-xyz-firebase-adminsdk-36s2d-bd6e795bf3.json"
// )

// func main() {
// 	ctx := context.Background()
// 	rpcClient, err := rpc.Dial("https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// details, err := clients.GetSafeHarborAdoptions("PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3", rpcClient)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// fmt.Println(details)
// 	sa := option.WithCredentialsFile(firebaseCredFile)
// 	app, err := firebase.NewApp(ctx, nil, sa)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize Firebase app: %v", err)
// 	}

// 	firestoreClient, err := app.Firestore(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize Firestore client: %v", err)
// 	}

// 	adoptions, err := clients.GetSafeHarborAdoptions(apiKey, rpcClient)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, adoption := range adoptions {
// 		fmt.Println(adoption)
// 		// Check if this adoption already exists in Firestore
// 		exists, err := checkAdoptionExists(ctx, firestoreClient, adoption.ProtocolName)
// 		if err != nil {
// 			log.Printf("Error checking adoption existence: %v", err)
// 			continue
// 		}

// 		fmt.Println(exists)

// 		// if !exists {
// 		// 	err = processNewAdoption(ctx, rpcClient, firestoreClient, adoption)
// 		// 	if err != nil {
// 		// 		log.Printf("Error processing new adoption for %s: %v", adoption.ProtocolName, err)
// 		// 	}
// 		// }
// 	}
// }

// func checkAdoptionExists(ctx context.Context, client *firestore.Client, protocolName string) (bool, error) {
// 	docs, err := client.Collection("safeHarborAgreements").Where("protocol.name", "==", protocolName).Documents(ctx).GetAll()
// 	if err != nil {
// 		return false, err
// 	}
// 	return len(docs) > 0, nil
// }
