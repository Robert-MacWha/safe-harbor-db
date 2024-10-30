# Safe Harbor Monitor CLI Tool

## Description

The Safe Harbor Monitor CLI tool continuously monitors blockchain events related to Safe Harbor agreements across multiple chains. It fetches events from Etherscan and processes newly created Safe Harbor agreements. The tool stores event details in Firestore. It is designed to work in parallel for different blockchain networks and registry configurations, updating Firestore with new agreements in real time.

## Usage

Run the CLI tool. Have an env file with the following variables:

-   FIREBASE_CREDENTIALS: A string of the json of the Firebase Credentials (literally copy and paste the json)

```
'{
  "type": "service_account",
  "project_id": "skylock-xyz",
  "private_key_id": "",
  "private_key": "",
  "client_email": "firebase-adminsdk-36s2d@skylock-xyz.iam.gserviceaccount.com",
  "client_id": "108051924350571600344",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-36s2d%40skylock-xyz.iam.gserviceaccount.com",
  "universe_domain": "googleapis.com"
}'
```

-   CHAIN_CONFIG: A dictionary of chainid, and APIKey & RPCURL for that chain

```
{
  "1": {
      "APIKey": "",
      "RPCURL": ""
  },
  "56": {
      "APIKey": "",
      "RPCURL": ""
  },
  "137": {
      "APIKey": "",
      "RPCURL": ""
  },
}
```

-   MAILGUN_API_KEY: The API key for Mailgun

## Example

```bash
go run cmd/lighthouse/main.go
```
