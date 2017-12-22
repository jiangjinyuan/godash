// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
    "time"

    "github.com/nargott/godash/chaincfg/chainhash"
    "github.com/nargott/godash/wire"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
    Version: 1,
    TxIn: []*wire.TxIn{
        {
            PreviousOutPoint: wire.OutPoint{
                Hash:  chainhash.Hash{},
                Index: 0xffffffff,
            },
            SignatureScript: []byte{
                0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x45, /* |.......E| */
                0x54, 0x68, 0x65, 0x20, 0x54, 0x69, 0x6d, 0x65, /* |The Time| */
                0x73, 0x20, 0x30, 0x33, 0x2f, 0x4a, 0x61, 0x6e, /* |s 03/Jan| */
                0x2f, 0x32, 0x30, 0x30, 0x39, 0x20, 0x43, 0x68, /* |/2009 Ch| */
                0x61, 0x6e, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x72, /* |ancellor| */
                0x20, 0x6f, 0x6e, 0x20, 0x62, 0x72, 0x69, 0x6e, /* | on brin| */
                0x6b, 0x20, 0x6f, 0x66, 0x20, 0x73, 0x65, 0x63, /* |k of sec|*/
                0x6f, 0x6e, 0x64, 0x20, 0x62, 0x61, 0x69, 0x6c, /* |ond bail| */
                0x6f, 0x75, 0x74, 0x20, 0x66, 0x6f, 0x72, 0x20, /* |out for |*/
                0x62, 0x61, 0x6e, 0x6b, 0x73,                   /* |banks| */
            },
            Sequence: 0xffffffff,
        },
    },
    TxOut: []*wire.TxOut{
        {
            Value: 0x12a05f200,
            PkScript: []byte{
                0x41, 0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, /* |A.g....U| */
                0x48, 0x27, 0x19, 0x67, 0xf1, 0xa6, 0x71, 0x30, /* |H'.g..q0| */
                0xb7, 0x10, 0x5c, 0xd6, 0xa8, 0x28, 0xe0, 0x39, /* |..\..(.9| */
                0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f, 0x61, /* |..yb...a| */
                0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, /* |..I..?L.| */
                0x38, 0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, /* |8..U....| */
                0x12, 0xde, 0x5c, 0x38, 0x4d, 0xf7, 0xba, 0x0b, /* |..\8M...| */
                0x8d, 0x57, 0x8a, 0x4c, 0x70, 0x2b, 0x6b, 0xf1, /* |.W.Lp+k.| */
                0x1d, 0x5f, 0xac,                               /* |._.| */
            },
        },
    },
    LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the DASH main
// network (genesis block).
var genesisHash = chainhash.Hash([chainhash.HashSize]byte{// Make go vet happy.
    0xb6, 0x7a, 0x40, 0xf3, 0xcd, 0x58, 0x04, 0x43,
    0x7a, 0x10, 0x8f, 0x10, 0x55, 0x33, 0x73, 0x9c,
    0x37, 0xe6, 0x22, 0x9b, 0xc1, 0xad, 0xca, 0xb3,
    0x85, 0x14, 0x0b, 0x59, 0xfd, 0x0f, 0x00, 0x00,
})

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the DASH main network.
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{// Make go vet happy.
    0xc7, 0x62, 0xa6, 0x56, 0x7f, 0x3c, 0xc0, 0x92,
    0xf0, 0x68, 0x4b, 0xb6, 0x2b, 0x7e, 0x00, 0xa8,
    0x48, 0x90, 0xb9, 0x90, 0xf0, 0x7c, 0xc7, 0x1a,
    0x6b, 0xb5, 0x8d, 0x64, 0xb9, 0x8e, 0x02, 0xe0,
})

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
    Header: wire.BlockHeader{
        Version:    1,
        PrevBlock:  chainhash.Hash{},         // DASH 00000ffd590b1485b3caadc19b22e6379c733355108f107a430458cdf3407ab6
        MerkleRoot: genesisMerkleRoot,        // DASH e0028eb9648db56b1ac77cf090b99048a8007e2bb64b68f092c03c7f56a662c7
        Timestamp:  time.Unix(0x52DB2D02, 0), // DASH Unix 1390095618
        Bits:       0x1e0ffff0,               // DASH
        Nonce:      0x121b062,                // 28917698 DASH
    },
    Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// regTestGenesisHash is the hash of the first block in the block chain for the
// regression test network (genesis block).
var regTestGenesisHash = chainhash.Hash([chainhash.HashSize]byte{// Make go vet happy.
    0x2e, 0x3d, 0xf2, 0x3e, 0xec, 0x5c, 0xd6, 0xa8,
    0x6e, 0xdd, 0x50, 0x95, 0x39, 0x02, 0x8e, 0x2c,
    0x3a, 0x3d, 0xc0, 0x53, 0x15, 0xeb, 0x28, 0xf2,
    0xba, 0xa4, 0x32, 0x18, 0xca, 0x08, 0x00, 0x00,
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
        PrevBlock:  chainhash.Hash{},         // DASH 000008ca1832a4baf228eb1553c03d3a2c8e02399550dd6ea8d65cec3ef23d2e
        MerkleRoot: regTestGenesisMerkleRoot, // DASH e0028eb9648db56b1ac77cf090b99048a8007e2bb64b68f092c03c7f56a662c7
        Timestamp:  time.Unix(1417713337, 0), // DASH 1417713337
        Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
        Nonce:      1096447,
    },
    Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet3GenesisHash is the hash of the first block in the block chain for the
// test network (version 3).
var testNet3GenesisHash = chainhash.Hash([chainhash.HashSize]byte{// Make go vet happy.
    0x2c, 0xbc, 0xf8, 0x3b, 0x62, 0x91, 0x3d, 0x56,
    0xf6, 0x05, 0xc0, 0xe5, 0x81, 0xa4, 0x88, 0x72,
    0x83, 0x94, 0x28, 0xc9, 0x2e, 0x5e, 0xb7, 0x6c,
    0xd7, 0xad, 0x94, 0xbc, 0xaf, 0x0b, 0x00, 0x00,
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
        PrevBlock:  chainhash.Hash{},          // 00000bafbc94add76cb75e2ec92894837288a481e5c005f6563d91623bf8bc2c
        MerkleRoot: testNet3GenesisMerkleRoot, // DASH genesisMerkleRoot
        Timestamp:  time.Unix(1390666206, 0),  // 2011-02-02 23:16:42 +0000 UTC
        Bits:       0x1e0ffff0,                //
        Nonce:      0xE627C9C3,                // 3861367235
    },
    Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}