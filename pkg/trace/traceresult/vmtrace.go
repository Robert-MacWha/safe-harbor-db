package traceresult

// VMTrace contains information on the full trace of the EVM's state during
// the execution of the given transaction, including for any sub-calls.
type VMTrace struct {
	Ops []*VMTraceOp `json:"ops"`
}

// VMTraceOp is a single operation executed on the EVM. If this opcode lead
// to sub-calls, they will be contained in the `sub` value.
type VMTraceOp struct {
	Op   string     `json:"op"`
	Cost int        `json:"cost"`
	Ex   *VMTraceEx `json:"ex"`
	Pc   int        `json:"pc"`
	Sub  *VMTrace   `json:"sub,omitempty"`
}

// VMTraceEx is one element of the vmTrace ops trace.
type VMTraceEx struct {
	Mem   *VMTraceMem   `json:"mem,omitempty"`
	Push  []string      `json:"push,omitempty"`
	Store *VMTraceStore `json:"store,omitempty"`
	Used  int           `json:"used,omitempty"`
}

// VMTraceMem is the state of memory
type VMTraceMem struct {
	Data string `json:"data"`
	Off  int    `json:"off"`
}

// VMTraceStore is one element of the vmTrace ops trace.
type VMTraceStore struct {
	Key string `json:"key"`
	Val string `json:"val"`
}
