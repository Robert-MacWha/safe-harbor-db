package trace

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

func TestTraceCallArgs_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Call   CallArgs
		Traces []Tracer
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Standard TraceCallArgs",
			fields: fields{
				Call: CallArgs{
					From:  &web3.Address{0x01},
					To:    &web3.Address{0x02},
					Value: &web3.BigInt{Int: big.NewInt(3)},
					Data:  &web3.Bytes{0x04},
				},
				Traces: []Tracer{
					VMTracer,
					CallTracer,
				},
			},
			args: args{
				data: []byte(
					//revive:disable-next-line:line-length-limit
					`[{"from":"0x0100000000000000000000000000000000000000","to":"0x0200000000000000000000000000000000000000","value":"0x3","data":"0x04"},["vmTrace","trace"]]`,
				),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TraceCallArgs{
				Call:   tt.fields.Call,
				Traces: tt.fields.Traces,
			}

			if err := tr.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("TraceCallArgs.UnmarshalJSON() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestTraceCallArgs_MarshalJSON(t *testing.T) {
	type fields struct {
		Call   CallArgs
		Traces []Tracer
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Standard TraceCallArgs",
			fields: fields{
				Call: CallArgs{
					From:  &web3.Address{0x01},
					To:    &web3.Address{0x02},
					Value: &web3.BigInt{Int: big.NewInt(3)},
					Data:  &web3.Bytes{0x04},
				},
				Traces: []Tracer{
					VMTracer,
					CallTracer,
				},
			},
			want: []byte(
				//revive:disable-next-line:line-length-limit
				`[{"from":"0x0100000000000000000000000000000000000000","to":"0x0200000000000000000000000000000000000000","value":"0x3","data":"0x04"},["vmTrace","trace"]]`,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TraceCallArgs{
				Call:   tt.fields.Call,
				Traces: tt.fields.Traces,
			}
			got, err := tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("TraceCallArgs.MarshalJSON() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TraceCallArgs.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
