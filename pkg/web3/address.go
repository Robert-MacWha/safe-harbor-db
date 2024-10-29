package web3

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// Address type is an alias for [20]byte.
type Address [AddressLength]byte

// AddressLength is the length of an address in bytes.
const AddressLength int = 20

// HexToAddress creates a new Address instance from a hexString.
func HexToAddress(hexStr string) (*Address, error) {
	str := strings.TrimPrefix(hexStr, "0x")
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return nil,
			fmt.Errorf("cound not decode hexStr %v, error: %v", hexStr, err)
	}

	return BytesToAddress(bytes)
}

// BytesToAddress creates a new Address instance from a byte slice.
func BytesToAddress(bytes []byte) (*Address, error) {
	if len(bytes) != AddressLength {
		return nil,
			fmt.Errorf("Address must be length 20, provided bytes invalid: %v", bytes)
	}

	var a Address
	copy(a[:], bytes)
	return &a, nil
}

// CommonToAddress creates a new web3.Address struct from a geth common.address.
func CommonToAddress(addr common.Address) *Address {
	var a Address
	copy(a[:], addr[:])
	return &a
}

// GetContractAddress Calculates the address for some new contract deployed by
// deployerAddress when their nonce is some value.
func GetContractAddress(deployerAddress *Address, nonce int) *Address {
	// Prepare data for RLP Encoding [address, uint256]
	dataToEncode := []interface{}{
		deployerAddress,
		big.NewInt(int64(nonce)).Bytes(),
	}

	// RLP encode the data
	encodedData, _ := rlp.EncodeToBytes(dataToEncode)

	// hash the result
	hash := sha3.NewLegacyKeccak256()
	hash.Write(encodedData)
	hashedResult := hash.Sum(nil)[12:]

	addr, _ := BytesToAddress(hashedResult)
	return addr
}

// ToCommon converts a web3.Address to a common.Address.
func (a Address) ToCommon() common.Address {
	return common.BytesToAddress(a[:])
}

// ToBytes converts a web3.Address to a byte slice.
func (a Address) ToBytes() []byte {
	return a[:]
}

// ToHex converts a web3.Address to a hexString with a prefix of 0x.
func (a Address) ToHex() string {
	return "0x" + hex.EncodeToString(a[:])
}

// MarshalJSON override.
func (a Address) MarshalJSON() ([]byte, error) {
	hexValue := hex.EncodeToString(a[:])
	hexValue = "0x" + hexValue
	return json.Marshal(hexValue)
}

// UnmarshalJSON override.
func (a *Address) UnmarshalJSON(data []byte) error {
	if string(data) == `""` || string(data) == `"0x"` {
		*a = Address{}
		return nil
	}
	var hexValue string
	err := json.Unmarshal(data, &hexValue)
	if err != nil {
		return fmt.Errorf("error unmarshalling address: %v", err)
	}

	addr, err := HexToAddress(hexValue)
	if err != nil {
		return fmt.Errorf("error converting string %v to address: %v", hexValue, err)
	}

	copy(a[:], addr[:])
	return nil
}

// UnmarshalText override
func (a *Address) UnmarshalText(input []byte) error {
	hexValue := string(input)
	addr, err := HexToAddress(hexValue)
	if err != nil {
		return fmt.Errorf("error converting string %v to address: %v", hexValue, err)
	}

	copy(a[:], addr[:])
	return nil
}

// UnmarshalYAML override
func (a *Address) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var hexValue string
	if err := unmarshal(&hexValue); err != nil {
		return fmt.Errorf("error unmarshalling address: %v", err)
	}

	addr, err := HexToAddress(hexValue)
	if err != nil {
		return fmt.Errorf("error converting address string to hex: %v", err)
	}

	copy(a[:], addr[:])
	return nil
}

// String implementation to print as a hexStr
func (a Address) String() string {
	val, exists := ERC20Addresses[a]
	if exists {
		return val
	}

	return a.ToHex()
}
