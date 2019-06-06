// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2018 The Flo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"errors"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/bitspill/flod/chaincfg/chainhash"
	"github.com/bitspill/flod/wire"
)

// These variables are the chain proof-of-work limit parameters for each default
// network.
var (
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// mainPowLimit is the highest proof of work value a Bitcoin block can
	// have for the main network.  It is the value 2^224 - 1.
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 236), bigOne)

	// regressionPowLimit is the highest proof of work value a Bitcoin block
	// can have for the regression test network.  It is the value 2^255 - 1.
	regressionPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)

	// testNet3PowLimit is the highest proof of work value a Bitcoin block
	// can have for the test network (version 3).  It is the value
	// 2^236 - 1.
	testNet3PowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 236), bigOne)

	// simNetPowLimit is the highest proof of work value a Bitcoin block
	// can have for the simulation test network.  It is the value 2^255 - 1.
	simNetPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)
)

// Checkpoint identifies a known good point in the block chain.  Using
// checkpoints allows a few optimizations for old blocks during initial download
// and also prevents forks from old blocks.
//
// Each checkpoint is selected based upon several factors.  See the
// documentation for blockchain.IsCheckpointCandidate for details on the
// selection criteria.
type Checkpoint struct {
	Height int32
	Hash   *chainhash.Hash
}

// DNSSeed identifies a DNS seed.
type DNSSeed struct {
	// Host defines the hostname of the seed.
	Host string

	// HasFiltering defines whether the seed supports filtering
	// by service flags (wire.ServiceFlag).
	HasFiltering bool
}

// ConsensusDeployment defines details related to a specific consensus rule
// change that is voted in.  This is part of BIP0009.
type ConsensusDeployment struct {
	// BitNumber defines the specific bit number within the block version
	// this particular soft-fork deployment refers to.
	BitNumber uint8

	// StartTime is the median block time after which voting on the
	// deployment starts.
	StartTime uint64

	// ExpireTime is the median block time after which the attempted
	// deployment expires.
	ExpireTime uint64
}

// Constants that define the deployment offset in the deployments field of the
// parameters for each deployment.  This is useful to be able to get the details
// of a specific deployment by name.
const (
	// DeploymentTestDummy defines the rule change deployment ID for testing
	// purposes.
	DeploymentTestDummy = iota

	// DeploymentCSV defines the rule change deployment ID for the CSV
	// soft-fork package. The CSV package includes the deployment of BIPS
	// 68, 112, and 113.
	DeploymentCSV

	// DeploymentSegwit defines the rule change deployment ID for the
	// Segregated Witness (segwit) soft-fork package. The segwit package
	// includes the deployment of BIPS 141, 142, 144, 145, 147 and 173.
	DeploymentSegwit

	// NOTE: DefinedDeployments must always come last since it is used to
	// determine how many defined deployments there currently are.

	// DefinedDeployments is the number of currently defined deployments.
	DefinedDeployments
)

