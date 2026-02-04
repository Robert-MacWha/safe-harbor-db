// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package safeharbor_v3

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

// SafeHarborV3MetaData contains all meta data concerning the SafeHarborV3 contract.
var SafeHarborV3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_legacyRegistry\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_adopters\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"SafeHarborRegistry__NoAgreement\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"legacyRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"migratedCount\",\"type\":\"uint256\"}],\"name\":\"LegacyDataMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"adopter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agreementAddress\",\"type\":\"address\"}],\"name\":\"SafeHarborAdoption\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agreementAddress\",\"type\":\"address\"}],\"name\":\"adoptSafeHarbor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_adopter\",\"type\":\"address\"}],\"name\":\"getAgreement\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// SafeHarborV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use SafeHarborV3MetaData.ABI instead.
var SafeHarborV3ABI = SafeHarborV3MetaData.ABI

// SafeHarborV3 is an auto generated Go binding around an Ethereum contract.
type SafeHarborV3 struct {
	SafeHarborV3Caller     // Read-only binding to the contract
	SafeHarborV3Transactor // Write-only binding to the contract
	SafeHarborV3Filterer   // Log filterer for contract events
}

// SafeHarborV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type SafeHarborV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeHarborV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeHarborV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeHarborV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeHarborV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeHarborV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeHarborV3Session struct {
	Contract     *SafeHarborV3     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeHarborV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeHarborV3CallerSession struct {
	Contract *SafeHarborV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// SafeHarborV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeHarborV3TransactorSession struct {
	Contract     *SafeHarborV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SafeHarborV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type SafeHarborV3Raw struct {
	Contract *SafeHarborV3 // Generic contract binding to access the raw methods on
}

// SafeHarborV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeHarborV3CallerRaw struct {
	Contract *SafeHarborV3Caller // Generic read-only contract binding to access the raw methods on
}

// SafeHarborV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeHarborV3TransactorRaw struct {
	Contract *SafeHarborV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeHarborV3 creates a new instance of SafeHarborV3, bound to a specific deployed contract.
func NewSafeHarborV3(address common.Address, backend bind.ContractBackend) (*SafeHarborV3, error) {
	contract, err := bindSafeHarborV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeHarborV3{SafeHarborV3Caller: SafeHarborV3Caller{contract: contract}, SafeHarborV3Transactor: SafeHarborV3Transactor{contract: contract}, SafeHarborV3Filterer: SafeHarborV3Filterer{contract: contract}}, nil
}

// NewSafeHarborV3Caller creates a new read-only instance of SafeHarborV3, bound to a specific deployed contract.
func NewSafeHarborV3Caller(address common.Address, caller bind.ContractCaller) (*SafeHarborV3Caller, error) {
	contract, err := bindSafeHarborV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeHarborV3Caller{contract: contract}, nil
}

// NewSafeHarborV3Transactor creates a new write-only instance of SafeHarborV3, bound to a specific deployed contract.
func NewSafeHarborV3Transactor(address common.Address, transactor bind.ContractTransactor) (*SafeHarborV3Transactor, error) {
	contract, err := bindSafeHarborV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeHarborV3Transactor{contract: contract}, nil
}

// NewSafeHarborV3Filterer creates a new log filterer instance of SafeHarborV3, bound to a specific deployed contract.
func NewSafeHarborV3Filterer(address common.Address, filterer bind.ContractFilterer) (*SafeHarborV3Filterer, error) {
	contract, err := bindSafeHarborV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeHarborV3Filterer{contract: contract}, nil
}

// bindSafeHarborV3 binds a generic wrapper to an already deployed contract.
func bindSafeHarborV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SafeHarborV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeHarborV3 *SafeHarborV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeHarborV3.Contract.SafeHarborV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeHarborV3 *SafeHarborV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeHarborV3.Contract.SafeHarborV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeHarborV3 *SafeHarborV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeHarborV3.Contract.SafeHarborV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeHarborV3 *SafeHarborV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeHarborV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeHarborV3 *SafeHarborV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeHarborV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeHarborV3 *SafeHarborV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeHarborV3.Contract.contract.Transact(opts, method, params...)
}

