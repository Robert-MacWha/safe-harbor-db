// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package safeharbor

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

// SafeharborMetaData contains all meta data concerning the Safeharbor contract.
var SafeharborMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_fallbackRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"NoAgreement\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"entity\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldDetails\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newDetails\",\"type\":\"address\"}],\"name\":\"SafeHarborAdoption\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"protocolName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contact\",\"type\":\"string\"}],\"internalType\":\"structContact[]\",\"name\":\"contactDetails\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetRecoveryAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"internalType\":\"enumChildContractScope\",\"name\":\"childContractScope\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structAccount[]\",\"name\":\"accounts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"internalType\":\"structChain[]\",\"name\":\"chains\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"bountyPercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bountyCapUSD\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"retainable\",\"type\":\"bool\"},{\"internalType\":\"enumIdentityRequirements\",\"name\":\"identity\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"diligenceRequirements\",\"type\":\"string\"}],\"internalType\":\"structBountyTerms\",\"name\":\"bountyTerms\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"agreementURI\",\"type\":\"string\"}],\"internalType\":\"structAgreementDetailsV1\",\"name\":\"details\",\"type\":\"tuple\"}],\"name\":\"adoptSafeHarbor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"adopter\",\"type\":\"address\"}],\"name\":\"getAgreement\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"protocolName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contact\",\"type\":\"string\"}],\"internalType\":\"structContact[]\",\"name\":\"contactDetails\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetRecoveryAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"internalType\":\"enumChildContractScope\",\"name\":\"childContractScope\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structAccount[]\",\"name\":\"accounts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"internalType\":\"structChain[]\",\"name\":\"chains\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"bountyPercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bountyCapUSD\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"retainable\",\"type\":\"bool\"},{\"internalType\":\"enumIdentityRequirements\",\"name\":\"identity\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"diligenceRequirements\",\"type\":\"string\"}],\"internalType\":\"structBountyTerms\",\"name\":\"bountyTerms\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"agreementURI\",\"type\":\"string\"}],\"internalType\":\"structAgreementDetailsV1\",\"name\":\"details\",\"type\":\"tuple\"}],\"name\":\"getTypedDataHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"protocolName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contact\",\"type\":\"string\"}],\"internalType\":\"structContact[]\",\"name\":\"contactDetails\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetRecoveryAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"internalType\":\"enumChildContractScope\",\"name\":\"childContractScope\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structAccount[]\",\"name\":\"accounts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"internalType\":\"structChain[]\",\"name\":\"chains\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"bountyPercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bountyCapUSD\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"retainable\",\"type\":\"bool\"},{\"internalType\":\"enumIdentityRequirements\",\"name\":\"identity\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"diligenceRequirements\",\"type\":\"string\"}],\"internalType\":\"structBountyTerms\",\"name\":\"bountyTerms\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"agreementURI\",\"type\":\"string\"}],\"internalType\":\"structAgreementDetailsV1\",\"name\":\"details\",\"type\":\"tuple\"}],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wantSigner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isSignatureValid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"protocolName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contact\",\"type\":\"string\"}],\"internalType\":\"structContact[]\",\"name\":\"contactDetails\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"assetRecoveryAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"internalType\":\"enumChildContractScope\",\"name\":\"childContractScope\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structAccount[]\",\"name\":\"accounts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"internalType\":\"structChain[]\",\"name\":\"chains\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"bountyPercentage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bountyCapUSD\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"retainable\",\"type\":\"bool\"},{\"internalType\":\"enumIdentityRequirements\",\"name\":\"identity\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"diligenceRequirements\",\"type\":\"string\"}],\"internalType\":\"structBountyTerms\",\"name\":\"bountyTerms\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"agreementURI\",\"type\":\"string\"}],\"internalType\":\"structAgreementDetailsV1\",\"name\":\"details\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"internalType\":\"enumChildContractScope\",\"name\":\"childContractScope\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structAccount\",\"name\":\"account\",\"type\":\"tuple\"}],\"name\":\"validateAccount\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agreementAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"accountAddress\",\"type\":\"address\"},{\"internalType\":\"enumChildContractScope\",\"name\":\"childContractScope\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structAccount\",\"name\":\"account\",\"type\":\"tuple\"}],\"name\":\"validateAccountByAddress\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// SafeharborABI is the input ABI used to generate the binding from.
// Deprecated: Use SafeharborMetaData.ABI instead.
var SafeharborABI = SafeharborMetaData.ABI

