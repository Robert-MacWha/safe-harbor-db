package scan

import (
	"log/slog"
	"sync"
	"sync/atomic"
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
	NormalTxByAddress(address string, startBlock *int, endBlock *int, page int, offset int, desc bool) (txs []Tx, err error)
	InternalTxByAddress(address string, startBlock *int, endBlock *int, page int, offset int, desc bool) (txs []Tx, err error)
}

type RateLimitedClient struct {
	api *etherscan.Client

	waitCount atomic.Int64
	mu        sync.Mutex
}

func NewRateLimitedClient(api *etherscan.Client) *RateLimitedClient {
	return &RateLimitedClient{api: api}
}

func (c *RateLimitedClient) NormalTxByAddress(address string, startBlock *int, endBlock *int, page int, offset int, desc bool) (txs []Tx, err error) {
	c.waitCount.Add(1)
	c.mu.Lock()

	defer slog.Debug("RateLimitedClient", "waiting", c.waitCount.Load())
	defer c.waitCount.Add(-1)
	defer c.mu.Unlock()

	normalTxns, err := c.api.NormalTxByAddress(address, startBlock, endBlock, page, offset, desc)
	time.Sleep(200 * time.Millisecond)
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
	c.waitCount.Add(1)
	c.mu.Lock()

	defer slog.Debug("RateLimitedClient", "waiting", c.waitCount.Load())
	defer c.waitCount.Add(-1)
	defer c.mu.Unlock()

	internalTxns, err := c.api.InternalTxByAddress(address, startBlock, endBlock, page, offset, desc)
	time.Sleep(200 * time.Millisecond)
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
