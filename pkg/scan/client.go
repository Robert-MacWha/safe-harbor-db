package scan

import (
	"fmt"
	"time"

	"github.com/nanmu42/etherscan-api"
)

type Tx struct {
	BlockNumber int
	TimeStamp   etherscan.Time
	Hash        string
	From        string
	To          string
	Value       *etherscan.BigInt
	Input       string
}

type Client interface {
	ContractName(address string) (name string)
	ContractSource(address string) (source etherscan.ContractSource, err error)
	NormalTxByAddress(address string, startBlock *int, endBlock *int, page int, offset int, desc bool) (txs []Tx, err error)
	InternalTxByAddress(address string, startBlock *int, endBlock *int, page int, offset int, desc bool) (txs []Tx, err error)
}

type RateLimitedClient struct {
	api *etherscan.Client

	lastCalled time.Time
}

const rateLimit = 200 * time.Millisecond

func NewRateLimitedClient(api *etherscan.Client) *RateLimitedClient {
	return &RateLimitedClient{api: api}
}

func (c *RateLimitedClient) ContractName(address string) (name string) {
	source, err := c.ContractSource(address)
	if err != nil {
		return ""
	}

	return source.ContractName
}

func (c *RateLimitedClient) ContractSource(address string) (source etherscan.ContractSource, err error) {
	defer c.sleep()

	sources, err := c.api.ContractSource(address)
	if err != nil {
		return etherscan.ContractSource{}, err
	}

	if len(sources) == 0 {
		return etherscan.ContractSource{}, fmt.Errorf("no source found for contract %s", address)
	}

	return sources[0], nil
}

func (c *RateLimitedClient) NormalTxByAddress(address string, startBlock *int, endBlock *int, page int, offset int, desc bool) (txs []Tx, err error) {
	defer c.sleep()

	normalTxns, err := c.api.NormalTxByAddress(address, startBlock, endBlock, page, offset, desc)
	if err != nil {
		return nil, err
	}

	txs = make([]Tx, len(normalTxns))
	for i, txn := range normalTxns {
		txs[i] = Tx{
			BlockNumber: txn.BlockNumber,
			TimeStamp:   txn.TimeStamp,
			Hash:        txn.Hash,
			From:        txn.From,
			To:          txn.To,
			Value:       txn.Value,
			Input:       txn.Input,
		}
	}

	return txs, nil
}

func (c *RateLimitedClient) InternalTxByAddress(address string, startBlock *int, endBlock *int, page int, offset int, desc bool) (txs []Tx, err error) {
	defer c.sleep()

	internalTxns, err := c.api.InternalTxByAddress(address, startBlock, endBlock, page, offset, desc)
	if err != nil {
		return nil, err
	}

	txs = make([]Tx, len(internalTxns))
	for i, txn := range internalTxns {
		txs[i] = Tx{
			BlockNumber: txn.BlockNumber,
			TimeStamp:   txn.TimeStamp,
			Hash:        txn.Hash,
			From:        txn.From,
			To:          txn.To,
			Value:       txn.Value,
			Input:       txn.Input,
		}
	}

	return txs, nil
}

// Sleep sleeps for the remaining time until the next call can be made.
func (c *RateLimitedClient) sleep() {
	time.Sleep(rateLimit)
	// sleepTime := rateLimit - time.Since(c.lastCalled)
	// if sleepTime > 0 {
	// 	time.Sleep(sleepTime)
	// }
	// c.lastCalled = time.Now()
}
