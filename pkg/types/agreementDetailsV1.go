package types

import (
	"SHDB/pkg/contracts/adoptiondetails"
	"SHDB/pkg/scan"
	"fmt"
)

const (
	ChildContractScopeNone         ChildContractScope = "None"
	ChildContractScopeExistingOnly ChildContractScope = "ExistingOnly"
	ChildContractScopeAll          ChildContractScope = "All"
)

const (
	IdentityAnonymous   Identity = "Anonymous"
	IdentityPsedonymous Identity = "Pseudonymous"
	IdentityNamed       Identity = "Named"
)

var (
	ChildContractScopes = []ChildContractScope{ChildContractScopeNone, ChildContractScopeExistingOnly, ChildContractScopeAll}
	Identities          = []Identity{IdentityAnonymous, IdentityPsedonymous, IdentityNamed}
)

type ChildContractScope string
type Identity string

type AgreementDetailsV1 struct {
	Name         string        `firestore:"name"`
	Contact      string        `firestore:"contact"`
	Chains       []ChainV1     `firestore:"chains"`
	BountyTerms  BountyTermsV1 `firestore:"bountyTerms"`
	AgreementURI string        `firestore:"agreementURI"`
}

type ContactV1 struct {
	Name    string `firestore:"name"`
	Contact string `firestore:"contact"`
}

type ChainV1 struct {
	ID                   int         `firestore:"id"`
	AssetRecoveryAddress string      `firestore:"assetRecoveryAddress"`
	Accounts             []AccountV1 `firestore:"accounts"`
}

type AccountV1 struct {
	Name               string             `firestore:"name"`
	Address            string             `firestore:"address"`
	ChildContractScope ChildContractScope `firestore:"childContractScope"`
	Signature          string             `firestore:"signature"`
	Children           []ChildAccountV1   `firestore:"children"`
}

type ChildAccountV1 struct {
	Name    string `firestore:"name"`
	Address string `firestore:"address"`
}

type BountyTermsV1 struct {
	BountyCapUSD          int      `firestore:"bountyCapUSD"`
	BountyPercentage      int      `firestore:"bountyPercentage"`
	DiligenceRequirements string   `firestore:"diligenceRequirements"`
	Identity              Identity `firestore:"identity"`
	Retainable            bool     `firestore:"retainable"`
}

func (v *AgreementDetailsV1) FromRawAgreementDetails(d *adoptiondetails.AgreementDetailsV1) error {
	v.Name = d.ProtocolName
	v.Contact = formatContactDetails(d.ContactDetails)
	v.Chains = make([]ChainV1, len(d.Chains))
	v.BountyTerms = BountyTermsV1{}
	if err := v.BountyTerms.fromRawBountyTerms(d.BountyTerms); err != nil {
		return fmt.Errorf("fromRawBountyTerms: %w", err)
	}

	for i, chain := range d.Chains {
		v.Chains[i] = ChainV1{}
		if err := v.Chains[i].fromRawChain(chain); err != nil {
			return fmt.Errorf("fromRawChain: %w", err)
		}
	}

	return nil
}

func (v *AgreementDetailsV1) TryNameAddresses(client scan.Client) {
	for i, chain := range v.Chains {
		for j, account := range chain.Accounts {
			name := client.ContractName(account.Address)
			account.Name = name
			v.Chains[i].Accounts[j] = account
		}
	}
}

func (c *ChainV1) fromRawChain(chain adoptiondetails.Chain) error {
	c.ID = int(chain.Id.Int64())
	c.AssetRecoveryAddress = chain.AssetRecoveryAddress.Hex()
	c.Accounts = make([]AccountV1, len(chain.Accounts))
	for i, account := range chain.Accounts {
		c.Accounts[i] = AccountV1{}
		if err := c.Accounts[i].fromRawAccount(account); err != nil {
			return fmt.Errorf("fromRawAccount: %w", err)
		}
	}

	return nil
}

func (a *AccountV1) fromRawAccount(account adoptiondetails.Account) error {
	if int(account.ChildContractScope) >= len(ChildContractScopes) {
		return fmt.Errorf("invalid child contract scope: %d", account.ChildContractScope)
	}

	a.Name = ""
	a.Address = account.AccountAddress.Hex()
	a.ChildContractScope = ChildContractScopes[account.ChildContractScope]
	a.Signature = string(account.Signature)
	a.Children = []ChildAccountV1{}

	return nil
}

func (b *BountyTermsV1) fromRawBountyTerms(bounty adoptiondetails.BountyTerms) error {
	if int(bounty.Identity) >= len(Identities) {
		return fmt.Errorf("invalid identity: %d", bounty.Identity)
	}

	b.BountyCapUSD = int(bounty.BountyPercentage.Int64())
	b.BountyPercentage = int(bounty.BountyPercentage.Int64())
	b.DiligenceRequirements = bounty.DiligenceRequirements
	b.Identity = Identities[bounty.Identity]
	b.Retainable = bounty.Retainable

	return nil
}

func formatContactDetails(d []adoptiondetails.Contact) string {
	s := ""
	for _, c := range d {
		s += fmt.Sprintf("%s: %s\n", c.Name, c.Contact)
	}

	return s
}
