package firewall

import (
	"time"

	"SHDB/pkg/web3"

	"cloud.google.com/go/firestore"
)

type ChildContractScope string

type AddressLocationType string

const (
	ChildContractScopeNone ChildContractScope = "None"
	ChildContractScopeAll  ChildContractScope = "All"
)

const (
	// AddressLocationTypeTopic is when the address is in the topic
	AddressLocationTypeTopic AddressLocationType = "Topic"
	// AddressLocationTypeData is when the address is in the data
	AddressLocationTypeData AddressLocationType = "Data"
)

// FirewallAgreement represents the main struct for the firewall agreement data.
type FirewallAgreement struct {
	Protocol    *firestore.DocumentRef `firestore:"protocol"` // Reference to Protocol document
	Chains      []Chain                `firestore:"chains"`
	BountyTerms BountyTerms            `firestore:"bountyTerms"`
	CreatedAt   time.Time              `firestore:"createdAt"`
}

// Chain represents a chain in the firewall agreement.
type Chain struct {
	AssetRecoveryAddress web3.Address `firestore:"assetRecoveryAddress"`
	ChainID              int          `firestore:"chainID"`
	Accounts             []Account    `firestore:"accounts"`
}

// Account represents an account in the firewall agreement.
type Account struct {
	Name               string             `firestore:"name"`
	Address            web3.Address       `firestore:"address"`
	ChildContractScope ChildContractScope `firestore:"childContractScope"`
	Children           []ChildAccount     `firestore:"children"`
	RegisteredEvents   []RegisteredEvent  `firestore:"registeredEvents"`
}

// ChildAccount represents a child account within an account.
type ChildAccount struct {
	Name    string `firestore:"name"`
	Address string `firestore:"address"`
}

// RegisteredEvent represents a registered event in the firewall agreement.
type RegisteredEvent struct {
	Topic0              web3.Hash           `firestore:"topic0"`
	AddressLocationType AddressLocationType `firestore:"addressLocationType"` // "Topic" or "Data"
	AddressLocation     int                 `firestore:"addressLocation"`
}

// BountyTerms represents the bounty terms in the firewall agreement.
type BountyTerms struct {
	BountyPercentage int `firestore:"bountyPercentage"`
	BountyCapUSD     int `firestore:"bountyCapUSD"`
}
