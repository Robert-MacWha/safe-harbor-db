package web3

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// Hash type represents a hashed value stored as a [32]byte
type Hash [HashLength]byte

// HashLength is the const length of a hash value
const HashLength = 32

// HexToHash Converts a hexString, with or without prefix, to a Hash
func HexToHash(hexStr string) (*Hash, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
	if err != nil {
		return nil, err
	}
	return BytesToHash(b)
}

// BytesToHash Converts a []byte to a Hash
func BytesToHash(b []byte) (*Hash, error) {
	if len(b) != HashLength {
		return nil,
			fmt.Errorf("Hash must be length 32, provided bytes invalid: %v", b)
	}

	var h Hash
	copy(h[:], b)
	return &h, nil
}

// ToCommon converts a Hash to a geth common.Hash
func (h Hash) ToCommon() common.Hash {
	return common.BytesToHash(h[:])
}

// CommonToHash converts a geth common.Hash to a Hash
func CommonToHash(hash common.Hash) *Hash {
	var h Hash
	copy(h[:], hash[:])
	return &h
}

// ToHex converts a web3.Hash to a hexString with a prefix of 0x.
func (h Hash) ToHex() string {
	return "0x" + hex.EncodeToString(h[:])
}

// MarshalJSON override.
func (h *Hash) MarshalJSON() ([]byte, error) {
	hexValue := hex.EncodeToString(h[:])
	hexValue = "0x" + hexValue
	return json.Marshal(hexValue)
}

// UnmarshalJSON override.
func (h *Hash) UnmarshalJSON(data []byte) error {
	var hexValue string
	err := json.Unmarshal(data, &hexValue)
	if err != nil {
		return fmt.Errorf("error unmarshalling hash: %v", err)
	}

	hash, err := HexToHash(hexValue)
	if err != nil {
		return fmt.Errorf("error converting hash string to hex: %v", err)
	}

	copy(h[:], hash[:])
	return nil
}

// UnmarshalText override
func (h *Hash) UnmarshalText(input []byte) error {
	return h.UnmarshalJSON(input)
}

// String implementation to print as a hexStr
func (h Hash) String() string {
	return h.ToHex()
}