// GetAgreement is a free data retrieval call binding the contract method 0x295c1fed.
//
// Solidity: function getAgreement(address _adopter) view returns(address)
func (_SafeHarborV3 *SafeHarborV3Caller) GetAgreement(opts *bind.CallOpts, _adopter common.Address) (common.Address, error) {
	var out []interface{}
	err := _SafeHarborV3.contract.Call(opts, &out, "getAgreement", _adopter)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAgreement is a free data retrieval call binding the contract method 0x295c1fed.
//
// Solidity: function getAgreement(address _adopter) view returns(address)
func (_SafeHarborV3 *SafeHarborV3Session) GetAgreement(_adopter common.Address) (common.Address, error) {
	return _SafeHarborV3.Contract.GetAgreement(&_SafeHarborV3.CallOpts, _adopter)
}

// GetAgreement is a free data retrieval call binding the contract method 0x295c1fed.
//
// Solidity: function getAgreement(address _adopter) view returns(address)
func (_SafeHarborV3 *SafeHarborV3CallerSession) GetAgreement(_adopter common.Address) (common.Address, error) {
	return _SafeHarborV3.Contract.GetAgreement(&_SafeHarborV3.CallOpts, _adopter)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_SafeHarborV3 *SafeHarborV3Caller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SafeHarborV3.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_SafeHarborV3 *SafeHarborV3Session) Version() (string, error) {
	return _SafeHarborV3.Contract.Version(&_SafeHarborV3.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_SafeHarborV3 *SafeHarborV3CallerSession) Version() (string, error) {
	return _SafeHarborV3.Contract.Version(&_SafeHarborV3.CallOpts)
}

// AdoptSafeHarbor is a paid mutator transaction binding the contract method 0x344fbd20.
//
// Solidity: function adoptSafeHarbor(address _agreementAddress) returns()
func (_SafeHarborV3 *SafeHarborV3Transactor) AdoptSafeHarbor(opts *bind.TransactOpts, _agreementAddress common.Address) (*types.Transaction, error) {
	return _SafeHarborV3.contract.Transact(opts, "adoptSafeHarbor", _agreementAddress)
}

// AdoptSafeHarbor is a paid mutator transaction binding the contract method 0x344fbd20.
//
// Solidity: function adoptSafeHarbor(address _agreementAddress) returns()
func (_SafeHarborV3 *SafeHarborV3Session) AdoptSafeHarbor(_agreementAddress common.Address) (*types.Transaction, error) {
	return _SafeHarborV3.Contract.AdoptSafeHarbor(&_SafeHarborV3.TransactOpts, _agreementAddress)
}

// AdoptSafeHarbor is a paid mutator transaction binding the contract method 0x344fbd20.
//
// Solidity: function adoptSafeHarbor(address _agreementAddress) returns()
func (_SafeHarborV3 *SafeHarborV3TransactorSession) AdoptSafeHarbor(_agreementAddress common.Address) (*types.Transaction, error) {
	return _SafeHarborV3.Contract.AdoptSafeHarbor(&_SafeHarborV3.TransactOpts, _agreementAddress)
}

// SafeHarborV3LegacyDataMigratedIterator is returned from FilterLegacyDataMigrated and is used to iterate over the raw logs and unpacked data for LegacyDataMigrated events raised by the SafeHarborV3 contract.
type SafeHarborV3LegacyDataMigratedIterator struct {
	Event *SafeHarborV3LegacyDataMigrated // Event containing the contract specifics and raw log

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
func (it *SafeHarborV3LegacyDataMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafeHarborV3LegacyDataMigrated)
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
		it.Event = new(SafeHarborV3LegacyDataMigrated)
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
func (it *SafeHarborV3LegacyDataMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafeHarborV3LegacyDataMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafeHarborV3LegacyDataMigrated represents a LegacyDataMigrated event raised by the SafeHarborV3 contract.
type SafeHarborV3LegacyDataMigrated struct {
	LegacyRegistry common.Address
	MigratedCount  *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterLegacyDataMigrated is a free log retrieval operation binding the contract event 0xd875cedf460af369381881dd6a86abfd9417752e8731620115d44a3c5dda165f.
//
// Solidity: event LegacyDataMigrated(address indexed legacyRegistry, uint256 migratedCount)
func (_SafeHarborV3 *SafeHarborV3Filterer) FilterLegacyDataMigrated(opts *bind.FilterOpts, legacyRegistry []common.Address) (*SafeHarborV3LegacyDataMigratedIterator, error) {

	var legacyRegistryRule []interface{}
	for _, legacyRegistryItem := range legacyRegistry {
		legacyRegistryRule = append(legacyRegistryRule, legacyRegistryItem)
	}

	logs, sub, err := _SafeHarborV3.contract.FilterLogs(opts, "LegacyDataMigrated", legacyRegistryRule)
	if err != nil {
		return nil, err
	}
	return &SafeHarborV3LegacyDataMigratedIterator{contract: _SafeHarborV3.contract, event: "LegacyDataMigrated", logs: logs, sub: sub}, nil
}

// WatchLegacyDataMigrated is a free log subscription operation binding the contract event 0xd875cedf460af369381881dd6a86abfd9417752e8731620115d44a3c5dda165f.
//
// Solidity: event LegacyDataMigrated(address indexed legacyRegistry, uint256 migratedCount)
func (_SafeHarborV3 *SafeHarborV3Filterer) WatchLegacyDataMigrated(opts *bind.WatchOpts, sink chan<- *SafeHarborV3LegacyDataMigrated, legacyRegistry []common.Address) (event.Subscription, error) {

	var legacyRegistryRule []interface{}
	for _, legacyRegistryItem := range legacyRegistry {
		legacyRegistryRule = append(legacyRegistryRule, legacyRegistryItem)
	}

	logs, sub, err := _SafeHarborV3.contract.WatchLogs(opts, "LegacyDataMigrated", legacyRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafeHarborV3LegacyDataMigrated)
				if err := _SafeHarborV3.contract.UnpackLog(event, "LegacyDataMigrated", log); err != nil {
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

// ParseLegacyDataMigrated is a log parse operation binding the contract event 0xd875cedf460af369381881dd6a86abfd9417752e8731620115d44a3c5dda165f.
//
// Solidity: event LegacyDataMigrated(address indexed legacyRegistry, uint256 migratedCount)
func (_SafeHarborV3 *SafeHarborV3Filterer) ParseLegacyDataMigrated(log types.Log) (*SafeHarborV3LegacyDataMigrated, error) {
	event := new(SafeHarborV3LegacyDataMigrated)
	if err := _SafeHarborV3.contract.UnpackLog(event, "LegacyDataMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafeHarborV3SafeHarborAdoptionIterator is returned from FilterSafeHarborAdoption and is used to iterate over the raw logs and unpacked data for SafeHarborAdoption events raised by the SafeHarborV3 contract.
type SafeHarborV3SafeHarborAdoptionIterator struct {
	Event *SafeHarborV3SafeHarborAdoption // Event containing the contract specifics and raw log

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
func (it *SafeHarborV3SafeHarborAdoptionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafeHarborV3SafeHarborAdoption)
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
		it.Event = new(SafeHarborV3SafeHarborAdoption)
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
func (it *SafeHarborV3SafeHarborAdoptionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafeHarborV3SafeHarborAdoptionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafeHarborV3SafeHarborAdoption represents a SafeHarborAdoption event raised by the SafeHarborV3 contract.
type SafeHarborV3SafeHarborAdoption struct {
	Adopter          common.Address
	AgreementAddress common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSafeHarborAdoption is a free log retrieval operation binding the contract event 0x1d458e1beb6e779286fe07bce764d5e0dfc5d75eb6caef4b19f4f2ce9c1a6d3e.
//
// Solidity: event SafeHarborAdoption(address indexed adopter, address agreementAddress)
func (_SafeHarborV3 *SafeHarborV3Filterer) FilterSafeHarborAdoption(opts *bind.FilterOpts, adopter []common.Address) (*SafeHarborV3SafeHarborAdoptionIterator, error) {

	var adopterRule []interface{}
	for _, adopterItem := range adopter {
		adopterRule = append(adopterRule, adopterItem)
	}

	logs, sub, err := _SafeHarborV3.contract.FilterLogs(opts, "SafeHarborAdoption", adopterRule)
	if err != nil {
		return nil, err
	}
	return &SafeHarborV3SafeHarborAdoptionIterator{contract: _SafeHarborV3.contract, event: "SafeHarborAdoption", logs: logs, sub: sub}, nil
}

// WatchSafeHarborAdoption is a free log subscription operation binding the contract event 0x1d458e1beb6e779286fe07bce764d5e0dfc5d75eb6caef4b19f4f2ce9c1a6d3e.
//
// Solidity: event SafeHarborAdoption(address indexed adopter, address agreementAddress)
func (_SafeHarborV3 *SafeHarborV3Filterer) WatchSafeHarborAdoption(opts *bind.WatchOpts, sink chan<- *SafeHarborV3SafeHarborAdoption, adopter []common.Address) (event.Subscription, error) {

	var adopterRule []interface{}
	for _, adopterItem := range adopter {
		adopterRule = append(adopterRule, adopterItem)
	}

	logs, sub, err := _SafeHarborV3.contract.WatchLogs(opts, "SafeHarborAdoption", adopterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafeHarborV3SafeHarborAdoption)
				if err := _SafeHarborV3.contract.UnpackLog(event, "SafeHarborAdoption", log); err != nil {
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

// ParseSafeHarborAdoption is a log parse operation binding the contract event 0x1d458e1beb6e779286fe07bce764d5e0dfc5d75eb6caef4b19f4f2ce9c1a6d3e.
//
// Solidity: event SafeHarborAdoption(address indexed adopter, address agreementAddress)
func (_SafeHarborV3 *SafeHarborV3Filterer) ParseSafeHarborAdoption(log types.Log) (*SafeHarborV3SafeHarborAdoption, error) {
	event := new(SafeHarborV3SafeHarborAdoption)
	if err := _SafeHarborV3.contract.UnpackLog(event, "SafeHarborAdoption", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
