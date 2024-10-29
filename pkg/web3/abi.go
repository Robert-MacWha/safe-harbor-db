package web3

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// LoadABIs loads the ABIs from a given directory path
func LoadABIs(abis []string, dirPath string) (map[string]abi.ABI, error) {
	response := make(map[string]abi.ABI)
	for _, name := range abis {

		abiFilePath := dirPath + "/" + name + ".json"

		contractAbi, err := LoadABI(abiFilePath)
		if err != nil {
			return nil, fmt.Errorf("error loading ABI: %w", err)
		}

		response[name] = contractAbi
	}

	return response, nil
}

// LoadABI loads the ABI from a given file path
func LoadABI(abiFilePath string) (abi.ABI, error) {
	abiData, err := os.ReadFile(abiFilePath)
	if err != nil {
		return abi.ABI{},
			fmt.Errorf("error reading ABI file: %w", err)
	}

	// Parse the JSON ABI data
	var contractAbi abi.ABI
	err = json.Unmarshal(abiData, &contractAbi)
	if err != nil {
		return abi.ABI{},
			fmt.Errorf("Error parsing ABI JSON: %v", err)
	}

	return contractAbi, nil
}

// EncodeABI encodes the ABI for a given contract and function signature
func EncodeABI(
	contract abi.ABI,
	sig []byte,
	inputs []interface{},
) ([]byte, error) {

	method, err := contract.MethodById(sig)
	if err != nil {
		return nil, err
	}

	result, err := method.Inputs.PackValues(inputs)
	if err != nil {
		return nil, err
	}

	result = append(sig, result...)

	return result, nil
}

// DecodeABI decodes the ABI for a given contract and function signature
func DecodeABI(contract abi.ABI, input []byte) ([]interface{}, error) {

	method, err := contract.MethodById(input[:4])
	if err != nil {
		return nil, err
	}

	inputsResult, err := method.Inputs.UnpackValues(input[4:])
	if err != nil {
		return nil, err
	}

	return inputsResult, nil
}
