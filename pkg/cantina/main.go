package cantina

import (
	"SHDB/pkg/types"
	"encoding/json"
	"fmt"
	"net/http"
)

const API_URL = "https://cantina.xyz/api/v0/bounties"

type Client struct {
}

type bounty struct {
	Id         string             `json:"id"`
	Name       string             `json:"name"`
	Url        string             `json:"url"`
	SafeHarbor *safeHarborDetails `json:"safeHarbor"`
}

type safeHarborDetails struct {
	Id              string          `json:"id"`
	Cap             string          `json:"cap"`
	Reward          string          `json:"reward"`
	ReturnAddresses []returnAddress `json:"return_addresses"`
	Assets          []asset         `json:"assets"`
}

type returnAddress struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`
}

type asset struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewClient() (*Client, error) {
	return &Client{}, nil
}

func (c *Client) GetAgreements() ([]types.SafeHarborAgreementCantinaV1, error) {
	// Fetch bounties[] from cantina API
	resp, err := http.Get(API_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cResp []bounty
	err = json.NewDecoder(resp.Body).Decode(&cResp)
	if err != nil {
		return nil, err
	}

	var agreements []types.SafeHarborAgreementCantinaV1
	for _, b := range cResp {
		if b.SafeHarbor == nil {
			continue
		}

		agreements = append(agreements, b.toSafeHarborAgreement())
	}

	return agreements, nil
}

func (c *Client) GetAgreement(id string) (types.SafeHarborAgreementCantinaV1, error) {
	agreements, err := c.GetAgreements()
	if err != nil {
		return types.SafeHarborAgreementCantinaV1{}, err
	}

	for _, agreement := range agreements {
		if agreement.Slug == "cantina-"+id {
			return agreement, nil
		}
	}

	return types.SafeHarborAgreementCantinaV1{}, fmt.Errorf("agreement with id %s not found", id)
}

func (b bounty) toSafeHarborAgreement() types.SafeHarborAgreementCantinaV1 {
	recoveryAddresses := make([]types.CantinaRecoveryAddressV1, 0, len(b.SafeHarbor.ReturnAddresses))
	for _, ra := range b.SafeHarbor.ReturnAddresses {
		recoveryAddress := types.CantinaRecoveryAddressV1{
			Address: ra.Address,
			Chain:   ra.Chain,
		}
		recoveryAddresses = append(recoveryAddresses, recoveryAddress)
	}

	assets := make([]types.CantinaAssetsV1, 0, len(b.SafeHarbor.Assets))
	for _, a := range b.SafeHarbor.Assets {
		asset := types.CantinaAssetsV1{
			Name:        a.Name,
			Description: a.Description,
		}
		assets = append(assets, asset)
	}

	safeHarborAgreement := types.SafeHarborAgreementCantinaV1{
		SafeHarborAgreementBase: types.SafeHarborAgreementBase{
			AdoptionProposalURI: "",
			Slug:                "cantina-" + b.Id,
			Version:             types.CantinaV1,
		},
		AgreementDetails: types.CantinaDetailsV1{
			Name:    b.Name,
			Contact: fmt.Sprintf("https://cantina.xyz/bounties/%s", b.Id),
			BountyTerms: types.BountyTermsV1{
				Retainable:            false,
				Identity:              types.IdentityNamed,
				DiligenceRequirements: "Diligence performed by Cantina as described on their platform",
			},
			RecoveryAddresses: recoveryAddresses,
			Assets:            assets,
			CantinaUrl:        fmt.Sprintf("https://cantina.xyz/bounties/%s", b.Id),
		},
	}

	return safeHarborAgreement
}
