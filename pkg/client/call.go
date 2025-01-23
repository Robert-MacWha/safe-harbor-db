package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nanmu42/etherscan-api"
)

// Define a call object that matches the structure
type call struct {
	Type    string           `json:"type"`
	From    common.Address   `json:"from"`
	To      *common.Address  `json:"to,omitempty"`
	Gas     etherscan.BigInt `json:"gas"`
	GasUsed etherscan.BigInt `json:"gasUsed"`
	Input   string           `json:"input"`
	Output  string           `json:"output"`
	Value   etherscan.BigInt `json:"value"`
	Calls   []call           `json:"calls,omitempty"`
}

// Calls holds the root call
type Calls struct {
	Calls []call `json:"calls"`
}

// flatten recursively flattens a call and its subcalls
func (c *call) Flatten() []call {
	var result []call

	result = append(result, *c)
	for _, subCall := range c.Calls {
		result = append(result, subCall.Flatten()...)
	}

	return result
}

// flatten recursively flattens a debugResult and its calls
func (r *Calls) Flatten() []call {
	var result []call

	for _, call := range r.Calls {
		result = append(result, call.Flatten()...)
	}

	return result
}
