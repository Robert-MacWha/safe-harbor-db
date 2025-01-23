package client

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient interface {
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	BlockNumber(ctx context.Context) (uint64, error)
	DebugTraceTransaction(hash common.Hash) (*Calls, error)
}

type RateLimitedClient struct {
	c *ethclient.Client

	lastCalled time.Time
}

const rateLimit = 75 * time.Millisecond

func NewRateLimitedClient(c *ethclient.Client) *RateLimitedClient {
	return &RateLimitedClient{c: c}
}

func Dial(url string) (*RateLimitedClient, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	return NewRateLimitedClient(client), nil
}

func (r *RateLimitedClient) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	defer r.sleep()

	return r.c.TransactionByHash(ctx, hash)
}

func (r *RateLimitedClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	defer r.sleep()

	return r.c.TransactionReceipt(ctx, txHash)
}

func (r *RateLimitedClient) BlockNumber(ctx context.Context) (uint64, error) {
	defer r.sleep()

	return r.c.BlockNumber(ctx)
}

func (r *RateLimitedClient) DebugTraceTransaction(hash common.Hash) (*Calls, error) {
	defer r.sleep()

	var result Calls
	err := r.c.Client().Call(&result, "debug_traceTransaction", hash.String(), map[string]interface{}{
		"tracer": "callTracer",
		"tracerConfig": map[string]interface{}{
			"onlyTopCall": false,
		},
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Sleep sleeps for the remaining time until the next call can be made.
func (r *RateLimitedClient) sleep() {
	sleepTime := rateLimit - time.Since(r.lastCalled)
	if sleepTime > 0 {
		time.Sleep(sleepTime)
	}
	r.lastCalled = time.Now()
}
