# SHDB

## Adding private modules

To go get the arianrhod module, make sure your terminal is logged into github via ssh, then run `export GOPRIVATE=github.com/Skylock-ai`. From there you can run `go mod tidy` as normal.

## Infrastructure

SHDB infrastructure is managed by ansible playbooks and github workflows. See [infrastructure](./infrastructure/README.md) for more info.

# TODO: Setup workflow to fetch adoption from immunefi

# TODO: Fix bug where polymarket has no children (and place limit on # of children provided to frontend)
