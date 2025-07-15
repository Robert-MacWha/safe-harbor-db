package immunefi

import (
	"SHDB/pkg/config"
	"SHDB/pkg/types"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Client struct {
	buildId string
}

type basicImmunefiBountyResp struct {
	PageProps struct {
		Bounties []basicImmunefiBounty `json:"bounties"`
	} `json:"pageProps"`
}

type basicImmunefiBounty struct {
	ContentfulId string `json:"contentfulId"`
	Id           string `json:"id"`
	Project      string `json:"project"`
	Tags         struct {
		General []string `json:"general"` // Safe Harbor may be one of the tags
	} `json:"tags"`
}

type detailedImmunefiBountyResp struct {
	PageProps struct {
		Bounty  detailedImmunefiBounty  `json:"bounty"`
		Project detailedImmunefiProject `json:"project"`
	} `json:"pageProps"`
}

type detailedImmunefiBounty struct {
	Slug    string                        `json:"slug"`
	Project string                        `json:"project"`
	Assets  []detailedImmunefiBountyAsset `json:"assets"`
}

type detailedImmunefiBountyAsset struct {
	Type        string `json:"type"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type detailedImmunefiProject struct {
	IsSafeHarborActive bool `json:"isSafeHarborActive"`
}

func NewClient() (*Client, error) {
	// Fetch the build ID with an http request
	buildId, err := getCurrentBuildId()
	if err != nil {
		return nil, fmt.Errorf("getCurrentBuildId: %w", err)
	}

	slog.Info("Immunefi client initialized", "buildId", buildId)

	return &Client{
		buildId: buildId,
	}, nil
}

// GetAgreements fetches all immunefi protocols that have safe harbor agreements
// and returns their slugs.
func (c *Client) GetAgreements() ([]types.SafeHarborAgreementImmunefiV1, error) {
	basicBounties, err := c.getBasicBounties()
	if err != nil {
		return nil, fmt.Errorf("getBasicBounties: %w", err)
	}

	detailedBounties := make([]detailedImmunefiBounty, 0)
	for _, bounty := range basicBounties {
		for _, tag := range bounty.Tags.General {
			if strings.Contains(tag, "Safe Harbor") {
				detailedBounty, err := c.getDetailedBounty(bounty.Id)
				if err != nil {
					slog.Error("Failed to get detailed bounty", "bountyId", bounty.Id, "error", err)
					continue
				}
				detailedBounties = append(detailedBounties, *detailedBounty)
			}
		}
	}

	agreementDetails := make([]types.SafeHarborAgreementImmunefiV1, 0, len(detailedBounties))
	for _, bounty := range detailedBounties {
		agreementDetails = append(agreementDetails, bounty.toSafeHarborAgreement())
	}

	return agreementDetails, nil
}

func (c *Client) GetAgreement(protocolId string) (*types.SafeHarborAgreementImmunefiV1, error) {
	detailedBounty, err := c.getDetailedBounty(protocolId)
	if err != nil {
		return nil, fmt.Errorf("getDetailedBounty(%s): %w", protocolId, err)
	}

	safeHarborAgreement := detailedBounty.toSafeHarborAgreement()
	return &safeHarborAgreement, nil
}

func (c *Client) getBasicBounties() ([]basicImmunefiBounty, error) {
	url := fmt.Sprintf("https://immunefi.com/_next/data/%s/bug-bounty.json", c.buildId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var iResp basicImmunefiBountyResp
	err = json.NewDecoder(resp.Body).Decode(&iResp)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	return iResp.PageProps.Bounties, nil
}

func (c *Client) getDetailedBounty(bountyId string) (*detailedImmunefiBounty, error) {
	url := fmt.Sprintf("https://immunefi.com/_next/data/%s/bug-bounty/%s/scope.json", c.buildId, bountyId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var iResp detailedImmunefiBountyResp
	err = json.NewDecoder(resp.Body).Decode(&iResp)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	if !iResp.PageProps.Project.IsSafeHarborActive {
		return nil, fmt.Errorf("project is not safe harbor active")
	}

	return &iResp.PageProps.Bounty, nil
}

func (d detailedImmunefiBounty) toSafeHarborAgreement() types.SafeHarborAgreementImmunefiV1 {
	chains := map[int]types.ImmunefiChainV1{}

	safeHarborAgreement := types.SafeHarborAgreementImmunefiV1{
		SafeHarborAgreementBase: types.SafeHarborAgreementBase{
			AdoptionProposalURI: fmt.Sprintf("https://immunefi.com/bug-bounty/%s/safe-harbor/", d.Slug),
			Slug:                "immunefi-" + d.Slug,
			Version:             types.ImmunefiV1,
		},
		AgreementDetails: types.ImmunefiDetailsV1{
			Name:    d.Project,
			Contact: fmt.Sprintf("https://immunefi.com/bug-bounty/%s/safe-harbor/", d.Slug),
			BountyTerms: types.BountyTermsV1{
				Retainable:            false,
				Identity:              types.IdentityNamed,
				DiligenceRequirements: "Diligence performed by Immunefi, including creating an account on their website and submitting a report through their platform",
			},
		},
	}

	for _, asset := range d.Assets {
		addAssetToChain(asset, chains)
	}

	for _, chain := range chains {
		safeHarborAgreement.AgreementDetails.Chains = append(safeHarborAgreement.AgreementDetails.Chains, chain)
	}

	return safeHarborAgreement
}

func addAssetToChain(asset detailedImmunefiBountyAsset, chains map[int]types.ImmunefiChainV1) {
	if asset.Type != "smart_contract" {
		return
	}

	address, chainId, err := config.ParseAddressFromURL(asset.URL)
	if err != nil {
		slog.Warn("Failed to parse address from URL", "url", asset.URL, "error", err)
		return
	}

	account := types.ImmunefiAccountV1{
		Name:    asset.Description,
		Address: address,
	}

	chain, exists := chains[chainId]
	if exists {
		chain.Accounts = append(chain.Accounts, account)
	} else {
		chain = types.ImmunefiChainV1{
			ID:       chainId,
			Accounts: []types.ImmunefiAccountV1{account},
		}
	}

	chains[chainId] = chain
}

func getCurrentBuildId() (string, error) {
	resp, err := http.Get("https://immunefi.com/bug-bounty")
	if err != nil {
		return "", fmt.Errorf("http.Get: %w", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("io.ReadAll: %w", err)
	}

	body := string(bodyBytes)
	if start := strings.Index(body, `"buildId":"`); start != -1 {
		start += len(`"buildId":"`)
		if end := strings.Index(body[start:], `"`); end != -1 {
			return body[start : start+end], nil
		}
	}

	return "", fmt.Errorf("buildId not found in response")
}