// Params defines a Bitcoin network by its parameters.  These parameters may be
// used by Bitcoin applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name string

	// Net defines the magic bytes used to identify the network.
	Net wire.BitcoinNet

	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string

	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds []DNSSeed

	// GenesisBlock defines the first block of the chain.
	GenesisBlock *wire.MsgBlock

	// GenesisHash is the starting block hash.
	GenesisHash *chainhash.Hash

	// PowLimit defines the highest allowed proof of work value for a block
	// as a uint256.
	PowLimit *big.Int

	// PowLimitBits defines the highest allowed proof of work value for a
	// block in compact form.
	PowLimitBits uint32

	// These fields define the block heights at which the specified softfork
	// BIP became active.
	BIP0034Height int32
	BIP0065Height int32
	BIP0066Height int32

	// CoinbaseMaturity is the number of blocks required before newly mined
	// coins (coinbase transactions) can be spent.
	CoinbaseMaturity uint16

	// SubsidyReductionInterval is the interval of blocks before the subsidy
	// is reduced.
	SubsidyReductionInterval int32

	// TargetTimespan is the desired amount of time that should elapse
	// before the block difficulty requirement is examined to determine how
	// it should be changed in order to maintain the desired block
	// generation rate.
	TargetTimespan time.Duration

	// TargetTimePerBlock is the desired amount of time to generate each
	// block.
	TargetTimePerBlock time.Duration

	// RetargetAdjustmentFactor is the adjustment factor used to limit
	// the minimum and maximum amount of adjustment that can occur between
	// difficulty retargets.
	RetargetAdjustmentFactorUp   int64
	RetargetAdjustmentFactorDown int64

	// ReduceMinDifficulty defines whether the network should reduce the
	// minimum required difficulty after a long enough period of time has
	// passed without finding a block.  This is really only useful for test
	// networks and should not be set on a main network.
	ReduceMinDifficulty bool

	// MinDiffReductionTime is the amount of time after which the minimum
	// required difficulty should be reduced when a block hasn't been found.
	//
	// NOTE: This only applies if ReduceMinDifficulty is true.
	MinDiffReductionTime time.Duration

	// GenerateSupported specifies whether or not CPU mining is allowed.
	GenerateSupported bool
	
	// PowNoRetargeting specifies wether the block retargeting calculation is performed.
	PowNoRetargeting bool

	// Checkpoints ordered from oldest to newest.
	Checkpoints []Checkpoint

	// These fields are related to voting on consensus rule changes as
	// defined by BIP0009.
	//
	// RuleChangeActivationThreshold is the number of blocks in a threshold
	// state retarget window for which a positive vote for a rule change
	// must be cast in order to lock in a rule change. It should typically
	// be 95% for the main network and 75% for test networks.
	//
	// MinerConfirmationWindow is the number of blocks in each threshold
	// state retarget window.
	//
	// Deployments define the specific consensus rule changes to be voted
	// on.
	RuleChangeActivationThreshold uint32
	MinerConfirmationWindow       uint32
	Deployments                   [DefinedDeployments]ConsensusDeployment

	// Mempool parameters
	RelayNonStdTxs bool

	// Human-readable part for Bech32 encoded segwit addresses, as defined
	// in BIP 173.
	Bech32HRPSegwit string

	// Address encoding magics
	PubKeyHashAddrID        byte // First byte of a P2PKH address
	ScriptHashAddrID        byte // First byte of a P2SH address
	PrivateKeyID            byte // First byte of a WIF private key
	WitnessPubKeyHashAddrID byte // First byte of a P2WPKH address
	WitnessScriptHashAddrID byte // First byte of a P2WSH address

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID [4]byte
	HDPublicKeyID  [4]byte

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType uint32
}

