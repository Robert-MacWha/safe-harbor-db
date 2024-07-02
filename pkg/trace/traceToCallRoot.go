package trace

import (
	"fmt"
	"math/big"

	"SHDB/pkg/trace/traceresult"

	"github.com/Skylock-ai/Arianrhod/pkg/blackwood/call"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

// ErrEmptyTrace is the error used when ParityToCallRoot is called with no traces.
var ErrEmptyTrace = fmt.Errorf("empty trace")

// TODO merge tracestatediff.go's Handle action with this one

// ParityToCallRoot converts Parity trace to a nested Call.Call representation.
func ParityToCallRoot(parityTrace []traceresult.CallTrace) (*call.Call, error) {
	// This method is tricky because reth sometimes returns traces in a different
	// order than they were called. For example, the return order for the trace
	// addresses might be [(), (0), (0, 25), (1), ...].
	//
	// Because of this, we make the assumption that parent traces will always appear
	// before children, but we don't make the assumption that children come in order.

	if len(parityTrace) == 0 {
		return nil, ErrEmptyTrace
	}

	// root call = first trace
	root, err := newCallFromParity(parityTrace[0])
	if err != nil {
		return nil, fmt.Errorf("parityToCall: %w", err)
	}

	for _, parity := range parityTrace[1:] {
		call, err := newCallFromParity(parity)
		if err != nil {
			return nil, fmt.Errorf("parityToCall: %w", err)
		}

		err = root.Inject(call, parity.TraceAddress)
		if err != nil {
			return nil, fmt.Errorf("root.Insert: %w", err)
		}
	}

	return root, nil
}

func newCallFromParity(parityCall traceresult.CallTrace) (*call.Call, error) {
	var traceCall *call.Call
	var err error

	switch parityCall.Type {
	case "CALL":
		traceCall, err = handleCallAction(parityCall)
	case "CREATE":
		traceCall, err = handleCreateAction(parityCall)
	case "SUICIDE":
		traceCall, err = handleSuicideAction(parityCall)
	default:
		return nil, fmt.Errorf("unknown Trace type")
	}

	if err != nil {
		return nil, err
	}

	traceCall.Children = make([]*call.Call, parityCall.Subtraces)
	return traceCall, err
}

func handleCallAction(parityCall traceresult.CallTrace) (*call.Call, error) {
	traceCall := &call.Call{
		From:      parityCall.Action.Call.From,
		To:        parityCall.Action.Call.To,
		Gas:       parityCall.Action.Call.Gas,
		Value:     parityCall.Action.Call.Value,
		CallType:  parityCall.Action.Call.CallType,
		Input:     parityCall.Action.Call.Input,
		Output:    parityCall.Result.Output,
		Error:     parityCall.Error,
		TraceType: parityCall.Type,
	}

	return traceCall, nil
}

func handleCreateAction(parityCall traceresult.CallTrace) (*call.Call, error) {
	if parityCall.Result.Address == nil {
		// Edge case where the contract creation errors out that to address is nil
		// Ex: 0x88425792ff00b7d3c87341dd1f17a6b502bbce3338255a863f26c740fa50de99
		return nil, fmt.Errorf("Error. To address of contract creation is nil")
	}

	traceCall := &call.Call{
		From:      parityCall.Action.Create.From,
		To:        *parityCall.Result.Address, //Create Address = to
		Gas:       parityCall.Action.Create.Gas,
		Value:     parityCall.Action.Create.Value,
		CallType:  "",                              // Note: No Create CallType
		Input:     parityCall.Action.Create.Init,   // init is the input for create
		Output:    *web3.NewBytes(make([]byte, 0)), // There's no Output in Creates
		Error:     parityCall.Error,
		TraceType: parityCall.Type,
	}

	return traceCall, nil
}

// Example:0x7d2296bcb936aa5e2397ddf8ccba59f54a178c3901666b49291d880369dbcf31
func handleSuicideAction(parityCall traceresult.CallTrace) (*call.Call, error) {
	traceCall := &call.Call{
		From:      parityCall.Action.Suicide.Address,
		To:        parityCall.Action.Suicide.RefundAddress, //Refund is to
		Gas:       *web3.NewBigInt(big.NewInt(0)),          // No gas in suicide
		Value:     parityCall.Action.Suicide.Balance,
		CallType:  "",                              // No Standard CallType for suicide
		Input:     *web3.NewBytes(make([]byte, 0)), // No input data for suicide
		Output:    *web3.NewBytes(make([]byte, 0)), // There's no Output in suicide
		Error:     parityCall.Error,
		TraceType: parityCall.Type,
	}
	return traceCall, nil
}
