package chaincfg

import (
	"time"

	"github.com/bitspill/flod/chaincfg/chainhash"
	"github.com/bitspill/flod/wire"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 2,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x53, 0x6c, 0x61, 0x73, 0x68, 0x64, 0x6f, 0x74, /* | Slashdot | */
				0x20, 0x2d, 0x20, 0x31, 0x37, 0x20, 0x4a, 0x75, /* |  - 17 Ju | */
				0x6e, 0x65, 0x20, 0x32, 0x30, 0x31, 0x33, 0x20, /* | ne 2013  | */
				0x2d, 0x20, 0x53, 0x61, 0x75, 0x64, 0x69, 0x20, /* | - Saudi  | */
				0x41, 0x72, 0x61, 0x62, 0x69, 0x61, 0x20, 0x53, /* | Arabia S | */
				0x65, 0x74, 0x20, 0x54, 0x6f, 0x20, 0x42, 0x61, /* | et To Ba | */
				0x6e, 0x20, 0x57, 0x68, 0x61, 0x74, 0x73, 0x41, /* | n WhatsA | */
				0x70, 0x70, 0x2c, 0x20, 0x53, 0x6b, 0x79, 0x70, /* | pp, Skyp | */
				0x65, /* | e | */
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{
				0x41, 0x04, 0x01, 0x84, 0x71, 0x0f, 0xa6, 0x89,
				0xad, 0x50, 0x23, 0x69, 0x0c, 0x80, 0xf3, 0xa4,
				0x9c, 0x8f, 0x13, 0xf8, 0xd4, 0x5b, 0x8c, 0x85,
				0x7f, 0xbc, 0xbc, 0x8b, 0xc4, 0xa8, 0xe4, 0xd3,
				0xeb, 0x4b, 0x10, 0xf4, 0xd4, 0x60, 0x4f, 0xa0,
				0x8d, 0xce, 0x60, 0x1a, 0xaf, 0x0f, 0x47, 0x02,
				0x16, 0xfe, 0x1b, 0x51, 0x85, 0x0b, 0x4a, 0xcf,
				0x21, 0xb1, 0x79, 0xc4, 0x50, 0x70, 0xac, 0x7b,
				0x03, 0xa9, 0xac,
			},
		},
	},
	FloData:  []byte("text:Florincoin genesis block"),
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main
// network (genesis block).
var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xea, 0x1c, 0x3e, 0xff, 0x7c, 0x4b, 0x27, 0x67,
	0xba, 0x56, 0x1c, 0xdd, 0xfc, 0xec, 0xd7, 0x41,
	0x90, 0x5c, 0xea, 0x38, 0x5d, 0xc3, 0x78, 0xe2,
	0x08, 0x07, 0xf9, 0x9d, 0x1c, 0x78, 0xc7, 0x09,
})

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the main network.
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x00, 0x23, 0x2b, 0xab, 0xa7, 0x29, 0x1c, 0x9d,
	0x84, 0x4b, 0xb2, 0x48, 0x67, 0xaa, 0xfe, 0x45,
	0x3e, 0xa7, 0xe2, 0x90, 0x68, 0x56, 0x12, 0x55,
	0x2d, 0x59, 0x5a, 0xdc, 0x8d, 0x0c, 0x0f, 0x73,
})

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(0x51bf408c, 0), // 2009-01-03 18:15:05 +0000 UTC
		Bits:       0x1e0ffff0,               // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      0x3b9c81a4,               // 2083236893
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// regTestGenesisHash is the hash of the first block in the block chain for the
// regression test network (genesis block).
var regTestGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xd7, 0xe0, 0xde, 0xa8, 0xb9, 0xd5, 0x17, 0x89,
	0xdb, 0x36, 0x57, 0x8d, 0x6f, 0x56, 0xa4, 0x5e,
	0x93, 0x61, 0xb1, 0x24, 0x1d, 0x9a, 0xb5, 0x03,
	0x11, 0xcb, 0x6d, 0xca, 0x26, 0xfa, 0x42, 0xec,
})

// regTestGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the regression test network.  It is the same as the merkle root for
// the main network.
var regTestGenesisMerkleRoot = genesisMerkleRoot

// regTestGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the regression test network.
var regTestGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: regTestGenesisMerkleRoot, // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1371387277, 0), // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      0,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet3GenesisHash is the hash of the first block in the block chain for the
// test network (version 3).
var testNet3GenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x31, 0x37, 0xe0, 0x27, 0x8a, 0x19, 0x2a, 0xda,
	0xf0, 0xa1, 0xa7, 0x22, 0x6c, 0x96, 0x7a, 0x80,
	0xe8, 0x7f, 0x6b, 0x03, 0x7c, 0x36, 0x39, 0x3a,
	0x5e, 0x4b, 0xc3, 0x36, 0x62, 0xc8, 0x7b, 0x9b,
})

// testNet3GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 3).  It is the same as the merkle root
// for the main network.
var testNet3GenesisMerkleRoot = genesisMerkleRoot

// testNet3GenesisBlock defines the genesis block of the block chain which
// serves as the public transaction ledger for the test network (version 3).
var testNet3GenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},          // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: testNet3GenesisMerkleRoot, // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1371387277, 0),  // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x1e0ffff0,                // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      1000580675,                // 414098458
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// simNetGenesisHash is the hash of the first block in the block chain for the
// simulation test network.
var simNetGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xdf, 0xfa, 0x0a, 0xe6, 0xaa, 0x0d, 0x9b, 0xe3,
	0xc9, 0x03, 0xa7, 0x25, 0xea, 0x8c, 0x61, 0x61,
	0x5c, 0x48, 0xa6, 0xa1, 0x60, 0x30, 0x71, 0x80,
	0xb5, 0xb8, 0x46, 0x91, 0xfb, 0xf0, 0xae, 0xf2,
})

// simNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the simulation test network.  It is the same as the merkle root for
// the main network.
var simNetGenesisMerkleRoot = genesisMerkleRoot

// simNetGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the simulation test network.
var simNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: simNetGenesisMerkleRoot,  // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1401292357, 0), // 2014-05-28 15:52:37 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      2,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}
