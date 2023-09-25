package main

import (
	"fmt"
	"math/big"
	"time"
	"os"
	"io/ioutil"

	ec "github.com/daweth/gevm/core"
	//	"github.com/daweth/gevm/logger"
	//	"github.com/daweth/gevm/state"
	//	"github.com/daweth/gevm/types"
	//	"github.com/daweth/gevm/vm"
	"github.com/ethereum/go-ethereum/core"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/ethdb/pebble"
	"github.com/ethereum/go-ethereum/params"
)

var (
	testAddress  = common.HexToAddress("alice")
	toAddress    = common.HexToAddress("bob")
	amount       = big.NewInt(1)
	accountNonce = uint64(0)
	gasLimit     = uint64(100000)
	gasUsed      = uint64(1)
	codeStr      = "0x6060604052341561000f57600080fd5b60b18061001d6000396000f300606060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063c6888fa1146044575b600080fd5b3415604e57600080fd5b606260048080359060200190919050506078565b6040518082815260200191505060405180910390f35b60006007820290509190505600a165627a7a72305820c4ac950a92caa9944a7e07e030542e9ed7db92631adcc234d86a105c853b81a20029"
	blobHashes   = []common.Hash{}
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}	
func loadBin(filename string) []byte {
	code, err := ioutil.ReadFile(filename)
	must(err)
	return hexutil.MustDecode("0x" + string(code))
}
func loadAbi(filename string) abi.ABI {
	abiFile, err := os.Open(filename)
	must(err)
	defer abiFile.Close()
	abiObj, err := abi.JSON(abiFile)
	must(err)
	return abiObj
}

func main() {
	binFilePath := "./sum.bin"
	abiFilePath := "./sum.abi"
	data := loadBin(binFilePath)
	abiObj := loadAbi(abiFilePath)

	alice, err := testAddress.MarshalText()
	must(err)
	bob, err := toAddress.MarshalText()
	must(err)

	fmt.Println("Alice Addr=", alice)
	fmt.Println("Bob Addr=", bob)

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
		To:                &toAddress,
		From:              testAddress,
		Nonce:             uint64(1),
		Value:             amount,
		GasLimit:          gasLimit,
		GasPrice:          big.NewInt(0),
		GasFeeCap:         big.NewInt(0),
		GasTipCap:         big.NewInt(0),
		Data:              data,
		AccessList:        types.AccessList{},
		BlobGasFeeCap:     big.NewInt(0),
		BlobHashes:        blobHashes,
		SkipAccountChecks: false,
	}

	cc := ChainContext{}
	btx := ec.NewEVMBlockContext(&header, cc, &testAddress)
	ctx := ec.NewEVMTxContext(&message)

	pbl, err := pebble.New("gevm-db", 0, 0, "gevm", false, false)
	must(err)

	rdb := rawdb.NewDatabase(pbl)
	db := state.NewDatabaseWithConfig(rdb, nil)

	statedb, err := state.New(common.Hash{}, db, nil)

	statedb.GetOrNewStateObject(testAddress)
	statedb.GetOrNewStateObject(toAddress)
	statedb.AddBalance(testAddress, big.NewInt(1e18))
	testBalance := statedb.GetBalance(testAddress)
	fmt.Println("testBalance =", testBalance)
	must(err)

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

	evm := vm.NewEVM(btx, ctx, statedb, chainConfig, vmConfig)

	// creating the contract

	contractRef := vm.AccountRef(testAddress)
	contractCode, _, gasLeftOver, vmerr := evm.Create(contractRef, data, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	statedb.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftOver))
	testBalance = statedb.GetBalance(testAddress)
	fmt.Println("after contract creation, testBalance=", testBalance, contractCode)

	// calling the contract

	inttypes, err := abi.NewType("uint", "uint256", []abi.ArgumentMarshaling{})
	must(err)

	abiCall := abi.NewMethod(
		"multiply",
		"",
		abi.FunctionType(1),
		"view",
		false,
		false,
		[]abi.Argument{abi.Argument{Name: "a", Type: inttypes}},
		[]abi.Argument{abi.Argument{Name: "d", Type: inttypes}},
	)

	pm := common.BigToHash(big.NewInt(-10)).Hex()
	fmt.Println(pm)

	inputstr := hexutil.Encode(abiCall.ID) + pm[2:]
	input, err := hexutil.Decode((inputstr))
	must(err)

	fmt.Println("begin to exec contract")
	statedb.SetCode(testAddress, contractCode)
	outputs, gasLeft, vmerr := evm.Call(contractRef, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	statedb.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeft))
	testBalance = statedb.GetBalance(testAddress)
	fmt.Println("after call contract, testBalance =", testBalance)
	fmt.Printf("Output %#v\n", hexutil.Encode(outputs))

}

type ChainContext struct{}

func (cc ChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	fmt.Println("(cc ChainContext) GetHeader (hash common.Hash, number uint64)")
	return nil
}
