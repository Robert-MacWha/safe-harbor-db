package main

import (
	"SHDB/pkg/cantina"
	"SHDB/pkg/config"
	adoptv2 "SHDB/pkg/contracts/adoptiondetailsv2"
	"SHDB/pkg/contracts/safeharbor"
	"SHDB/pkg/contracts/safeharbor_v3"
	"SHDB/pkg/deduab"
	"SHDB/pkg/defiliama"
	"SHDB/pkg/firebase"
	"SHDB/pkg/scan"
	"SHDB/pkg/telegram"
	"SHDB/pkg/types"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

// TODO: Add refresh-contract-names command
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
					&cli.StringFlag{
						Name:  "bugBounty",
						Usage: "Protocol's bug bounty program URL",
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
				Name:   "add-adoption-v2",
				Usage:  "Add a v2 adoption to the database",
				Action: runAddAdoptionV2,
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
					&cli.StringFlag{
						Name:  "bugBounty",
						Usage: "Protocol's bug bounty program URL",
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
				Name:   "add-adoption-v3",
				Usage:  "Add a v3 adoption to the database",
				Action: runAddAdoptionV3,
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
					&cli.StringFlag{
						Name:  "bugBounty",
						Usage: "Protocol's bug bounty program URL",
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
				Name:   "add-cantina-adoption",
				Usage:  "Adds an adoption from Cantina to the database",
				Action: runAddCantinaAdoption,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "cantina-id",
						Usage:    "Cantina ID of the adopting protocol",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "defiliama-slug",
						Usage:    "Defiliama slug of the adopting protocol",
						Required: true,
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
				Usage:  "Refresh the child contracts of all agreements in the database",
				Action: runRefreshChildContracts,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "Run in production mode. If false, writes to a test collection",
						Value: false,
					},
				},
			},
			{
				Name:   "refresh-adoption-dates",
				Usage:  "Refreshes the adoption dates of all agreements in the database",
				Action: runRefreshAdoptionDates,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "Run in production mode. If false, writes to a test collection",
						Value: false,
					},
				},
			},
			{
				Name:   "check-new-adoptions",
				Usage:  "Checks for new adoptions and sends notifications",
				Action: runCheckNewAdoptions,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "Run in production mode. If false, reads from a test collection",
						Value: false,
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "Only check and log, don't send notifications",
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
		if os.Getenv("DEV") != "1" {
			telegram.SendNotification(
				fmt.Sprintf("Error running SHDB: %s", err),
				os.Getenv("TELEGRAM_BOT_TOKEN"),
				os.Getenv("TELEGRAM_CHAT_ID"),
			)
		}

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
	bugBounty := cCtx.String("bugBounty")
	txhash := common.HexToHash(txHashStr)

	protocolCol, agreementCol := getCollectionNames(cCtx)

	err := addAdoptionOnchain(
		protocolCol,
		agreementCol,
		slug,
		chain,
		txhash,
		adoptionProposalUri,
		bugBounty,
		force,
		types.SealV1,
		&AgreementDetailsV1Fetcher{chain: chain},
	)
	if err != nil {
		return fmt.Errorf("addAdoptionOnchain: %w", err)
	}

	return nil
}

func runAddAdoptionV2(cCtx *cli.Context) error {
	_ = godotenv.Load()

	slug := cCtx.String("slug")
	chain := cCtx.Int("chain")
	txHashStr := cCtx.String("txhash")
	force := cCtx.Bool("force")
	adoptionProposalUri := cCtx.String("adoptionProposalUri")
	bugBounty := cCtx.String("bugBounty")
	txhash := common.HexToHash(txHashStr)

	protocolCol, agreementCol := getCollectionNames(cCtx)

	err := addAdoptionOnchain(protocolCol,
		agreementCol,
		slug,
		chain,
		txhash,
		adoptionProposalUri,
		bugBounty,
		force,
		types.SealV2,
		&AgreementDetailsV2Fetcher{chain: chain},
	)
	if err != nil {
		return fmt.Errorf("addAdoptionOnchain: %w", err)
	}

	return nil
}

