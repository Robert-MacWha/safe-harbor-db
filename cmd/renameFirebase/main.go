package main

import (
	"SHDB/pkg/protocol"
	safeharbor "SHDB/pkg/safeHarbor"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/option"
)

const (
	SAFE_HARBOR_COLLECTION = "safeHarborAgreements"
	PROTOCOLS_COLLECTION   = "protocols"
)

// RenameDocument renames a document in the specified collection by copying the data to a new document with a new ID and deleting the old document.
func RenameDocument(client *firestore.Client, collectionName string, oldID string, newID string) error {
	ctx := context.Background()

	// Get reference to the old document
	oldDocRef := client.Collection(collectionName).Doc(oldID)

	// Read the old document data
	docSnap, err := oldDocRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get document with ID %s: %w", oldID, err)
	}

	// Check if document exists
	if !docSnap.Exists() {
		return fmt.Errorf("document with ID %s does not exist", oldID)
	}

	// Extract data from the old document
	docData := docSnap.Data()

	// Create a new document with the new ID and copy the data
	newDocRef := client.Collection(collectionName).Doc(newID)
	_, err = newDocRef.Set(ctx, docData)
	if err != nil {
		return fmt.Errorf("failed to create new document with ID %s: %w", newID, err)
	}

	// Delete the old document
	_, err = oldDocRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document with ID %s: %w", oldID, err)
	}

	return nil
}

// Run is the main action function for the CLI
func Run(c *cli.Context) error {
	// Get the flags
	credsPath := c.String("creds")
	setSafeHarbor := c.Bool("setSafeHarbor")
	oldID := c.String("oldID")
	newID := c.String("protocolID")

	fmt.Println(oldID, newID)

	// Initialize Firestore client
	credsFile, err := os.Open(credsPath)
	if err != nil {
		return fmt.Errorf("failed to open Firestore credentials file: %w", err)
	}
	defer credsFile.Close()

	firestoreClient, err := newFirestoreClient(credsFile)
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}
	defer firestoreClient.Close()

	fmt.Printf("Renaming document: %s -> %s\n", oldID, newID)
	err = RenameDocument(firestoreClient, SAFE_HARBOR_COLLECTION, oldID, newID)
	if err != nil {
		return fmt.Errorf("failed to rename document: %w", err)
	}

	// Optionally set Safe Harbor Agreement reference if the flag is enabled
	if setSafeHarbor {
		err = safeharbor.SetProtocol(firestoreClient, newID)
		if err != nil {
			return fmt.Errorf("failed to set Safe Harbor Agreement reference: %w", err)
		}
		err = protocol.SetSafeHarborAgreement(firestoreClient, newID)
		if err != nil {
			return fmt.Errorf("failed to set protocol reference: %w", err)
		}
	}

	fmt.Println("Operation completed successfully.")
	return nil
}

func main() {
	// Define the CLI app
	app := &cli.App{
		Name:  "rename",
		Usage: "A simple CLI app to rename Firestore documents and set Safe Harbor Agreement",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "oldID",
				Aliases:  []string{"old"},
				Usage:    "Old Firestore (ID)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "protocolID",
				Aliases:  []string{"p"},
				Usage:    "Protocol's Firestore Document Ref",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "creds",
				Aliases:  []string{"c"},
				Usage:    "Path to Firestore credentials file",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "setSafeHarbor",
				Usage: "Set the Safe Harbor Agreement reference",
				Value: true,
			},
		},
		Action: Run,
	}

	// Run the CLI app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
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
