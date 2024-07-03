package clients

import (
	"fmt"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum/rpc"
)

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

func (f *Firewall) ToMonitored(rpcClient *rpc.Client, apiKey string) (*Monitored, error) {
	monitored := &Monitored{
		MockData: f.MockData,
	}

	for _, account := range f.Accounts {
		accountM := AccountM{
			Name:    account.Name,
			Address: account.Address.String(),
			Chains:  make([]int, len(account.ChainIDs)),
		}

		for i, chainID := range account.ChainIDs {
			accountM.Chains[i] = int(chainID)
		}

		// Fill name for the account if it's not already set
		if accountM.Name == "" {
			err := fillName(&accountM, apiKey)
			if err != nil {
				return nil, fmt.Errorf("failed to fill name for account %s: %w", account.Address.String(), err)
			}
		}

		// Find children if necessary
		if account.ChildContractScope == ChildContractScopeAll || len(account.RegisteredEvents) > 0 {
			children, err := findAllChildren(
				account.Address,
				account.ChainIDs,
				account.ChildContractScope,
				account.RegisteredEvents,
				rpcClient,
				apiKey,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to find children for account %s: %w", account.Address.String(), err)
			}
			accountM.Children = children
		}

		monitored.Addresses = append(monitored.Addresses, accountM)
	}

	return monitored, nil
}
