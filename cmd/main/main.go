package main

import (
	"SHDB/pkg/call"
	"SHDB/pkg/config"
	"SHDB/pkg/contracts/safeharbor"
	"SHDB/pkg/defiliama"
	"SHDB/pkg/firebase"
	"SHDB/pkg/scan"
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/nanmu42/etherscan-api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const EtherscanErrNoTxns string = "etherscan server: No transactions found"

const (
	ProtocolsCollection = "protocols-test"
	AgreementCollection = "safeHarborAgreement-test"
)

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
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "Force the command to run even if the protocol already exists in the database",
						Value:   false,
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
	txhash := common.HexToHash(txHashStr)

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

	err = addAdoption(eClient, fClient, slug, chain, txhash, force)
	if err != nil {
		return fmt.Errorf("addAdoption: %w", err)
	}

	err = refreshChildContracts(fClient, slug, chainCfg)
	if err != nil {
		return fmt.Errorf("refreshChildContracts: %w", err)
	}

	return nil
}

func runRefreshTvl(cCtx *cli.Context) error {
	_ = godotenv.Load()

	//* Load config
	slug := cCtx.String("slug")

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	//* Refresh single protocol
	if slug != "all" {
		return refreshTvl(fClient, slug)
	}

	//* Refresh all protocols
	protocols := fClient.Collection(ProtocolsCollection)
	documents, err := protocols.Documents(context.Background()).GetAll()
	if err != nil {
		return fmt.Errorf("firestore.GetAll: %w", err)
	}

	for _, doc := range documents {
		slug := doc.Ref.ID
		err = refreshTvl(fClient, slug)
		if err != nil {
			return fmt.Errorf("refreshTvl: %w", err)
		}
	}

	return nil
}

func runRefreshChildContracts(cCtx *cli.Context) error {
	_ = godotenv.Load()

	//* Load config
	slug := cCtx.String("slug")

	fClient, err := firebase.NewFirestoreClient()
	if err != nil {
		return fmt.Errorf("NewFirestoreClient: %w", err)
	}

	chainCfg, err := config.LoadChainCfg()
	if err != nil {
		return fmt.Errorf("LoadChainCfg: %w", err)
	}

	err = refreshChildContracts(fClient, slug, chainCfg)
	if err != nil {
		return fmt.Errorf("refreshChildContracts: %w", err)
	}

	return nil
}

func addAdoption(eClient *ethclient.Client, fClient *firestore.Client, slug string, chain int, txhash common.Hash, force bool) error {
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
		return fmt.Errorf("defiliama.GetProtocol: %w", err)
	}

	//* Upload protocol & adoption to firestore if not already present
	// Check if protocol exists
	protocolDocRef := fClient.Collection(ProtocolsCollection).Doc(protocol.Slug)
	protocolDoc, err := protocolDocRef.Get(context.Background())
	if err != nil && status.Code(err) != codes.NotFound {
		return fmt.Errorf("firestore.Get: %w", err)
	}

	exists := protocolDoc.Exists()
	if !exists && !force {
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
		ContactDetails:      firebase.FormatContactDetails(agreement.ContactDetails),
		Chains:              firebase.FormatChains(agreement.Chains),
		BountyTerms:         firebase.FormatBountyTerms(agreement.BountyTerms),
		CreatedAt:           txbody.Time(),
		CreatedBlock:        int(receipt.BlockNumber.Int64()),
	}

	// Upload safe harbor adoption
	slog.Info("Uploading adoption", "adoption", txhash.String())
	agreementDocRef := fClient.Collection(AgreementCollection).Doc(txhash.String())
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

