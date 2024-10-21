package clients

import (
	"SHDB/pkg/etherscan"
	"SHDB/pkg/trace"
	"fmt"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum/rpc"
)

func findAllChildren(
	address web3.Address,
	chainIDs []int64,
	childContractScope ChildContractScope,
	events []EventFW,
	rpcClient *rpc.Client,
	apiKey string,
) ([]AccountM, error) {
	var children []AccountM
	childrenMap := make(map[web3.Address]map[int64]bool)

	for _, chainID := range chainIDs {
		// Handle fetching of child contract addresses from events
		for _, event := range events {
			addressesFromEvent, err := getAllEventAddresses(
				chainID,
				apiKey,
				address,
				event.Topic0,
				event.AddressLocationType,
				event.AddressLocation,
				0, // Start from block 0
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get addresses from events: %w", err)
			}
			for _, addr := range addressesFromEvent {
				if chainIDs, exists := childrenMap[addr]; exists {
					chainIDs[chainID] = true
				} else {
					childrenMap[addr] = map[int64]bool{chainID: true}
				}
			}
		}

		// Handle fetching of child contract addresses from contract creation
		if childContractScope == ChildContractScopeAll {
			childAddresses, err := getAllSubContractAddresses(
				rpcClient,
				chainID,
				apiKey,
				address,
				0, // Start from block 0
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get child contract addresses: %w", err)
			}
			for _, childAddr := range childAddresses {
				if chainIDs, exists := childrenMap[childAddr]; exists {
					chainIDs[chainID] = true
				} else {
					childrenMap[childAddr] = map[int64]bool{chainID: true}
				}
			}
		}
	}

	for childAddr, chainIDsMap := range childrenMap {
		chainIDsSlice := []int{}
		for chainID := range chainIDsMap {
			chainIDsSlice = append(chainIDsSlice, int(chainID))
		}

		childAccount := AccountM{
			Address: childAddr.String(),
			Chains:  chainIDsSlice,
		}

		// Fill name for the child account
		err := fillName(&childAccount, apiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to fill name for child account %s: %w", childAddr.String(), err)
		}

		// // Recursively find all children
		// if childContractScope == ChildContractScopeAll {
		// 	chainIDsList := make([]int64, len(chainIDsSlice))
		// 	for i, chainID := range chainIDsSlice {
		// 		chainIDsList[i] = int64(chainID)
		// 	}
		// 	childrenOfChild, err := findAllChildren(childAddr, chainIDsList, childContractScope, events, rpcClient, apiKey)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("failed to find all children: %w", err)
		// 	}
		// 	childAccount.Children = childrenOfChild
		// }

		children = append(children, childAccount)
	}

	return children, nil
}

func fillName(account *AccountM, apiKey string) error {
	for _, chainID := range account.Chains {
		address, err := web3.HexToAddress(account.Address)
		if err != nil {
			return fmt.Errorf("failed to convert address: %w", err)
		}
		result, err := etherscan.FetchSourceCode(int64(chainID), apiKey, *address)
		if err != nil {
			account.Name = ""
			continue
		}

		account.Name = result.ContractName
		fmt.Println("Name: ", account.Name)
		return nil
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
