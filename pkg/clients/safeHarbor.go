package clients

import (
	"SHDB/pkg/etherscan"
	"context"
	"math/big"
	"strings"

	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type AgreementDetailsV1 struct {
	ProtocolName   string
	ContactDetails string
	Chains         []ChainSH
	BountyTerms    BountyTermsSH
	AgreementURI   string
}

type ChainSH struct {
	AssetRecoveryAddress common.Address
	Accounts             []AccountSH
	ID                   *big.Int
}

type AccountSH struct {
	AccountAddress     common.Address
	ChildContractScope uint8
	Signature          []byte
}

type BountyTermsSH struct {
	BountyPercentage *big.Int
	BountyCapUSD     *big.Int
	Verification     uint8
}

const contractABI = `[{"inputs":[{"components":[{"internalType":"string","name":"protocolName","type":"string"},{"internalType":"string","name":"contactDetails","type":"string"},{"components":[{"internalType":"address","name":"assetRecoveryAddress","type":"address"},{"components":[{"internalType":"address","name":"accountAddress","type":"address"},{"internalType":"enum ChildContractScope","name":"childContractScope","type":"uint8"},{"internalType":"bytes","name":"signature","type":"bytes"}],"internalType":"struct Account[]","name":"accounts","type":"tuple[]"},{"internalType":"uint256","name":"id","type":"uint256"}],"internalType":"struct Chain[]","name":"chains","type":"tuple[]"},{"components":[{"internalType":"uint256","name":"bountyPercentage","type":"uint256"},{"internalType":"uint256","name":"bountyCapUSD","type":"uint256"},{"internalType":"enum IdentityVerification","name":"verification","type":"uint8"}],"internalType":"struct BountyTerms","name":"bountyTerms","type":"tuple"},{"internalType":"string","name":"agreementURI","type":"string"}],"internalType":"struct AgreementDetailsV1","name":"_details","type":"tuple"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"getDetails","outputs":[{"components":[{"internalType":"string","name":"protocolName","type":"string"},{"internalType":"string","name":"contactDetails","type":"string"},{"components":[{"internalType":"address","name":"assetRecoveryAddress","type":"address"},{"components":[{"internalType":"address","name":"accountAddress","type":"address"},{"internalType":"enum ChildContractScope","name":"childContractScope","type":"uint8"},{"internalType":"bytes","name":"signature","type":"bytes"}],"internalType":"struct Account[]","name":"accounts","type":"tuple[]"},{"internalType":"uint256","name":"id","type":"uint256"}],"internalType":"struct Chain[]","name":"chains","type":"tuple[]"},{"components":[{"internalType":"uint256","name":"bountyPercentage","type":"uint256"},{"internalType":"uint256","name":"bountyCapUSD","type":"uint256"},{"internalType":"enum IdentityVerification","name":"verification","type":"uint8"}],"internalType":"struct BountyTerms","name":"bountyTerms","type":"tuple"},{"internalType":"string","name":"agreementURI","type":"string"}],"internalType":"struct AgreementDetailsV1","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"}]`

func GetSafeHarborAdoptions(apiKey string, rpcClient *rpc.Client) ([]Client, error) {
	const (
		chainID               = int64(11155111)
		factoryAddress        = "0x2697beaf5c0ddd825b952e0d5d41918c01f85dbb"
		agreementCreatedTopic = "0xfb9c334c719c97ecac9e4d31dec8572d1e2cf193a6af229da967437a30dc7010"
	)

	address, err := web3.HexToAddress(factoryAddress)
	if err != nil {
		return nil, err
	}

	topic0, err := web3.HexToHash(agreementCreatedTopic)
	if err != nil {
		return nil, err
	}

	logs, err := etherscan.FetchLogs(chainID, apiKey, *address, *topic0, 0)
	if err != nil {
		return nil, err
	}

	entitySafeHarborAdoption := make(map[web3.Address]web3.Address)
	for _, log := range logs {
		if len(log.Topics) < 2 || len(log.Data) < 64 {
			continue
		}

		entity, err := web3.BytesToAddress(log.Topics[1][12:])
		if err != nil {
			return nil, err
		}

		safeHarbor, err := web3.BytesToAddress(log.Data[44:64])
		if err != nil {
			return nil, err
		}

		entitySafeHarborAdoption[*entity] = *safeHarbor
	}

	addresses := make([]web3.Address, 0, len(entitySafeHarborAdoption))
	for _, address := range entitySafeHarborAdoption {
		addresses = append(addresses, address)
	}

	return fetchAgreementDetails(addresses, rpcClient)
}

func uint8ToChildContractScope(value uint8) ChildContractScope {
	switch value {
	case 0:
		return ChildContractScopeNone
	case 1:
		return ChildContractScopeAll
	default:
		return ChildContractScopeNone // Default to None if unknown value
	}
}

func uint8ToIdentityVerification(value uint8) IdentityVerification {
	switch value {
	case 0:
		return Retainable
	case 1:
		return Immunefi
	case 2:
		return Bugcrowd
	case 3:
		return Hackerone
	default:
		return Retainable // Default to Retainable if unknown value
	}
}

func fetchAgreementDetails(addresses []web3.Address, rpcClient *rpc.Client) ([]Client, error) {
	client := ethclient.NewClient(rpcClient)
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	var clients []Client

	for _, address := range addresses {
		commonAddress := address.ToCommon()
		callData, err := parsedABI.Pack("getDetails")
		if err != nil {
			return nil, err
		}

		data, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &commonAddress,
			Data: callData,
		}, nil)
		if err != nil {
			return nil, err
		}

		var result struct {
			AgreementDetails AgreementDetailsV1
		}

		if err := parsedABI.UnpackIntoInterface(&result, "getDetails", data); err != nil {
			return nil, err
		}

		// Convert AgreementDetailsV1 to Client
		newClient := Client{
			ProtocolName:   result.AgreementDetails.ProtocolName,
			AgreementURI:   result.AgreementDetails.AgreementURI,
			ContactDetails: result.AgreementDetails.ContactDetails,
			BountyTerms: BountyTerms{
				BountyPercentage: int(result.AgreementDetails.BountyTerms.BountyPercentage.Int64()),
				BountyCapUSD:     int(result.AgreementDetails.BountyTerms.BountyCapUSD.Int64()),
				Retainable:       uint8ToIdentityVerification(result.AgreementDetails.BountyTerms.Verification),
			},
		}

		// Convert Chains
		for _, chain := range result.AgreementDetails.Chains {
			newChain := Chain{
				AssetRecoveryAddress: web3.Address(chain.AssetRecoveryAddress),
				ChainID:              chain.ID.Int64(),
			}
			newClient.Chains = append(newClient.Chains, newChain)

			// Convert Accounts
			for _, account := range chain.Accounts {
				newAccount := Account{
					Address:            web3.Address(account.AccountAddress),
					ChildContractScope: uint8ToChildContractScope(account.ChildContractScope),
					SignatureValid:     len(account.Signature) > 0,
					ChainIDs:           []int64{chain.ID.Int64()},
					LastQuery:          make(map[int64]int),
				}
				newClient.Accounts = append(newClient.Accounts, newAccount)
			}
		}

		clients = append(clients, newClient)
	}

	return clients, nil
}
