# Safe Harbor CLI Tool

## Description

The Safe Harbor CLI tool processes Safe Harbor agreements for a specific transaction and uploads the agreement details to Firestore. It will grab the scope details from the chain, process the children if necessary, then upload the Safe Harbor Agreement to Firestore.

Warning: If there are sub-contracts that need to be processed, it might take a while. So please be patient!

## Usage

Run the CLI tool with the following flags:

## Flags

1. **`--config`** (`-c`) - **Required**: Path to the chain config JSON file.
2. **`--txHash`** (`-t`) - **Required**: The Ethereum transaction hash.
3. **`--safeHarborAddress`** (`-s`) - **Required**: The Safe Harbor contract address.
4. **`--deployer`** (`-d`) - **Required**: The deployer address (the address that deployed the Safe Harbor contract).
5. **`--chainId`** (`-i`) - **Required**: The blockchain chain ID.
6. **`--blockNumber`** (`-b`) - **Required**: The block number where the transaction occurred.
7. **`--protocol`** (`-p`) - **Required**: The protocol name for Firestore (document reference).
8. **`--creds`** (`-f`) - **Required**: Path to the Firestore credentials file.
9. **`--setProtocol`** (`-sp`) - **Optional**: Boolean flag to set Safe Harbor Agreement reference in the protocol. Default: `true`.

## Example

```bash
go run cmd/manualSafeHarbor/main.go \
  --config=/path/to/chainConfigs.json \
  --txHash=0x62a554a7a8f8a7ab49f41b4df5b72eea6ca30680adc2f61f608d3c2d47296685 \
  --safeHarborAddress=0x2f6748580b200b9b2ace5774edc2657ff7ccc56b \
  --deployer=0x566345a70d70ce724cc1a441dca748b6b6c31265 \
  --chainId=137 \
  --blockNumber=63210447 \
  --protocol=polymarket \
  --creds=/path/to/firestore-creds.json \
  --setProtocol=true
```
