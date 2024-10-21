package flow

import (
	"SHDB/pkg/protocol"
	safeharbor "SHDB/pkg/safeHarbor"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// ProcessProtocol processes and uploads a protocol to Firestore and optionally sets the Safe Harbor Agreement reference.
func ProcessProtocol(
	protocolName string,
	credsPath string,
	setSafeHarbor bool,
) error {
	// Open the credentials file
	credsFile, err := os.Open(credsPath)
	if err != nil {
		return fmt.Errorf("failed to open Firestore credentials file: %w", err)
	}
	defer credsFile.Close()

	// Create Firestore client
	firestoreClient, err := newFirestoreClient(credsFile)
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}
	defer firestoreClient.Close()

	// Fetch the protocol
	resultProtocol, err := protocol.GetProtocol(protocolName)
	if err != nil {
		return fmt.Errorf("failed to get protocol: %w", err)
	}

	// Generate a lowercase protocol ID
	protocolID := strings.ToLower(protocolName)

	// Upload protocol to Firestore
	err = resultProtocol.Upload(firestoreClient, protocolID)
	if err != nil {
		return fmt.Errorf("failed to upload protocol to Firestore: %w", err)
	}

	// Optionally set the Safe Harbor Agreement reference
	if setSafeHarbor {
		err = protocol.SetSafeHarborAgreement(firestoreClient, protocolID)
		if err != nil {
			return fmt.Errorf("failed to set Safe Harbor Agreement reference: %w", err)
		}
		err = safeharbor.SetProtocol(firestoreClient, protocolID)
		if err != nil {
			return fmt.Errorf("failed to set protocol reference: %w", err)
		}
	}

	log.Printf("Successfully processed protocol: %s", protocolName)
	return nil
}

// newFirestoreClient creates a Firestore client using credentials
func newFirestoreClient(creds io.Reader) (*firestore.Client, error) {
	ctx := context.Background()

	// Read credentials file
	credsData, err := io.ReadAll(creds)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	// Create Firestore client
	client, err := firestore.NewClient(ctx, "skylock-xyz", option.WithCredentialsJSON(credsData))
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}

	return client, nil
}