func runAddAdoptionV3(cCtx *cli.Context) error {
	_ = godotenv.Load()

	slug := cCtx.String("slug")
	chain := cCtx.Int("chain")
	txHashStr := cCtx.String("txhash")
	force := cCtx.Bool("force")
	adoptionProposalUri := cCtx.String("adoptionProposalUri")
	bugBounty := cCtx.String("bugBounty")
	txhash := common.HexToHash(txHashStr)

	protocolCol, agreementCol := getCollectionNames(cCtx)

	err := addAdoptionOnchain(
		protocolCol,
		agreementCol,
		slug,
		chain,
		txhash,
		adoptionProposalUri,
		bugBounty,
		force,
		types.SealV2,
		&AgreementDetailsV3Fetcher{chain: chain},
	)
	if err != nil {
		return fmt.Errorf("addAdoptionOnchain: %w", err)
	}

	return nil
}

func runAddCantinaAdoption(cCtx *cli.Context) error {
	_ = godotenv.Load()

	//* Load config
	cantinaID := cCtx.String("cantina-id")
	defiliamaSlug := cCtx.String("defiliama-slug")
	force := cCtx.Bool("force")

	protocolCol, agreementCol := getCollectionNames(cCtx)

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	cClient, err := cantina.NewClient()
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	err = addCantinaAdoption(cClient, fClient, protocolCol, agreementCol, cantinaID, defiliamaSlug, force)
	if err != nil {
		return fmt.Errorf("addCantinaAdoption: %w", err)
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
		slog.Info("Refreshing TVL for protocol", "protocol", slug)
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
	_, agreementCol := getCollectionNames(cCtx)

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	//* Refresh all protocols
	protocols := fClient.Collection(agreementCol)
	documents, err := protocols.Documents(context.Background()).GetAll()
	if err != nil {
		return fmt.Errorf("firestore.GetAll: %w", err)
	}

	for _, doc := range documents {
		slug := doc.Ref.ID
		slog.Info("Refreshing child contracts for agreement", "slug", slug)
		err = refreshChildContracts(fClient, agreementCol, slug)
		if err != nil {
			slog.Warn("refreshChildContracts", "error", err)
		}
	}

	return nil
}

func runRefreshAdoptionDates(cCtx *cli.Context) error {
	_ = godotenv.Load()

	_, agreementCol := getCollectionNames(cCtx)

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	agreements := fClient.Collection(agreementCol)
	documents, err := agreements.Documents(context.Background()).GetAll()
	if err != nil {
		return fmt.Errorf("firestore.GetAll: %w", err)
	}

	for _, doc := range documents {
		slug := doc.Ref.ID
		slog.Info("Refreshing adoption date for agreement", "slug", slug)
		err = refreshAdoptionDate(fClient, agreementCol, slug)
		if err != nil {
			slog.Warn("refreshChildContracts", "error", err)
		}
	}

	return nil
}

func runCheckNewAdoptions(cCtx *cli.Context) error {
	_ = godotenv.Load()

	//* Load config
	dryRun := cCtx.Bool("dry-run")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	_, agreementCol := getCollectionNames(cCtx)

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	cClient, err := cantina.NewClient()
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	//* Get new adoptions
	newCantinaAdoptions, err := checkNewCantinaAdoptions(cClient, fClient, agreementCol)
	if err != nil {
		return err
	}

	newOnchainAdoptions, err := checkNewOnchainAdoptions(fClient, agreementCol)
	if err != nil {
		return fmt.Errorf("checkNewOnchainAdoptions: %w", err)
	}

	var newAdoptions []string
	newAdoptions = append(newAdoptions, newCantinaAdoptions...)
	newAdoptions = append(newAdoptions, newOnchainAdoptions...)

	if len(newAdoptions) == 0 {
		slog.Info("No new adoptions found")
		return nil
	}

	//* Send TG notification
	message := "🚨 *New Safe Harbor Adoptions Found*\n\n"
	for _, slug := range newAdoptions {
		message += fmt.Sprintf("- `%s`\n", slug)
	}
	message += "\nRun the `add-*-adoption` github action to add them to the database."

	slog.Info("New adoptions found", "count", len(newAdoptions), "message", message)
	if !dryRun {
		err := telegram.SendNotification(message, botToken, chatID)
		if err != nil {
			return fmt.Errorf("SendNotification: %w", err)
		}
	}

	return nil
}

type DetailFetcher interface {
	GetDetails(txhash common.Hash) (common.Address, interface{}, error)
}

// Universal adoption addter for v1, 2, and 3
func addAdoptionOnchain(
	protocolCol string,
	agreementCol string,
	slug string,
	chain int,
	txhash common.Hash,
	adoptionProposalUri string,
	bugBounty string,
	force bool,
	version types.SafeHarborVersion,
	detailFetcher DetailFetcher,
) error {
	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	eClient, err := getChainClient(chain)
	if err != nil {
		return fmt.Errorf("getChainClient: %w", err)
	}

	protocol, err := defiliama.GetProtocol(slug)
	if err != nil {
		slog.Error("defiliama.GetProtocol", "slug", slug, "error", err)
		slog.Info("Falling back to no Defiliama protocol data")

		protocol = types.Protocol{
			Slug:      slug,
			Name:      slug,
			BugBounty: bugBounty,
		}
	}
	protocol.BugBounty = bugBounty

	// Fetch safe harbor adoption from blockchain
	txbody, _, err := eClient.TransactionByHash(context.Background(), txhash)
	if err != nil {
		return fmt.Errorf("rpc.TransactionByHash: %w", err)
	}

	receipt, err := eClient.TransactionReceipt(context.Background(), txhash)
	if err != nil {
		return fmt.Errorf("rpc.TransactionReceipt: %w", err)
	}

	sender, err := eClient.TransactionSender(context.Background(), txbody, receipt.BlockHash, receipt.TransactionIndex)
	if err != nil {
		return fmt.Errorf("types.Sender: %w", err)
	}

	agreementAddress, details, err := detailFetcher.GetDetails(txhash)
	if err != nil {
		return fmt.Errorf("GetDetails: %w", err)
	}

	protocolDocRef, err := uploadProtocol(fClient, protocol, protocolCol, force)
	if err != nil {
		return fmt.Errorf("uploadProtocol: %w", err)
	}

	// Create safe harbor adoption
	fAdoption := types.SafeHarborAgreementGeneric{
		SafeHarborAgreementBase: types.SafeHarborAgreementBase{
			AdoptionProposalURI: adoptionProposalUri,
			Protocol:            protocolDocRef,
			Slug:                "onchain-" + txhash.String(),
			Version:             version,
		},
		AgreementDetails:    details,
		AgreementAddress:    agreementAddress.String(),
		CreatedAt:           txbody.Time(),
		Creator:             sender.String(),
		RegistryTransaction: txhash.String(),
		RegistryChainID:     chain,
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
	protocol.SafeHarborAgreement = agreementDocRef
	_, err = protocolDocRef.Set(context.Background(), protocol)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	slog.Info("Successfully added adoption to database")
	// Send Telegram notification
	telegramMessage := "🚨 New Safe Harbor Adoption\n\n"
	telegramMessage += fmt.Sprintf("Protocol: %s\n", slug)
	telegramMessage += fmt.Sprintf("URL: https://safe-harbor-d9e89.web.app/database/%s\n", slug)

	// err = telegram.SendNotification(telegramMessage, os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("TELEGRAM_CHAT_ID"))
	// if err != nil {
	// 	slog.Error("Failed to send Telegram notification", "error", err)
	// }

	return nil
}

type AgreementDetailsV1Fetcher struct {
	chain int
}

func (f *AgreementDetailsV1Fetcher) GetDetails(txhash common.Hash) (common.Address, interface{}, error) {
	eClient, err := getChainClient(f.chain)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("getChainClient: %w", err)
	}

	agreementAddress, rawAgreement, err := safeharbor.GetAgreement(txhash, eClient)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("getAgreement: %w", err)
	}

	agreementDetails := types.AgreementDetailsV1{}
	err = agreementDetails.FromRawAgreementDetails(rawAgreement)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("FromRawAgreementDetails: %w", err)
	}

	scanClient, err := getScanClient(f.chain)
	if err == nil {
		slog.Info("Naming addresses...")
		agreementDetails.TryNameAddresses(scanClient)
	} else {
		slog.Warn("getScanClient", "error", err)
	}

	return *agreementAddress, agreementDetails, nil
}

