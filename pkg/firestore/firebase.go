package firestore

import (
	"SHDB/pkg/clients"
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func initFirestore() (*firestore.Client, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("/path/to/your/firebase-adminsdk.json")
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

// WriteProtocolInfoToFirestore writes a ProtocolInfo struct to the 'protocols' collection
func WriteProtocolInfoToFirestore(client *firestore.Client, info clients.ProtocolInfo) (*firestore.DocumentRef, error) {
	ctx := context.Background()

	_, err := client.Collection("protocols").Doc(info.Name).Set(ctx, map[string]interface{}{
		"name":           info.Name,
		"website":        info.Website,
		"icon":           info.Icon,
		"tvl":            info.TVL,
		"category":       info.Category,
		"contactDetails": info.ContactDetails,
	})

	if err != nil {
		return nil, err
	}

	docRef := client.Collection("protocols").Doc(info.Name)

	return docRef, nil
}

// WriteAgreementDetailsToFirestore writes an AgreementDetailsV1 struct to the 'safeHarborAgreements' collection
func WriteAgreementDetailsToFirestore(client *firestore.Client, protocolName string, details clients.AgreementDetailsV1, protocolRef *firestore.DocumentRef) (*firestore.DocumentRef, error) {
	ctx := context.Background()

	// docRef, _, err := client.Collection("safeHarborAgreements").Add(ctx, map[string]interface{}{
	_, err := client.Collection("safeHarborAgreements").Doc(protocolName).Set(ctx, map[string]interface{}{
		"protocol":            protocolRef,
		"registryTransaction": details.RegistryTransaction,
		"registryExplorerURL": details.RegistryExplorerURL,
		"agreementURI":        details.AgreementURI,
		"contactDetails":      details.ContactDetails,
		"chains":              chainsToFirestoreArray(details.Chains),
		"bountyTerms": map[string]interface{}{
			"bountyPercentage": details.BountyTerms.BountyPercentage.Int64(),
			"bountyCapUSD":     details.BountyTerms.BountyCapUSD.Int64(),
			"verification":     string(details.BountyTerms.Verification),
		},
		"createdAt": firestore.ServerTimestamp,
		"entity":    details.Entity.Hex(),
	})

	if err != nil {
		return nil, err
	}

	docRef := client.Collection("safeHarborAgreements").Doc(protocolName)

	return docRef, nil
}

// WriteFirewallToFirestore writes a Firewall struct to the 'firewallAgreements' collection
func WriteFirewallToFirestore(client *firestore.Client, protocolName string, firewall clients.Firewall, protocolRef *firestore.DocumentRef) (*firestore.DocumentRef, error) {
	ctx := context.Background()

	// docRef, _, err := client.Collection("firewallAgreements").Add(ctx, map[string]interface{}{
	_, err := client.Collection("firewallAgreements").Doc(protocolName).Set(ctx, map[string]interface{}{
		"protocol":    protocolRef,
		"chains":      firewallChainsToFirestoreArray(firewall.Chains),
		"accounts":    firewallAccountsToFirestoreArray(firewall.Accounts),
		"bountyTerms": firewallBountyTermsToFirestoreMap(firewall.BountyTerms),
		"createdAt":   firestore.ServerTimestamp,
	})

	if err != nil {
		return nil, err
	}

	docRef := client.Collection("firewallAgreements").Doc(protocolName)

	return docRef, nil
}

// WriteMonitoredToFirestore writes a Monitored struct to the 'monitorableAddresses' collection
func WriteMonitoredToFirestore(client *firestore.Client, protocolName string, monitored clients.Monitored, protocolRef, safeHarborRef, firewallRef *firestore.DocumentRef) (*firestore.DocumentRef, error) {
	ctx := context.Background()

	_, err := client.Collection("monitorableAddresses").Doc(protocolName).Set(ctx, map[string]interface{}{
		"protocol":            protocolRef,
		"safeHarborAgreement": safeHarborRef,
		"firewallAgreement":   firewallRef,
		"addresses":           monitoredAccountsToFirestoreArray(monitored.Addresses),
	})

	if err != nil {
		return nil, err
	}

	ref := client.Collection("monitorableAddresses").Doc(protocolName)

	return ref, nil
}

// Helper functions

func chainsToFirestoreArray(chains []clients.ChainSH) []map[string]interface{} {
	result := make([]map[string]interface{}, len(chains))
	for i, chain := range chains {
		result[i] = map[string]interface{}{
			"assetRecoveryAddress": chain.AssetRecoveryAddress.Hex(),
			"id":                   chain.ID.Int64(),
			"accounts":             accountsSHToFirestoreArray(chain.Accounts),
		}
	}
	return result
}

func accountsSHToFirestoreArray(accounts []clients.AccountSH) []map[string]interface{} {
	result := make([]map[string]interface{}, len(accounts))
	for i, account := range accounts {
		result[i] = map[string]interface{}{
			"accountAddress":     account.AccountAddress.Hex(),
			"childContractScope": string(account.ChildContractScope),
			"signature":          account.Signature,
		}
	}
	return result
}

func firewallChainsToFirestoreArray(chains []clients.ChainFW) []map[string]interface{} {
	result := make([]map[string]interface{}, len(chains))
	for i, chain := range chains {
		result[i] = map[string]interface{}{
			"assetRecoveryAddress": chain.AssetRecoveryAddress.String(),
			"chainID":              chain.ChainID,
		}
	}
	return result
}

func firewallAccountsToFirestoreArray(accounts []clients.AccountFW) []map[string]interface{} {
	result := make([]map[string]interface{}, len(accounts))
	for i, account := range accounts {
		result[i] = map[string]interface{}{
			"name":                  account.Name,
			"address":               account.Address.String(),
			"registeredEvents":      eventsToFirestoreArray(account.RegisteredEvents),
			"includeChildAddresses": string(account.ChildContractScope),
			"chainIDs":              account.ChainIDs,
		}
	}
	return result
}

func eventsToFirestoreArray(events []clients.EventFW) []map[string]interface{} {
	result := make([]map[string]interface{}, len(events))
	for i, event := range events {
		result[i] = map[string]interface{}{
			"topic0":              event.Topic0.String(),
			"addressLocationType": string(event.AddressLocationType),
			"addressLocation":     event.AddressLocation,
		}
	}
	return result
}

func firewallBountyTermsToFirestoreMap(terms clients.BountyTermsFW) map[string]interface{} {
	return map[string]interface{}{
		"bountyPercentage": terms.BountyPercentage,
		"bountyCapUSD":     terms.BountyCapUSD,
	}
}

func monitoredAccountsToFirestoreArray(accounts []clients.AccountM) []map[string]interface{} {
	result := make([]map[string]interface{}, len(accounts))
	for i, account := range accounts {
		result[i] = map[string]interface{}{
			"name":     account.Name,
			"address":  account.Address,
			"chains":   account.Chains,
			"children": monitoredAccountsToFirestoreArray(account.Children),
		}
	}
	return result
}
