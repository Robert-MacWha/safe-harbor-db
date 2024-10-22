# Safe Harbor Monitor CLI Tool

## Description

The Safe Harbor Monitor CLI tool continuously monitors blockchain events related to Safe Harbor agreements across multiple chains. It fetches events from Etherscan and processes newly created Safe Harbor agreements. The tool stores event details in Firestore. It is designed to work in parallel for different blockchain networks and registry configurations, updating Firestore with new agreements in real time.

## Usage

Run the CLI tool with the following flags:

### Flags

1. **`--config`** (`-c`) - **Required**: Path to the registry configuration JSON file.
2. **`--chainConfigs`** (`-cc`) - **Required**: Path to the blockchain chain configurations JSON file.
3. **`--creds`** (`-f`) - **Required**: Path to the Firestore credentials file.

## Example

```bash
go run main.go \
  --config=/path/to/registryConfigs.json \
  --chainConfigs=/path/to/chainConfigs.json \
  --creds=/path/to/firestore-creds.json
```
