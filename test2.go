// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
// 	"github.com/ethereum/go-ethereum/rpc"
// )

// // Define a Call object that matches the structure
// type Call struct {
// 	Type    string        `json:"type"`
// 	From    web3.Address  `json:"from"`
// 	To      *web3.Address `json:"to,omitempty"`
// 	Gas     web3.BigInt   `json:"gas"`
// 	GasUsed web3.BigInt   `json:"gasUsed"`
// 	Input   string        `json:"input"`
// 	Output  string        `json:"output"`
// 	Value   web3.BigInt   `json:"value"`
// 	Calls   []Call        `json:"calls,omitempty"` // Subcalls are also of type Call
// }

// // Result holds the root call
// type Result struct {
// 	Calls []Call `json:"calls"`
// }

// func main() {
// 	fmt.Println("hi")
// 	// Connect to the Ethereum node (Assuming Geth is running with --rpcapi debug)
// 	client, err := rpc.DialContext(context.Background(), "https://frequent-sleek-frost.ethereum-holesky.quiknode.pro/43e48929064f5df35ecab07516a344defb2e0e69")
// 	if err != nil {
// 		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
// 	}

// 	fmt.Println("hi")

// 	// Replace this with the transaction hash you want to trace
// 	txHash := "0x3d4e5452c80286519be1709340bae6bbb3e6415e3e6dc1be46ede588040ed4bb"

// 	// Prepare the tracer parameters
// 	params := map[string]interface{}{
// 		"tracer": "callTracer",
// 		"tracerConfig": map[string]interface{}{
// 			"onlyTopCall": false, // Change to true if you want only the top-level call
// 			"timeout":     "30s", // Timeout configuration
// 		},
// 	}

// 	// Define a variable to hold the result
// 	var result Result

// 	// Perform the RPC call to debug_traceTransaction
// 	err = client.CallContext(context.Background(), &result, "debug_traceTransaction", txHash, params)
// 	if err != nil {
// 		log.Fatalf("Failed to trace transaction: %v", err)
// 	}

// 	for _, call := range result.Calls {
// 		if call.Type == "CREATE" && call.To != nil {
// 			fmt.Println("Contract created at address:", call)

// 		}
// 	}

// 	// Output the result
// 	fmt.Printf("Transaction Trace: %+v\n", result)
// }
