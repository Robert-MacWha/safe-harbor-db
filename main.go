package main

import (
	"SHDB/pkg/clients"
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/rpc"

	firestoreClients "SHDB/pkg/firestore"

	firebase "firebase.google.com/go"
	"github.com/Skylock-ai/Arianrhod/pkg/types/web3"
	"google.golang.org/api/option"
)

const (
	apiKey           = "PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3"
	rpcEndpoint      = "https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/"
	checkInterval    = 1 * time.Hour
	firebaseCredFile = "/home/dwu/Arianrhod/testdata/skylock-xyz-firebase-adminsdk-36s2d-265672c820.json"
)

func addressOrPanic(address string) web3.Address {
	addr, err := web3.HexToAddress(address)
	if err != nil {
		panic(err)
	}
	return *addr
}

func hashOrPanic(hash string) web3.Hash {
	h, err := web3.HexToHash(hash)
	if err != nil {
		panic(err)
	}
	return *h
}

func main() {
	ctx := context.Background()

	rpcClient, err := rpc.Dial("https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/")
	if err != nil {
		panic(err)
	}

	// Initialize Firestore client
	sa := option.WithCredentialsFile(firebaseCredFile)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	// firewallClient := clients.Firewall{
	// 	ProtocolName: "Native",
	// 	Accounts: []clients.AccountFW{
	// 		{
	// 			Address: addressOrPanic("0x85b0f66e83515ff4e825dfcaa58e040e08278ef9"),
	// 			RegisteredEvents: []clients.EventFW{
	// 				{
	// 					Topic0:              hashOrPanic("0xcf4eef09472ed02a0130e7d806d47fbd0db559056897eb16edb78506bdf45eed"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     12,
	// 				},
	// 				{
	// 					Topic0:              hashOrPanic("0x82cc5194ed1e660e8a6a4b2c99f7305283b16d8db2a1b2c9bb440096a4a07435"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     12,
	// 				},
	// 				{
	// 					Topic0:              hashOrPanic("0x583aa3202641e70abf6c4e526dd0bf713aa272b67f557ed30f2e40a37d985a8e"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     108,
	// 				},
	// 			},
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 			Name:               "Native PoolFactory",
	// 		},
	// 	},
	// }

	// // "0x8392f6669292fa56123f71949b52d883ae57e225", // Treasury
	// // 	"0x9e2b6378ee8ad2a4a95fe481d63caba8fb0ebbf9", // Dev Multisig
	// // 	"0x5C6374a2ac4EBC38DeA0Fc1F8716e5Ea1AdD94dd", // alUSDAlchemist
	// // 	"0x062Bf725dC4cDF947aa79Ca2aaCCD4F385b13b5c", // alETHAlchemist
	// // 	"0x9735f7d3ea56b454b24ffd74c58e9bd85cfad31b", // alUSD AMO
	// // 	"0x06378717d86b8cd2dba58c87383da1eda92d3495", // alUSDFRAXBP AMO
	// // 	"0xe761bf731a06fe8259fee05897b2687d56933110", // alETH AMO
	// // 	"0x9fb54d1F6F506Feb4c65B721bE931e59BB538c63", // alETH/fraxETH AMO
	// // 	"0xeE69BD81Bd056339368c97c4B2837B4Dc4b796E7", // USDTransmuterB
	// // 	"0xb039eA6153c827e59b620bDCd974F7bbFe68214A", // USDYearnVaultAdapter
	// // 	"0x6Fe02BE0EC79dCF582cBDB936D7037d2eB17F661", // USDYearnVaultAdapterTransmuterB
	// // 	"0x9FD9946E526357B35D95Bcb4b388614be4cFd4AC", // ETHTransmuter
	// // 	"0xf8317BD5F48B6fE608a52B48C856D3367540B73B", // ETHAlchemist
	// // 	"0x546E6711032Ec744A7708D4b7b283A210a85B3BC", // ETHYearnVaultAdapter
	// // 	"0x6d75657771256C7a8CB4d475fDf5047B70160132", // ETHYearnVaultAdapterB
	// // 	"0xb4E7cc74e004F95AEe7565a97Dbfdea9c1761b24", // WETHMigration
	// // 	"0x72A7cb4d5daB8E9Ba23f30DBE8E72Bc854a9945A", // DAIMigration
	// // 	"0xAB8e74017a8Cc7c15FFcCd726603790d26d7DeCa", // Staking
	// // 	"0xA840C73a004026710471F727252a9a2800a5197F", // Transmuter (DAI) alUSD
	// // 	"0x49930AD9eBbbc0EB120CCF1a318c3aE5Bb24Df55", // Transmuter (USDC) alUSD
	// // 	"0xfC30820ba6d045b95D13a5B8dF4fB0E6B5bdF5b9", // Transmuter (USDT) alUSD
	// // 	"0xE107Fa35D775C77924926C0292a9ec1FC14262b2", // Transmuter (FRAX) alUSD
	// // 	"0x1EEd2DbeB9fc23Ab483F447F38F289cA15f79Bac", // TransmuterBuffer alUSD
	// // 	"0x03323143a5f0D0679026C2a9fB6b0391e4D64811", // Transmuter (WETH) alETH
	// // 	"0xbc2FB245594a68c927C930FBE2d00680A8C90B9e", // TransmuterBuffer alETH
	// // 	"0xBE1C919cA137299715e9c929BC7126Af14f76091", // DefiLlama alUSD+FRAXBP AMO

	// firewallClient := clients.Firewall{
	// 	ProtocolName: "Alchemix",
	// 	Accounts: []clients.AccountFW{
	// 		{
	// 			Name:               "Treasury",
	// 			Address:            addressOrPanic("0x8392f6669292fa56123f71949b52d883ae57e225"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Dev Multisig",
	// 			Address:            addressOrPanic("0x9e2b6378ee8ad2a4a95fe481d63caba8fb0ebbf9"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "STD Controller",
	// 			Address:            addressOrPanic("0x3216d2a52f0094aa860ca090bc5c335de36e6273"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "alUSD Alchemist",
	// 			Address:            addressOrPanic("0x5C6374a2ac4EBC38DeA0Fc1F8716e5Ea1AdD94dd"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "alETH Alchemist",
	// 			Address:            addressOrPanic("0x062Bf725dC4cDF947aa79Ca2aaCCD4F385b13b5c"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "alUSD AMO",
	// 			Address:            addressOrPanic("0x9735f7d3ea56b454b24ffd74c58e9bd85cfad31b"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "alUSD FRAXBP AMO",
	// 			Address:            addressOrPanic("0x06378717d86b8cd2dba58c87383da1eda92d3495"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "alETH AMO",
	// 			Address:            addressOrPanic("0xe761bf731a06fe8259fee05897b2687d56933110"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "alETH FRAXETH AMO",
	// 			Address:            addressOrPanic("0x9fb54d1F6F506Feb4c65B721bE931e59BB538c63"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "USD TransmuterB",
	// 			Address:            addressOrPanic("0xeE69BD81Bd056339368c97c4B2837B4Dc4b796E7"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "USD YearnVaultAdapter",
	// 			Address:            addressOrPanic("0xb039eA6153c827e59b620bDCd974F7bbFe68214A"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "USD YearnVaultAdapter TransmuterB",
	// 			Address:            addressOrPanic("0x6Fe02BE0EC79dCF582cBDB936D7037d2eB17F661"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "ETH Transmuter",
	// 			Address:            addressOrPanic("0x9FD9946E526357B35D95Bcb4b388614be4cFd4AC"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "ETH Alchemist",
	// 			Address:            addressOrPanic("0xf8317BD5F48B6fE608a52B48C856D3367540B73B"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "ETH YearnVaultAdapter",
	// 			Address:            addressOrPanic("0x546E6711032Ec744A7708D4b7b283A210a85B3BC"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "ETH YearnVaultAdapterB",
	// 			Address:            addressOrPanic("0x6d75657771256C7a8CB4d475fDf5047B70160132"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "WETH Migration",
	// 			Address:            addressOrPanic("0xb4E7cc74e004F95AEe7565a97Dbfdea9c1761b24"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "DAI Migration",
	// 			Address:            addressOrPanic("0x72A7cb4d5daB8E9Ba23f30DBE8E72Bc854a9945A"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Staking",
	// 			Address:            addressOrPanic("0xAB8e74017a8Cc7c15FFcCd726603790d26d7DeCa"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Transmuter (DAI) alUSD",
	// 			Address:            addressOrPanic("0xA840C73a004026710471F727252a9a2800a5197F"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Transmuter (USDC) alUSD",
	// 			Address:            addressOrPanic("0x49930AD9eBbbc0EB120CCF1a318c3aE5Bb24Df55"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Transmuter (USDT) alUSD",
	// 			Address:            addressOrPanic("0xfC30820ba6d045b95D13a5B8dF4fB0E6B5bdF5b9"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Transmuter (FRAX) alUSD",
	// 			Address:            addressOrPanic("0xE107Fa35D775C77924926C0292a9ec1FC14262b2"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "TransmuterBuffer alUSD",
	// 			Address:            addressOrPanic("0x1EEd2DbeB9fc23Ab483F447F38F289cA15f79Bac"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Transmuter (WETH) alETH",
	// 			Address:            addressOrPanic("0x03323143a5f0D0679026C2a9fB6b0391e4D64811"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "TransmuterBuffer alETH",
	// 			Address:            addressOrPanic("0xbc2FB245594a68c927C930FBE2d00680A8C90B9e"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "AlETH lockbox",
	// 			Address:            addressOrPanic("0x9141776017D6A8a8522f913fddFAcAe3e84a7CDb"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "AlUSD lockbox",
	// 			Address:            addressOrPanic("0x2930cda830b206c84ae8d4ca3f77ec0eaa77a14b"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 	},
	// }

	// // // Market Factory V1
	// // 	{"0x27b1dAcd74688aF24a64BD3C9C1B143118740784", []string{
	// // 		// CreateNewMarket
	// // 		"0x166ae5f55615b65bbd9a2496e98d4e4d78ca15bd6127c0fe2dc27b76f6c03143",
	// // 	},
	// // 		[]addressExtractor{
	// // 			extractAddressFromTopic1, // Market
	// // 		}},
	// // 	// Market Factory V3
	// // 	{"0x1A6fCc85557BC4fB7B534ed835a03EF056552D52", []string{
	// // 		// CreateNewMarket
	// // 		"0xae811fae25e2770b6bd1dcb1475657e8c3a976f91d1ebf081271db08eef920af",
	// // 	}, []addressExtractor{
	// // 		extractAddressFromTopic1, // Market
	// // 	}},
	// // 	// Yield Contract V1
	// // 	{"0x70ee0A6DB4F5a2Dc4d9c0b57bE97B9987e75BAFD", []string{
	// // 		// CreateYieldContract
	// // 		"0xaa79d8f17776adeaa316c5411b72e8b0057d064974fa8748f32492ecaa22ecd1",
	// // 	}, []addressExtractor{
	// // 		extractAddressFromTopic1,     // SY
	// // 		extractAddressFromData12To32, // PT
	// // 		extractAddressFromData44To64, // YT
	// // 	}},
	// // 	// Yield Contract V3
	// // 	{"0xdF3601014686674e53d1Fa52F7602525483F9122", []string{
	// // 		// CreateYieldContract
	// // 		"0xaa79d8f17776adeaa316c5411b72e8b0057d064974fa8748f32492ecaa22ecd1",
	// // 	}, []addressExtractor{
	// // 		extractAddressFromTopic1,     // SY
	// // 		extractAddressFromData12To32, // PT
	// // 		extractAddressFromData44To64, // YT
	// // 	}},

	// firewallClient := clients.Firewall{
	// 	ProtocolName: "Pendle",
	// 	Accounts: []clients.AccountFW{
	// 		{
	// 			Address: addressOrPanic("0x27b1dAcd74688aF24a64BD3C9C1B143118740784"),
	// 			RegisteredEvents: []clients.EventFW{
	// 				{
	// 					Topic0:              hashOrPanic("0x166ae5f55615b65bbd9a2496e98d4e4d78ca15bd6127c0fe2dc27b76f6c03143"),
	// 					AddressLocationType: clients.AddressLocationTypeTopic,
	// 					AddressLocation:     1,
	// 				},
	// 			},
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 			Name:               "Market Factory V1",
	// 		},
	// 		{
	// 			Address: addressOrPanic("0x1A6fCc85557BC4fB7B534ed835a03EF056552D52"),
	// 			RegisteredEvents: []clients.EventFW{
	// 				{
	// 					Topic0:              hashOrPanic("0xae811fae25e2770b6bd1dcb1475657e8c3a976f91d1ebf081271db08eef920af"),
	// 					AddressLocationType: clients.AddressLocationTypeTopic,
	// 					AddressLocation:     1,
	// 				},
	// 			},
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 			Name:               "Market Factory V3",
	// 		},
	// 		{
	// 			Address: addressOrPanic("0x70ee0A6DB4F5a2Dc4d9c0b57bE97B9987e75BAFD"),
	// 			RegisteredEvents: []clients.EventFW{
	// 				{
	// 					Topic0:              hashOrPanic("0xaa79d8f17776adeaa316c5411b72e8b0057d064974fa8748f32492ecaa22ecd1"),
	// 					AddressLocationType: clients.AddressLocationTypeTopic,
	// 					AddressLocation:     1,
	// 				},
	// 				{
	// 					Topic0:              hashOrPanic("0xaa79d8f17776adeaa316c5411b72e8b0057d064974fa8748f32492ecaa22ecd1"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     12,
	// 				},
	// 				{
	// 					Topic0:              hashOrPanic("0xaa79d8f17776adeaa316c5411b72e8b0057d064974fa8748f32492ecaa22ecd1"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     44,
	// 				},
	// 			},
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 			Name:               "Yield Contract V1",
	// 		},
	// 		{
	// 			Address: addressOrPanic("0xdF3601014686674e53d1Fa52F7602525483F9122"),
	// 			RegisteredEvents: []clients.EventFW{
	// 				{
	// 					Topic0:              hashOrPanic("0xaa79d8f17776adeaa316c5411b72e8b0057d064974fa8748f32492ecaa22ecd1"),
	// 					AddressLocationType: clients.AddressLocationTypeTopic,
	// 					AddressLocation:     1,
	// 				},
	// 				{
	// 					Topic0:              hashOrPanic("0xaa79d8f17776adeaa316c5411b72e8b0057d064974fa8748f32492ecaa22ecd1"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     12,
	// 				},
	// 				{
	// 					Topic0:              hashOrPanic("0xaa79d8f17776adeaa316c5411b72e8b0057d064974fa8748f32492ecaa22ecd1"),
	// 					AddressLocationType: clients.AddressLocationTypeData,
	// 					AddressLocation:     44,
	// 				},
	// 			},
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 			Name:               "Yield Contract V3",
	// 		},
	// 		{
	// 			Name:               "vePendle",
	// 			Address:            addressOrPanic("0x4f30A9D41B80ecC5B94306AB4364951AE3170210"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "router",
	// 			Address:            addressOrPanic("0x888888888889758F76e7103c6CbF23ABbF58F946"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 	},
	// }

	// // 	"0x1BFFaBc6dFcAfB4177046db6686e3F135E8Bc732", // aveQI
	// // "0x8549ba7f483afb13b8321830d6f07f30f0a2f1de", // reward distributor
	// // "0x3FEACf904b152b1880bDE8BF04aC9Eb636fEE4d8", // Main safe
	// // "0x3182E6856c3B59C39114416075770Ec9DC9Ff436", // QiDao Guardians (UMH)
	// // "0x594F17028522BF85e830b689973682967E0DbcBc", // Revenue Managers
	// // "0x9d3c8a651e48e4D89ca5D1553035A4BE3c17cFe6", // Token Managers
	// // "0x60d133c666919B54a3254E0d3F14332cB783B733", // YLEMVT Vault
	// // "0xEcbd32bD581e241739be1763DFE7a8fFcC844ae1", // YEEMVT Vault
	// // "0x98eb27E5F24FB83b7D129D789665b08C258b4cCF", // WEMVT Vault
	// // "0x8C45969aD19D297c9B85763e90D0344C6E2ac9d1", // WBMVT Vault
	// // "0xcc61Ee649A95F2E2f0830838681f839BDb7CB823", // SCSEMVT Vault
	// // "0x82E90EB7034C1DF646bD06aFb9E67281AAb5ed28", // YCSEMVT Vault
	// // "0x67411793c5dcf9abc5a8d113ddd0e596cd5ba3e7", // StakeDAO Curve stETH
	// // "0xD1a6F422ceFf5a39b764e340Fd1bCd46C0744F83", // Yearn Curve stETH
	// // "0x86f78d3cbCa0636817AD9e27a44996C738Ec4932", // Beefy Convex Curve stETH
	// // "0xCA3EB45FB186Ed4e75B9B22A514fF1d4abAdD123", // CRV Vault
	// // "0x4ce4c542d96ce1872fea4fa3fbb2e7ae31862bad", // cbETH Vault
	// // "0x5773e8953cf60f495eb3c2db45dd753b5c4b7473", // stETH Vault
	// // "0x954ac12c339c60eafbb32213b15af3f7c7a0dec2", // LDO Vault
	// // "0xEd8a2759B0f8ea0f33225C86cB726fa9C6E030A4", // Performance Fee Management
	// // "0xf2833F5E72207D1Da1EEE7F8395Fb5f49895BBb4", // Stake DAO ETH Strategy Performance Fee Tokens
	// // "0xE9D954a9A6A1a61bc1120970f84CDd76562c4a0c", // Yearn Curve stETH Performance Fee Tokens
	// // "0x3c82A9514327A93928108e9F00D89877F4beB6e3", // Beefy Convex Curve stETH Performance Fee Tokens
	// // "0x97451025De0beef64c1A454bcF995de6FB8e0f2A", // cbETH Performance Fee Tokens
	// // "0x9414e766E8B59473599b9968aAf52CDCd07f59a9", // stETH Performance Fee Tokens

	// firewallClient := clients.Firewall{
	// 	ProtocolName: "QiDAO",
	// 	Accounts: []clients.AccountFW{
	// 		{
	// 			Name:               "aveQI",
	// 			Address:            addressOrPanic("0x1BFFaBc6dFcAfB4177046db6686e3F135E8Bc732"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Reward Distributor",
	// 			Address:            addressOrPanic("0x8549ba7f483afb13b8321830d6f07f30f0a2f1de"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Main Safe",
	// 			Address:            addressOrPanic("0x3FEACf904b152b1880bDE8BF04aC9Eb636fEE4d8"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "QiDao Guardians (UMH)",
	// 			Address:            addressOrPanic("0x3182E6856c3B59C39114416075770Ec9DC9Ff436"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Revenue Managers",
	// 			Address:            addressOrPanic("0x594F17028522BF85e830b689973682967E0DbcBc"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Token Managers",
	// 			Address:            addressOrPanic("0x9d3c8a651e48e4D89ca5D1553035A4BE3c17cFe6"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "YLEMVT Vault",
	// 			Address:            addressOrPanic("0x60d133c666919B54a3254E0d3F14332cB783B733"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "YEEMVT Vault",
	// 			Address:            addressOrPanic("0xEcbd32bD581e241739be1763DFE7a8fFcC844ae1"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "WEMVT Vault",
	// 			Address:            addressOrPanic("0x98eb27E5F24FB83b7D129D789665b08C258b4cCF"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "WBMVT Vault",
	// 			Address:            addressOrPanic("0x8C45969aD19D297c9B85763e90D0344C6E2ac9d1"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "SCSEMVT Vault",
	// 			Address:            addressOrPanic("0xcc61Ee649A95F2E2f0830838681f839BDb7CB823"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "YCSEMVT Vault",
	// 			Address:            addressOrPanic("0x82E90EB7034C1DF646bD06aFb9E67281AAb5ed28"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "StakeDAO Curve stETH",
	// 			Address:            addressOrPanic("0x67411793c5dcf9abc5a8d113ddd0e596cd5ba3e7"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Yearn Curve stETH",
	// 			Address:            addressOrPanic("0xD1a6F422ceFf5a39b764e340Fd1bCd46C0744F83"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Beefy Convex Curve stETH",
	// 			Address:            addressOrPanic("0x86f78d3cbCa0636817AD9e27a44996C738Ec4932"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "CRV Vault",
	// 			Address:            addressOrPanic("0xCA3EB45FB186Ed4e75B9B22A514fF1d4abAdD123"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "cbETH Vault",
	// 			Address:            addressOrPanic("0x4ce4c542d96ce1872fea4fa3fbb2e7ae31862bad"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "stETH Vault",
	// 			Address:            addressOrPanic("0x5773e8953cf60f495eb3c2db45dd753b5c4b7473"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "LDO Vault",
	// 			Address:            addressOrPanic("0x954ac12c339c60eafbb32213b15af3f7c7a0dec2"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Performance Fee Management",
	// 			Address:            addressOrPanic("0xEd8a2759B0f8ea0f33225C86cB726fa9C6E030A4"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Stake DAO ETH Strategy Performance Fee Tokens",
	// 			Address:            addressOrPanic("0xf2833F5E72207D1Da1EEE7F8395Fb5f49895BBb4"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Yearn Curve stETH Performance Fee Tokens",
	// 			Address:            addressOrPanic("0xE9D954a9A6A1a61bc1120970f84CDd76562c4a0c"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Beefy Convex Curve stETH Performance Fee Tokens",
	// 			Address:            addressOrPanic("0x3c82A9514327A93928108e9F00D89877F4beB6e3"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "cbETH Performance Fee Tokens",
	// 			Address:            addressOrPanic("0x97451025De0beef64c1A454bcF995de6FB8e0f2A"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "stETH Performance Fee Tokens",
	// 			Address:            addressOrPanic("0x9414e766E8B59473599b9968aAf52CDCd07f59a9"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 	},
	// }

	// // "0x3b8F6D6970a24A58b52374C539297ae02A3c4Ae4", //Mainnet Exchange
	// // 	"0x0c378fb17e87b180256a87e3f671cd83bf3236db", //Locked RBX

	// firewallClient := clients.Firewall{
	// 	ProtocolName: "RabbitX",
	// 	Accounts: []clients.AccountFW{
	// 		{
	// 			Name:               "Mainnet Exchange",
	// 			Address:            addressOrPanic("0x3b8F6D6970a24A58b52374C539297ae02A3c4Ae4"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Locked RBX",
	// 			Address:            addressOrPanic("0x0c378fb17e87b180256a87e3f671cd83bf3236db"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 	},
	// }

	// "0x4c406C068106375724275Cbff028770C544a1333", // Emerald (scETH)
	// 	"0x096697720056886b905D0DEB0f06AfFB8e4665E5", // Opal (scUSDC)
	// 	"0xdb369eEB33fcfDCd1557E354dDeE7d6cF3146A11", // Amber (scLUSD)
	// 	"0x0a36f9565c6fb862509ad8d148941968344a55d8", // Rewards Tracker

	// firewallClient := clients.Firewall{
	// 	ProtocolName: "Sandclock",
	// 	Accounts: []clients.AccountFW{
	// 		{
	// 			Name:               "Emerald (scETH)",
	// 			Address:            addressOrPanic("0x4c406C068106375724275Cbff028770C544a1333"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Opal (scUSDC)",
	// 			Address:            addressOrPanic("0x096697720056886b905D0DEB0f06AfFB8e4665E5"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Amber (scLUSD)",
	// 			Address:            addressOrPanic("0xdb369eEB33fcfDCd1557E354dDeE7d6cF3146A11"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 		{
	// 			Name:               "Rewards Tracker",
	// 			Address:            addressOrPanic("0x0a36f9565c6fb862509ad8d148941968344a55d8"),
	// 			ChildContractScope: clients.ChildContractScopeNone,
	// 			ChainIDs:           []int64{1},
	// 		},
	// 	},
	// }

	// "0xB878DC600550367e14220d4916Ff678fB284214F", // Factory
	// 		[]string{
	// 			"0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9", // Pair Created
	// 			"0xeda679f3434de965730a28b8b694f2a348c09a2c1bb3e226633b6af24841adc1", // Pair Added
	// 		},
	// 		[]addressExtractor{
	// 			extractAddressFromData12To32, // pair
	// 		},

	firewallClient := clients.Firewall{
		ProtocolName: "Smardex",
		Accounts: []clients.AccountFW{
			{
				Address: addressOrPanic("0xB878DC600550367e14220d4916Ff678fB284214F"),
				RegisteredEvents: []clients.EventFW{
					{
						Topic0:              hashOrPanic("0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9"),
						AddressLocationType: clients.AddressLocationTypeTopic,
						AddressLocation:     12,
					},
					{
						Topic0:              hashOrPanic("0xeda679f3434de965730a28b8b694f2a348c09a2c1bb3e226633b6af24841adc1"),
						AddressLocationType: clients.AddressLocationTypeData,
						AddressLocation:     12,
					},
				},
				ChildContractScope: clients.ChildContractScopeNone,
				ChainIDs:           []int64{1},
				Name:               "Factory",
			},
		},
	}

	protocolInfo, err := clients.GetProtocolInfo(firewallClient.ProtocolName)
	if err != nil {
		panic(err)
	}

	// Generate Monitored data
	monitored, err := firewallClient.ToMonitored(rpcClient, apiKey)
	if err != nil {
		panic(err)
	}

	// Write protocol info to Firestore
	protocolRef, err := firestoreClients.WriteProtocolInfoToFirestore(firestoreClient, *protocolInfo)
	if err != nil {
		panic(err)
	}

	// Write agreement details to Firestore
	firewallRef, err := firestoreClients.WriteFirewallToFirestore(firestoreClient, protocolInfo.Name, firewallClient, protocolRef)
	if err != nil {
		panic(err)
	}

	// Write Monitored data to Firestore
	_, err = firestoreClients.WriteMonitoredToFirestore(firestoreClient, protocolInfo.Name, *monitored, protocolRef, nil, firewallRef)
	if err != nil {
		panic(err)
	}

	// chainID := int64(1)
	// apiKey := "PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3"
	// address, err := web3.HexToAddress("0xdf3601014686674e53d1fa52f7602525483f9122")
	// if err != nil {
	// 	panic(err)
	// }
	// startBlock := 19017136
	// _, err = etherscan.FetchRegularTransactions(chainID, apiKey, *address, startBlock)
	// if err != nil {
	// 	panic(err)
	// }
}

