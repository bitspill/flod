package chaincfg

import (
	"time"

	"github.com/bitspill/flod/wire"
)

// MainNetParams defines the network parameters for the main Litecoin network.
var LtcMainNetParams = Params{
	Name:        "mainnet",
	Net:         wire.LtcMainNet,
	DefaultPort: "9333",
	DNSSeeds: []DNSSeed{
		{"seed-a.litecoin.loshan.co.uk", true},
		{"dnsseed.thrasher.io", true},
		{"dnsseed.litecointools.com", false},
		{"dnsseed.litecoinpool.org", false},
		{"dnsseed.koin-project.com", false},
	},

	// Chain parameters
	GenesisBlock:                 &genesisBlock,
	GenesisHash:                  &genesisHash,
	PowLimit:                     mainPowLimit,
	PowLimitBits:                 504365055,
	BIP0034Height:                710000,
	BIP0065Height:                918684,
	BIP0066Height:                811879,
	CoinbaseMaturity:             100,
	SubsidyReductionInterval:     840000,
	TargetTimespan:               (time.Hour * 24 * 3) + (time.Hour * 12), // 3.5 days
	TargetTimePerBlock:           (time.Minute * 2) + (time.Second * 30),  // 2.5 minutes
	RetargetAdjustmentFactorUp:   4,                                       // 25% less, 400% more
	RetargetAdjustmentFactorDown: 4,
	ReduceMinDifficulty:          false,
	MinDiffReductionTime:         0,
	GenerateSupported:            false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		{1500, newHashFromStr("841a2965955dd288cfa707a755d05a54e45f8bd476835ec9af4402a2b59a2967")},
		{4032, newHashFromStr("9ce90e427198fc0ef05e5905ce3503725b80e26afd35a987965fd7e3d9cf0846")},
		{8064, newHashFromStr("eb984353fc5190f210651f150c40b8a4bab9eeeff0b729fcb3987da694430d70")},
		{16128, newHashFromStr("602edf1859b7f9a6af809f1d9b0e6cb66fdc1d4d9dcd7a4bec03e12a1ccd153d")},
		{23420, newHashFromStr("d80fdf9ca81afd0bd2b2a90ac3a9fe547da58f2530ec874e978fce0b5101b507")},
		{50000, newHashFromStr("69dc37eb029b68f075a5012dcc0419c127672adb4f3a32882b2b3e71d07a20a6")},
		{80000, newHashFromStr("4fcb7c02f676a300503f49c764a89955a8f920b46a8cbecb4867182ecdb2e90a")},
		{120000, newHashFromStr("bd9d26924f05f6daa7f0155f32828ec89e8e29cee9e7121b026a7a3552ac6131")},
		{161500, newHashFromStr("dbe89880474f4bb4f75c227c77ba1cdc024991123b28b8418dbbf7798471ff43")},
		{179620, newHashFromStr("2ad9c65c990ac00426d18e446e0fd7be2ffa69e9a7dcb28358a50b2b78b9f709")},
		{240000, newHashFromStr("7140d1c4b4c2157ca217ee7636f24c9c73db39c4590c4e6eab2e3ea1555088aa")},
		{383640, newHashFromStr("2b6809f094a9215bafc65eb3f110a35127a34be94b7d0590a096c3f126c6f364")},
		{409004, newHashFromStr("487518d663d9f1fa08611d9395ad74d982b667fbdc0e77e9cf39b4f1355908a3")},
		{456000, newHashFromStr("bf34f71cc6366cd487930d06be22f897e34ca6a40501ac7d401be32456372004")},
		{638902, newHashFromStr("15238656e8ec63d28de29a8c75fcf3a5819afc953dcd9cc45cecc53baec74f38")},
		{721000, newHashFromStr("198a7b4de1df9478e2463bd99d75b714eab235a2e63e741641dc8a759a9840e5")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 6048, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       8064, //
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1485561600, // January 28, 2017 UTC
			ExpireTime: 1517356801, // January 31st, 2018 UTC
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1485561600, // January 28, 2017 UTC
			ExpireTime: 1517356801, // January 31st, 2018 UTC.
		},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "ltc", // always ltc for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x30, // starts with L
	ScriptHashAddrID:        0x32, // starts with M
	PrivateKeyID:            0xB0, // starts with 6 (uncompressed) or T (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 2,
}

func init() {
	mustRegister(&LtcMainNetParams)
}
