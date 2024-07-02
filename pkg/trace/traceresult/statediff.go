package traceresult

import (
	"encoding/json"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
)

// StateDiffAccount contains the information on altered Ethereum state due to
// the execution of the given transaction.
type StateDiffAccount struct {
	Balance StateDiffAccountValue               `json:"balance"`
	Code    StateDiffAccountValue               `json:"code"`
	Nonce   StateDiffAccountValue               `json:"nonce"`
	Storage map[web3.Hash]StateDiffAccountValue `json:"storage"`
}

// StateDiffAccountValue contains information on a single state difference from
// a stateDiff trace
type StateDiffAccountValue struct {
	Changed bool
	From    web3.Bytes `json:"from"`
	To      web3.Bytes `json:"to"`
}

// UnmarshalJSON override
func (s *StateDiffAccountValue) UnmarshalJSON(data []byte) error {
	// Checks if the value is a unchanged value
	if data[1] == byte('=') {
		s.Changed = false
		return nil
	}

	var temp struct {
		Change struct {
			From web3.Bytes `json:"from"`
			To   web3.Bytes `json:"to"`
		} `json:"*"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	s.From = temp.Change.From
	s.To = temp.Change.To
	s.Changed = true

	return nil
}
