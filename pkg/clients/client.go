package clients

import (
	"SHDB/pkg/etherscan"
	"SHDB/pkg/trace"
	"fmt"

	"github.com/Skylock-ai/Arianrhod/pkg/source"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client represents a client that can be queried for addresses.
type Client struct {
	ProtocolName string
	AgreementURI string

	Chains         []Chain
	Accounts       []Account
	ContactDetails string
	BountyTerms    BountyTerms

	Website        string
	Date           string
	Icon           string
	TVL            int
	Category       string
	FirewallClient bool
}

// Chain represents information about protected accounts by chain
type Chain struct {
	AssetRecoveryAddress web3.Address
	ChainID              int64
}

// Account represents information about the EOA or contract or child contracts to include.
type Account struct {
	Name    string
	Address web3.Address

	RegisteredEvents   []Event
	ChildContractScope ChildContractScope
	Children           []Account

	SignatureValid bool
	ChainIDs       []int64

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
type Event struct {
	Topic0              web3.Hash
	AddressLocationType AddressLocationType
	AddressLocation     int
}

type IdentityVerification string

const (
	// Retainable is when the identity is retainable
	Retainable IdentityVerification = "Retainable"
	// Immunefi is when the identity is verified through Immunefi
	Immunefi IdentityVerification = "Immunefi"
	// Bugcrowd is when the identity is verified through Bugcrowd
	Bugcrowd IdentityVerification = "Bugcrowd"
	// Hackerone is when the identity is verified through Hackerone
	Hackerone IdentityVerification = "Hackerone"
)

type BountyTerms struct {
	BountyPercentage int
	BountyCapUSD     int
	Retainable       IdentityVerification
}

// ChildContractScope represents the scope of child contracts to include.
type ChildContractScope string

const (
	// ChildContractScopeNone is when no child contracts are included
	ChildContractScopeNone ChildContractScope = "None"
	// ChildContractScopeAll is when all child contracts, both existing and new, are included
	ChildContractScopeAll ChildContractScope = "All"
)

func (c *Client) GetClientInformation(
	rpcClient *rpc.Client,
	apiKey string,
) error {
	for i := range c.Accounts {
		account := &c.Accounts[i]
		if err := account.findAllChildren(rpcClient, apiKey); err != nil {
			return fmt.Errorf("failed to find all children: %w", err)
		}

		if err := account.fillNamesAndChildren(apiKey); err != nil {
			return fmt.Errorf("failed to fill names and children: %w", err)
		}
	}

	_, url, tvl, category, _, logo, err := GetProtocolInfo(c.ProtocolName)
	if err != nil {
		// This is to be expected as the naming on Defillama isn't consistent
		fmt.Println("Error, could not fetch protocol info for ", c.ProtocolName, err)
		return nil
	}

	c.Website = url
	c.TVL = int(tvl)
	c.Category = category
	c.Icon = logo

	return nil
}

func (a *Account) findAllChildren(
	rpcClient *rpc.Client,
	apiKey string,
) error {
	childrenMap := make(map[web3.Address]map[int64]bool)

	for _, chainID := range a.ChainIDs {
		startBlock, ok := a.LastQuery[chainID]
		if !ok {
			startBlock = 0
		}

		for _, event := range a.RegisteredEvents {

			addressesFromEvent, err := getAllEventAddresses(
				chainID,
				apiKey,
				a.Address,
				event.Topic0,
				event.AddressLocationType,
				event.AddressLocation,
				startBlock,
			)
			if err != nil {
				return fmt.Errorf("failed to get addresses from events: %w", err)
			}
			for _, addr := range addressesFromEvent {
				if chainIDs, exists := childrenMap[addr]; exists {
					// If the address is in the map, add the current chain ID
					chainIDs[chainID] = true
				} else {
					// If the address isn't in the map, add it with the current chain ID
					childrenMap[addr] = map[int64]bool{chainID: true}
				}
			}
		}

		// Handle fetching of child contract addresses
		if a.ChildContractScope == ChildContractScopeAll {

			childAddresses, err := getAllSubContractAddresses(
				rpcClient,
				chainID,
				apiKey,
				a.Address,
				startBlock,
			)
			if err != nil {
				return fmt.Errorf("failed to get child contract addresses: %w", err)
			}
			for _, childAddr := range childAddresses {
				if chainIDs, exists := childrenMap[childAddr]; exists {
					// If the address is in the map, add the current chain ID
					chainIDs[chainID] = true
				} else {
					// If the address isn't in the map, add it with the current chain ID
					childrenMap[childAddr] = map[int64]bool{chainID: true}
				}
			}
		}
		// ! FIX ME: Temprorary cause in a rush to finish but:
		// This RPC is for 1 single chain, we need different RPCS for each chain :skull:
		rpcSource := source.RPCBackend{Client: rpcClient}
		blockNumber, err := rpcSource.GetBlockNumber()
		if err != nil {
			return fmt.Errorf("failed to get block number: %w", err)
		}

		a.LastQuery[chainID] = int(blockNumber)
	}

	for childAddr, chainIDs := range childrenMap {
		chainIDsSlice := []int64{}
		for chainID := range chainIDs {
			chainIDsSlice = append(chainIDsSlice, chainID)
		}

		a.Children = append(a.Children, Account{
			Address:  childAddr,
			ChainIDs: chainIDsSlice,
		})

		// Recursively find all children
		if a.ChildContractScope == ChildContractScopeAll {
			err := a.Children[len(a.Children)-1].findAllChildren(rpcClient, apiKey)
			if err != nil {
				return fmt.Errorf("failed to find all children: %w", err)
			}
		}
	}

	return nil
}

func (a *Account) fillNames(
	apiKey string,
) error {
	// errors := []error{}
	for _, chainID := range a.ChainIDs {
		result, err := etherscan.FetchSourceCode(chainID, apiKey, a.Address)
		if err != nil {
			a.Name = ""
			// if err
			// errors = append(errors, fmt.Errorf("failed to fetch source code: %w", err))
			continue
		}

		a.Name = result.ContractName
		fmt.Println("Name: ", a.Name)
		return nil
	}

	return nil
	// return errors[0]
}

func (a *Account) fillNamesAndChildren(apiKey string) error {
	err := a.fillNames(apiKey)
	if err != nil {
		return fmt.Errorf("failed to fill names: %w", err)
	}

	// for _, child := range a.Children {
	for i := range a.Children {
		child := &a.Children[i]
		err := child.fillNamesAndChildren(apiKey)
		if err != nil {
			return fmt.Errorf("failed to fill names: %w", err)
		}
	}

	return nil
}

// getAllEventAddresses fetches logs and extracts addresses from the log data.
func getAllEventAddresses(
	chainID int64,
	apiKey string,
	address web3.Address,
	topic0 web3.Hash,
	addressLocationType AddressLocationType,
	addressLocation int,
	startBlock int,
) ([]web3.Address, error) {
	logs, err := etherscan.FetchLogs(chainID, apiKey, address, topic0, startBlock)
	if err != nil {
		return nil, fmt.Errorf("error fetching logs: %w", err)
	}

	var addresses []web3.Address
	for _, log := range logs {
		var address *web3.Address

		if addressLocationType == AddressLocationTypeData {

			if len(log.Data) < (addressLocation + 20 - 1) {
				continue
			}

			extractedAddress := log.Data[addressLocation : addressLocation+20]
			address, err = web3.BytesToAddress(extractedAddress)
			if err != nil {
				return nil, fmt.Errorf("error converting bytes to address: %w", err)
			}
		} else if addressLocationType == AddressLocationTypeTopic {
			if len(log.Topics) < (addressLocation + 1) {
				continue
			}

			extractedAddress := log.Topics[addressLocation][12:]
			address, err = web3.BytesToAddress(extractedAddress)
			if err != nil {
				return nil, fmt.Errorf("error converting bytes to address: %w", err)
			}
		}

		addresses = append(addresses, *address)
	}

	uniqueAddresses := make(map[web3.Address]bool)
	for _, address := range addresses {
		uniqueAddresses[address] = true
	}

	addresses = []web3.Address{}
	for address := range uniqueAddresses {
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// getAllSubContractAddresses fetches both regular and internal transactions
// and gets the contract addresses.
func getAllSubContractAddresses(
	rpcClient *rpc.Client,
	chainID int64,
	apiKey string,
	address web3.Address,
	startBlock int,
) ([]web3.Address, error) {
	var addresses []web3.Address

	// Fetch regular transactions
	regularTransactions, err := etherscan.FetchRegularTransactions(chainID, apiKey, address, startBlock)
	if err != nil {
		return nil, fmt.Errorf("error fetching regular transactions: %w", err)
	}

	// Fetch internal transactions
	internalTransactions, err := etherscan.FetchInternalTransactions(chainID, apiKey, address, startBlock)
	if err != nil {
		return nil, fmt.Errorf("error fetching internal transactions: %w", err)
	}

	var txHashes []web3.Hash
	for _, tx := range regularTransactions {
		txHashes = append(txHashes, tx.Hash)
	}
	for _, tx := range internalTransactions {
		txHashes = append(txHashes, tx.Hash)
	}

	var traces []trace.Tracer
	traces = append(traces, trace.CallTracer)

	for _, txHash := range txHashes {
		traceResults, err := trace.TraceTransaction(rpcClient, &txHash, traces)

		if err != nil {
			return nil, fmt.Errorf("error tracing transaction: %w", err)
		}

		traceCall := traceResults[0].Trace

		for _, call := range traceCall {
			if call.Action.Create == nil {
				continue
			}

			if call.Action.Create.From != address {
				continue
			}

			if call.Result.Address == nil {
				continue
			}

			addresses = append(addresses, *call.Result.Address)
		}
	}

	uniqueAddresses := make(map[web3.Address]bool)
	for _, address := range addresses {
		uniqueAddresses[address] = true
	}

	addresses = []web3.Address{}
	for address := range uniqueAddresses {
		addresses = append(addresses, address)
	}

	return addresses, nil
}