// MainNetParams defines the network parameters for the main Bitcoin network.
var MainNetParams = Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "7312",
	DNSSeeds: []DNSSeed{
		{"flodns.oip.fun", false},
		{"flodns.oip.li", false},
		{"node.oip.fun", false},
	},

	// Chain parameters
	GenesisBlock:                 &genesisBlock,
	GenesisHash:                  &genesisHash,
	PowLimit:                     mainPowLimit,
	PowLimitBits:                 0x1d00ffff,
	BIP0034Height:                1679161, // 490a10507efe42b89104408787088b7c43310cc230310201b5f57dac6f513b8b
	BIP0065Height:                1679161, // 490a10507efe42b89104408787088b7c43310cc230310201b5f57dac6f513b8b
	BIP0066Height:                1679161, // 490a10507efe42b89104408787088b7c43310cc230310201b5f57dac6f513b8b
	CoinbaseMaturity:             100,
	SubsidyReductionInterval:     800000,
	TargetTimespan:               time.Second * 40 * 6, // 40 seconds * 6 blocks
	TargetTimePerBlock:           time.Second * 40,     // 40 seconds
	RetargetAdjustmentFactorUp:   2,
	RetargetAdjustmentFactorDown: 3,
	ReduceMinDifficulty:          false,
	MinDiffReductionTime:         0,
	GenerateSupported:            false,
	PowNoRetargeting:             false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		{8002, newHashFromStr("73bc3b16d99bbf797f396c9532f80c3b73bb21304280de2efbc5edcb75739234")},
		{18001, newHashFromStr("5a7a4821aa4fc7ee3dea2f8319e9fa4d991a8c6762e79cb624c64e4cf1031582")},
		{38002, newHashFromStr("4962437c6d0a450f44c1e40cd38ff220f8122af1517e1329f1abd07fb7791e40")},
		{160002, newHashFromStr("478d381c92298614c3a05fb934a4fffc4d3e5b573efbba9b3e8b2ce8d26a0f8f")},
		{208001, newHashFromStr("2bb3f8b2d5081aefa0af9f5d8de42bd73a5d89eebf78aa7421cd63dc40a56d4c")},
		{270001, newHashFromStr("74988a3179ae6bbc5986e63f71bafc855202502b07e4d9331015eee82df80860")},
		{290036, newHashFromStr("145994381e5e4f0e5674adc1ace9a03b670838792f6bd6b650c80466453c2da3")},
		{344665, newHashFromStr("40fe36d8dec357aa529b6b1d99b2989a37ed8c7b065a0e3345cd15a751b9c1ad")},
		{400236, newHashFromStr("f9a4b8e21d410539e45ff3f11c28dee8966de7edffc45fd02dd1a5f4e7d4ef38")},
		// { 415000, newHashFromStr("16ef8ab98a7300039a5755d5bdc00e31dada9d2f1c440ff7928f43c4ea41c0a8")},
		{420937, newHashFromStr("48a75e4687021ec0dda2031439de50b61933e197a4e1a1185d131cc2b59b8444")},
		{425606, newHashFromStr("62c8d811b1a49f6fdaffded704dc48b1c98d6f8dd736d8afb96c9b097774a85e")},
		{508694, newHashFromStr("65cde197e9118e5164c4dcdcdc6fcfaf8c0de605d569cefd56aa220e7739da6a")},
		{696454, newHashFromStr("8cfb75684405e22f8f69522ec11f1e5206758e37f25db13880548f69fe6f1976")},
		{955000, newHashFromStr("b5517a50aee6af59eb0ab4ee3262bcbaf3f6672b9301cdd3302e4bab491e7526")},
		{1505017, newHashFromStr("d38b306850bb26a5c98400df747d4391bb4e359e95e20dc79b50063ed3c5bfa7")},
		{1678879, newHashFromStr("1e874e2852e8dfb3553f0e6c02dcf70e9f5697effa31385d25a3c88fe26676fc")},
		{1678909, newHashFromStr("4c5a1040e337a542e6717904c8346bd72151fc34c390dff7b5cf23dcedc5058a")},
		{1679162, newHashFromStr("b32c64fb80a4196ff3e1be883db10629e1d7cd27c00ef0b5d1fe54af481fc10f")},
		{1796633, newHashFromStr("c2da8b936a7f2c0de02aa0c6c45f3d971ebad78655255a945b0e60b62f27d445")},
		{2094558, newHashFromStr("946616c88286f32bfac15868456d87a86f8611e1f9b56594b81e46831ce43f81")},
		{2532181, newHashFromStr("cacd5149aaed1088ae1db997a741210b0525e941356104120f182f3159931c79")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 6048, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       8064, //
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1522562766, // May 1st, 2016
			ExpireTime: 1554098766, // May 1st, 2017
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1522562766, // November 15, 2016 UTC
			ExpireTime: 1554098766, // November 15, 2017 UTC.
		},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "flo", // always bc for main net

	// Address encoding magics
	// ToDo: bitspill parameters
	PubKeyHashAddrID: 35, // starts with 1
	//ScriptHashAddrID:        8, // starts with 3
	ScriptHashAddrID: 94,  // starts with 3
	PrivateKeyID:     163, // starts with 5 (uncompressed) or K (compressed)
	//PrivateKeyID:            176, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0,
}

// RegressionNetParams defines the network parameters for the regression test
// Bitcoin network.  Not to be confused with the test Bitcoin network (version
// 3), this network is sometimes simply called "testnet".
var RegressionNetParams = Params{
	Name:        "regtest",
	Net:         wire.TestNet,
	DefaultPort: "17412",
	DNSSeeds:    []DNSSeed{},

	// Chain parameters
	GenesisBlock:                 &regTestGenesisBlock,
	GenesisHash:                  &regTestGenesisHash,
	PowLimit:                     regressionPowLimit,
	PowLimitBits:                 0x207fffff,
	CoinbaseMaturity:             100,
	BIP0034Height:                100000000, // Not active - Permit ver 1 blocks
	BIP0065Height:                1351,      // Used by regression tests
	BIP0066Height:                1251,      // Used by regression tests
	SubsidyReductionInterval:     150,
	TargetTimespan:               time.Second * 40 * 6, // 40 seconds * 6 blocks
	TargetTimePerBlock:           time.Second * 40,     // 40 seconds
	RetargetAdjustmentFactorUp:   2,                   // 25% less, 400% more
	RetargetAdjustmentFactorDown: 3,
	ReduceMinDifficulty:          true,
	MinDiffReductionTime:         time.Second * 80, // TargetTimePerBlock * 2
	GenerateSupported:            true,
	PowNoRetargeting:             true,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 108, // 75%  of MinerConfirmationWindow
	MinerConfirmationWindow:       144, // Faster than normal for regtest (144 instead of 2016)
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires.
		},
	},

	// Mempool parameters
	RelayNonStdTxs: true,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "rflo", // always bcrt for reg test net

	// Address encoding magics
	PubKeyHashAddrID: 115,
	ScriptHashAddrID: 198,
	PrivateKeyID:     239, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,
}