func refreshTvl(fClient *firestore.Client, slug string) error {
	//* Confirm protocol exists
	protocolDocRef := fClient.Collection(ProtocolsCollection).Doc(slug)
	protocolDoc, err := protocolDocRef.Get(context.Background())
	if err != nil {
		return fmt.Errorf("firestore.Get: %w", err)
	}

	if !protocolDoc.Exists() {
		return fmt.Errorf("protocol not found in firestore")
	}

	//* Fetch and update TVL
	protocol, err := defiliama.GetProtocol(slug)
	if err != nil {
		return fmt.Errorf("defiliama.GetProtocol: %w", err)
	}

	tvl := protocol.TVL
	data := map[string]interface{}{
		"tvl": tvl,
	}

	_, err = protocolDocRef.Set(context.Background(), data, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	slog.Info("Successfully refreshed TVL")

	return nil
}

func refreshChildContracts(
	fClient *firestore.Client,
	slug string,
	chainCfg map[int]config.ChainCfg,
) error {
	protocolDocRef := fClient.Collection(ProtocolsCollection).Doc(slug)
	protocolDoc, err := protocolDocRef.Get(context.Background())
	if err != nil {
		return fmt.Errorf("firestore.Get: %w", err)
	}

	var protocol firebase.Protocol
	err = protocolDoc.DataTo(&protocol)
	if err != nil {
		return fmt.Errorf("firestore.DataTo: %w", err)
	}

	agreementId := protocol.SafeHarborAgreement.ID
	agreementDocRef := fClient.Collection(AgreementCollection).Doc(agreementId)
	agreementDoc, err := agreementDocRef.Get(context.Background())
	if err != nil {
		return fmt.Errorf("firestore.Get: %w", err)
	}

	var agreement firebase.SafeHarborAgreement
	err = agreementDoc.DataTo(&agreement)
	if err != nil {
		return fmt.Errorf("firestore.DataTo: %w", err)
	}

	totalChildren := 0
	for i, chain := range agreement.Chains {
		if _, exists := chainCfg[chain.ID]; !exists {
			slog.Info("Chain not found in chain config, skipping", "chain", chain.ID)
			continue
		}

		eClient, err := ethclient.Dial(chainCfg[chain.ID].RpcUrl)
		if err != nil {
			return fmt.Errorf("rpc.Dial: %w", err)
		}

		sClient := scan.NewRateLimitedClient(etherscan.NewCustomized(etherscan.Customization{
			Timeout: 15 * time.Second,
			Key:     chainCfg[chain.ID].ScanKey,
			BaseURL: chainCfg[chain.ID].ScanUrl,
			Verbose: false,
		}))

		block, err := eClient.BlockNumber(context.Background())
		if err != nil {
			return fmt.Errorf("rpc.BlockNumber: %w", err)
		}

		lastIndexed := chain.LastIndexedBlock
		for j, account := range chain.Accounts {
			log := slog.Default().With("account", account.Address)
			if account.ChildContractScope == firebase.ChildContractScopeNone {
				log.Debug("Skipping account", "reason", "ChildContractScopeNone")
				continue
			}

			var b = int(block)
			var endBlock *int = &b
			if account.ChildContractScope == firebase.ChildContractScopeExistingOnly {
				endBlock = &agreement.CreatedBlock
			}

			log.Debug("Refreshing child contracts", "scope", account.ChildContractScope, "from", lastIndexed, "to", endBlock)
			var accountAddr = common.HexToAddress(account.Address)
			children, err := getAccountChildren(eClient, sClient, accountAddr, lastIndexed, endBlock)
			if err != nil {
				return fmt.Errorf(
					"getAccountChildren(accountAddr=%v, chain=%v, agreement=%v): %w",
					accountAddr,
					chain.ID,
					agreement.AgreementAddress,
					err,
				)
			}

			agreement.Chains[i].Accounts[j].Children = children
			totalChildren += len(children)
		}

		agreement.Chains[i].LastIndexedBlock = int(block)
		slog.Info(
			"Successfully refreshed chain",
			"agreement", agreement.AgreementAddress,
			"chain", chain.ID,
			"start", lastIndexed,
			"end", block,
		)
	}

	// Upload updated agreement to firestore
	_, err = agreementDocRef.Set(context.Background(), agreement)
	if err != nil {
		return fmt.Errorf("firestore.Set: %w", err)
	}

	slog.Info("Successfully refreshed child contracts", "agreement", agreement.AgreementAddress, "total", totalChildren)

	return nil
}

func getAccountChildren(
	eClient *ethclient.Client,
	sClient scan.Client,
	accountAddr common.Address,
	startBlock int,
	endBlock *int,
) (children []firebase.ChildAccount, err error) {
	// Get all txhashes that interacted with this account
	txnHashes, err := getAccountTxns(sClient, accountAddr, startBlock, endBlock)
	if err != nil {
		return nil, fmt.Errorf("getAccountTxns: %w", err)
	}

	if len(txnHashes) == 0 {
		return []firebase.ChildAccount{}, nil
	}
	slog.Debug("Found transactions", "account", accountAddr, "count", len(txnHashes))

	// Trace all transactions to find subcontracts
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	children = []firebase.ChildAccount{}
	for _, hash := range txnHashes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newChildren := getTxnChildren(eClient, sClient, hash, accountAddr, startBlock, endBlock)

			mu.Lock()
			children = append(children, newChildren...)
			mu.Unlock()
		}()
	}

	wg.Wait()
	if len(children) != 0 {
		slog.Debug("Found child contracts", "account", accountAddr, "count", len(children))
	}

	return children, nil
}

