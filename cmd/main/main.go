package main

import (
	"SHDB/pkg/config"
	"SHDB/pkg/contracts/safeharbor"
	"SHDB/pkg/defiliama"
	"SHDB/pkg/firebase"
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ProtocolCollectionName  = "protocols"
	AgreementCollectionName = "safeHarborAgreements"
)

const EtherscanErrNoTxns string = "etherscan server: No transactions found"

func main() {
	app := &cli.App{
		Name:  "shdb",
		Usage: "SHDB is a cli manager for Skylock's firestore safe harbor database",
		Commands: []*cli.Command{
			{
				Name:   "add-adoption",
				Usage:  "Add an adoption to the database",
				Action: runAddAdoption,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "slug",
						Usage:    "Defiliama slug of the adopting protocol",
						Required: true,
					},
					&cli.Int64Flag{
						Name:     "chain",
						Usage:    "Chain ID of the adoption transaction",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "txhash",
						Usage:    "Transaction hash of the safe harbor adoption",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "adoptionProposalUri",
						Usage: "URI of the adoption proposal",
					},
					&cli.BoolFlag{
						Name:  "force",
						Usage: "Force the command to run even if the protocol already exists in the database",
						Value: false,
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "Run in production mode. If false, writes to a test collection",
						Value: false,
					},
				},
			},
			{
				Name:   "refresh-tvl",
				Usage:  "Refresh the TVL of a protocol in the database",
				Action: runRefreshTvl,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "slug",
						Usage:    "Defiliama slug of the protocol. If 'all', refreshes all protocols",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "Run in production mode. If false, writes to a test collection",
						Value: false,
					},
				},
			},
			{
				Name:   "refresh-child-contracts",
				Usage:  "Refresh the child contracts of an agreement in the database",
				Action: runRefreshChildContracts,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "slug",
						Usage:    "Defiliama slug of the protocol. If 'all', refreshes all protocols",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "Run in production mode. If false, writes to a test collection",
						Value: false,
					},
				},
			},
		},
	}

	slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	err := app.Run(os.Args)
	if err != nil {
		slog.Error("app.Run", "error", err)
		os.Exit(1)
	}
}

func runAddAdoption(cCtx *cli.Context) error {
	_ = godotenv.Load()

	//* Load config
	slug := cCtx.String("slug")
	chain := cCtx.Int("chain")
	txHashStr := cCtx.String("txhash")
	force := cCtx.Bool("force")
	adoptionProposalUri := cCtx.String("adoptionProposalUri")
	txhash := common.HexToHash(txHashStr)

	protocolCol, agreementCol := getCollectionNames(cCtx)

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	chainCfg, err := config.LoadChainCfg()
	if err != nil {
		return fmt.Errorf("LoadChainCfg: %w", err)
	}

	if _, exists := chainCfg[chain]; !exists {
		return fmt.Errorf("chain ID not found in chain config: %d", chain)
	}

	eClient, err := ethclient.Dial(chainCfg[chain].RpcUrl)
	if err != nil {
		return fmt.Errorf("rpc.Dial: %w", err)
	}

	err = addAdoption(eClient, fClient, protocolCol, agreementCol, slug, chain, txhash, adoptionProposalUri, force)
	if err != nil {
		return fmt.Errorf("addAdoption: %w", err)
	}

	return nil
}

func runRefreshTvl(cCtx *cli.Context) error {
	_ = godotenv.Load()

	//* Load config
	slug := cCtx.String("slug")

	protocolCol, _ := getCollectionNames(cCtx)

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	//* Refresh single protocol
	if slug != "all" {
		return refreshTvl(fClient, protocolCol, slug)
	}

	//* Refresh all protocols
	protocols := fClient.Collection(protocolCol)
	documents, err := protocols.Documents(context.Background()).GetAll()
	if err != nil {
		return fmt.Errorf("firestore.GetAll: %w", err)
	}

	for _, doc := range documents {
		slug := doc.Ref.ID
		slog.Info("Refreshing TVL for protocol", "slug", slug)
		err = refreshTvl(fClient, protocolCol, slug)
		if err != nil {
			slog.Warn("refreshTvl", "error", err)
		}
	}

	return nil
}

func runRefreshChildContracts(cCtx *cli.Context) error {
	_ = godotenv.Load()

	//* Load config
	slug := cCtx.String("slug")
	protocolCol, agreementCol := getCollectionNames(cCtx)

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	chainCfg, err := config.LoadChainCfg()
	if err != nil {
		return fmt.Errorf("LoadChainCfg: %w", err)
	}

	//* Refresh single protocol
	if slug != "all" {
		return refreshChildContracts(fClient, protocolCol, agreementCol, slug, chainCfg)
	}

	//* Refresh all protocols
	protocols := fClient.Collection(protocolCol)
	documents, err := protocols.Documents(context.Background()).GetAll()
	if err != nil {
		return fmt.Errorf("firestore.GetAll: %w", err)
	}

	for _, doc := range documents {
		slug := doc.Ref.ID
		err = refreshChildContracts(fClient, protocolCol, agreementCol, slug, chainCfg)
		if err != nil {
			slog.Warn("refreshChildContracts", "error", err)
		}
	}

	return nil
}

