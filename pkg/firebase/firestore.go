package firebase

import (
	"SHDB/pkg/contracts/adoptiondetails"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func NewFirestoreClient() (*firestore.Client, error) {
	creds := os.Getenv("FIREBASE_CREDENTIALS")
	if creds == "" {
		return nil, fmt.Errorf("missing FIREBASE_CREDENTIALS env")
	}

	ctx := context.Background()
	creds = strings.Trim(creds, "'")
	client, err := firestore.NewClient(ctx, "skylock-xyz", option.WithCredentialsJSON([]byte(creds)))
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}

	return client, nil
}

func FormatContactDetails(d []adoptiondetails.Contact) string {
	s := ""
	for _, c := range d {
		s += fmt.Sprintf("%s: %s\n", c.Name, c.Contact)
	}

	return s
}

func FormatChains(c []adoptiondetails.Chain) []Chain {
	chains := make([]Chain, len(c))

	for i, chain := range c {
		chains[i] = Chain{
			AssetRecoveryAddress: chain.AssetRecoveryAddress.String(),
			ID:                   int(chain.Id.Int64()),
			Accounts:             FormatAccounts(chain.Accounts),
		}
	}

	return chains
}

func FormatAccounts(a []adoptiondetails.Account) []Account {
	accounts := make([]Account, len(a))

	for i, account := range a {
		if int(account.ChildContractScope) >= len(ChildContractScopes) {
			log.Fatalf("invalid child contract scope: %d", account.ChildContractScope)
		}

		accounts[i] = Account{
			Address:            account.AccountAddress.String(),
			ChildContractScope: ChildContractScopes[account.ChildContractScope],
			Signature:          string(account.Signature),
		}
	}

	return accounts
}

func FormatBountyTerms(b adoptiondetails.BountyTerms) BountyTerms {
	if int(b.Identity) >= len(Identities) {
		log.Fatalf("invalid identity: %d", b.Identity)
	}

	return BountyTerms{
		BountyPercentage:      int(b.BountyPercentage.Int64()),
		BountyCapUSD:          int(b.BountyCapUSD.Int64()),
		Retainable:            b.Retainable,
		Identity:              Identities[b.Identity],
		DiligenceRequirements: b.DiligenceRequirements,
	}
}
