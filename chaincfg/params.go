// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
    "errors"
    "math"
    "math/big"
    "strings"
    "time"

    "github.com/nargott/godash/chaincfg/chainhash"
    "github.com/nargott/godash/wire"
)

// These variables are the chain proof-of-work limit parameters for each default
// network.
var (
    // bigOne is 1 represented as a big.Int.  It is defined here to avoid
    // the overhead of creating it multiple times.
    bigOne = big.NewInt(1)

    // mainPowLimit is the highest proof of work value a Bitcoin block can
    // have for the main network.  It is the value 2^224 - 1.
    mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)

    // regressionPowLimit is the highest proof of work value a Bitcoin block
    // can have for the regression test network.  It is the value 2^255 - 1.
    regressionPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)

    // testNet3PowLimit is the highest proof of work value a Bitcoin block
    // can have for the test network (version 3).  It is the value
    // 2^224 - 1.
    testNet3PowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)

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
    // soft-fork package. The CSV package includes the depolyment of BIPS
    // 68, 112, and 113.
    DeploymentCSV

    // DeploymentSegwit defines the rule change deployment ID for the
    // Segragated Witness (segwit) soft-fork package. The segwit package
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
    Net wire.DASHNet

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
    RetargetAdjustmentFactor int64

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
    DefaultPort: "9999",
    DNSSeeds: []DNSSeed{
        {"dnsseed.dash.org", true},
        {"dnsseed.dashdot.io", true},
        {"dnsseed.masternode.io", false},
        {"dnsseed.dashpay.io", true},
    },

    // Chain parameters
    GenesisBlock:             &genesisBlock,
    GenesisHash:              &genesisHash,
    PowLimit:                 mainPowLimit,
    PowLimitBits:             0x1d00ffff, //? DASH 00000fffff000000000000000000000000000000000000000000000000000000
    BIP0034Height:            1, // DASH 000007d91d1254d60e2dd1ae580383070a4ddffa4c64c2eeb4a2f9ecc0414343
    BIP0065Height:            388381, // 000000000000000004c2b624ed5d7756c508d90fd0da2c7c679febfa6c4735f0
    BIP0066Height:            363725, // 00000000000000000379eaa19dce8c9b722d46ae6a57c2f1a988119488b50931
    CoinbaseMaturity:         100,
    SubsidyReductionInterval: 210240,
    TargetTimespan:           24 * 60 * 60,      // Dash: 1 day
    TargetTimePerBlock:       time.Second * 150, // Dash: 2.5 minutes
    RetargetAdjustmentFactor: 4,                 // 25% less, 400% more
    ReduceMinDifficulty:      false,
    MinDiffReductionTime:     0,
    GenerateSupported:        false,

    // Checkpoints ordered from oldest to newest for DASH
    Checkpoints: []Checkpoint{
        {4991, newHashFromStr("000000003b01809551952460744d5dbb8fcbd6cbae3c220267bf7fa43f837367")},
        {9918, newHashFromStr("00000000213e229f332c0ffbe34defdaa9e74de87f2d8d1f01af8d121c3c170b")},
        {16912, newHashFromStr("00000000075c0d10371d55a60634da70f197548dbbfa4123e12abfcbc5738af9")},
        {23912, newHashFromStr("0000000000335eac6703f3b1732ec8b2f89c3ba3a7889e5767b090556bb9a276")},
        {35457, newHashFromStr("0000000000b0ae211be59b048df14820475ad0dd53b9ff83b010f71a77342d9f")},
        {45479, newHashFromStr("000000000063d411655d590590e16960f15ceea4257122ac430c6fbe39fbf02d")},
        {55895, newHashFromStr("0000000000ae4c53a43639a4ca027282f69da9c67ba951768a20415b6439a2d7")},
        {68899, newHashFromStr("0000000000194ab4d3d9eeb1f2f792f21bb39ff767cb547fe977640f969d77b7")},
        {74619, newHashFromStr("000000000011d28f38f05d01650a502cc3f4d0e793fbc26e2a2ca71f07dc3842")},
        {75095, newHashFromStr("0000000000193d12f6ad352a9996ee58ef8bdc4946818a5fec5ce99c11b87f0d")},
        {88805, newHashFromStr("00000000001392f1652e9bf45cd8bc79dc60fe935277cd11538565b4a94fa85f")},
        {107996, newHashFromStr("00000000000a23840ac16115407488267aa3da2b9bc843e301185b7d17e4dc40")},
        {137993, newHashFromStr("00000000000cf69ce152b1bffdeddc59188d7a80879210d6e5c9503011929c3c")},
        {167996, newHashFromStr("000000000009486020a80f7f2cc065342b0c2fb59af5e090cd813dba68ab0fed")},
        {207992, newHashFromStr("00000000000d85c22be098f74576ef00b7aa00c05777e966aff68a270f1e01a5")},
        {312645, newHashFromStr("0000000000059dcb71ad35a9e40526c44e7aae6c99169a9e7017b7d84b1c2daf")},
        {407452, newHashFromStr("000000000003c6a87e73623b9d70af7cd908ae22fee466063e4ffc20be1d2dbc")},
        {523412, newHashFromStr("000000000000e54f036576a10597e0e42cc22a5159ce572f999c33975e121d4d")},
        {523930, newHashFromStr("0000000000000bccdb11c2b1cfb0ecab452abf267d89b7f46eaf2d54ce6e652c")},
        {750000, newHashFromStr("00000000000000b4181bbbdddbae464ce11fede5d0292fb63fdede1e7c8ab21c")},
    },

    // Consensus rule change deployments.
    //
    // The miner confirmation window is defined as:
    //   target proof of work timespan / target proof of work spacing
    RuleChangeActivationThreshold: 1916, // 95% of 2016
    MinerConfirmationWindow:       2016, // nPowTargetTimespan / nPowTargetSpacing
    Deployments: [DefinedDeployments]ConsensusDeployment{
        DeploymentTestDummy: {
            BitNumber:  28,
            StartTime:  1199145601, // January 1, 2008 UTC
            ExpireTime: 1230767999, // December 31, 2008 UTC
        },
        DeploymentCSV: {
            BitNumber:  0,
            StartTime:  1486252800, // Feb 5th, 2017
            ExpireTime: 1517788800, // Feb 5th, 2018
        },
        DeploymentSegwit: {
            BitNumber:  1,
            StartTime:  1508025600, // Oct 15th, 2017
            ExpireTime: 1539561600, // Oct 15th, 2018
        },
    },

    // Mempool parameters
    RelayNonStdTxs: false,

    // Human-readable part for Bech32 encoded segwit addresses, as defined in
    // BIP 173.
    Bech32HRPSegwit: "bc", // always bc for main net

    // Address encoding magics
    PubKeyHashAddrID:        0x4c, // Dash addresses start with 'X'
    ScriptHashAddrID:        0x10, // Dash script addresses start with '7'
    PrivateKeyID:            0xcc, // Dash private keys start with '7' or 'X'
    WitnessPubKeyHashAddrID: 0x06, // starts with p2
    WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

    // BIP32 hierarchical deterministic extended key magics (DASH = Bitcoin)
    HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
    HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

    // BIP44 coin type used in the hierarchical deterministic path for
    // address generation.
    HDCoinType: 5, //for DASH
}