type AgreementDetailsV2Fetcher struct {
	chain int
}

func (f *AgreementDetailsV2Fetcher) GetDetails(txhash common.Hash) (common.Address, interface{}, error) {
	eClient, err := getChainClient(f.chain)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("getChainClient: %w", err)
	}

	agreementAddress, err := safeharbor.GetAgreementAddress(txhash, eClient)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("GetAgreementAddress: %w", err)
	}

	v2Contract, err := adoptv2.NewAdoptiondetails(*agreementAddress, eClient)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("adoptiondetailsv2.NewAdoptiondetails: %w", err)
	}
	rawDetails, err := v2Contract.GetDetails(nil)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("adoptiondetailsv2.GetDetails: %w", err)
	}

	details := types.AgreementDetailsV2{}
	details.FromRawAgreementDetails(rawDetails)

	// Best-effort naming for EVM chains using CAIP-2 eip155:<id>
	details.TryNameAddressesByCAIP2(func(chainID int) (scan.Client, error) {
		return getScanClient(chainID)
	})

	return *agreementAddress, details, nil
}

type AgreementDetailsV3Fetcher struct {
	chain int
}

func (f *AgreementDetailsV3Fetcher) GetDetails(txhash common.Hash) (common.Address, interface{}, error) {
	eClient, err := getChainClient(f.chain)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("getChainClient: %w", err)
	}

	agreementAddress, err := safeharbor_v3.GetAgreementAddress(txhash, eClient)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("GetAgreementAddress: %w", err)
	}

	v2Contract, err := adoptv2.NewAdoptiondetails(*agreementAddress, eClient)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("NewSafeHarbor: %w", err)
	}

	rawDetails, err := v2Contract.GetDetails(nil)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("GetAgreementDetails: %w", err)
	}

	details := types.AgreementDetailsV2{}
	details.FromRawAgreementDetails(rawDetails)

	// Best-effort naming for EVM chains using CAIP-2 eip155:<id>
	details.TryNameAddressesByCAIP2(func(chainID int) (scan.Client, error) {
		return getScanClient(chainID)
	})

	return *agreementAddress, details, nil
}

