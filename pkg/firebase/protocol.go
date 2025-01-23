package firebase

import "cloud.google.com/go/firestore"

type Protocol struct {
	Name                string                 `firestore:"name"`
	Slug                string                 `firestore:"slug"`
	Website             string                 `firestore:"website"`
	Icon                string                 `firestore:"icon"`
	TVL                 float64                `firestore:"tvl"`
	Category            string                 `firestore:"category"`
	ContactDetails      string                 `firestore:"contactDetails"`
	SafeHarborAgreement *firestore.DocumentRef `firestore:"safeHarborAgreement"` // Reference to SafeHarborAgreement document
	FirewallAgreement   *firestore.DocumentRef `firestore:"firewallAgreement"`   // Reference to FirewallAgreement document
	Owner               *firestore.DocumentRef `firestore:"owner"`               // Reference to Owner document
}
