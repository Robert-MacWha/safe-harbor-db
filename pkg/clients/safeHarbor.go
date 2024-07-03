package clients

import (
	"SHDB/pkg/etherscan"
	"context"
	"fmt"
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
	ProtocolName        string
	ContactDetails      string
	Chains              []ChainSH
	BountyTerms         BountyTermsSH
	AgreementURI        string
	RegistryTransaction string
	RegistryExplorerURL string
	CreatedAt           int
	Entity              common.Address
	MockData            bool
}

type ChainSH struct {
	AssetRecoveryAddress common.Address
	Accounts             []AccountSH
	ID                   *big.Int
}

type AccountSH struct {
	AccountAddress     common.Address
	ChildContractScope ChildContractScope
	Signature          []byte
}

type ChildContractScope string

const (
	ChildContractScopeNone ChildContractScope = "None"
	ChildContractScopeAll  ChildContractScope = "All"
)

type BountyTermsSH struct {
	BountyPercentage *big.Int
	BountyCapUSD     *big.Int
	Verification     IdentityVerification
}

type IdentityVerification string

const (
	Retainable IdentityVerification = "Retainable"
	Immunefi   IdentityVerification = "Immunefi"
	Bugcrowd   IdentityVerification = "Bugcrowd"
	Hackerone  IdentityVerification = "Hackerone"
)

var explorerPrefixes = map[int64]string{
	1:        "https://etherscan.io/tx/",
	11155111: "https://sepolia.etherscan.io/tx/",
}

const contractABI = `[{"inputs":[{"components":[{"internalType":"string","name":"protocolName","type":"string"},{"internalType":"string","name":"contactDetails","type":"string"},{"components":[{"internalType":"address","name":"assetRecoveryAddress","type":"address"},{"components":[{"internalType":"address","name":"accountAddress","type":"address"},{"internalType":"enum ChildContractScope","name":"childContractScope","type":"uint8"},{"internalType":"bytes","name":"signature","type":"bytes"}],"internalType":"struct Account[]","name":"accounts","type":"tuple[]"},{"internalType":"uint256","name":"id","type":"uint256"}],"internalType":"struct Chain[]","name":"chains","type":"tuple[]"},{"components":[{"internalType":"uint256","name":"bountyPercentage","type":"uint256"},{"internalType":"uint256","name":"bountyCapUSD","type":"uint256"},{"internalType":"enum IdentityVerification","name":"verification","type":"uint8"}],"internalType":"struct BountyTerms","name":"bountyTerms","type":"tuple"},{"internalType":"string","name":"agreementURI","type":"string"}],"internalType":"struct AgreementDetailsV1","name":"_details","type":"tuple"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"getDetails","outputs":[{"components":[{"internalType":"string","name":"protocolName","type":"string"},{"internalType":"string","name":"contactDetails","type":"string"},{"components":[{"internalType":"address","name":"assetRecoveryAddress","type":"address"},{"components":[{"internalType":"address","name":"accountAddress","type":"address"},{"internalType":"enum ChildContractScope","name":"childContractScope","type":"uint8"},{"internalType":"bytes","name":"signature","type":"bytes"}],"internalType":"struct Account[]","name":"accounts","type":"tuple[]"},{"internalType":"uint256","name":"id","type":"uint256"}],"internalType":"struct Chain[]","name":"chains","type":"tuple[]"},{"components":[{"internalType":"uint256","name":"bountyPercentage","type":"uint256"},{"internalType":"uint256","name":"bountyCapUSD","type":"uint256"},{"internalType":"enum IdentityVerification","name":"verification","type":"uint8"}],"internalType":"struct BountyTerms","name":"bountyTerms","type":"tuple"},{"internalType":"string","name":"agreementURI","type":"string"}],"internalType":"struct AgreementDetailsV1","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"}]`

func GetSafeHarborAdoptions(apiKey string, rpcClient *rpc.Client) ([]AgreementDetailsV1, error) {
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
	txHashes := make(map[web3.Address]web3.Hash)
	blockTimes := make(map[web3.Address]web3.BigInt)

	// Logs are returned in reverse chronological order, thus we use maps
	// to store the latest log for each entity
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
		txHashes[*entity] = log.TransactionHash
		blockTimes[*entity] = log.TimeStamp
	}

	return fetchAgreementDetails(entitySafeHarborAdoption, txHashes, blockTimes, chainID, rpcClient)
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

