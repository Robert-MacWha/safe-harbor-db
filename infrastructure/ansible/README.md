# Infrastructure

SHDB's infrastructure is based on [SEP-6](https://github.com/Skylock-ai/SEP/blob/main/decisions/SEP-6.md)

## Building

SHDB's build workflow is triggered on the creation of releases. The workflow builds a docker container for all CMDs and deploys them to the lighthouse docker registry in github.

The build workflow should be re-run whenever code changes are made. [Semantic versioning](https://semver.org/) is used.

## Deployment

SHDB's deploy workflow is triggered manually. Deploy workflows should exist per-application. They connect to the prod server and execute the coresponding ansible playbook and create the application's environment.

The deploy workflow should be re-run whenever secrets are updated or a new version has been built.

## Secret Management

Secrets are stored within AWS secret manager and are populated by ansible scripts during the deployment workflow. Under no condition should secrets be stored in plaintext in the github repository.

## Tasks

The current tasks include:

-   lighthouse

### Lighthouse

[Lighthouse](/cmd/lighthouse/README.md) maintaines the firestore `safeHarborAgreements` collection based on on-chain events and issues a notification whenever new agreements are adopted.
