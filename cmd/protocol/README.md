# Protocol CLI Tool

## Description

The Protocol CLI tool is a simple command-line application for populating protocol data in Firestore. It fetches protocol details and optionally sets the Safe Harbor Agreement reference in the protocol's Firestore document.

## Usage

Run the CLI tool with the following flags:

## Flags

1. **`--protocol`** (`-p`) - **Required**: The protocol's Firestore document reference.
2. **`--creds`** (`-c`) - **Required**: Path to Firestore credentials file (JSON).
3. **`--setSafeHarbor`** (`-sp`) - **Optional**: Boolean flag to set the Safe Harbor Struct reference and the Protocol references to each other in Firestore. Default: `true`.

## Example

```bash
go run cmd/protocol/main.go --protocol=polymarket --creds=/path/to/firestore-creds.json --setSafeHarbor=false
```

This command will fetch the polymarket protocol, upload it to Firestore, and set will not set the Safe Harbor Agreement reference in the protocol's Firestore document.

Not setting the references to each other should be used when the other document isn't created yet (in this case the Safe Harbor Agreement document).
