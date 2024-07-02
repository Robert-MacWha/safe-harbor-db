package traceresult

import (
	"encoding/json"
	"fmt"

	"github.com/Skylock-ai/Arianrhod/pkg/types/tracetypes"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

// Action contains all possible call actions that a CallTrace can execute
type Action struct {
	Call    *CallAction
	Create  *CreateAction
	Suicide *SuicideAction
	Reward  *RewardAction
}

// CallTrace represnts a single call's trace using the party trace
// style. This struct is used when receiving a calltrace from the
// trace_call/trace_callMany endpoints.
type CallTrace struct {
	// Action will hold a value for one of the four structs: CallAction,
	// CreateAction, SuicideAction, or RewardAction
	Action       Action               `json:"action"`
	Error        string               `json:"error,omitempty"`
	Result       CallResult           `json:"result"`
	Subtraces    int                  `json:"subtraces"`
	TraceAddress []int                `json:"traceAddress"`
	Type         tracetypes.TraceType `json:"type"`
}

// CallAction is part of the ParityTrace struct.
type CallAction struct {
	From     web3.Address        `json:"from"`
	CallType tracetypes.CallType `json:"callType"`
	Gas      web3.BigInt         `json:"gas"`
	Input    web3.Bytes          `json:"input"`
	To       web3.Address        `json:"to"`
	Value    web3.BigInt         `json:"value"`
}

// CreateAction is part of the ParityTrace struct.
type CreateAction struct {
	From  web3.Address `json:"from"`
	Gas   web3.BigInt  `json:"gas"`
	Init  web3.Bytes   `json:"init"`
	Value web3.BigInt  `json:"value"`
}

// SuicideAction is part of the ParityTrace struct.
type SuicideAction struct {
	Address       web3.Address `json:"address"`
	RefundAddress web3.Address `json:"refundAddress"`
	Balance       web3.BigInt  `json:"balance"`
}

// RewardAction is part of the ParityTrace struct.
type RewardAction struct {
	Author     web3.Address `json:"author"`
	RewardType string       `json:"rewardType"`
	Value      web3.BigInt  `json:"value,omitempty"`
}

// CallResult is part of the ParityTrace struct.
type CallResult struct {
	GasUsed web3.BigInt   `json:"gasUsed,omitempty"`
	Output  web3.Bytes    `json:"output,omitempty"`
	Address *web3.Address `json:"address,omitempty"`
}

// UnmarshalJSON override
func (ct *CallTrace) UnmarshalJSON(data []byte) error {
	// define a temporary struct to avoid infinite recursion
	type Alias CallTrace
	aux := &struct {
		Action json.RawMessage `json:"action"`
		*Alias
	}{
		Alias: (*Alias)(ct),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	traceType := aux.Type
	switch {
	case traceType == tracetypes.TraceTypeCall:
		ct.Action.Call = &CallAction{}
		return json.Unmarshal(aux.Action, ct.Action.Call)
	case traceType == tracetypes.TraceTypeCreate:
		ct.Action.Create = &CreateAction{}
		return json.Unmarshal(aux.Action, ct.Action.Create)
	case traceType == tracetypes.TraceTypeSuicide:
		ct.Action.Suicide = &SuicideAction{}
		return json.Unmarshal(aux.Action, ct.Action.Suicide)
	default:
		return fmt.Errorf(
			"error unmarshalling CallTrace: unrecognized traceType %v",
			traceType,
		)
	}
}
