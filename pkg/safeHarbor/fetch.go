package safeharbor

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const contractABI = `[{"inputs":[{"components":[{"internalType":"string","name":"protocolName","type":"string"},{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"contact","type":"string"}],"internalType":"struct Contact[]","name":"contactDetails","type":"tuple[]"},{"components":[{"internalType":"address","name":"assetRecoveryAddress","type":"address"},{"components":[{"internalType":"address","name":"accountAddress","type":"address"},{"internalType":"enum ChildContractScope","name":"childContractScope","type":"uint8"},{"internalType":"bytes","name":"signature","type":"bytes"}],"internalType":"struct Account[]","name":"accounts","type":"tuple[]"},{"internalType":"uint256","name":"id","type":"uint256"}],"internalType":"struct Chain[]","name":"chains","type":"tuple[]"},{"components":[{"internalType":"uint256","name":"bountyPercentage","type":"uint256"},{"internalType":"uint256","name":"bountyCapUSD","type":"uint256"},{"internalType":"bool","name":"retainable","type":"bool"},{"internalType":"enum IdentityRequirements","name":"identity","type":"uint8"},{"internalType":"string","name":"diligenceRequirements","type":"string"}],"internalType":"struct BountyTerms","name":"bountyTerms","type":"tuple"},{"internalType":"string","name":"agreementURI","type":"string"}],"internalType":"struct AgreementDetailsV1","name":"_details","type":"tuple"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"getDetails","outputs":[{"components":[{"internalType":"string","name":"protocolName","type":"string"},{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"contact","type":"string"}],"internalType":"struct Contact[]","name":"contactDetails","type":"tuple[]"},{"components":[{"internalType":"address","name":"assetRecoveryAddress","type":"address"},{"components":[{"internalType":"address","name":"accountAddress","type":"address"},{"internalType":"enum ChildContractScope","name":"childContractScope","type":"uint8"},{"internalType":"bytes","name":"signature","type":"bytes"}],"internalType":"struct Account[]","name":"accounts","type":"tuple[]"},{"internalType":"uint256","name":"id","type":"uint256"}],"internalType":"struct Chain[]","name":"chains","type":"tuple[]"},{"components":[{"internalType":"uint256","name":"bountyPercentage","type":"uint256"},{"internalType":"uint256","name":"bountyCapUSD","type":"uint256"},{"internalType":"bool","name":"retainable","type":"bool"},{"internalType":"enum IdentityRequirements","name":"identity","type":"uint8"},{"internalType":"string","name":"diligenceRequirements","type":"string"}],"internalType":"struct BountyTerms","name":"bountyTerms","type":"tuple"},{"internalType":"string","name":"agreementURI","type":"string"}],"internalType":"struct AgreementDetailsV1","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"}]`

func FetchAgreementDetails(
	blockNumber big.Int,
	registryTransaction string,
	registryChainId string,
	entity string,
	safeHarborAddress common.Address,
	client *ethclient.Client,
) (*SafeHarborAgreement, string, error) {
	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Pack the call data for getDetails()
	callData, err := parsedABI.Pack("getDetails")
	if err != nil {
		return nil, "", fmt.Errorf("failed to pack call data: %w", err)
	}

	// Call the contract
	data, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &safeHarborAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to call contract: %w", err)
	}

	// Define the result struct inline
	var result struct {
		AgreementDetails struct {
			ProtocolName   string
			ContactDetails []struct {
				Name    string
				Contact string
			}
			Chains []struct {
				AssetRecoveryAddress common.Address
				Accounts             []struct {
					AccountAddress     common.Address
					ChildContractScope uint8
					Signature          []byte
				}
				ID *big.Int
			}
			BountyTerms struct {
				BountyPercentage      *big.Int
				BountyCapUSD          *big.Int
				Retainable            bool
				Identity              uint8
				DiligenceRequirements string
			}
			AgreementURI string
		}
	}

	// Unpack the result into our struct
	err = parsedABI.UnpackIntoInterface(&result, "getDetails", data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to unpack result: %w", err)
	}

	// Getting Block timestamp of txBody
	block, err := client.BlockByNumber(context.Background(), &blockNumber)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get block: %w", err)
	}

	// Map the result to your SafeHarborAgreement struct
	agreementDetails := &SafeHarborAgreement{
		RegistryTransaction: registryTransaction,
		RegistryChainId:     registryChainId,
		AgreementAddress:    safeHarborAddress.String(),
		Entity:              entity,
		AgreementURI:        result.AgreementDetails.AgreementURI,
		ContactDetails:      formatContactDetails(result.AgreementDetails.ContactDetails),
		Chains:              mapChains(result.AgreementDetails.Chains),
		BountyTerms: BountyTerms{
			BountyPercentage:      int(result.AgreementDetails.BountyTerms.BountyPercentage.Int64()),
			BountyCapUSD:          int(result.AgreementDetails.BountyTerms.BountyCapUSD.Int64()),
			Retainable:            result.AgreementDetails.BountyTerms.Retainable,
			Identity:              identityEnumToString(result.AgreementDetails.BountyTerms.Identity),
			DiligenceRequirements: result.AgreementDetails.BountyTerms.DiligenceRequirements,
		},
		CreatedAt: time.Unix(int64(block.Time()), 0),
	}

	return agreementDetails, result.AgreementDetails.ProtocolName, nil
}

// Utility functions

// formatContactDetails converts Contact array to string format
func formatContactDetails(contactDetails []struct {
	Name    string
	Contact string
}) string {
	contactStr := ""
	for _, contact := range contactDetails {
		contactStr += fmt.Sprintf("%s: %s\n", contact.Name, contact.Contact)
	}
	return contactStr
}

// mapChains converts the result chains to the Chain struct
func mapChains(resultChains []struct {
	AssetRecoveryAddress common.Address
	Accounts             []struct {
		AccountAddress     common.Address
		ChildContractScope uint8
		Signature          []byte
	}
	ID *big.Int
}) []Chain {
	var chains []Chain
	for _, resultChain := range resultChains {
		newChain := Chain{
			AssetRecoveryAddress: resultChain.AssetRecoveryAddress.String(),
			ID:                   strconv.Itoa(int(resultChain.ID.Int64())),
			Accounts:             mapAccounts(resultChain.Accounts),
		}
		chains = append(chains, newChain)
	}
	return chains
}

// mapAccounts converts the result accounts to Account struct
func mapAccounts(resultAccounts []struct {
	AccountAddress     common.Address
	ChildContractScope uint8
	Signature          []byte
}) []Account {
	var accounts []Account
	for _, resultAccount := range resultAccounts {
		newAccount := Account{
			Address:            resultAccount.AccountAddress.String(),
			ChildContractScope: childContractScopeEnumToString(resultAccount.ChildContractScope),
			Signature:          string(resultAccount.Signature),
		}
		accounts = append(accounts, newAccount)
	}
	return accounts
}

// childContractScopeEnumToString converts ChildContractScope enum to string
func childContractScopeEnumToString(scope uint8) ChildContractScope {
	switch scope {
	case 0:
		return ChildContractScopeNone
	case 1:
		return ChildContractScopeExistingOnly
	case 2:
		return ChildContractScopeAll
	default:
		return ChildContractScopeNone
	}
}

// identityEnumToString converts Identity enum to string
func identityEnumToString(identity uint8) string {
	switch identity {
	case 0:
		return "Anonymous"
	case 1:
		return "Pseudonymous"
	case 2:
		return "Named"
	default:
		return "Unknown"
	}
}
