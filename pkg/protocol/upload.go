package protocol

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

const (
	SAFE_HARBOR_COLLECTION = "safeHarborAgreements"
	PROTOCOLS_COLLECTION   = "protocols"
)

type CollectionClient interface {
	Collection(string) *firestore.CollectionRef
}

func (p *Protocol) Upload(client CollectionClient, protocolID string) error {
	_, err := client.Collection("protocols").Doc(protocolID).Set(context.Background(), p)
	return err
}

func SetSafeHarborAgreement(client CollectionClient, protocolID string) error {
	safeHarborRef := client.Collection(SAFE_HARBOR_COLLECTION).Doc(protocolID)
	docRef := client.Collection(PROTOCOLS_COLLECTION).Doc(protocolID)
	_, err := docRef.Update(context.Background(), []firestore.Update{
		{
			Path: "safeHarborAgreement", Value: safeHarborRef,
		},
	})
	if err != nil {
		return fmt.Errorf("docRef.Update(): %w", err)
	}
	return nil
}
