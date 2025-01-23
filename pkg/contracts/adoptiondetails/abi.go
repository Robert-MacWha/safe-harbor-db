// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package adoptiondetails

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Account is an auto generated low-level Go binding around an user-defined struct.
type Account struct {
	AccountAddress     common.Address
	ChildContractScope uint8
	Signature          []byte
}

// AgreementDetailsV1 is an auto generated low-level Go binding around an user-defined struct.
type AgreementDetailsV1 struct {
	ProtocolName   string
	ContactDetails []Contact
	Chains         []Chain
	BountyTerms    BountyTerms
	AgreementURI   string
}

// BountyTerms is an auto generated low-level Go binding around an user-defined struct.
type BountyTerms struct {
	BountyPercentage      *big.Int
	BountyCapUSD          *big.Int
	Retainable            bool
	Identity              uint8
	DiligenceRequirements string
}

// Chain is an auto generated low-level Go binding around an user-defined struct.
type Chain struct {
	AssetRecoveryAddress common.Address
	Accounts             []Account
	Id                   *big.Int
}

// Contact is an auto generated low-level Go binding around an user-defined struct.
type Contact struct {
	Name    string
	Contact string
}

// AdoptiondetailsMetaData contains all meta data concerning the Adoptiondetails contract.
var AdoptiondetailsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"protocolName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contact\",\"type\":\"string\"}],\"internalType\":\"structContact[]\",\"name\":\"contactDetails\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetRecoveryAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"internalType\":\"enumChildContractScope\",\"name\":\"childContractScope\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structAccount[]\",\"name\":\"accounts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"internalType\":\"structChain[]\",\"name\":\"chains\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"bountyPercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bountyCapUSD\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"retainable\",\"type\":\"bool\"},{\"internalType\":\"enumIdentityRequirements\",\"name\":\"identity\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"diligenceRequirements\",\"type\":\"string\"}],\"internalType\":\"structBountyTerms\",\"name\":\"bountyTerms\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"agreementURI\",\"type\":\"string\"}],\"internalType\":\"structAgreementDetailsV1\",\"name\":\"_details\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getDetails\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"protocolName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contact\",\"type\":\"string\"}],\"internalType\":\"structContact[]\",\"name\":\"contactDetails\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetRecoveryAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"internalType\":\"enumChildContractScope\",\"name\":\"childContractScope\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structAccount[]\",\"name\":\"accounts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"internalType\":\"structChain[]\",\"name\":\"chains\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"bountyPercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bountyCapUSD\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"retainable\",\"type\":\"bool\"},{\"internalType\":\"enumIdentityRequirements\",\"name\":\"identity\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"diligenceRequirements\",\"type\":\"string\"}],\"internalType\":\"structBountyTerms\",\"name\":\"bountyTerms\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"agreementURI\",\"type\":\"string\"}],\"internalType\":\"structAgreementDetailsV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// AdoptiondetailsABI is the input ABI used to generate the binding from.
// Deprecated: Use AdoptiondetailsMetaData.ABI instead.
var AdoptiondetailsABI = AdoptiondetailsMetaData.ABI

// Adoptiondetails is an auto generated Go binding around an Ethereum contract.
type Adoptiondetails struct {
	AdoptiondetailsCaller     // Read-only binding to the contract
	AdoptiondetailsTransactor // Write-only binding to the contract
	AdoptiondetailsFilterer   // Log filterer for contract events
}

// AdoptiondetailsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AdoptiondetailsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdoptiondetailsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AdoptiondetailsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdoptiondetailsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AdoptiondetailsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdoptiondetailsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AdoptiondetailsSession struct {
	Contract     *Adoptiondetails  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AdoptiondetailsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AdoptiondetailsCallerSession struct {
	Contract *AdoptiondetailsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AdoptiondetailsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AdoptiondetailsTransactorSession struct {
	Contract     *AdoptiondetailsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AdoptiondetailsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AdoptiondetailsRaw struct {
	Contract *Adoptiondetails // Generic contract binding to access the raw methods on
}

// AdoptiondetailsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AdoptiondetailsCallerRaw struct {
	Contract *AdoptiondetailsCaller // Generic read-only contract binding to access the raw methods on
}

// AdoptiondetailsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AdoptiondetailsTransactorRaw struct {
	Contract *AdoptiondetailsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAdoptiondetails creates a new instance of Adoptiondetails, bound to a specific deployed contract.
func NewAdoptiondetails(address common.Address, backend bind.ContractBackend) (*Adoptiondetails, error) {
	contract, err := bindAdoptiondetails(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Adoptiondetails{AdoptiondetailsCaller: AdoptiondetailsCaller{contract: contract}, AdoptiondetailsTransactor: AdoptiondetailsTransactor{contract: contract}, AdoptiondetailsFilterer: AdoptiondetailsFilterer{contract: contract}}, nil
}

// NewAdoptiondetailsCaller creates a new read-only instance of Adoptiondetails, bound to a specific deployed contract.
func NewAdoptiondetailsCaller(address common.Address, caller bind.ContractCaller) (*AdoptiondetailsCaller, error) {
	contract, err := bindAdoptiondetails(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AdoptiondetailsCaller{contract: contract}, nil
}

// NewAdoptiondetailsTransactor creates a new write-only instance of Adoptiondetails, bound to a specific deployed contract.
func NewAdoptiondetailsTransactor(address common.Address, transactor bind.ContractTransactor) (*AdoptiondetailsTransactor, error) {
	contract, err := bindAdoptiondetails(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AdoptiondetailsTransactor{contract: contract}, nil
}

// NewAdoptiondetailsFilterer creates a new log filterer instance of Adoptiondetails, bound to a specific deployed contract.
func NewAdoptiondetailsFilterer(address common.Address, filterer bind.ContractFilterer) (*AdoptiondetailsFilterer, error) {
	contract, err := bindAdoptiondetails(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AdoptiondetailsFilterer{contract: contract}, nil
}

// bindAdoptiondetails binds a generic wrapper to an already deployed contract.
func bindAdoptiondetails(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AdoptiondetailsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Adoptiondetails *AdoptiondetailsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Adoptiondetails.Contract.AdoptiondetailsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Adoptiondetails *AdoptiondetailsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Adoptiondetails.Contract.AdoptiondetailsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Adoptiondetails *AdoptiondetailsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Adoptiondetails.Contract.AdoptiondetailsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Adoptiondetails *AdoptiondetailsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Adoptiondetails.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Adoptiondetails *AdoptiondetailsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Adoptiondetails.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Adoptiondetails *AdoptiondetailsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Adoptiondetails.Contract.contract.Transact(opts, method, params...)
}

// GetDetails is a free data retrieval call binding the contract method 0xfbbf93a0.
//
// Solidity: function getDetails() view returns((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string))
func (_Adoptiondetails *AdoptiondetailsCaller) GetDetails(opts *bind.CallOpts) (AgreementDetailsV1, error) {
	var out []interface{}
	err := _Adoptiondetails.contract.Call(opts, &out, "getDetails")

	if err != nil {
		return *new(AgreementDetailsV1), err
	}

	out0 := *abi.ConvertType(out[0], new(AgreementDetailsV1)).(*AgreementDetailsV1)

	return out0, err

}

// GetDetails is a free data retrieval call binding the contract method 0xfbbf93a0.
//
// Solidity: function getDetails() view returns((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string))
func (_Adoptiondetails *AdoptiondetailsSession) GetDetails() (AgreementDetailsV1, error) {
	return _Adoptiondetails.Contract.GetDetails(&_Adoptiondetails.CallOpts)
}

// GetDetails is a free data retrieval call binding the contract method 0xfbbf93a0.
//
// Solidity: function getDetails() view returns((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string))
func (_Adoptiondetails *AdoptiondetailsCallerSession) GetDetails() (AgreementDetailsV1, error) {
	return _Adoptiondetails.Contract.GetDetails(&_Adoptiondetails.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_Adoptiondetails *AdoptiondetailsCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Adoptiondetails.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_Adoptiondetails *AdoptiondetailsSession) Version() (string, error) {
	return _Adoptiondetails.Contract.Version(&_Adoptiondetails.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_Adoptiondetails *AdoptiondetailsCallerSession) Version() (string, error) {
	return _Adoptiondetails.Contract.Version(&_Adoptiondetails.CallOpts)
}