// TestNet3Params defines the network parameters for the test Bitcoin network
// (version 3).  Not to be confused with the regression test network, this
// network is sometimes simply called "testnet".
var TestNet3Params = Params{
	Name:        "testnet3",
	Net:         wire.TestNet3,
	DefaultPort: "17312",
	DNSSeeds: []DNSSeed{
		{"testnet.oip.fun", false},
	},

	// Chain parameters
	GenesisBlock:                 &testNet3GenesisBlock,
	GenesisHash:                  &testNet3GenesisHash,
	PowLimit:                     testNet3PowLimit,
	PowLimitBits:                 0x1e0fffff,
	BIP0034Height:                33600, // 4ac31d938531317c065405a9b23478c8c99204ff17fc294cb09821e2c2b42e65
	BIP0065Height:                33600, // 4ac31d938531317c065405a9b23478c8c99204ff17fc294cb09821e2c2b42e65
	BIP0066Height:                33600, // 4ac31d938531317c065405a9b23478c8c99204ff17fc294cb09821e2c2b42e65
	CoinbaseMaturity:             100,
	SubsidyReductionInterval:     800000,
	TargetTimespan:               time.Second * 40 * 6, // 40 seconds * 6 blocks
	TargetTimePerBlock:           time.Second * 40,     // 40 seconds
	RetargetAdjustmentFactorUp:   2,
	RetargetAdjustmentFactorDown: 3,
	ReduceMinDifficulty:          true,
	MinDiffReductionTime:         time.Second * 80, // TargetTimePerBlock * 2
	GenerateSupported:            false,
	PowNoRetargeting:             false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		{2056, newHashFromStr("d3334db071731beaa651f10624c2fea1a5e8c6f9e50f0e602f86262938374148")},
		{230000, newHashFromStr("c2e6451240a580c3bfa5ddbfad1b001f8655e7c51d5c32c123e16f69c2d2b539")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 600, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       800,
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1483228800, // January 1, 2017
			ExpireTime: 1530446401, // July 1, 2018
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1483228800, // January 1, 2017
			ExpireTime: 1530446401, // July 1, 2018
		},
	},

	// Mempool parameters
	RelayNonStdTxs: true,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "tb", // always tb for test net

	// Address encoding magics
	PubKeyHashAddrID: 115, // starts with m or n
	ScriptHashAddrID: 198, // starts with 2
	// 58
	WitnessPubKeyHashAddrID: 0x03, // starts with QW
	WitnessScriptHashAddrID: 0x28, // starts with T7n
	PrivateKeyID:            239,  // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x01, 0x34, 0x40, 0xe2}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x01, 0x34, 0x3c, 0x23}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,
}

// SimNetParams defines the network parameters for the simulation test Bitcoin
// network.  This network is similar to the normal test network except it is
// intended for private use within a group of individuals doing simulation
// testing.  The functionality is intended to differ in that the only nodes
// which are specifically specified are used to create the network rather than
// following normal discovery rules.  This is important as otherwise it would
// just turn into another public testnet.
var SimNetParams = Params{
	Name:        "simnet",
	Net:         wire.SimNet,
	DefaultPort: "18555",
	DNSSeeds:    []DNSSeed{}, // NOTE: There must NOT be any seeds.

	// Chain parameters
	GenesisBlock:                 &simNetGenesisBlock,
	GenesisHash:                  &simNetGenesisHash,
	PowLimit:                     simNetPowLimit,
	PowLimitBits:                 0x207fffff,
	BIP0034Height:                0, // Always active on simnet
	BIP0065Height:                0, // Always active on simnet
	BIP0066Height:                0, // Always active on simnet
	CoinbaseMaturity:             100,
	SubsidyReductionInterval:     210000,
	TargetTimespan:               time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:           time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactorUp:   4,                   // 25% less, 400% more
	RetargetAdjustmentFactorDown: 4,
	ReduceMinDifficulty:          true,
	MinDiffReductionTime:         time.Minute * 20, // TargetTimePerBlock * 2
	GenerateSupported:            true,
	PowNoRetargeting:             true,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 75, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       100,
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires.
		},
	},

	// Mempool parameters
	RelayNonStdTxs: true,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "sb", // always sb for sim net

	// Address encoding magics
	PubKeyHashAddrID:        0x3f, // starts with S
	ScriptHashAddrID:        0x7b, // starts with s
	PrivateKeyID:            0x64, // starts with 4 (uncompressed) or F (compressed)
	WitnessPubKeyHashAddrID: 0x19, // starts with Gg
	WitnessScriptHashAddrID: 0x28, // starts with ?

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x20, 0xb9, 0x00}, // starts with sprv
	HDPublicKeyID:  [4]byte{0x04, 0x20, 0xbd, 0x3a}, // starts with spub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 115, // ASCII for s
}

