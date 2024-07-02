package trace

import (
	"SHDB/pkg/trace/traceresult"

	"github.com/Skylock-ai/Arianrhod/pkg/types/primitives"
)

// PruneRevertedCalls prunes all reverted calls, and the children of these reverted
// calls, from a []CallTrace
func PruneRevertedCalls(trace []traceresult.CallTrace) (pruned []traceresult.CallTrace) {
	// loop over the trace backwards and find all
	pruned = trace

	for i := len(trace) - 1; i >= 0; i-- {
		if trace[i].Error != "" {
			pruned = pruneChildCall(i, pruned)
		}
	}

	return pruned
}

// pruneChildCall prunes the provided parent and all children, returning the
// pruned []CallTrace.
func pruneChildCall(
	parent int, trace []traceresult.CallTrace,
) (pruned []traceresult.CallTrace) {
	pruned = trace

	if len(trace) <= parent {
		return pruned
	}

	callAddress := trace[parent].TraceAddress
	for i := len(trace) - 1; i >= parent; i-- {
		child := trace[i]
		if primitives.IntSlice().StartsWith(child.TraceAddress, callAddress) {
			pruned = append(pruned[:i], pruned[i+1:]...)
		}
	}

	return pruned
}
