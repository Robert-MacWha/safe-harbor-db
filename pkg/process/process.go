package process

import (
	"SHDB/pkg/etherscan"
	"SHDB/pkg/firewall"
	"context"
	"fmt"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum/rpc"
)

// Define a call object that matches the structure
type call struct {
	Type    string        `json:"type"`
	From    web3.Address  `json:"from"`
	To      *web3.Address `json:"to,omitempty"`
	Gas     web3.BigInt   `json:"gas"`
	GasUsed web3.BigInt   `json:"gasUsed"`
	Input   string        `json:"input"`
	Output  string        `json:"output"`
	Value   web3.BigInt   `json:"value"`
	Calls   []call        `json:"calls,omitempty"` // Subcalls are also of type Call
}

// debugResult holds the root call
type debugResult struct {
	Calls []call `json:"calls"`
}

func GetNameOrEmpty(address web3.Address, chainID int64, apiKey string) string {
	result, err := etherscan.FetchSourceCode(chainID, apiKey, address)
	if err != nil {
		fmt.Println("Failed to fetch source code", "error", err)
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

	// To avoid tracing the same transactions multiple times we cache the function signatures
	noSubContractFunctionSignatures := make(map[string]bool)

	// Prepare the tracer parameters
	params := map[string]interface{}{
		"tracer": "callTracer",
		"tracerConfig": map[string]interface{}{
			"onlyTopCall": false, // Change to true if you want only the top-level call
			"timeout":     "60s", // Timeout configuration
		},
	}

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

		var result debugResult
		// Perform the RPC call to debug_traceTransaction
		err = rpcClient.CallContext(context.Background(), &result, "debug_traceTransaction", txHash, params)
		if err != nil {
			fmt.Println("Failed to trace transaction", "error", err)
			return nil, fmt.Errorf("error tracing transaction: %w", err)
		}

		count := 0

		for _, call := range result.Calls {
			if call.Type == "CREATE" && call.From == address && call.To != nil {
				count++
				addresses = append(addresses, *call.To)
			}
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