var (
	// ErrDuplicateNet describes an error where the parameters for a Bitcoin
	// network could not be set due to the network already being a standard
	// network or previously-registered into this package.
	ErrDuplicateNet = errors.New("duplicate Bitcoin network")

	// ErrUnknownHDKeyID describes an error where the provided id which
	// is intended to identify the network for a hierarchical deterministic
	// private extended key is not registered.
	ErrUnknownHDKeyID = errors.New("unknown hd private extended key bytes")
)

var (
	registeredNets       = make(map[wire.BitcoinNet]struct{})
	pubKeyHashAddrIDs    = make(map[byte]struct{})
	scriptHashAddrIDs    = make(map[byte]struct{})
	bech32SegwitPrefixes = make(map[string]struct{})
	hdPrivToPubKeyIDs    = make(map[[4]byte][]byte)
)

// String returns the hostname of the DNS seed in human-readable form.
func (d DNSSeed) String() string {
	return d.Host
}

// Register registers the network parameters for a Bitcoin network.  This may
// error with ErrDuplicateNet if the network is already registered (either
// due to a previous Register call, or the network being one of the default
// networks).
//
// Network parameters should be registered into this package by a main package
// as early as possible.  Then, library packages may lookup networks or network
// parameters based on inputs and work regardless of the network being standard
// or not.
func Register(params *Params) error {
	if _, ok := registeredNets[params.Net]; ok {
		return ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	pubKeyHashAddrIDs[params.PubKeyHashAddrID] = struct{}{}
	scriptHashAddrIDs[params.ScriptHashAddrID] = struct{}{}
	hdPrivToPubKeyIDs[params.HDPrivateKeyID] = params.HDPublicKeyID[:]

	// A valid Bech32 encoded segwit address always has as prefix the
	// human-readable part for the given net followed by '1'.
	bech32SegwitPrefixes[params.Bech32HRPSegwit+"1"] = struct{}{}
	return nil
}

// mustRegister performs the same function as Register except it panics if there
// is an error.  This should only be called from package init functions.
func mustRegister(params *Params) {
	if err := Register(params); err != nil {
		panic("failed to register network: " + err.Error())
	}
}

// IsPubKeyHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-pubkey-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsScriptHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsPubKeyHashAddrID(id byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

// IsScriptHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-script-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsPubKeyHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsScriptHashAddrID(id byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

// IsBech32SegwitPrefix returns whether the prefix is a known prefix for segwit
// addresses on any default or registered network.  This is used when decoding
// an address string into a specific address type.
func IsBech32SegwitPrefix(prefix string) bool {
	prefix = strings.ToLower(prefix)
	_, ok := bech32SegwitPrefixes[prefix]
	return ok
}

// HDPrivateKeyToPublicKeyID accepts a private hierarchical deterministic
// extended key id and returns the associated public key id.  When the provided
// id is not registered, the ErrUnknownHDKeyID error will be returned.
func HDPrivateKeyToPublicKeyID(id []byte) ([]byte, error) {
	if len(id) != 4 {
		return nil, ErrUnknownHDKeyID
	}

	var key [4]byte
	copy(key[:], id)
	pubBytes, ok := hdPrivToPubKeyIDs[key]
	if !ok {
		return nil, ErrUnknownHDKeyID
	}

	return pubBytes, nil
}

// newHashFromStr converts the passed big-endian hex string into a
// chainhash.Hash.  It only differs from the one available in chainhash in that
// it panics on an error since it will only (and must only) be called with
// hard-coded, and therefore known good, hashes.
func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		// Ordinarily I don't like panics in library code since it
		// can take applications down without them having a chance to
		// recover which is extremely annoying, however an exception is
		// being made in this case because the only way this can panic
		// is if there is an error in the hard-coded hashes.  Thus it
		// will only ever potentially panic on init and therefore is
		// 100% predictable.
		panic(err)
	}
	return hash
}

func init() {
	// Register all default networks when the package is initialized.
	mustRegister(&MainNetParams)
	mustRegister(&TestNet3Params)
	mustRegister(&RegressionNetParams)
	mustRegister(&SimNetParams)
}
