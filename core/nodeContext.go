package core

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	// logger "github.com/daweth/gevm/logger"
	// "github.com/daweth/gevm/state"
	"github.com/daweth/gevm/types"
	"github.com/daweth/gevm/util"
	"github.com/daweth/gevm/vm"

	"github.com/daweth/gevm/gevmtypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"

	gstate "github.com/ethereum/go-ethereum/core/state"
	gtypes "github.com/ethereum/go-ethereum/core/types"
	gvm "github.com/ethereum/go-ethereum/core/vm"
	glogger "github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/ethdb/pebble"
	"github.com/ethereum/go-ethereum/params"
)

type NodeCtx struct {
	Accounts []common.Address // accounts that this node handles.
	StateDB  *gstate.StateDB
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
	db := gstate.NewDatabaseWithConfig(rdb, nil)
	statedb, err := gstate.New(common.Hash{}, db, nil)

	// fill database with addresses
	for i := 0; i < len(accounts); i++ {
		fmt.Println("seeding the balance of the new account", accounts[i])
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
		Number:      big.NewInt(18437836),
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
		AccessList:        gtypes.AccessList{},
		BlobGasFeeCap:     big.NewInt(0),
		BlobHashes:        []common.Hash{},
		SkipAccountChecks: false,
	}

	cc := EmptyChainContext{}
	btx := NewEVMBlockContext(&header, cc, &accounts[0])
	ctx := NewEVMTxContext(&message)

	// create structLogger (for EVM config)
	chainConfig := params.TestChainConfig
	logConfig := glogger.Config{
		EnableMemory:     true,
		DisableStack:     true,
		DisableStorage:   false,
		EnableReturnData: true,
		Debug:            true,
		Limit:            0,
		Overrides:        chainConfig,
	}
	logger := glogger.NewStructLogger(&logConfig)
	vmConfig := gvm.Config{
		Tracer:                  logger,
		NoBaseFee:               true,
		EnablePreimageRecording: false,
		ExtraEips:               []int{},
	}
	vmcfg := util.ConvertGConfigToConfig(vmConfig, logConfig)
	// create new EVM
	evm := vm.NewEVM(btx, ctx, statedb, chainConfig, vmcfg)

	return NodeCtx{
		Accounts: accounts,
		StateDB:  statedb,
		Evm:      evm,
	}

}

var (
	gasLimit = uint64(1000000000000)
	gasUsed  = uint64(1)
	admin    = common.HexToAddress("alice")
	account1 = common.HexToAddress("bob")
)

func Default() NodeCtx {
	return NewNodeContext(gasLimit, gasUsed, admin, account1)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (n *NodeCtx) HandleTransaction(txn gevmtypes.Transaction) ([]byte, uint64) {
	if n == nil {
		// return empty data types if node context does not exist
		return []byte(""), 0
	}
	value := big.NewInt(0).SetUint64(txn.Value)
	outputs, gasLeft, vmerr := n.Evm.Call(StringToContractRef(txn.From), StringToAddress(txn.To), []byte(txn.Data), txn.Gas, value)
	must(vmerr)
	return outputs, gasLeft
}

func (n *NodeCtx) HandleCreateTransaction(txn gevmtypes.Transaction) ([]byte, uint64) {
	if n == nil {
		// return empty data types if node context does not exist
		return []byte(""), 0
	}

	value := big.NewInt(0).SetUint64(txn.Value)
	_, contractAddress, gasLeft, vmerr := n.Evm.Create(StringToContractRef(txn.From), []byte(txn.Data), txn.Gas, value)
	must(vmerr)
	return contractAddress[:], gasLeft
}

// HELPER FUNCTIONS
func StringToContractRef(s string) vm.ContractRef {
	hex := common.HexToAddress(s)
	return vm.AccountRef(hex)
}

func StringToAddress(s string) common.Address {
	return common.HexToAddress(s)
}

func StringToRawMessage(s string) json.RawMessage {
	msg := json.RawMessage([]byte(fmt.Sprintf("{\"data\": \"%v\"}", s)))
	return msg
}
