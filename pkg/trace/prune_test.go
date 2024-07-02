package trace

import (
	"reflect"
	"testing"

	"SHDB/pkg/trace/traceresult"
)

func Test_pruneChildCall(t *testing.T) {
	type args struct {
		callIndex int
		trace     []traceresult.CallTrace
	}
	tests := []struct {
		name       string
		args       args
		wantPruned []traceresult.CallTrace
	}{
		{
			"Base",
			args{
				callIndex: 2,
				trace: []traceresult.CallTrace{
					{TraceAddress: []int{}},
					{TraceAddress: []int{0}},
					{TraceAddress: []int{1}},
					{TraceAddress: []int{1, 0}},
					{TraceAddress: []int{1, 1}},
					{TraceAddress: []int{1, 1, 0}},
					{TraceAddress: []int{2}},
				},
			},
			[]traceresult.CallTrace{
				{TraceAddress: []int{}},
				{TraceAddress: []int{0}},
				{TraceAddress: []int{2}},
			},
		},
		{
			"No Prune",
			args{
				callIndex: 3,
				trace: []traceresult.CallTrace{
					{TraceAddress: []int{}},
					{TraceAddress: []int{0}},
				},
			},
			[]traceresult.CallTrace{
				{TraceAddress: []int{}},
				{TraceAddress: []int{0}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPruned := pruneChildCall(tt.args.callIndex, tt.args.trace)
			if !reflect.DeepEqual(gotPruned, tt.wantPruned) {
				t.Errorf("pruneChildCall() = %v, want %v", gotPruned, tt.wantPruned)
			}
		})
	}
}

func Test_pruneRevertedCalls(t *testing.T) {
	type args struct {
		trace []traceresult.CallTrace
	}
	tests := []struct {
		name       string
		args       args
		wantPruned []traceresult.CallTrace
	}{
		{
			"Base",
			args{
				trace: []traceresult.CallTrace{
					{TraceAddress: []int{}, Error: ""},
					{TraceAddress: []int{0}, Error: ""},
					{TraceAddress: []int{1}, Error: "Reverted"},
					{TraceAddress: []int{1, 0}, Error: ""},
					{TraceAddress: []int{1, 1}, Error: ""},
					{TraceAddress: []int{1, 1, 0}, Error: ""},
					{TraceAddress: []int{2}, Error: ""},
				},
			},
			[]traceresult.CallTrace{
				{TraceAddress: []int{}, Error: ""},
				{TraceAddress: []int{0}, Error: ""},
				{TraceAddress: []int{2}, Error: ""},
			},
		},
		{
			"No Errors",
			args{
				trace: []traceresult.CallTrace{
					{TraceAddress: []int{}, Error: ""},
					{TraceAddress: []int{0}, Error: ""},
					{TraceAddress: []int{1}, Error: ""},
					{TraceAddress: []int{1, 0}, Error: ""},
					{TraceAddress: []int{2}, Error: ""},
				},
			},
			[]traceresult.CallTrace{
				{TraceAddress: []int{}, Error: ""},
				{TraceAddress: []int{0}, Error: ""},
				{TraceAddress: []int{1}, Error: ""},
				{TraceAddress: []int{1, 0}, Error: ""},
				{TraceAddress: []int{2}, Error: ""},
			},
		},
		{
			"Multiple Errors",
			args{
				trace: []traceresult.CallTrace{
					{TraceAddress: []int{}, Error: ""},
					{TraceAddress: []int{0}, Error: ""},
					{TraceAddress: []int{1}, Error: "Reverted"},
					{TraceAddress: []int{1, 0}, Error: ""},
					{TraceAddress: []int{1, 1}, Error: "Reverted"},
					{TraceAddress: []int{1, 1, 0}, Error: ""},
					{TraceAddress: []int{2}, Error: "Reverted"},
				},
			},
			[]traceresult.CallTrace{
				{TraceAddress: []int{}, Error: ""},
				{TraceAddress: []int{0}, Error: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPruned := PruneRevertedCalls(tt.args.trace)
			if !reflect.DeepEqual(gotPruned, tt.wantPruned) {
				t.Errorf("pruneRevertedCalls() = %v, want %v", gotPruned, tt.wantPruned)
			}
		})
	}
}
