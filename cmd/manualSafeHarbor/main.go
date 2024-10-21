package main

import (
	"SHDB/pkg/flow"
	safeharbor "SHDB/pkg/safeHarbor"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/option"
)

func main() {
	// Define the CLI app
	app := &cli.App{
		Name:  "safe-harbor-agreement",
		Usage: "CLI to process Safe Harbor adoptions for a specific safe harbor adoption and store in Firestore",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to the chain config JSON file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "txHash",
				Aliases:  []string{"t"},
				Usage:    "Transaction hash",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "safeHarborAddress",
				Aliases:  []string{"s"},
				Usage:    "Safe Harbor Address",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "deployer",
				Aliases:  []string{"d"},
				Usage:    "Deployer address",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "chainId",
				Aliases:  []string{"i"},
				Usage:    "Chain ID",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "blockNumber",
				Aliases:  []string{"b"},
				Usage:    "Block number",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "protocol",
				Aliases:  []string{"p"},
				Usage:    "Protocol name for Firestore",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "creds",
				Aliases:  []string{"f"},
				Usage:    "Path to Firestore credentials file",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "setProtocol",
				Aliases: []string{"sp"},
				Usage:   "Set Safe Harbor Agreement reference in protocol",
				Value:   true, // Default to true
			},
		},
		Action: Run,
	}

	// Run the CLI app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}

// Run is the main function that gets called when the CLI app is run
func Run(c *cli.Context) error {
	// Load chain configs from JSON
	chainConfigs, err := loadChainConfigs(c.String("config"))
	if err != nil {
		return fmt.Errorf("failed to load chain configs: %w", err)
	}

	// Parse command-line inputs
	txHash, err := web3.HexToHash(c.String("txHash"))
	if err != nil {
		return fmt.Errorf("failed to parse txHash: %w", err)
	}

	safeHarborAddress, err := web3.HexToAddress(c.String("safeHarborAddress"))
	if err != nil {
		return fmt.Errorf("failed to parse safeHarborAddress: %w", err)
	}

	deployer, err := web3.HexToAddress(c.String("deployer"))
	if err != nil {
		return fmt.Errorf("failed to parse deployer: %w", err)
	}

	chainId := c.Int("chainId")
	blockNumber := big.NewInt(int64(c.Int("blockNumber")))

	// Load Firestore credentials and create Firestore client
	credsPath := c.String("creds")
	credsFile, err := os.Open(credsPath)
	if err != nil {
		return fmt.Errorf("failed to open Firestore credentials file: %w", err)
	}
	defer credsFile.Close()

	firestoreClient, err := newFirestoreClient(credsFile)
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}
	defer firestoreClient.Close()

	// Call the flow.ProcessSafeHarborAgreement function
	err = flow.ProcessSafeHarborAgreement(
		chainConfigs,
		*txHash,
		*safeHarborAddress,
		*deployer,
		chainId,
		blockNumber,
		c.String("protocol"),
		firestoreClient,
		c.Bool("setProtocol"),
	)
	if err != nil {
		return fmt.Errorf("failed to process Safe Harbor agreement: %w", err)
	}

	fmt.Printf("Successfully processed protocol: %s\n", c.String("protocol"))
	return nil
}

// loadChainConfigs loads the chain configurations from a JSON file
func loadChainConfigs(filePath string) (map[int64]safeharbor.ChainConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open chain config file: %w", err)
	}
	defer file.Close()

	var chainConfigs map[int64]safeharbor.ChainConfig
	err = json.NewDecoder(file).Decode(&chainConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode chain config JSON: %w", err)
	}

	return chainConfigs, nil
}

// newFirestoreClient creates a Firestore client using credentials
func newFirestoreClient(creds io.Reader) (*firestore.Client, error) {
	ctx := context.Background()

	// Read credentials file
	credsData, err := io.ReadAll(creds)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	// Create Firestore client
	client, err := firestore.NewClient(ctx, "skylock-xyz", option.WithCredentialsJSON(credsData))
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}

	return client, nil
}
