package web3

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Bytes acts as a type alias for []byte that can easily be marshalled and
// unmarshalled as a hexstring
type Bytes []byte

// NewBytes creates a new Bytes object from a byte slice
func NewBytes(d []byte) *Bytes {
	b := make(Bytes, len(d))
	copy(b, d)

	return &b
}

// MarshalJSON override
func (b Bytes) MarshalJSON() ([]byte, error) {
	hexstr := "0x" + hex.EncodeToString(b)
	return json.Marshal(hexstr)
}

// UnmarshalJSON override
func (b *Bytes) UnmarshalJSON(data []byte) error {
	var hexValue string
	err := json.Unmarshal(data, &hexValue)
	if err != nil {
		return fmt.Errorf("error unmarshalling ByteSlice: %v", err)
	}

	hexValue = strings.TrimPrefix(hexValue, "0x")
	if len(hexValue)%2 != 0 {
		hexValue = "0" + hexValue
	}

	bytes, err := hex.DecodeString(hexValue)
	if err != nil {
		return fmt.Errorf("error decoding hexStr: %v", err)
	}
	*b = Bytes(bytes)

	return nil
}

// String implementation to print as a hexStr
func (b Bytes) String() string {
	return "0x" + hex.EncodeToString(b)
}