// RegressionNetParams defines the network parameters for the regression test
// Bitcoin network.  Not to be confused with the test Bitcoin network (version
// 3), this network is sometimes simply called "testnet".
var RegressionNetParams = Params{
    Name:        "regtest",
    Net:         wire.TestNet,
    DefaultPort: "19994",
    DNSSeeds:    []DNSSeed{},

    // Chain parameters
    GenesisBlock:             &regTestGenesisBlock,
    GenesisHash:              &regTestGenesisHash,
    PowLimit:                 regressionPowLimit,
    PowLimitBits:             0x207fffff,
    CoinbaseMaturity:         100,
    BIP0034Height:            100000000, // Not active - Permit ver 1 blocks
    BIP0065Height:            1351,      // Used by regression tests
    BIP0066Height:            1251,      // Used by regression tests
    SubsidyReductionInterval: 150,
    TargetTimespan:           time.Hour * 24 * 1, // DASH 1 day
    TargetTimePerBlock:       time.Second * 150,    // DASH 2.5 minutes
    RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
    ReduceMinDifficulty:      true,
    MinDiffReductionTime:     time.Minute * 20, // TargetTimePerBlock * 2
    GenerateSupported:        true,

    // Checkpoints ordered from oldest to newest.
    Checkpoints: nil,

    // Consensus rule change deployments.
    //
    // The miner confirmation window is defined as:
    //   target proof of work timespan / target proof of work spacing
    RuleChangeActivationThreshold: 108, // 75%  of MinerConfirmationWindow
    MinerConfirmationWindow:       144,
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
    Bech32HRPSegwit: "tb", // always tb for test net

    // Address encoding magics
    PubKeyHashAddrID: 0x8c, // Regtest Dash addresses start with 'y'
    ScriptHashAddrID: 0x13, // Regtest Dash script addresses start with '8' or '9'
    PrivateKeyID:     0xef, // starts with 9 (uncompressed) or c (compressed)

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
    DefaultPort: "19999",
    DNSSeeds: []DNSSeed{
        {"testnet-seed.dashdot.io", true},
        {"test.dnsseed.masternode.io", true},
    },

    // Chain parameters
    GenesisBlock:             &testNet3GenesisBlock,
    GenesisHash:              &testNet3GenesisHash,
    PowLimit:                 testNet3PowLimit,
    PowLimitBits:             0x1d00ffff,
    BIP0034Height:            1,  // 0000047d24635e347be3aaaeb66c26be94901a2f962feccd4f95090191f208c1
    BIP0065Height:            581885, // 00000000007f6655f22f98e72ed80d8b06dc761d5da09df0fa1dc4be4f861eb6
    BIP0066Height:            330776, // 000000002104c8c45e99a8853285a3b592602a3ccde2b832481da85e9e4ba182
    CoinbaseMaturity:         100,
    SubsidyReductionInterval: 210240,
    TargetTimespan:           time.Hour * 24 * 1, // DASH 1 day
    TargetTimePerBlock:       time.Second * 150,    // DASH 2.5 minutes
    RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
    ReduceMinDifficulty:      true,
    MinDiffReductionTime:     time.Minute * 20, // TargetTimePerBlock * 2
    GenerateSupported:        false,

    // Checkpoints ordered from oldest to newest.
    Checkpoints: []Checkpoint{
        {261, newHashFromStr("00000c26026d0815a7e2ce4fa270775f61403c040647ff2c3091f99e894a4618")},
        {1999, newHashFromStr("00000052e538d27fa53693efe6fb6892a0c1d26c0235f599171c48a3cce553b1")},
        {2999, newHashFromStr("0000024bc3f4f4cb30d29827c13d921ad77d2c6072e586c7f60d83c2722cdcc5")},
    },

    // Consensus rule change deployments.
    //
    // The miner confirmation window is defined as:
    //   target proof of work timespan / target proof of work spacing
    RuleChangeActivationThreshold: 1512, // 75% of MinerConfirmationWindow
    MinerConfirmationWindow:       2016,
    Deployments: [DefinedDeployments]ConsensusDeployment{
        DeploymentTestDummy: {
            BitNumber:  28,
            StartTime:  1199145601, // January 1, 2008 UTC
            ExpireTime: 1230767999, // December 31, 2008 UTC
        },
        DeploymentCSV: {
            BitNumber:  0,
            StartTime:  1456790400, // March 1st, 2016
            ExpireTime: 1493596800, // May 1st, 2017
        },
        DeploymentSegwit: {
            BitNumber:  1,
            StartTime:  1462060800, // May 1, 2016 UTC
            ExpireTime: 1493596800, // May 1, 2017 UTC.
        },
    },

    // Mempool parameters
    RelayNonStdTxs: true,

    // Human-readable part for Bech32 encoded segwit addresses, as defined in
    // BIP 173.
    Bech32HRPSegwit: "tb", // always tb for test net

    // Address encoding magics
    PubKeyHashAddrID:        0x8c, // Testnet Dash addresses start with 'y'
    ScriptHashAddrID:        0x13, // Testnet Dash script addresses start with '8' or '9'
    WitnessPubKeyHashAddrID: 0x03, // starts with QW
    WitnessScriptHashAddrID: 0x28, // starts with T7n
    PrivateKeyID:            0xef, // starts with 9 (uncompressed) or c (compressed)

    // BIP32 hierarchical deterministic extended key magics
    HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
    HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

    // BIP44 coin type used in the hierarchical deterministic path for
    // address generation.
    HDCoinType: 1,
}

var (
    // ErrDuplicateNet describes an error where the parameters for a Bitcoin
    // network could not be set due to the network already being a standard
    // network or previously-registered into this package.
    ErrDuplicateNet = errors.New("duplicate DASH network")

    // ErrUnknownHDKeyID describes an error where the provided id which
    // is intended to identify the network for a hierarchical deterministic
    // private extended key is not registered.
    ErrUnknownHDKeyID = errors.New("unknown hd private extended key bytes")
)

var (
    registeredNets       = make(map[wire.DASHNet]struct{})
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
}
