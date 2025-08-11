// Code generated-like - minimal binding for V2 getDetails. Do not edit without updating abi.json accordingly.

package adoptiondetailsv2

import (
	_ "embed"
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = abi.ConvertType
)

// Account mirrors the V2 ABI struct.
type Account struct {
	AccountAddress     string
	ChildContractScope uint8
}

// Chain mirrors the V2 ABI struct.
type Chain struct {
	AssetRecoveryAddress string
	Accounts             []Account
	Caip2ChainId         string
}

// BountyTerms mirrors the V2 ABI struct.
type BountyTerms struct {
	BountyPercentage      *big.Int
	BountyCapUSD          *big.Int
	Retainable            bool
	Identity              uint8
	DiligenceRequirements string
	AggregateBountyCapUSD *big.Int
}

// Contact mirrors the V2 ABI struct.
type Contact struct {
	Name    string
	Contact string
}

// AgreementDetailsV2 mirrors the V2 ABI struct.
type AgreementDetailsV2 struct {
	ProtocolName   string
	ContactDetails []Contact
	Chains         []Chain
	BountyTerms    BountyTerms
	AgreementURI   string
}

//go:embed abi.json
var adoptionDetailsV2ABI string

// AdoptiondetailsV2MetaData contains the minimal ABI (getDetails) for v2.
var AdoptiondetailsV2MetaData = &bind.MetaData{
	ABI: adoptionDetailsV2ABI,
}

// AdoptiondetailsV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use AdoptiondetailsV2MetaData.ABI instead.
var AdoptiondetailsV2ABI = AdoptiondetailsV2MetaData.ABI

// Adoptiondetails is a minimal binding around the V2 adoption details contract.
type Adoptiondetails struct {
	AdoptiondetailsCaller
}

// AdoptiondetailsCaller is a read-only binding around an Ethereum contract.
type AdoptiondetailsCaller struct {
	contract *bind.BoundContract
}

// NewAdoptiondetails creates a new instance bound to a specific deployed contract.
func NewAdoptiondetails(address common.Address, backend bind.ContractBackend) (*Adoptiondetails, error) {
	contract, err := bindAdoptiondetails(address, backend)
	if err != nil {
		return nil, err
	}
	return &Adoptiondetails{AdoptiondetailsCaller: AdoptiondetailsCaller{contract: contract}}, nil
}

// bindAdoptiondetails binds a generic wrapper to an already deployed contract.
func bindAdoptiondetails(address common.Address, caller bind.ContractCaller) (*bind.BoundContract, error) {
	// Parse ABI directly from embedded JSON to avoid any MetaData parsing issues
	parsed, err := abi.JSON(strings.NewReader(adoptionDetailsV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, nil, nil), nil
}

// GetDetails is a free data retrieval call binding the contract method getDetails.
func (_Adoptiondetails *AdoptiondetailsCaller) GetDetails(opts *bind.CallOpts) (AgreementDetailsV2, error) {
	var out []interface{}
	err := _Adoptiondetails.contract.Call(opts, &out, "getDetails")
	if err != nil {
		return *new(AgreementDetailsV2), err
	}
	out0 := *abi.ConvertType(out[0], new(AgreementDetailsV2)).(*AgreementDetailsV2)
	return out0, err
}
