package process

import (
	"SHDB/pkg/etherscan"
	"SHDB/pkg/firewall"
	"SHDB/pkg/trace"
	"fmt"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

func GetNameOrEmpty(address web3.Address, chainID int64, apiKey string) string {
	result, err := etherscan.FetchSourceCode(chainID, apiKey, address)
	if err != nil {
		log.Error("Failed to fetch source code", "error", err)
		return ""
	}

	return result.ContractName
}

// GetAllEventAddresses fetches logs and extracts addresses from the log data.
func GetAllEventAddresses(
	chainID int64,
	apiKey string,
	address web3.Address,
	topic0 web3.Hash,
	addressLocationType firewall.AddressLocationType,
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

		if addressLocationType == firewall.AddressLocationTypeData {

			if len(log.Data) < (addressLocation + 20 - 1) {
				continue
			}

			extractedAddress := log.Data[addressLocation : addressLocation+20]
			address, err = web3.BytesToAddress(extractedAddress)
			if err != nil {
				return nil, fmt.Errorf("error converting bytes to address: %w", err)
			}
		} else if addressLocationType == firewall.AddressLocationTypeTopic {
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

// GetAllSubContractAddresses fetches both regular and internal transactions
// and gets the contract addresses.
func GetAllSubContractAddresses(
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

	// To avoid tracing the same transactions multiple times we cache the function signatures
	noSubContractFunctionSignatures := make(map[string]bool)

	for _, txHash := range txHashes {
		var txBody web3.TxBody
		err = rpcClient.Call(&txBody, "eth_getTransactionByHash", txHash.ToHex())
		if err != nil {
			return nil, fmt.Errorf("error getting transaction by hash: %v", err)
		}

		var signature string

		if len(txBody.Input.String()) == 0 {
			continue
		} else if len(txBody.Input.String()) < 10 {
			signature = txBody.Input.String()
		} else {
			signature = txBody.Input.String()[:10]
		}

		if noSubContractFunctionSignatures[signature] {
			continue
		}

		traceResults, err := trace.TraceTransaction(rpcClient, &txHash, traces)

		if err != nil {
			return nil, fmt.Errorf("error tracing transaction: %w", err)
		}

		traceCall := traceResults[0].Trace

		count := 0
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
			count++
		}

		if count == 0 {
			noSubContractFunctionSignatures[signature] = true
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