func addCantinaAdoption(
	cClient *cantina.Client,
	fClient *firestore.Client,
	protocolCol string,
	agreementCol string,
	cantinaID string,
	defiliamaSlug string,
	force bool,
) error {
	protocol, err := defiliama.GetProtocol(defiliamaSlug)
	if err != nil {
		slog.Error("defiliama.GetProtocol", "slug", defiliamaSlug, "error", err)
		slog.Info("Falling back to no Defiliama protocol data")

		protocol = types.Protocol{
			Slug: defiliamaSlug,
			Name: defiliamaSlug,
		}
	}

	// Fetch safe harbor agreement from Cantina
	agreement, err := cClient.GetAgreement(cantinaID)
	if err != nil {
		return fmt.Errorf("GetAgreement: %w", err)
	}

	//* Upload protocol & adoption to firestore if not already present
	protocolDocRef, err := uploadProtocol(fClient, protocol, protocolCol, force)
	if err != nil {
		return fmt.Errorf("uploadProtocol: %w", err)
	}
	agreement.Protocol = protocolDocRef

	// Upload safe harbor adoption
	slug := "cantina-" + cantinaID
	slog.Info("Uploading adoption", "adoption", slug)
	agreementDocRef := fClient.Collection(agreementCol).Doc(slug)
	_, err = agreementDocRef.Set(context.Background(), agreement)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	// Update protocol with safe harbor agreement reference
	slog.Info("Updating protocol with safe harbor agreement reference")
	protocol.SafeHarborAgreement = agreementDocRef
	_, err = protocolDocRef.Set(context.Background(), protocol)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	slog.Info("Successfully added adoption to database")
	// Send Telegram notification
	telegramMessage := "🚨 New Cantina Safe Harbor Adoption\n\n"
	telegramMessage += fmt.Sprintf("Protocol: %s\n", slug)
	telegramMessage += fmt.Sprintf("URL: https://safe-harbor-d9e89.web.app/%s\n", slug)

	err = telegram.SendNotification(telegramMessage, os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {
		slog.Error("Failed to send Telegram notification", "error", err)
	}
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
	var protocol types.Protocol
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

	data := map[string]any{
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
	agreementCol string,
	slug string,
) error {
	var dClient = deduab.NewClient()
	var doc, err = fClient.Collection(agreementCol).Doc(slug).Get(context.Background())
	if err != nil {
		return fmt.Errorf("firestore.Get: %w", err)
	}

	var version types.AgreementVersion
	if err := doc.DataTo(&version); err != nil {
		return fmt.Errorf("firestore.DataTo(version): %w", err)
	}

	switch version.Version {
	case types.SealV1:
		err := refreshChildContractSealV1(doc, dClient)
		if err != nil {
			return fmt.Errorf("refreshChildContractSealV1: %w", err)
		}

	case types.ImmunefiV1:
		return nil // ImmunefiV1 does not have child contracts to refresh
	default:
		slog.Warn("refreshChildContracts not implemented", "version", version.Version, "slug", slug)
		return nil
	}

	return nil
}

func refreshAdoptionDate(
	fClient *firestore.Client,
	agreementCol string,
	slug string,
) error {
	doc, err := fClient.Collection(agreementCol).Doc(slug).Get(context.Background())
	if err != nil {
		return fmt.Errorf("firestore.Get: %w", err)
	}

	var version types.AgreementVersion
	if err := doc.DataTo(&version); err != nil {
		return fmt.Errorf("firestore.DataTo(version): %w", err)
	}

	switch version.Version {
	case types.SealV1:
		var agreement types.SafeHarborAgreementV1
		if err := doc.DataTo(&agreement); err != nil {
			return fmt.Errorf("firestore.DataTo(seal v1): %w", err)
		}

		eClient, err := getChainClient(agreement.RegistryChainID)
		if err != nil {
			return fmt.Errorf("getChainClient: %w", err)
		}

		tx := agreement.RegistryTransaction
		txHash := common.HexToHash(tx)
		txBody, err := eClient.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return fmt.Errorf("rpc.TransactionByHash: %w", err)
		}

		agreement.CreatedBlock = int(txBody.BlockNumber.Int64())

		block, err := eClient.BlockByNumber(context.Background(), txBody.BlockNumber)
		if err != nil {
			return fmt.Errorf("rpc.BlockByNumber: %w", err)
		}

		agreement.CreatedAt = time.UnixMilli(int64(block.Time() * 1000))

		// Update the document with the value of agreement
		_, err = doc.Ref.Set(context.Background(), agreement)
		if err != nil {
			return fmt.Errorf("firestore.Set: %w", err)
		}
	case types.ImmunefiV1:
		return nil
	default:
		slog.Error("refreshAdoptionDate not implemented", "version", version.Version, "slug", slug)
		return nil
	}

	return nil
}

func checkNewCantinaAdoptions(cClient *cantina.Client, fClient *firestore.Client, agreementCol string) ([]string, error) {
	agreements, err := cClient.GetAgreements()
	if err != nil {
		return nil, fmt.Errorf("cantina.GetAgreements: %w", err)
	}

	agreementDocs, err := fClient.Collection(agreementCol).Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore.GetAll: %w", err)
	}

	existingAdoptions := make(map[string]bool)
	for _, doc := range agreementDocs {
		existingAdoptions[doc.Ref.ID] = true
	}

	var newAdoptions []string
	for _, agreement := range agreements {
		if _, exists := existingAdoptions[agreement.Slug]; !exists {
			newAdoptions = append(newAdoptions, agreement.Slug)
		}
	}
	return newAdoptions, nil
}

