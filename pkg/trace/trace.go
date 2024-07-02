package trace

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"

	"SHDB/pkg/trace/traceresult"

	"github.com/Skylock-ai/Arianrhod/pkg/blackwood"
	"github.com/Skylock-ai/Arianrhod/pkg/types/tracetypes"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"

	"github.com/ethereum/go-ethereum/rpc"
)

// TraceTransaction traces the provided transaction as a call at the end
// of the block before the txHash was originally deployed.  Returns the
// traceResultfor the selected `tracers`
func TraceTransaction(
	rpcClient *rpc.Client,
	txHash *web3.Hash,
	tracers []Tracer,
) (traceResult []*traceresult.TraceResult, err error) {
	callArgs, hexBlock, err := GetTransactionAsCall(
		rpcClient,
		txHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	traceResults, err := TraceCallMany(
		rpcClient,
		[]TraceCallArgs{*callArgs},
		hexBlock,
		tracers,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to trace transaction: %v", err)
	}

	return traceResults, nil
}

// TraceCall is a fully-managed interface used to trace a transaction, either
// the original or the replay.
func TraceCall(
	bw *blackwood.Blackwood, txHashStr string, rpcEndpoint string, replay bool,
) (result *traceresult.TraceResult, err error) {
	txHash, err := web3.HexToHash(txHashStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing txHashString: %v", err)
	}

	rpcClient, err := rpc.Dial(rpcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error dialing rpc endpoint: %v", err)
	}

	tracers := []Tracer{CallTracer}

	//* Original trace
	if !replay {
		traceResults, err := TraceTransaction(rpcClient, txHash, tracers)
		if err != nil {
			return nil, fmt.Errorf("error tracing transaction: %v", err)
		}
		traceResult := traceResults[0]

		return traceResult, nil
	}

	traceResult, _, err := TraceReplayTransaction(
		bw,
		rpcClient, txHash, web3.DefaultDeployer,
		web3.DefaultBeneficiary, tracers, web3.DefaultBeneficiary,
	)
	if err != nil {
		return nil, fmt.Errorf("error tracing transaction: %v", err)
	}

	return traceResult, nil
}

// TraceVM is a fully-managed interface used to trace a transaction's vmtrace,
// either the original or the replay.
func TraceVM(
	bw *blackwood.Blackwood, txHashStr string, rpcEndpoint string, replay bool,
) (result *traceresult.TraceResult, err error) {
	txHash, err := web3.HexToHash(txHashStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing txHashString: %v", err)
	}

	rpcClient, err := rpc.Dial(rpcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error dialing rpc endpoint: %v", err)
	}

	tracers := []Tracer{VMTracer}

	//* Original trace
	if !replay {
		traceResults, err := TraceTransaction(rpcClient, txHash, tracers)
		if err != nil {
			return nil, fmt.Errorf("error tracing transaction: %v", err)
		}
		traceResult := traceResults[0]

		return traceResult, nil
	}

	//* Replay trace

	traceResult, _, err := TraceReplayTransaction(
		bw, rpcClient, txHash, web3.DefaultDeployer, web3.DefaultDeployer,
		tracers, web3.DefaultDeployer,
	)
	if err != nil {
		return nil, fmt.Errorf("error tracing transaction: %v", err)
	}

	return traceResult, nil
}

// TraceReplayTransaction attempts to replay a transaction at the end of the
// prior block using Blackwood & Kali to create the response contract.
//
// Returns the traceresult.TraceResult for the call made to the deployed
// response contract from `deployer` traced using the provided tracers, and
// the map of tokens received by the `checkAddress`.
func TraceReplayTransaction(
	bw *blackwood.Blackwood,
	rpcClient *rpc.Client,
	txHash *web3.Hash,
	deployer *web3.Address,
	beneficiary *web3.Address,
	tracers []Tracer,
	checkAddress *web3.Address,
) (traceResult *traceresult.TraceResult, ERCResult map[web3.Address]web3.BigInt, err error) {
	txBody, err := GetTransaction(rpcClient, txHash)
	if err != nil {
		return nil, nil, err
	}

	var hexResCallBlock string
	if txBody.Block == nil {
		hexResCallBlock = "latest"
	} else {
		resCallBlock := int(txBody.Block.Int64() - 1)
		hexResCallBlock = fmt.Sprintf("0x%x", resCallBlock)
	}

	resCallArgs, _, err := GetResponseContractCallArgs(
		bw,
		rpcClient,
		txHash,
		deployer,
		beneficiary,
		tracers,
	)
	if err != nil {
		return nil,
			nil,
			fmt.Errorf("error getting response contract call args: %v", err)
	}

	callTraces, err := TraceCallMany(
		rpcClient,
		resCallArgs,
		hexResCallBlock,
		nil,
	)

	if err != nil {
		return nil,
			nil,
			fmt.Errorf("error tracing response contract: %v", err)
	}
	replayCallTrace := callTraces[len(callTraces)-1]

	if len(replayCallTrace.Trace) != 0 && replayCallTrace.Trace[0].Error != "" {
		return replayCallTrace,
			nil,
			fmt.Errorf(
				"error tracing response contract: %v",
				replayCallTrace.Trace[0].Error,
			)
	}

	ERCResult, err = TraceReplayTransactionERCResult(
		rpcClient,
		resCallArgs,
		hexResCallBlock,
		checkAddress,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting ERCResults: %v", err)
	}

	return replayCallTrace, ERCResult, nil
}

// TraceCallMany traces the provided calls using the trace_callMany RPC method
// and returns the unmarshalled results.
//
// If tracers []Tracer is provided, this tracer will be used in place of the
// tracers attached to calls.
// Block should be the hex of the block number (with 0x) or "latest"
func TraceCallMany(
	rpcClient *rpc.Client,
	calls []TraceCallArgs,
	block string,
	tracers []Tracer,
) (traceResult []*traceresult.TraceResult, err error) {
	if tracers != nil {
		for i := range calls {
			calls[i].Traces = tracers
		}
	}

	err = rpcClient.Call(&traceResult, "trace_callMany", calls, block)
	if err != nil {
		return nil, fmt.Errorf("error tracing calls: %v", err)
	}

	return traceResult, nil
}

// TraceFromTxBody traces the provided transaction body and returns the
// unmarshalled results.
func TraceFromTxBody(
	rpcClient *rpc.Client,
	txBody web3.TxBody,
) (traceResult *traceresult.TraceResult, err error) {
	tracers := []Tracer{CallTracer}
	callArgs := &TraceCallArgs{
		Call: CallArgs{
			From:  &txBody.From,
			To:    txBody.To,
			Value: &txBody.Value,
			Data:  &txBody.Input,
		},
		Traces: []Tracer{},
	}
	var hexBlock string
	if txBody.Block == nil {
		hexBlock = "latest"
	} else {
		block := int(txBody.Block.Int64() - 1)
		hexBlock = fmt.Sprintf("0x%x", block)
	}

	traceResults, err := TraceCallMany(
		rpcClient,
		[]TraceCallArgs{*callArgs},
		hexBlock,
		tracers,
	)

	if err != nil {
		return nil, fmt.Errorf("error tracing transaction: %v", err)
	}

	if len(traceResults) == 0 {
		return nil, fmt.Errorf("no trace results returned")
	}

	traceResult = traceResults[0]

	return traceResult, err
}

// TraceERC20Balance traces the provided calls and returns the balance for some
// `targetAddress` across all ERC-20 tokens it interacted with, and Ethereum.
//
// revive:disable:cyclomatic High Complexity due to nested for loops, actual
// calculations are not complex.
func TraceERC20Balance(
	rpcClient *rpc.Client,
	calls []TraceCallArgs,
	block string,
	targetAddress web3.Address,
) (balances map[web3.Address]web3.BigInt, err error) {
	//* Trace the calls & get the addresses of all relevant ERC20s
	callTraces, err := TraceCallMany(rpcClient, calls, block, []Tracer{CallTracer})
	if err != nil {
		return nil, fmt.Errorf("error calling callTrace: %v", err)
	}

	last := len(callTraces) - 1
	erc20s := getERC20AddressesFromCallTrace(
		callTraces[last].Trace, targetAddress,
	)

	//* Get ERC20 Balances
	ercBalances, err := GetERC20BalanceForAddress(
		rpcClient, calls, targetAddress, block, erc20s,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching erc20 balances for address: %v", err)
	}

	//* Get ETH Balance
	ethBalance := GetETHBalanceFromTrace(callTraces[last].Trace, targetAddress)

	ercBalances[*web3.ETHAddress] = ethBalance

	return ercBalances, nil
}

// revive:enable:cyclomatic

// TraceReplayTransactionERCResult traces a replay transaction and returns the
// ERCResult erc20 map.
//
// TODO: Only check the ERCResult for the deployer address once Blackwood & Kali
// support Balance Exfiltration
func TraceReplayTransactionERCResult(
	rpcClient *rpc.Client,
	resCallArgs []TraceCallArgs,
	resCallBlock string,
	checkAddress *web3.Address,
) (ERCResult map[web3.Address]web3.BigInt, err error) {
	ERCResult, err = TraceERC20Balance(
		rpcClient,
		resCallArgs,
		resCallBlock,
		*checkAddress,
	)

	if err != nil {
		return nil, fmt.Errorf("error tracing contract erc20 balances: %v", err)
	}

	return ERCResult, nil
}

// getERC20AddressesFromCallTrace gets the set of all ERC20 addresses who
// experienced a state change for the target address within a given []CallTrace.
//
//revive:disable:cyclomatic Many multi-param checks, but all clearly documented
func getERC20AddressesFromCallTrace(
	call []traceresult.CallTrace,
	targetAddr web3.Address,
) (ercs []web3.Address) {
	const (
		transferFuncSig     = "a9059cbb"
		transferfromFuncSig = "23b872dd"
		mintFuncSig         = "40c10f19"
		burnFuncSig         = "42966c68"
		transferFuncLen     = 64
		transferfromFuncLen = 96
		mintFuncLen         = 64
		burnFuncLen         = 32
		functionSigLen      = 4
	)

	ercMap := map[web3.Address]bool{}

	// loop over each trace within the call
	for _, trace := range call {
		if !trace.Type.IsCall() {
			continue
		}

		if len(trace.Action.Call.Input) < functionSigLen {
			continue
		}

		fSig := hex.EncodeToString(trace.Action.Call.Input[:functionSigLen])
		fData := trace.Action.Call.Input[functionSigLen:]

		// check if the trace was a transfer / transferfrom to / from our address.
		// if so, record the erc20's address.
		switch fSig {
		case transferFuncSig:

			toAddr, _ := web3.BytesToAddress(fData[12:32])

			targetAssociated := *toAddr == targetAddr || trace.Action.Call.From == targetAddr
			if targetAssociated {
				ercMap[trace.Action.Call.To] = true
			}

		case transferfromFuncSig:
			spenderAddr, _ := web3.BytesToAddress(fData[12:32])
			toAddr, _ := web3.BytesToAddress(fData[44:64])

			interactedWithTarget := *spenderAddr == targetAddr || *toAddr == targetAddr
			if interactedWithTarget {
				ercMap[trace.Action.Call.To] = true
			}

		case mintFuncSig:
			receiver, _ := web3.BytesToAddress(fData[12 : 12+20])

			interactedWithTarget := *receiver == targetAddr ||
				trace.Action.Call.From == targetAddr

			if interactedWithTarget {
				ercMap[trace.Action.Call.To] = true
			}

		case burnFuncSig:
			if trace.Action.Call.From == targetAddr {
				ercMap[trace.Action.Call.To] = true
			}
		}

	}

	ercs = make([]web3.Address, 0, len(ercMap))
	for k := range ercMap {
		if k == call[0].Action.Call.To {
			continue
		}
		ercs = append(ercs, k)
	}

	// always add weth, usdc, usdt, dai
	ercs = append(ercs, *web3.WETHAddress)
	ercs = append(ercs, *web3.USDCAddress)
	ercs = append(ercs, *web3.USDTAddress)
	ercs = append(ercs, *web3.DAIAddress)
	return ercs
}

//revive:enable:cyclomatic

// GetERC20BalanceForAddress returns the erc20 balances for some target
// address after a set of calls have been traced.
func GetERC20BalanceForAddress(
	rpcClient *rpc.Client,
	preCalls []TraceCallArgs,
	targetAddress web3.Address,
	block string,
	erc20s []web3.Address,
) (ercBalances map[web3.Address]web3.BigInt, err error) {
	//* Add the erc20 balance calls to the preCalls
	ercCalls, err := AddERC20BalanceCalls(erc20s, preCalls, targetAddress)
	if err != nil {
		return nil, fmt.Errorf("error adding ERC-20 balance calls: %v", err)
	}

	ercCallTrace, err := TraceCallMany(rpcClient, ercCalls, block, []Tracer{CallTracer})
	if err != nil {
		return nil, fmt.Errorf("error tracing ERC-20 balance calls: %v", err)
	}

	//* Get the balance mapping
	ercBalanceMap := ProcessERC20BalanceMap(ercCallTrace, len(preCalls))

	return ercBalanceMap, nil
}

// ProcessERC20BalanceMap processes the provided erc20 balance map and returns
func ProcessERC20BalanceMap(
	ercCallTrace []*traceresult.TraceResult,
	numberPreCalls int,
) map[web3.Address]web3.BigInt {
	//* Get the balance mapping
	ercBalanceMap := map[web3.Address]web3.BigInt{}
	for _, call := range ercCallTrace[numberPreCalls:] {
		tokenAddr := call.Trace[0].Action.Call.To
		tokenBal := big.NewInt(0).SetBytes(call.Trace[0].Result.Output)
		revert := call.Trace[0].Error

		// edge case where token balance is greater than 2^256 (it's an error message)
		var twoTo256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil) // nolint:mnd
		if tokenBal.Cmp(twoTo256) > 0 {
			slog.Debug(
				"Token balance overflow inside of getERC20BalanceForAddress",
			)
			continue
		}

		if revert != "" {
			continue
		}
		if tokenBal.String() != "0" {
			bigTokenBal := web3.NewBigInt(tokenBal)
			if bigTokenBal.Cmp(web3.NewBigInt(big.NewInt(10))) >= 0 {
				ercBalanceMap[tokenAddr] = *bigTokenBal
			}

		}
	}

	return ercBalanceMap
}

