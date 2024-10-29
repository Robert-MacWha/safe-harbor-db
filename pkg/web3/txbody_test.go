package web3

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxBody_MarshalJSON(t *testing.T) {
	tx := TxBody{
		From: Address{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		},
		To: &Address{
			1, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
		},
		Hash:                 "0x010203",
		Gas:                  *IntToBig(1000),
		GasPrice:             *IntToBig(2000),
		Input:                Bytes{41, 42, 43, 44},
		Value:                *IntToBig(3000),
		Block:                IntToBig(37),
		Nonce:                *IntToBig(21),
		ChainID:              *IntToBig(5),
		MaxFeePerGas:         IntToBig(31),
		MaxPriorityFeePerGas: IntToBig(32),
	}

	data, err := json.Marshal(&tx)
	assert.NoError(t, err)

	//revive:disable
	expectedJSON := `{"from":"0x0102030405060708090a0b0c0d0e0f1011121314","to":"0x01161718191a1b1c1d1e1f202122232425262728","hash":"0x010203","gas":"0x3e8","gasPrice":"0x7d0","input":"0x292a2b2c","value":"0xbb8", "blockNumber": "0x25", "nonce": "0x15", "chainId": "0x5", "maxFeePerGas": "0x1f", "maxPriorityFeePerGas": "0x20"}`
	//revive:enable

	var expectedMap map[string]interface{}
	var actualMap map[string]interface{}

	err = json.Unmarshal([]byte(expectedJSON), &expectedMap)
	assert.NoError(t, err)

	err = json.Unmarshal(data, &actualMap)
	assert.NoError(t, err)

	assert.Equal(t, expectedMap, actualMap)
}

func TestTxBody_UnmarshalJSON(t *testing.T) {
	//revive:disable
	jsonData := `{"from":"0x0102030405060708090a0b0c0d0e0f1011121314","to":"0x01161718191a1b1c1d1e1f202122232425262728","hash":"0x010203","gas":"0x3e8","gasPrice":"0x7d0","input":"0x292a2b2c","value":"0xbb8", "blockNumber": "0x25", "nonce": "0x15", "chainId": "0x5", "maxFeePerGas": "0x1f", "maxPriorityFeePerGas": "0x20"}`
	//revive:enable

	var tx TxBody
	err := json.Unmarshal([]byte(jsonData), &tx)
	assert.NoError(t, err)

	expectedTX := TxBody{
		From: Address{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		},
		To: &Address{
			1, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
		},
		Hash:                 "0x010203",
		Gas:                  *IntToBig(1000),
		GasPrice:             *IntToBig(2000),
		Input:                Bytes{41, 42, 43, 44},
		Value:                *IntToBig(3000),
		Block:                IntToBig(37),
		Nonce:                *IntToBig(21),
		ChainID:              *IntToBig(5),
		MaxFeePerGas:         IntToBig(31),
		MaxPriorityFeePerGas: IntToBig(32),
	}

	assert.Equal(t, expectedTX, tx)
}