func checkNewOnchainAdoptions(fClient *firestore.Client, agreementCol string) ([]string, error) {
	existingAdoptions := make(map[string]bool)
	agreementDocs, err := fClient.Collection(agreementCol).Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore.GetAll: %w", err)
	}

	for _, doc := range agreementDocs {
		existingAdoptions[doc.Ref.ID] = true
	}

	// Fetch adoptions from all chains
	var allNewAdoptions = make(map[string]bool)
	for chainID := range config.SafeHarborV1Registries {
		slog.Info("Checking chain for new adoptions...", "chainID", chainID)
		adoptions, err := checkChainForAdoptions(chainID)
		if err != nil {
			slog.Warn("checkChainForAdoptions", "chainID", chainID, "error", err)
			continue
		}

		for _, txHash := range adoptions {
			if _, exists := existingAdoptions[txHash]; exists {
				continue
			}

			chained_txhash := fmt.Sprintf("chain={%d} txhash={%s}", chainID, txHash)
			allNewAdoptions[chained_txhash] = true
		}
		slog.Info("Checked chain for new adoptions", "chainID", chainID, "found", len(adoptions))
	}

	var newAdoptions []string
	for txHash := range allNewAdoptions {
		newAdoptions = append(newAdoptions, txHash)
	}

	return newAdoptions, nil
}