// AddERC20BalanceCalls adds the erc20 balance calls to the provided calls.
func AddERC20BalanceCalls(
	erc20s []web3.Address,
	preCalls []TraceCallArgs,
	targetAddress web3.Address,

) ([]TraceCallArgs, error) {
	ercCallData, err := hex.DecodeString("70a08231000000000000000000000000")
	if err != nil {
		return nil, fmt.Errorf("error generating ERC-20 balance calldata: %v", err)
	}

	ercCalls := []TraceCallArgs{}
	ercCalls = append(ercCalls, preCalls...)

	for i := range erc20s {
		ercCalls = append(ercCalls,
			TraceCallArgs{
				Call: CallArgs{
					From: &targetAddress,
					To:   &erc20s[i],
					Data: web3.NewBytes(append(ercCallData, targetAddress.ToBytes()...)),
				},
				Traces: []Tracer{CallTracer},
			},
		)
	}

	return ercCalls, nil
}

// GetETHBalanceFromTrace gets the ETH balance of the provided address from the
// provided trace.
func GetETHBalanceFromTrace(
	responseTrace []traceresult.CallTrace, address web3.Address,
) (balance web3.BigInt) {
	ethBalance := big.NewInt(0)
	prunedCallTrace := PruneRevertedCalls(responseTrace)

	for _, subTrace := range prunedCallTrace {
		if !subTrace.Type.IsCall() ||
			subTrace.Action.Call.CallType != tracetypes.CallTypeCall ||
			subTrace.Action.Call.Value.String() == "0" ||
			subTrace.Error != "" {
			continue
		}

		if subTrace.Action.Call.To == address {
			ethBalance.Add(ethBalance, subTrace.Action.Call.Value.Int)
		}
		if subTrace.Action.Call.From == address {
			ethBalance.Sub(ethBalance, subTrace.Action.Call.Value.Int)
		}

	}

	return *web3.NewBigInt(ethBalance)
}

func ercMapToStr(m map[web3.Address]web3.BigInt) map[string]web3.BigInt {
	// Convert the map into a map with string keys
	res := make(map[string]web3.BigInt)
	for k, v := range m {
		res[k.String()] = v
	}

	return res
}