// Safeharbor is an auto generated Go binding around an Ethereum contract.
type Safeharbor struct {
	SafeharborCaller     // Read-only binding to the contract
	SafeharborTransactor // Write-only binding to the contract
	SafeharborFilterer   // Log filterer for contract events
}

// SafeharborCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeharborCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeharborTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeharborTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeharborFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeharborFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeharborSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeharborSession struct {
	Contract     *Safeharbor       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeharborCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeharborCallerSession struct {
	Contract *SafeharborCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// SafeharborTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeharborTransactorSession struct {
	Contract     *SafeharborTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SafeharborRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeharborRaw struct {
	Contract *Safeharbor // Generic contract binding to access the raw methods on
}

// SafeharborCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeharborCallerRaw struct {
	Contract *SafeharborCaller // Generic read-only contract binding to access the raw methods on
}

// SafeharborTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeharborTransactorRaw struct {
	Contract *SafeharborTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeharbor creates a new instance of Safeharbor, bound to a specific deployed contract.
func NewSafeharbor(address common.Address, backend bind.ContractBackend) (*Safeharbor, error) {
	contract, err := bindSafeharbor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Safeharbor{SafeharborCaller: SafeharborCaller{contract: contract}, SafeharborTransactor: SafeharborTransactor{contract: contract}, SafeharborFilterer: SafeharborFilterer{contract: contract}}, nil
}

// NewSafeharborCaller creates a new read-only instance of Safeharbor, bound to a specific deployed contract.
func NewSafeharborCaller(address common.Address, caller bind.ContractCaller) (*SafeharborCaller, error) {
	contract, err := bindSafeharbor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeharborCaller{contract: contract}, nil
}

// NewSafeharborTransactor creates a new write-only instance of Safeharbor, bound to a specific deployed contract.
func NewSafeharborTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeharborTransactor, error) {
	contract, err := bindSafeharbor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeharborTransactor{contract: contract}, nil
}

// NewSafeharborFilterer creates a new log filterer instance of Safeharbor, bound to a specific deployed contract.
func NewSafeharborFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeharborFilterer, error) {
	contract, err := bindSafeharbor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeharborFilterer{contract: contract}, nil
}

// bindSafeharbor binds a generic wrapper to an already deployed contract.
func bindSafeharbor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SafeharborMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Safeharbor *SafeharborRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Safeharbor.Contract.SafeharborCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Safeharbor *SafeharborRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safeharbor.Contract.SafeharborTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Safeharbor *SafeharborRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Safeharbor.Contract.SafeharborTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Safeharbor *SafeharborCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Safeharbor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Safeharbor *SafeharborTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Safeharbor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Safeharbor *SafeharborTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Safeharbor.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Safeharbor *SafeharborCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Safeharbor.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Safeharbor *SafeharborSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Safeharbor.Contract.DOMAINSEPARATOR(&_Safeharbor.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Safeharbor *SafeharborCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Safeharbor.Contract.DOMAINSEPARATOR(&_Safeharbor.CallOpts)
}