func checkChainForAdoptions(chainId int) ([]string, error) {
	eClient, err := getChainClient(chainId)
	if err != nil {
		return nil, fmt.Errorf("getChainClient: %w", err)
	}

	contractAddress := config.SafeHarborV1Registries[chainId]

	filter, err := safeharbor.NewSafeharborFilterer(contractAddress, eClient)
	if err != nil {
		return nil, fmt.Errorf("NewSafeharborFilterer: %w", err)
	}

	latestBlock, err := eClient.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("BlockByNumber: %w", err)
	}

	var adoptions []string

	//? 100k blocks should be enough to cover new adoptions since last check on all chains
	fromBlock := latestBlock - 100_000
	// Scan in 1000 block increments
	for start := fromBlock; start < latestBlock; start += 1000 {
		end := start + 1000
		end = min(end, latestBlock)

		// slog.Debug("Scanning blocks for adoptions...", "chainID", chainId, "start", start, "end", end)
		opts := &bind.FilterOpts{
			Start: start,
			End:   &end,
		}

		// Get all adoption events
		iter, err := filter.FilterSafeHarborAdoption(opts, nil)
		if err != nil {
			return nil, fmt.Errorf("FilterSafeHarborAdoption: %w", err)
		}
		defer iter.Close()

		for iter.Next() {
			event := iter.Event
			txHash := event.Raw.TxHash.Hex()
			adoptions = append(adoptions, txHash)
		}

		if iter.Error() != nil {
			return nil, fmt.Errorf("event iteration error: %w", iter.Error())
		}
	}

	return adoptions, nil
}

func refreshChildContractSealV1(doc *firestore.DocumentSnapshot, dClient *deduab.Client) error {
	var agreement types.SafeHarborAgreementV1
	if err := doc.DataTo(&agreement); err != nil {
		return fmt.Errorf("firestore.DataTo(seal v1): %w", err)
	}

	for i, chain := range agreement.AgreementDetails.Chains {
		for j, account := range chain.Accounts {
			var endBlock int

			switch account.ChildContractScope {
			case types.ChildContractScopeAll:
				endBlock = 2147483647 // Nax for API
			case types.ChildContractScopeExistingOnly:
				endBlock = agreement.CreatedBlock
			default:
				continue
			}

			childContracts, err := dClient.GetDeployed(account.Address, endBlock, 32767, endBlock, 200)
			if err != nil {
				slog.Warn("GetDeployed", "address", account.Address, "error", err)
			}

			account.Children = make([]types.ChildAccountV1, 0, len(childContracts))
			for _, child := range childContracts {
				account.Children = append(account.Children, types.ChildAccountV1{
					Address: child.Address,
					Name:    child.Name,
				})
			}

			agreement.AgreementDetails.Chains[i].Accounts[j] = account
		}
	}

	_, err := doc.Ref.Set(context.Background(), agreement)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	return nil
}

func uploadProtocol(fClient *firestore.Client, protocol types.Protocol, protocolCol string, force bool) (*firestore.DocumentRef, error) {
	// Check if protocol exists
	protocolDocRef := fClient.Collection(protocolCol).Doc(protocol.Slug)
	protocolDoc, err := protocolDocRef.Get(context.Background())
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, fmt.Errorf("firestore.Get: %w", err)
	}

	exists := protocolDoc.Exists()
	if exists && !force {
		return nil, fmt.Errorf("protocol already exists in firestore. Use --force to overwrite")
	}

	slog.Info("Uploading protocol", "protocol", protocol.Name)
	_, err = protocolDocRef.Set(context.Background(), protocol)
	if err != nil {
		return nil, fmt.Errorf("firestore.Set: %w", err)
	}

	return protocolDocRef, nil
}

func getChainClient(chain int) (*ethclient.Client, error) {
	rpcUrlEnv := fmt.Sprintf("RPC_URL_%d", chain)
	rpcUrl := os.Getenv(rpcUrlEnv)
	if rpcUrl == "" {
		return nil, fmt.Errorf("Unsupported chain ID %d", chain)
	}

	eClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("rpc.Dial: %w", err)
	}
	return eClient, nil
}

func getScanClient(chain int) (scan.Client, error) {
	apiKey := os.Getenv("SCAN_KEY")
	scanClient := scan.NewRateLimitedClient(apiKey, chain)
	return scanClient, nil
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
