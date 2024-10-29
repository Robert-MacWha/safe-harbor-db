package web3

import "math/big"

const weiToEth int64 = 1e18

// WeiToEth converts a balance in weth to a balance in eth
func WeiToEth(b BigInt) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(b.Int), big.NewFloat(float64(weiToEth)))
}
