package safeharbor

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

const (
	SAFE_HARBOR_COLLECTION = "safeHarborAgreements"
	PROTOCOLS_COLLECTION   = "protocols"
)

func (s SafeHarborAgreement) Upload(client *firestore.Client, protocolID string) error {
	_, err := client.Collection(SAFE_HARBOR_COLLECTION).Doc(protocolID).Set(context.Background(), s)
	if err != nil {
		return fmt.Errorf("client.Collection().Doc().Set(): %w", err)
	}

	return nil

}

func SetProtocol(client *firestore.Client, protocolID string) error {
	safeHarborRef := client.Collection(SAFE_HARBOR_COLLECTION).Doc(protocolID)
	protocolRef := client.Collection(PROTOCOLS_COLLECTION).Doc(protocolID)
	_, err := safeHarborRef.Update(context.Background(), []firestore.Update{
		{
			Path: "protocol", Value: protocolRef,
		},
	})
	if err != nil {
		return fmt.Errorf("protocolRef.Update(): %w", err)
	}

	return nil

}