func addAdoption(
	eClient *ethclient.Client,
	fClient *firestore.Client,
	protocolCol string,
	agreementCol string,
	slug string,
	chain int,
	txhash common.Hash,
	adoptionProposalUri string,
	force bool,
) error {
	txbody, _, err := eClient.TransactionByHash(context.Background(), txhash)
	if err != nil {
		return fmt.Errorf("rpc.TransactionByHash: %w", err)
	}

	receipt, err := eClient.TransactionReceipt(context.Background(), txhash)
	if err != nil {
		return fmt.Errorf("rpc.TransactionReceipt: %w", err)
	}

	sender, err := types.Sender(types.NewLondonSigner(big.NewInt(int64(chain))), txbody)
	if err != nil {
		return fmt.Errorf("types.Sender: %w", err)
	}

	agreementAddress, agreement, err := safeharbor.GetAgreement(txhash, eClient)
	if err != nil {
		return fmt.Errorf("getAgreement: %w", err)
	}

	// Fetch defiliama data for protocol
	protocol, err := defiliama.GetProtocol(slug)
	if err != nil {
		return fmt.Errorf("defiliama.GetProtocol(slug=%v): %w", slug, err)
	}

	//* Upload protocol & adoption to firestore if not already present
	// Check if protocol exists
	protocolDocRef := fClient.Collection(protocolCol).Doc(protocol.Slug)
	protocolDoc, err := protocolDocRef.Get(context.Background())
	if err != nil && status.Code(err) != codes.NotFound {
		return fmt.Errorf("firestore.Get: %w", err)
	}

	exists := protocolDoc.Exists()
	if exists && !force {
		return fmt.Errorf("protocol already exists in firestore. Use --force to overwrite")
	}

	// Upload protocol
	fProtocol := firebase.Protocol{
		Name:     protocol.Name,
		Slug:     protocol.Slug,
		Website:  protocol.Website,
		Icon:     protocol.Icon,
		TVL:      protocol.TVL,
		Category: protocol.Category,
	}
	slog.Info("Uploading protocol", "protocol", fProtocol.Name)
	_, err = protocolDocRef.Set(context.Background(), fProtocol)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	// Create safe harbor adoption
	fAdoption := firebase.SafeHarborAgreement{
		Protocol:            protocolDocRef,
		RegistryTransaction: txhash.String(),
		RegistryChainId:     fmt.Sprintf("%d", chain),
		AgreementAddress:    agreementAddress.String(),
		Entity:              sender.String(),
		AgreementURI:        agreement.AgreementURI,
		AdoptionProposalURI: adoptionProposalUri,
		ContactDetails:      firebase.FormatContactDetails(agreement.ContactDetails),
		Chains:              firebase.FormatChains(agreement.Chains),
		BountyTerms:         firebase.FormatBountyTerms(agreement.BountyTerms),
		CreatedAt:           txbody.Time(),
		CreatedBlock:        int(receipt.BlockNumber.Int64()),
	}

	// Upload safe harbor adoption
	slog.Info("Uploading adoption", "adoption", txhash.String())
	agreementDocRef := fClient.Collection(agreementCol).Doc(txhash.String())
	_, err = agreementDocRef.Set(context.Background(), fAdoption)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	// Update protocol with safe harbor agreement reference
	slog.Info("Updating protocol with safe harbor agreement reference")
	fProtocol.SafeHarborAgreement = agreementDocRef
	_, err = protocolDocRef.Set(context.Background(), fProtocol)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	slog.Info("Successfully added adoption to database")
	return nil
}

func refreshTvl(
	fClient *firestore.Client,
	protocolCol string,
	slug string,
) error {
	//* Confirm protocol exists
	protocolDocRef := fClient.Collection(protocolCol).Doc(slug)
	protocolDoc, err := protocolDocRef.Get(context.Background())
	if err != nil {
		return fmt.Errorf("firestore.Get: %w", err)
	}

	if !protocolDoc.Exists() {
		return fmt.Errorf("protocol not found in firestore")
	}

	//* Get the slug from the document
	var protocol firebase.Protocol
	err = protocolDoc.DataTo(&protocol)
	if err != nil {
		return fmt.Errorf("firestore.DataTo: %w", err)
	}

	slug = protocol.Slug

	//* Fetch and update TVL
	tvl, err := defiliama.GetTvl(protocol.Slug)
	if err != nil {
		return fmt.Errorf("defiliama.GetProtocol(slug=%v): %w", slug, err)
	}

	data := map[string]interface{}{
		"tvl": tvl,
	}

	_, err = protocolDocRef.Set(context.Background(), data, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	slog.Info("Successfully refreshed TVL", "protocol", slug, "tvl", tvl)

	return nil
}

func refreshChildContracts(
	fClient *firestore.Client,
	protocolCol string,
	agreementCol string,
	slug string,
	chainCfg map[int]config.ChainCfg,
) error {
	_, _, _, _, _ = fClient, protocolCol, agreementCol, slug, chainCfg
	return fmt.Errorf("refreshChildContracts not implemented - issue #16")
}

func getCollectionNames(cCtx *cli.Context) (protocolCol string, agreementCol string) {
	protocolCol = ProtocolCollectionName
	agreementCol = AgreementCollectionName

	if !cCtx.Bool("prod") {
		slog.Warn("Running in test mode")
		protocolCol = "test_" + protocolCol
		agreementCol = "test_" + agreementCol
	}

	return protocolCol, agreementCol
}
