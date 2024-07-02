package trace

import (
	"encoding/json"
	"fmt"
)

// TraceCallArgs represents the parameters for a trace_call RPC request
type TraceCallArgs struct {
	Call   CallArgs
	Traces []Tracer
}

// MarshalJSON override
func (t TraceCallArgs) MarshalJSON() ([]byte, error) {
	values := []interface{}{t.Call, t.Traces}

	return json.Marshal(values)
}

// UnmarshalJSON override
func (t *TraceCallArgs) UnmarshalJSON(data []byte) error {
	// Create a temporary structure to hold the data in the same format it was marshaled
	var temp []json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Check if the length of the temp array is as expected (2 in this case)
	if len(temp) != 2 {
		return fmt.Errorf("expected length of array is 2, got %d", len(temp))
	}

	// Unmarshal the first element of the array into t.Call
	if err := json.Unmarshal(temp[0], &t.Call); err != nil {
		return err
	}

	// Unmarshal the second element of the array into t.Traces
	return json.Unmarshal(temp[1], &t.Traces)
}
