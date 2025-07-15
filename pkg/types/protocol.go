package types

import "cloud.google.com/go/firestore"

// Protocol
type Protocol struct {
	Name                string                 `firestore:"name"`
	Website             string                 `firestore:"website"`
	Slug                string                 `firestore:"slug"`
	Icon                string                 `firestore:"icon"`
	TVL                 float64                `firestore:"tvl"`
	Category            string                 `firestore:"category"`
	ContactDetails      string                 `firestore:"contactDetails"`
	SafeHarborAgreement *firestore.DocumentRef `firestore:"safeHarborAgreement,omitempty"`
}
