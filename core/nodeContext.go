package gevm

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/ethdb/pebble"
	"github.com/ethereum/go-ethereum/params"
)

type NodeCtx struct {
	Accounts []common.Address // accounts that this node handles.
	StateDB  *state.StateDB
	Evm      *vm.EVM
}

type NodeParams struct {
	gasLimit uint64
	gasUsed  uint64
}

func NewNodeContext(gasLimit uint64, gasUsed uint64, accounts ...common.Address) NodeCtx {
	pbl, err := pebble.New("gevm-db", 0, 0, "gevm", false, false)
	must(err)
	rdb := rawdb.NewDatabase(pbl)
	db := state.NewDatabaseWithConfig(rdb, nil)
	statedb, err := state.New(common.Hash{}, db, nil)

	// fill database with addresses
	for i := 1; i < len(accounts); i++ {
		statedb.GetOrNewStateObject(accounts[i])
		statedb.AddBalance(accounts[i], big.NewInt(1e18))
	}

	header := types.Header{
		ParentHash:  common.Hash{},
		UncleHash:   common.Hash{},
		Coinbase:    common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.BytesToBloom([]byte("daweth")),
		Difficulty:  big.NewInt(1),
		Number:      big.NewInt(1),
		GasLimit:    gasLimit,
		GasUsed:     gasUsed,
		Time:        uint64(time.Now().Unix()),
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       types.EncodeNonce(1),
	}

	message := core.Message{
		To:                &accounts[0],
		From:              accounts[1],
		Nonce:             uint64(1),
		Value:             big.NewInt(1),
		GasLimit:          gasLimit,
		GasPrice:          big.NewInt(0),
		GasFeeCap:         big.NewInt(0),
		GasTipCap:         big.NewInt(0),
		Data:              []byte(""),
		AccessList:        types.AccessList{},
		BlobGasFeeCap:     big.NewInt(0),
		BlobHashes:        []common.Hash{},
		SkipAccountChecks: false,
	}

	cc := EmptyChainContext{}
	btx := NewEVMBlockContext(&header, cc, &accounts[0])
	ctx := NewEVMTxContext(&message)
	// create structLogger (for EVM config)
	chainConfig := params.TestChainConfig
	logConfig := logger.Config{
		EnableMemory:     true,
		DisableStack:     true,
		DisableStorage:   false,
		EnableReturnData: true,
		Debug:            true,
		Limit:            0,
		Overrides:        chainConfig,
	}
	logger := logger.NewStructLogger(&logConfig)
	vmConfig := vm.Config{
		Tracer:                  logger,
		NoBaseFee:               true,
		EnablePreimageRecording: false,
		ExtraEips:               []int{},
	}

	// create new EVM
	evm := vm.NewEVM(btx, ctx, statedb, chainConfig, vmConfig)

	return NodeCtx{
		Accounts: accounts,
		StateDB:  statedb,
		Evm:      evm,
	}

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