// import (
// 	"SHDB/pkg/clients"
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"cloud.google.com/go/firestore"
// 	firebase "firebase.google.com/go"
// 	"github.com/ethereum/go-ethereum/rpc"
// 	"google.golang.org/api/option"
// )

// const (
// 	apiKey           = "PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3"
// 	rpcEndpoint      = "https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/"
// 	checkInterval    = 1 * time.Hour
// 	firebaseCredFile = "/home/dwu/SHDB/skylock-xyz-firebase-adminsdk-36s2d-bd6e795bf3.json"
// )

// func main() {
// 	ctx := context.Background()
// 	rpcClient, err := rpc.Dial("https://aged-tiniest-uranium.ethereum-sepolia.quiknode.pro/c3b64aa1f9a4efb9642aba41db74bb288466a7c6/")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// details, err := clients.GetSafeHarborAdoptions("PVZZQKUUE28APGKU9T41FD45P8BQ8M4JV3", rpcClient)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// fmt.Println(details)
// 	sa := option.WithCredentialsFile(firebaseCredFile)
// 	app, err := firebase.NewApp(ctx, nil, sa)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize Firebase app: %v", err)
// 	}

// 	firestoreClient, err := app.Firestore(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize Firestore client: %v", err)
// 	}

// 	adoptions, err := clients.GetSafeHarborAdoptions(apiKey, rpcClient)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, adoption := range adoptions {
// 		fmt.Println(adoption)
// 		// Check if this adoption already exists in Firestore
// 		exists, err := checkAdoptionExists(ctx, firestoreClient, adoption.ProtocolName)
// 		if err != nil {
// 			log.Printf("Error checking adoption existence: %v", err)
// 			continue
// 		}

// 		fmt.Println(exists)

// 		// if !exists {
// 		// 	err = processNewAdoption(ctx, rpcClient, firestoreClient, adoption)
// 		// 	if err != nil {
// 		// 		log.Printf("Error processing new adoption for %s: %v", adoption.ProtocolName, err)
// 		// 	}
// 		// }
// 	}
// }

// func checkAdoptionExists(ctx context.Context, client *firestore.Client, protocolName string) (bool, error) {
// 	docs, err := client.Collection("safeHarborAgreements").Where("protocol.name", "==", protocolName).Documents(ctx).GetAll()
// 	if err != nil {
// 		return false, err
// 	}
// 	return len(docs) > 0, nil
// }
