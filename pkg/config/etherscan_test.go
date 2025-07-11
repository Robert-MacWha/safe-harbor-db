package config

import "testing"

func TestParseAddressFromURL(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string
		wantChainId int
		wantErr     bool
	}{
		{
			name: "Etherscan",
			args: args{
				raw: "https://etherscan.io/address/0x6ef3D766Dfe02Dc4bF04aAe9122EB9A0Ded25615",
			},
			wantAddress: "0x6ef3D766Dfe02Dc4bF04aAe9122EB9A0Ded25615",
			wantChainId: 1,
			wantErr:     false,
		},
		{
			name: "Unknown Chain",
			args: args{
				raw: "https://unknownscan.io/address/0x6ef3D766Dfe02Dc4bF04aAe9122EB9A0Ded25615",
			},
			wantAddress: "",
			wantChainId: 0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, gotChainId, err := ParseAddressFromURL(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAddressFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("ParseAddressFromURL() gotAddress = %v, want %v", gotAddress, tt.wantAddress)
			}
			if gotChainId != tt.wantChainId {
				t.Errorf("ParseAddressFromURL() gotChainId = %v, want %v", gotChainId, tt.wantChainId)
			}
		})
	}
}
