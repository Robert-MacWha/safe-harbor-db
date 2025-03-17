package firebase

import (
	"time"

	"cloud.google.com/go/firestore"
)

const (
	ChildContractScopeNone         ChildContractScope = "None"
	ChildContractScopeExistingOnly ChildContractScope = "ExistingOnly"
	ChildContractScopeAll          ChildContractScope = "All"
)

const (
	IdentityAnonymous   Identity = "Anonymous"
	IdentityPsedonymous Identity = "Pseudonymous"
	IdentityNamed       Identity = "Named"
)

var (
	ChildContractScopes = []ChildContractScope{ChildContractScopeNone, ChildContractScopeExistingOnly, ChildContractScopeAll}
	Identities          = []Identity{IdentityAnonymous, IdentityPsedonymous, IdentityNamed}
)

type ChildContractScope string
type Identity string

// SafeHarborAgreement represents the main struct for the agreement data.
type SafeHarborAgreement struct {
	Protocol            *firestore.DocumentRef `firestore:"protocol"`            // Reference to Protocol document
	RegistryTransaction string                 `firestore:"registryTransaction"` // Transaction hash
	RegistryChainId     string                 `firestore:"registryChainId"`
	AgreementAddress    string                 `firestore:"agreementAddress"`
	Entity              string                 `firestore:"entity"` // Creator EOA
	AgreementURI        string                 `firestore:"agreementURI"`
	AdoptionProposalURI string                 `firestore:"adoptionProposalURI"`
	ContactDetails      string                 `firestore:"contactDetails"`
	Chains              []Chain                `firestore:"chains"`
	BountyTerms         BountyTerms            `firestore:"bountyTerms"`
	CreatedAt           time.Time              `firestore:"createdAt"`
	CreatedBlock        int                    `firestore:"createdBlock"`
}

// Chain represents the chain details in the agreement.
type Chain struct {
	AssetRecoveryAddress string    `firestore:"assetRecoveryAddress"`
	ID                   int       `firestore:"id"`
	LastIndexedBlock     int       `firestore:"lastIndexedBlock"`
	Accounts             []Account `firestore:"accounts"`
}

// Account represents an account in the agreement.
type Account struct {
	Name               string             `firestore:"name"`
	Address            string             `firestore:"address"`
	ChildContractScope ChildContractScope `firestore:"childContractScope"`
	Children           []ChildAccount     `firestore:"children"`
	Signature          string             `firestore:"signature"`
}

// ChildAccount represents a child account within an account.
type ChildAccount struct {
	Name    string `firestore:"name"`
	Address string `firestore:"address"`
}

// BountyTerms represents the bounty terms in the agreement.
type BountyTerms struct {
	BountyPercentage      int      `firestore:"bountyPercentage"`
	BountyCapUSD          int      `firestore:"bountyCapUSD"`
	Retainable            bool     `firestore:"retainable"`
	Identity              Identity `firestore:"identity"`
	DiligenceRequirements string   `firestore:"diligenceRequirements"`
}