// getAccountTxns returns all transactions that interacted with an account, including internal transactions
//
// ? May return duplicate transactions, but they will be infrequent. Processing duplicates is probably more performant than
// ? trying to filter them out for large sets.
func getAccountTxns(sClient scan.Client, accountAddr common.Address, startBlock int, endBlock *int) ([]common.Hash, error) {
	var hashes []common.Hash

	// TODO: Merge these two loops if you can find a clean approach
	// TODO: Write smarter rate limit handling
	// Fetch regular transactions
	lastRegularIndexed := startBlock
	for {
		txns, err := sClient.NormalTxByAddress(accountAddr.String(), &lastRegularIndexed, endBlock, 0, 0, true)
		if err != nil && err.Error() == EtherscanErrNoTxns {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("etherscan.NormalTxByAddress: %w", err)
		}

		updatedLastIndexed := false
		for _, txn := range txns {
			hashes = append(hashes, common.HexToHash(txn.Hash))
			if txn.BlockNumber > lastRegularIndexed {
				lastRegularIndexed = txn.BlockNumber
				updatedLastIndexed = true
			}
		}

		// Break if no new transactions are found
		if !updatedLastIndexed {
			break
		}

	}

	// Fetch internal transactions
	lastInternalIndexed := startBlock
	for {
		txns, err := sClient.InternalTxByAddress(accountAddr.String(), &lastInternalIndexed, endBlock, 0, 0, true)
		if err != nil && err.Error() == EtherscanErrNoTxns {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("etherscan.InternalTxByAddress: %w", err)
		}

		updatedLastIndexed := false
		for _, txn := range txns {
			hashes = append(hashes, common.HexToHash(txn.Hash))
			if txn.BlockNumber > lastInternalIndexed {
				lastInternalIndexed = txn.BlockNumber
				updatedLastIndexed = true
			}
		}

		// Break if no new transactions are found
		if !updatedLastIndexed {
			break
		}

	}

	return hashes, nil
}

func getTxnChildren(
	eClient *ethclient.Client,
	sClient scan.Client,
	hash common.Hash,
	accountAddr common.Address,
	startBlock int,
	endBlock *int,
) []firebase.ChildAccount {
	calls, err := call.DebugTraceTransaction(eClient, hash)
	if err != nil {
		slog.Warn("Failed to debug trace transaction", "account", accountAddr, "hash", hash, "error", err)
		return nil
	}

	newChildren := []firebase.ChildAccount{}
	for _, call := range calls.Flatten() {
		if call.Type == "CREATE" && call.From == accountAddr && call.To != nil {
			recursiveChildren, err := getAccountChildren(eClient, sClient, *call.To, startBlock, endBlock)
			if err != nil {
				slog.Warn("Failed to get recursive children", "account", accountAddr, "error", err)
				return nil
			}

			newChildren = append(newChildren, firebase.ChildAccount{
				Address: call.To.String(),
			})
			newChildren = append(newChildren, recursiveChildren...)
		}
	}

	return newChildren
}