// GetAgreement is a free data retrieval call binding the contract method 0x295c1fed.
//
// Solidity: function getAgreement(address adopter) view returns(address)
func (_Safeharbor *SafeharborCaller) GetAgreement(opts *bind.CallOpts, adopter common.Address) (common.Address, error) {
	var out []interface{}
	err := _Safeharbor.contract.Call(opts, &out, "getAgreement", adopter)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAgreement is a free data retrieval call binding the contract method 0x295c1fed.
//
// Solidity: function getAgreement(address adopter) view returns(address)
func (_Safeharbor *SafeharborSession) GetAgreement(adopter common.Address) (common.Address, error) {
	return _Safeharbor.Contract.GetAgreement(&_Safeharbor.CallOpts, adopter)
}

// GetAgreement is a free data retrieval call binding the contract method 0x295c1fed.
//
// Solidity: function getAgreement(address adopter) view returns(address)
func (_Safeharbor *SafeharborCallerSession) GetAgreement(adopter common.Address) (common.Address, error) {
	return _Safeharbor.Contract.GetAgreement(&_Safeharbor.CallOpts, adopter)
}

// GetTypedDataHash is a free data retrieval call binding the contract method 0xc77e7b18.
//
// Solidity: function getTypedDataHash((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) view returns(bytes32)
func (_Safeharbor *SafeharborCaller) GetTypedDataHash(opts *bind.CallOpts, details AgreementDetailsV1) ([32]byte, error) {
	var out []interface{}
	err := _Safeharbor.contract.Call(opts, &out, "getTypedDataHash", details)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetTypedDataHash is a free data retrieval call binding the contract method 0xc77e7b18.
//
// Solidity: function getTypedDataHash((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) view returns(bytes32)
func (_Safeharbor *SafeharborSession) GetTypedDataHash(details AgreementDetailsV1) ([32]byte, error) {
	return _Safeharbor.Contract.GetTypedDataHash(&_Safeharbor.CallOpts, details)
}

// GetTypedDataHash is a free data retrieval call binding the contract method 0xc77e7b18.
//
// Solidity: function getTypedDataHash((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) view returns(bytes32)
func (_Safeharbor *SafeharborCallerSession) GetTypedDataHash(details AgreementDetailsV1) ([32]byte, error) {
	return _Safeharbor.Contract.GetTypedDataHash(&_Safeharbor.CallOpts, details)
}

// Hash is a free data retrieval call binding the contract method 0xbd78a34a.
//
// Solidity: function hash((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) pure returns(bytes32)
func (_Safeharbor *SafeharborCaller) Hash(opts *bind.CallOpts, details AgreementDetailsV1) ([32]byte, error) {
	var out []interface{}
	err := _Safeharbor.contract.Call(opts, &out, "hash", details)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Hash is a free data retrieval call binding the contract method 0xbd78a34a.
//
// Solidity: function hash((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) pure returns(bytes32)
func (_Safeharbor *SafeharborSession) Hash(details AgreementDetailsV1) ([32]byte, error) {
	return _Safeharbor.Contract.Hash(&_Safeharbor.CallOpts, details)
}

// Hash is a free data retrieval call binding the contract method 0xbd78a34a.
//
// Solidity: function hash((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) pure returns(bytes32)
func (_Safeharbor *SafeharborCallerSession) Hash(details AgreementDetailsV1) ([32]byte, error) {
	return _Safeharbor.Contract.Hash(&_Safeharbor.CallOpts, details)
}

// IsSignatureValid is a free data retrieval call binding the contract method 0xfd746e43.
//
// Solidity: function isSignatureValid(address wantSigner, bytes32 hash, bytes signature) view returns(bool)
func (_Safeharbor *SafeharborCaller) IsSignatureValid(opts *bind.CallOpts, wantSigner common.Address, hash [32]byte, signature []byte) (bool, error) {
	var out []interface{}
	err := _Safeharbor.contract.Call(opts, &out, "isSignatureValid", wantSigner, hash, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSignatureValid is a free data retrieval call binding the contract method 0xfd746e43.
//
// Solidity: function isSignatureValid(address wantSigner, bytes32 hash, bytes signature) view returns(bool)
func (_Safeharbor *SafeharborSession) IsSignatureValid(wantSigner common.Address, hash [32]byte, signature []byte) (bool, error) {
	return _Safeharbor.Contract.IsSignatureValid(&_Safeharbor.CallOpts, wantSigner, hash, signature)
}

// IsSignatureValid is a free data retrieval call binding the contract method 0xfd746e43.
//
// Solidity: function isSignatureValid(address wantSigner, bytes32 hash, bytes signature) view returns(bool)
func (_Safeharbor *SafeharborCallerSession) IsSignatureValid(wantSigner common.Address, hash [32]byte, signature []byte) (bool, error) {
	return _Safeharbor.Contract.IsSignatureValid(&_Safeharbor.CallOpts, wantSigner, hash, signature)
}

// ValidateAccount is a free data retrieval call binding the contract method 0x98e27ea9.
//
// Solidity: function validateAccount((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details, (address,uint8,bytes) account) view returns(bool)
func (_Safeharbor *SafeharborCaller) ValidateAccount(opts *bind.CallOpts, details AgreementDetailsV1, account Account) (bool, error) {
	var out []interface{}
	err := _Safeharbor.contract.Call(opts, &out, "validateAccount", details, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidateAccount is a free data retrieval call binding the contract method 0x98e27ea9.
//
// Solidity: function validateAccount((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details, (address,uint8,bytes) account) view returns(bool)
func (_Safeharbor *SafeharborSession) ValidateAccount(details AgreementDetailsV1, account Account) (bool, error) {
	return _Safeharbor.Contract.ValidateAccount(&_Safeharbor.CallOpts, details, account)
}

// ValidateAccount is a free data retrieval call binding the contract method 0x98e27ea9.
//
// Solidity: function validateAccount((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details, (address,uint8,bytes) account) view returns(bool)
func (_Safeharbor *SafeharborCallerSession) ValidateAccount(details AgreementDetailsV1, account Account) (bool, error) {
	return _Safeharbor.Contract.ValidateAccount(&_Safeharbor.CallOpts, details, account)
}

// ValidateAccountByAddress is a free data retrieval call binding the contract method 0xdb97c84b.
//
// Solidity: function validateAccountByAddress(address agreementAddress, (address,uint8,bytes) account) view returns(bool)
func (_Safeharbor *SafeharborCaller) ValidateAccountByAddress(opts *bind.CallOpts, agreementAddress common.Address, account Account) (bool, error) {
	var out []interface{}
	err := _Safeharbor.contract.Call(opts, &out, "validateAccountByAddress", agreementAddress, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidateAccountByAddress is a free data retrieval call binding the contract method 0xdb97c84b.
//
// Solidity: function validateAccountByAddress(address agreementAddress, (address,uint8,bytes) account) view returns(bool)
func (_Safeharbor *SafeharborSession) ValidateAccountByAddress(agreementAddress common.Address, account Account) (bool, error) {
	return _Safeharbor.Contract.ValidateAccountByAddress(&_Safeharbor.CallOpts, agreementAddress, account)
}

// ValidateAccountByAddress is a free data retrieval call binding the contract method 0xdb97c84b.
//
// Solidity: function validateAccountByAddress(address agreementAddress, (address,uint8,bytes) account) view returns(bool)
func (_Safeharbor *SafeharborCallerSession) ValidateAccountByAddress(agreementAddress common.Address, account Account) (bool, error) {
	return _Safeharbor.Contract.ValidateAccountByAddress(&_Safeharbor.CallOpts, agreementAddress, account)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_Safeharbor *SafeharborCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Safeharbor.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_Safeharbor *SafeharborSession) Version() (string, error) {
	return _Safeharbor.Contract.Version(&_Safeharbor.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_Safeharbor *SafeharborCallerSession) Version() (string, error) {
	return _Safeharbor.Contract.Version(&_Safeharbor.CallOpts)
}

// AdoptSafeHarbor is a paid mutator transaction binding the contract method 0x121e9ffe.
//
// Solidity: function adoptSafeHarbor((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) returns()
func (_Safeharbor *SafeharborTransactor) AdoptSafeHarbor(opts *bind.TransactOpts, details AgreementDetailsV1) (*types.Transaction, error) {
	return _Safeharbor.contract.Transact(opts, "adoptSafeHarbor", details)
}

// AdoptSafeHarbor is a paid mutator transaction binding the contract method 0x121e9ffe.
//
// Solidity: function adoptSafeHarbor((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) returns()
func (_Safeharbor *SafeharborSession) AdoptSafeHarbor(details AgreementDetailsV1) (*types.Transaction, error) {
	return _Safeharbor.Contract.AdoptSafeHarbor(&_Safeharbor.TransactOpts, details)
}

// AdoptSafeHarbor is a paid mutator transaction binding the contract method 0x121e9ffe.
//
// Solidity: function adoptSafeHarbor((string,(string,string)[],(address,(address,uint8,bytes)[],uint256)[],(uint256,uint256,bool,uint8,string),string) details) returns()
func (_Safeharbor *SafeharborTransactorSession) AdoptSafeHarbor(details AgreementDetailsV1) (*types.Transaction, error) {
	return _Safeharbor.Contract.AdoptSafeHarbor(&_Safeharbor.TransactOpts, details)
}

// SafeharborSafeHarborAdoptionIterator is returned from FilterSafeHarborAdoption and is used to iterate over the raw logs and unpacked data for SafeHarborAdoption events raised by the Safeharbor contract.
type SafeharborSafeHarborAdoptionIterator struct {
	Event *SafeharborSafeHarborAdoption // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SafeharborSafeHarborAdoptionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafeharborSafeHarborAdoption)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SafeharborSafeHarborAdoption)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SafeharborSafeHarborAdoptionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafeharborSafeHarborAdoptionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafeharborSafeHarborAdoption represents a SafeHarborAdoption event raised by the Safeharbor contract.
type SafeharborSafeHarborAdoption struct {
	Entity     common.Address
	OldDetails common.Address
	NewDetails common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSafeHarborAdoption is a free log retrieval operation binding the contract event 0xfb9c334c719c97ecac9e4d31dec8572d1e2cf193a6af229da967437a30dc7010.
//
// Solidity: event SafeHarborAdoption(address indexed entity, address oldDetails, address newDetails)
func (_Safeharbor *SafeharborFilterer) FilterSafeHarborAdoption(opts *bind.FilterOpts, entity []common.Address) (*SafeharborSafeHarborAdoptionIterator, error) {

	var entityRule []interface{}
	for _, entityItem := range entity {
		entityRule = append(entityRule, entityItem)
	}

	logs, sub, err := _Safeharbor.contract.FilterLogs(opts, "SafeHarborAdoption", entityRule)
	if err != nil {
		return nil, err
	}
	return &SafeharborSafeHarborAdoptionIterator{contract: _Safeharbor.contract, event: "SafeHarborAdoption", logs: logs, sub: sub}, nil
}

// WatchSafeHarborAdoption is a free log subscription operation binding the contract event 0xfb9c334c719c97ecac9e4d31dec8572d1e2cf193a6af229da967437a30dc7010.
//
// Solidity: event SafeHarborAdoption(address indexed entity, address oldDetails, address newDetails)
func (_Safeharbor *SafeharborFilterer) WatchSafeHarborAdoption(opts *bind.WatchOpts, sink chan<- *SafeharborSafeHarborAdoption, entity []common.Address) (event.Subscription, error) {

	var entityRule []interface{}
	for _, entityItem := range entity {
		entityRule = append(entityRule, entityItem)
	}

	logs, sub, err := _Safeharbor.contract.WatchLogs(opts, "SafeHarborAdoption", entityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafeharborSafeHarborAdoption)
				if err := _Safeharbor.contract.UnpackLog(event, "SafeHarborAdoption", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSafeHarborAdoption is a log parse operation binding the contract event 0xfb9c334c719c97ecac9e4d31dec8572d1e2cf193a6af229da967437a30dc7010.
//
// Solidity: event SafeHarborAdoption(address indexed entity, address oldDetails, address newDetails)
func (_Safeharbor *SafeharborFilterer) ParseSafeHarborAdoption(log types.Log) (*SafeharborSafeHarborAdoption, error) {
	event := new(SafeharborSafeHarborAdoption)
	if err := _Safeharbor.contract.UnpackLog(event, "SafeHarborAdoption", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
