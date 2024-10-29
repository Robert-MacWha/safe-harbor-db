package web3

// DAIAddress = 0x6b175474e89094c44da98b954eedeac495271d0f
var DAIAddress, _ = HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")

// USDCAddress = 0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48
var USDCAddress, _ = HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")

// USDTAddress = 0xdac17f958d2ee523a2206206994597c13d831ec7
var USDTAddress, _ = HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")

// WETHAddress = 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
var WETHAddress, _ = HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")

// OUSDAddress = 0x2a8e1e676ec238d8a992307b495b45b3feaa5e86
var OUSDAddress, _ = HexToAddress("0x2a8e1e676ec238d8a992307b495b45b3feaa5e86")

// ETHAddress is a placeholder address used to represent ethereum
var ETHAddress, _ = HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")

// DefaultDeployer is a default deployer address that can be used when deploying
// replay contracts.
var DefaultDeployer, _ = HexToAddress("0x736b796c6f636bEf4A8908c57fb6B9d1A4b94dE4")

// DefaultBeneficiary is a default beneficiary address that can be used when
// deploying replay contracts.
var DefaultBeneficiary, _ = HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

// DefaultContractAddress is a default contract address that can be used when
// deploying replay contracts.
var DefaultContractAddress, _ = HexToAddress("0xd0343d30b8a75b206fa061f741088183b967e6df")

// ERC20Tokens is a map of ERC20 token names to their addresses.
var ERC20Tokens = map[string]Address{
	"DAI":  *DAIAddress,
	"USDC": *USDCAddress,
	"USDT": *USDTAddress,
	"WETH": *WETHAddress,
	"ETH":  *ETHAddress,
}

// ERC20Addresses is a map of ERC20 addresses to their token names.
var ERC20Addresses = map[Address]string{
	*DAIAddress:  "DAI",
	*USDCAddress: "USDC",
	*USDTAddress: "USDT",
	*WETHAddress: "WETH",
	*ETHAddress:  "ETH",
}

// ERC20ToUSD contains the divisor to convert from a given ERC20 to USD
//
// nolint:mnd
var ERC20ToUSD = map[Address]float64{
	*USDCAddress: 1 / 1e6,
	*USDTAddress: 1 / 1e6,
	*DAIAddress:  1 / 1e18,
}
