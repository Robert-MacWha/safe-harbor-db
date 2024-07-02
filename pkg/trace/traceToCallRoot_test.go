package trace

import (
	"testing"

	"SHDB/pkg/trace/traceresult"

	"github.com/Skylock-ai/Arianrhod/pkg/blackwood/call"
	"github.com/Skylock-ai/Arianrhod/pkg/types/tracetypes"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

func TestParityToCallRoot(t *testing.T) {
	type args struct {
		parityTrace []traceresult.CallTrace
	}
	tests := []struct {
		name    string
		args    args
		want    *call.Call
		wantErr bool
	}{
		{
			"empty trace",
			args{
				[]traceresult.CallTrace{},
			},
			nil,
			true,
		},
		{
			"only root",
			args{
				[]traceresult.CallTrace{
					{
						Type: tracetypes.TraceTypeCall,
						Action: traceresult.Action{
							Call: &traceresult.CallAction{
								From:     web3.Address{0x01, 0x02, 0x3, 0x04},
								CallType: tracetypes.CallTypeCall,
								Gas:      *web3.IntToBig(10),
								Input:    web3.Bytes{0x09, 0x0A, 0x0B, 0x0C},
								To:       web3.Address{0x05, 0x06, 0x07, 0x08},
								Value:    *web3.IntToBig(20),
							},
						},
						Result: traceresult.CallResult{
							GasUsed: web3.BigInt{},
							Output:  web3.Bytes{0x0D, 0x0E, 0x0F, 0x10},
							Address: nil,
						},
						Subtraces:    0,
						TraceAddress: []int{0},
					},
				},
			},
			&call.Call{
				TraceType: tracetypes.TraceTypeCall,
				From:      web3.Address{0x01, 0x02, 0x3, 0x04},
				To:        web3.Address{0x05, 0x06, 0x07, 0x08},
				Gas:       *web3.IntToBig(10),
				Value:     *web3.IntToBig(20),
				CallType:  tracetypes.CallTypeCall,
				Input:     web3.Bytes{0x09, 0x0A, 0x0B, 0x0C},
				Output:    web3.Bytes{0x0D, 0x0E, 0x0F, 0x10},
				Error:     "",
				Children:  nil,
				Parent:    nil,
			},
			false,
		},
		{
			"ordered children",
			args{
				[]traceresult.CallTrace{
					{
						Type: tracetypes.TraceTypeCall,
						Action: traceresult.Action{
							Call: &traceresult.CallAction{
								From: web3.Address{0x11},
								To:   web3.Address{0x22},
							},
						},
						Subtraces:    2,
						TraceAddress: []int{},
					},
					{
						Type: tracetypes.TraceTypeCall,
						Action: traceresult.Action{
							Call: &traceresult.CallAction{
								From: web3.Address{0x22},
							},
						},
						Subtraces:    0,
						TraceAddress: []int{0},
					},
					{
						Type: tracetypes.TraceTypeCall,
						Action: traceresult.Action{
							Call: &traceresult.CallAction{
								From: web3.Address{0x22},
								To:   web3.Address{0x33},
							},
						},
						Subtraces:    1,
						TraceAddress: []int{1},
					},
					{
						Type: tracetypes.TraceTypeCall,
						Action: traceresult.Action{
							Call: &traceresult.CallAction{
								From: web3.Address{0x33},
								To:   web3.Address{0x44},
							},
						},
						Subtraces:    0,
						TraceAddress: []int{1, 0},
					},
				},
			},
			&call.Call{
				TraceType: tracetypes.TraceTypeCall,
				From:      web3.Address{0x11},
				To:        web3.Address{0x22},
				Children: []*call.Call{
					{
						TraceType: tracetypes.TraceTypeCall,
						From:      web3.Address{0x22},
					},
					{
						TraceType: tracetypes.TraceTypeCall,
						From:      web3.Address{0x22},
						To:        web3.Address{0x33},
						Children: []*call.Call{
							{
								TraceType: tracetypes.TraceTypeCall,
								From:      web3.Address{0x33},
								To:        web3.Address{0x44},
							},
						},
					},
				},
			},
			false,
		},
		{
			"unordered children",
			args{
				[]traceresult.CallTrace{
					{ // ()
						Type: tracetypes.TraceTypeCall,
						Action: traceresult.Action{
							Call: &traceresult.CallAction{
								From: web3.Address{0x00},
								To:   web3.Address{0x11},
							},
						},
						Subtraces:    2,
						TraceAddress: []int{},
					},
					{ // (), (1)
						Type: tracetypes.TraceTypeCall,
						Action: traceresult.Action{
							Call: &traceresult.CallAction{
								From: web3.Address{0x11},
								To:   web3.Address{0x22},
							},
						},
						Subtraces:    0,
						TraceAddress: []int{1},
					},
					{ // (), (0)
						Type: tracetypes.TraceTypeCall,
						Action: traceresult.Action{
							Call: &traceresult.CallAction{
								From: web3.Address{0x11},
								To:   web3.Address{0x33},
							},
						},
						Subtraces:    0,
						TraceAddress: []int{0},
					},
				},
			},
			&call.Call{
				TraceType: tracetypes.TraceTypeCall,
				From:      web3.Address{0x00},
				To:        web3.Address{0x11},
				Children: []*call.Call{
					{
						TraceType: tracetypes.TraceTypeCall,
						From:      web3.Address{0x11},
						To:        web3.Address{0x33},
					},
					{
						TraceType: tracetypes.TraceTypeCall,
						From:      web3.Address{0x11},
						To:        web3.Address{0x22},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		if tt.want != nil {
			tt.want.LinkParents(nil)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := ParityToCallRoot(tt.args.parityTrace)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParityToCallRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil && got == nil {
				return
			}

			if !tt.want.DeepEqual(got) {
				t.Errorf("ParityToCallRoot() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
