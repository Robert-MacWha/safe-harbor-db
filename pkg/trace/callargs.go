package trace

import (
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

// CallArgs represents a transaction to be traced
type CallArgs struct {
	From  *web3.Address `json:"from"`
	To    *web3.Address `json:"to,omitempty"`
	Value *web3.BigInt  `json:"value,omitempty"`
	Data  *web3.Bytes   `json:"data,omitempty"`
}
