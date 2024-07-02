package trace

// Tracer enum represents a trace type for trace_call* rpc endpoints
type Tracer string

const (
	// VMTracer traces each opcode execution during a call
	VMTracer Tracer = "vmTrace"

	// CallTracer traces each internal transaction during a call
	CallTracer Tracer = "trace"

	// StateDiffTracer traces the storage value changes at the end of a call
	StateDiffTracer Tracer = "stateDiff"
)
