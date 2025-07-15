package types

import (
	"time"

	"cloud.google.com/go/firestore"
)

type SafeHarborVersion string

const (
	SealV1     SafeHarborVersion = "seal-1"
	SealV2     SafeHarborVersion = "seal-2"
	ImmunefiV1 SafeHarborVersion = "immunefi-1"
)

type SafeHarborAgreementBase struct {
	AdoptionProposalURI string                 `firestore:"adoptionProposalURI"`
	Protocol            *firestore.DocumentRef `firestore:"protocol"`
	Slug                string                 `firestore:"slug"`
	Version             SafeHarborVersion      `firestore:"version"`
}

type SafeHarborAgreementV1 struct {
	SafeHarborAgreementBase
	AgreementAddress    string             `firestore:"agreementAddress"`
	AgreementDetails    AgreementDetailsV1 `firestore:"agreementDetails"`
	CreatedAt           time.Time          `firestore:"createdAt"`
	Creator             string             `firestore:"creator"`
	RegistryChainID     int                `firestore:"registryChainId"`
	RegistryTransaction string             `firestore:"registryTransaction"`
}

type SafeHarborAgreementImmunefiV1 struct {
	SafeHarborAgreementBase
	AgreementDetails ImmunefiDetailsV1 `firestore:"agreementDetails"`
}
