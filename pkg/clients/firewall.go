package clients

import "github.com/Skylock-ai/Arianrhod/pkg/types/web3"

// Firewall represents a Skylock Firewall client
type Firewall struct {
	ProtocolName string
	AgreementURI string

	Chains         []ChainFW
	Accounts       []AccountFW
	ContactDetails string
	BountyTerms    BountyTermsFW

	CreatedAt int

	MockData bool
}

// Chain represents information about protected accounts by chain
type ChainFW struct {
	AssetRecoveryAddress web3.Address
	ChainID              int64
}

// Account represents information about the EOA or contract or child contracts to include.
type AccountFW struct {
	Name    string
	Address web3.Address

	RegisteredEvents   []EventFW
	ChildContractScope ChildContractScope
	ChainIDs           []int64

	LastQuery map[int64]int //ChainID:Block
}

type AddressLocationType string

const (
	// AddressLocationTypeTopic is when the address is in the topic
	AddressLocationTypeTopic AddressLocationType = "Topic"
	// AddressLocationTypeData is when the address is in the data
	AddressLocationTypeData AddressLocationType = "Data"
)

// Event represents an event to register for a contract.
type EventFW struct {
	Topic0              web3.Hash
	AddressLocationType AddressLocationType
	AddressLocation     int
}

type BountyTermsFW struct {
	BountyPercentage int
	BountyCapUSD     int
}