func fetchAgreementDetails(
	entitySafeHarborAdoption map[web3.Address]web3.Address,
	txHashes map[web3.Address]web3.Hash,
	blockTimes map[web3.Address]web3.BigInt,
	chainID int64,
	rpcClient *rpc.Client,
) ([]AgreementDetailsV1, error) {
	client := ethclient.NewClient(rpcClient)
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	var agreementDetails []AgreementDetailsV1

	for entity, safeHarborAddress := range entitySafeHarborAdoption {
		commonAddress := safeHarborAddress.ToCommon()
		callData, err := parsedABI.Pack("getDetails")
		if err != nil {
			return nil, fmt.Errorf("failed to pack call data: %w", err)
		}

		data, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &commonAddress,
			Data: callData,
		}, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to call contract: %w", err)
		}

		var result struct {
			AgreementDetails struct {
				ProtocolName   string
				ContactDetails string
				Chains         []struct {
					AssetRecoveryAddress common.Address
					Accounts             []struct {
						AccountAddress     common.Address
						ChildContractScope uint8
						Signature          []byte
					}
					ID *big.Int
				}
				BountyTerms struct {
					BountyPercentage *big.Int
					BountyCapUSD     *big.Int
					Verification     uint8
				}
				AgreementURI string
			}
		}

		if err := parsedABI.UnpackIntoInterface(&result, "getDetails", data); err != nil {
			return nil, fmt.Errorf("failed to unpack result: %w", err)
		}

		// Convert the result to our new struct format
		agreement := AgreementDetailsV1{
			ProtocolName:        result.AgreementDetails.ProtocolName,
			ContactDetails:      result.AgreementDetails.ContactDetails,
			AgreementURI:        result.AgreementDetails.AgreementURI,
			Entity:              entity.ToCommon(),
			RegistryTransaction: txHashes[entity].String(),
			RegistryExplorerURL: fmt.Sprintf("%s%s", explorerPrefixes[chainID], txHashes[entity].String()),
			CreatedAt:           int(blockTimes[entity].Int64()),
		}

		// Convert Chains
		for _, chain := range result.AgreementDetails.Chains {
			chainSH := ChainSH{
				AssetRecoveryAddress: chain.AssetRecoveryAddress,
				ID:                   chain.ID,
			}
			for _, account := range chain.Accounts {
				accountSH := AccountSH{
					AccountAddress:     account.AccountAddress,
					ChildContractScope: uint8ToChildContractScope(account.ChildContractScope),
					Signature:          account.Signature,
				}
				chainSH.Accounts = append(chainSH.Accounts, accountSH)
			}
			agreement.Chains = append(agreement.Chains, chainSH)
		}

		// Convert BountyTerms
		agreement.BountyTerms = BountyTermsSH{
			BountyPercentage: result.AgreementDetails.BountyTerms.BountyPercentage,
			BountyCapUSD:     result.AgreementDetails.BountyTerms.BountyCapUSD,
			Verification:     uint8ToIdentityVerification(result.AgreementDetails.BountyTerms.Verification),
		}

		agreementDetails = append(agreementDetails, agreement)
	}

	return agreementDetails, nil
}

func (a *AgreementDetailsV1) ToMonitored(rpcClient *rpc.Client, apiKey string) (*Monitored, error) {
	monitored := &Monitored{
		MockData: a.MockData,
	}

	for _, chain := range a.Chains {
		for _, account := range chain.Accounts {
			chainID := chain.ID.Int64()
			accountM := AccountM{
				Address: account.AccountAddress.String(),
				Chains:  []int{int(chainID)},
			}

			// Fill name for the account
			err := fillName(&accountM, apiKey)
			if err != nil {
				return nil, fmt.Errorf("failed to fill name for account %s: %w", account.AccountAddress.String(), err)
			}

			// Find children if necessary
			if account.ChildContractScope == ChildContractScopeAll {
				children, err := findAllChildren(
					*web3.CommonToAddress(account.AccountAddress),
					[]int64{chainID},
					account.ChildContractScope,
					nil, // AgreementDetailsV1 doesn't have events, so pass nil
					rpcClient,
					apiKey,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to find children for account %s: %w", account.AccountAddress.String(), err)
				}
				accountM.Children = children
			}

			monitored.Addresses = append(monitored.Addresses, accountM)
		}
	}

	return monitored, nil
}
