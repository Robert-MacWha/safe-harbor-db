package main

import (
	"SHDB/pkg/flow"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	// Define the CLI app
	app := &cli.App{
		Name:  "protocol",
		Usage: "A simple CLI app to populate Protocol firebase struct",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "protocol",
				Aliases:  []string{"p"},
				Usage:    "Protocol's Firestore Document Ref",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "creds",
				Aliases:  []string{"c"},
				Usage:    "Path to Firestore credentials file",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "setSafeHarbor",
				Usage: "Set the Safe Harbor Agreement reference",
				Value: true,
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
	// Get the protocol name and credentials path from the flags
	protocolName := c.String("protocol")
	credsPath := c.String("creds")
	setSafeHarbor := c.Bool("setSafeHarbor")

	// Call the flow function to process the protocol
	err := flow.ProcessProtocol(protocolName, credsPath, setSafeHarbor)
	if err != nil {
		return err
	}

	return nil
}
