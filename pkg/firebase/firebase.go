package firebase

import (
	"SHDB/pkg/clients"
	"context"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func initFirestore() (*firestore.Client, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("/home/dwu/Arianrhod/skylock-xyz-firebase-adminsdk-36s2d-bd6e795bf3.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func clientToFirestoreMap(client clients.Client) map[string]interface{} {
	return map[string]interface{}{
		"ProtocolName":   client.ProtocolName,
		"AgreementURI":   client.AgreementURI,
		"ContactDetails": client.ContactDetails,
		"BountyTerms": map[string]interface{}{
			"BountyPercentage": client.BountyTerms.BountyPercentage,
			"BountyCapUSD":     client.BountyTerms.BountyCapUSD,
			"retainable":       string(client.BountyTerms.Retainable),
		},
		"Chains":         chainsToFirestoreArray(client.Chains),
		"Accounts":       accountsToFirestoreArray(client.Accounts),
		"Website":        client.Website,
		"Date":           client.Date,
		"Icon":           client.Icon,
		"TVL":            client.TVL,
		"Category":       client.Category,
		"FirewallClient": client.FirewallClient,
	}
}

func chainsToFirestoreArray(chains []clients.Chain) []map[string]interface{} {
	result := make([]map[string]interface{}, len(chains))
	for i, chain := range chains {
		result[i] = map[string]interface{}{
			"AssetRecoveryAddress": chain.AssetRecoveryAddress.String(),
			"ChainID":              chain.ChainID,
		}
	}
	return result
}

func accountsToFirestoreArray(accounts []clients.Account) []map[string]interface{} {
	result := make([]map[string]interface{}, len(accounts))
	for i, account := range accounts {
		lastQuery := make(map[string]int)
		for k, v := range account.LastQuery {
			lastQuery[strconv.FormatInt(k, 10)] = v
		}

		result[i] = map[string]interface{}{
			"Name":               account.Name,
			"Address":            account.Address.String(),
			"ChildContractScope": string(account.ChildContractScope),
			"SignatureValid":     account.SignatureValid,
			"ChainIDs":           account.ChainIDs,
			"LastQuery":          lastQuery,
			"RegisteredEvents":   eventsToFirestoreArray(account.RegisteredEvents),
			"Children":           accountsToFirestoreArray(account.Children),
		}
	}
	return result
}

func eventsToFirestoreArray(events []clients.Event) []map[string]interface{} {
	result := make([]map[string]interface{}, len(events))
	for i, event := range events {
		result[i] = map[string]interface{}{
			"Topic0":              event.Topic0.String(),
			"AddressLocationType": string(event.AddressLocationType),
			"AddressLocation":     event.AddressLocation,
		}
	}
	return result
}

func writeClientToFirestore(client clients.Client) error {
	ctx := context.Background()
	firestoreClient, err := initFirestore()
	if err != nil {
		return err
	}
	defer firestoreClient.Close()

	clientData := clientToFirestoreMap(client)

	// Create a unique ID by combining ProtocolName and a unique identifier
	docID := client.ProtocolName + "_" + generateUID()

	_, err = firestoreClient.Collection("firewall").Doc(docID).Set(ctx, clientData)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to generate a unique identifier
func generateUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())

}

func getClientByProtocolName(protocolName string) (*clients.Client, error) {
	ctx := context.Background()
	firestoreClient, err := initFirestore()
	if err != nil {
		return nil, err
	}
	defer firestoreClient.Close()

	iter := firestoreClient.Collection("firewall").Where("ProtocolName", "==", protocolName).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("no client found with protocol name: %s", protocolName)
	}

	// Assuming there's only one document per protocol name
	doc := docs[0]
	var client clients.Client
	err = doc.DataTo(&client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
