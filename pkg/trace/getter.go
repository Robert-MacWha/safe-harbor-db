package trace

import (
	"fmt"
	"strconv"

	"github.com/Skylock-ai/Arianrhod/pkg/blackwood"
	"github.com/Skylock-ai/Arianrhod/pkg/kali"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"

	"github.com/ethereum/go-ethereum/rpc"
)

// GetTransaction wraps the eth_getTransactionByHash RPC method and
// returns a *web3.TxBody.
func GetTransaction(
	rpcClient *rpc.Client,
	txHash *web3.Hash,
) (txBody *web3.TxBody, err error) {
	err = rpcClient.Call(&txBody, "eth_getTransactionByHash", txHash.ToHex())
	if err != nil {
		return nil, fmt.Errorf("error getting transaction by hash: %v", err)
	}

	if txBody == nil {
		return nil, fmt.Errorf("transaction not found and error was nil")
	}

	return txBody, nil
}

// GetTransactionAsCall returns a call representation of the transaction
// provided. CallArgs can be traced using `TraceCallMany.`
func GetTransactionAsCall(
	rpcClient *rpc.Client,
	txHash *web3.Hash,
) (callArgs *TraceCallArgs, block string, err error) {
	txBody, err := GetTransaction(rpcClient, txHash)
	if err != nil {
		return nil, "", err
	}

	callArgs = &TraceCallArgs{
		CallArgs{
			From:  &txBody.From,
			To:    txBody.To,
			Value: &txBody.Value,
			Data:  &txBody.Input,
		},
		[]Tracer{},
	}

	var hexBlock string
	if txBody.Block == nil {
		hexBlock = "latest"
	} else {
		hexBlock = fmt.Sprintf("0x%x", int(txBody.Block.Int64())-1)
	}

	return callArgs, hexBlock, nil
}

// GetResponseContractCallArgs returns the []TraceCallArgs used to trace
// deploying & executing a response smart contract.
func GetResponseContractCallArgs(
	bw *blackwood.Blackwood,
	rpcClient *rpc.Client,
	txHash *web3.Hash,
	deployer *web3.Address,
	beneficiary *web3.Address,
	tracers []Tracer,
) (traceCallArgs []TraceCallArgs, contractAddress *web3.Address, err error) {
	//* Get the response contract
	resContractBytecode, resContractAddress, err := GetResponseContract(
		bw,
		rpcClient,
		txHash,
		deployer,
		beneficiary,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("error generating response: %v", err)
	}

	//* Build response calls
	traceCallArgs = []TraceCallArgs{
		// Deploy response contract
		{
			CallArgs{
				From: deployer,
				Data: web3.NewBytes(resContractBytecode),
			},
			[]Tracer{},
		},
		// Call response contract
		{
			CallArgs{
				From: deployer,
				To:   resContractAddress,
			},
			tracers,
		},
	}

	return traceCallArgs, resContractAddress, nil
}

// GetNonce returns the current nonce for the provided address.
func GetNonce(rpcClient *rpc.Client, address web3.Address) (nonce int, err error) {
	var hexNonce string
	err = rpcClient.Call(
		&hexNonce,
		"eth_getTransactionCount",
		address.ToHex(),
		"latest",
	)
	if err != nil {
		return 0, fmt.Errorf("error fetching deployer nonce: %v", err)
	}

	nonce64, err := strconv.ParseInt(hexNonce[2:], 16, 64)

	return int(nonce64), err
}

// GetResponseContract traces a transaction and returns the initialization bytecode
// for a response contract deployed from `deployer`.
func GetResponseContract(
	bw *blackwood.Blackwood,
	rpcClient *rpc.Client,
	txHash *web3.Hash,
	deployer *web3.Address,
	beneficiary *web3.Address,
) (bytecode []byte, contractAddress *web3.Address, err error) {
	txBody, err := GetTransaction(rpcClient, txHash)
	if err != nil {
		return []byte{}, nil, fmt.Errorf("error getting transaction by hash: %v", err)
	}
	traceResult, err := TraceFromTxBody(rpcClient, *txBody)
	if err != nil {
		return []byte{}, nil, fmt.Errorf("error tracing transaction: %v", err)
	}
	rootCall, err := ParityToCallRoot(traceResult.Trace)
	if err != nil {
		return []byte{}, nil, fmt.Errorf("error converting trace to call root: %v", err)
	}

	//* Generate response contract blueprint
	deployerNonce, err := GetNonce(rpcClient, *deployer)
	if err != nil {
		return []byte{}, nil, err
	}

	contractAddress = web3.GetContractAddress(deployer, int(deployerNonce))

	blueprint, err := bw.CalldataBlueprint(
		*deployer,
		*contractAddress,
		*beneficiary,
		rootCall.From,
		rootCall.To,
		make(map[web3.Address]web3.Address),
		rootCall,
	)

	if err != nil {
		return []byte{}, nil, fmt.Errorf("error generating contract blueprint: %v", err)
	}

	bytecode, err = kali.CompileDefault(blueprint, contractAddress)
	if err != nil {
		return []byte{}, nil, fmt.Errorf("error compiling blueprint: %v", err)
	}

	return bytecode, contractAddress, nil
}
