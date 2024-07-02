package traceresult

import (
	"encoding/json"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

// TraceResult struct contains the unmarshalled response from a rpc
// `trace_call` or `trace_callMany` request. Based on the results RE:
// https://www.quicknode.com/docs/ethereum/trace_call
type TraceResult struct {
	Output    web3.Bytes                        `json:"output"`
	StateDiff map[web3.Address]StateDiffAccount `json:"stateDiff"`
	Trace     []CallTrace                       `json:"trace"`
	VMTrace   VMTrace                           `json:"vmTrace"`
}

// MarshalJSON override
func (t *TraceResult) MarshalJSON() ([]byte, error) {
	// Convert map[web3.Address]StateDiffAccount to map[string]StateDiffAccount
	transformedStateDiff := make(map[string]StateDiffAccount)
	for addr, account := range t.StateDiff {
		transformedStateDiff[addr.String()] = account
	}

	// Create a temporary struct to hold the converted data
	type Alias TraceResult
	tmp := &struct {
		StateDiff map[string]StateDiffAccount `json:"stateDiff"`
		*Alias
	}{
		StateDiff: transformedStateDiff,
		Alias:     (*Alias)(t),
	}

	// Marshal the temporary struct
	//nolint:staticcheck
	return json.Marshal(tmp)
}
