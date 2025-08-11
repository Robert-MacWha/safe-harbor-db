package types

import (
	"log/slog"
	"math/big"
	"strconv"
	"strings"

	"SHDB/pkg/contracts/adoptiondetailsv2"
	"SHDB/pkg/scan"
)

var childContractScopesV2 = []ChildContractScope{
	ChildContractScope("None"),
	ChildContractScope("ExistingOnly"),
	ChildContractScope("All"),
	ChildContractScope("FutureOnly"),
}

type AgreementDetailsV2 struct {
	Name         string        `firestore:"name"`
	Contact      string        `firestore:"contact"`
	Chains       []ChainV2     `firestore:"chains"`
	BountyTerms  BountyTermsV2 `firestore:"bountyTerms"`
	AgreementURI string        `firestore:"agreementURI"`
}

type ChainV2 struct {
	Caip2ChainId         string      `firestore:"caip2ChainId"`
	AssetRecoveryAddress string      `firestore:"assetRecoveryAddress"`
	Accounts             []AccountV2 `firestore:"accounts"`
}

type AccountV2 struct {
	Name               string             `firestore:"name"`
	Address            string             `firestore:"address"`
	ChildContractScope ChildContractScope `firestore:"childContractScope"`
}

type BountyTermsV2 struct {
	BountyCapUSD          int      `firestore:"bountyCapUSD"`
	BountyPercentage      int      `firestore:"bountyPercentage"`
	DiligenceRequirements string   `firestore:"diligenceRequirements"`
	Identity              Identity `firestore:"identity"`
	Retainable            bool     `firestore:"retainable"`
	AggregateBountyCapUSD int      `firestore:"aggregateBountyCapUSD"`
}

func (v *AgreementDetailsV2) FromRawAgreementDetails(d adoptiondetailsv2.AgreementDetailsV2) {
	v.Name = d.ProtocolName
	// Flatten contact details to a single string, consistent with V1 storage
	contact := ""
	for _, c := range d.ContactDetails {
		contact += c.Name + ": " + c.Contact + "\n"
	}
	v.Contact = contact
	v.AgreementURI = d.AgreementURI

	v.BountyTerms = BountyTermsV2{
		BountyCapUSD:          int(toBigInt(d.BountyTerms.BountyCapUSD)),
		BountyPercentage:      int(toBigInt(d.BountyTerms.BountyPercentage)),
		DiligenceRequirements: d.BountyTerms.DiligenceRequirements,
		Identity:              Identities[d.BountyTerms.Identity],
		Retainable:            d.BountyTerms.Retainable,
		AggregateBountyCapUSD: int(toBigInt(d.BountyTerms.AggregateBountyCapUSD)),
	}

	v.Chains = make([]ChainV2, len(d.Chains))
	for i, ch := range d.Chains {
		chain := ChainV2{
			Caip2ChainId:         ch.Caip2ChainId,
			AssetRecoveryAddress: ch.AssetRecoveryAddress,
			Accounts:             make([]AccountV2, len(ch.Accounts)),
		}
		for j, acc := range ch.Accounts {
			scopeIdx := int(acc.ChildContractScope)
			if scopeIdx >= len(childContractScopesV2) {
				scopeIdx = 0
			}
			chain.Accounts[j] = AccountV2{
				Name:               "",
				Address:            acc.AccountAddress,
				ChildContractScope: childContractScopesV2[scopeIdx],
			}
		}
		v.Chains[i] = chain
	}
}

// toBigInt converts interface{} coming from abi.ConvertType big-int fields to int64
func toBigInt(v interface{}) int64 {
	switch t := v.(type) {
	case *big.Int:
		return t.Int64()
	case big.Int:
		return t.Int64()
	default:
		return 0
	}
}

// TryNameAddressesByCAIP2 attempts to name contract addresses for EVM chains (eip155:<id>) only.
// It preserves the original string addresses and is a best-effort similar to V1.
func (v *AgreementDetailsV2) TryNameAddressesByCAIP2(getScan func(chainID int) (scan.Client, error)) {
	for i, chain := range v.Chains {
		if !strings.HasPrefix(chain.Caip2ChainId, "eip155:") {
			continue
		}
		idStr := strings.TrimPrefix(chain.Caip2ChainId, "eip155:")
		evmID, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Warn("invalid eip155 chain id", "caip2", chain.Caip2ChainId, "error", err)
			continue
		}
		client, err := getScan(evmID)
		if err != nil {
			slog.Warn("getScanClient", "chainID", evmID, "error", err)
			continue
		}
		for j, account := range chain.Accounts {
			name := client.ContractName(account.Address)
			account.Name = name
			v.Chains[i].Accounts[j] = account
			if name == "" {
				slog.Info("Naming address failed", "address", account.Address, "caip2", chain.Caip2ChainId)
			}
		}
	}
}
