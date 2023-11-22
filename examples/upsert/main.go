package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	ec "github.com/daweth/gevm/core"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	// gm "github.com/ethereum/go-ethereum/common/math"
)

var (
	testAddress      = common.HexToAddress("alice")
	toAddress    = common.HexToAddress("bob")
	uncreatedAddress = common.HexToAddress("vitalik")
	amount           = big.NewInt(1)
	accountNonce     = uint64(0)
	gasLimit         = uint64(1000000)
	gasUsed          = uint64(1)
	blobHashes       = []common.Hash{}
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
func getTPS(start time.Time, end time.Time) int64 {
	dur := end.Sub(start)
	sec, _ := time.ParseDuration("1s")

	return sec.Nanoseconds() / dur.Nanoseconds()
}

func main() {
	binFilePath := "./upsert.bin"
	abiFilePath := "./upsert.abi"
	data := loadBin(binFilePath)
	abiObj := loadAbi(abiFilePath)

	alice, err := testAddress.MarshalText()
	must(err)
	vitalik, err := uncreatedAddress.MarshalText()
	must(err)

	node := ec.NewNodeContext(gasLimit, gasUsed, testAddress, toAddress)
	fmt.Println("Alice Addr=", alice)
	fmt.Println("Vitalik Addr=", vitalik)

	// creating the contract
	fmt.Println("balance: ", node.StateDB.GetBalance(testAddress).Uint64())
	contractRef := vm.AccountRef(testAddress)
	contractCode, contractAddress, gasLeftOver, vmerr := node.Evm.Create(contractRef, data, node.StateDB.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	fmt.Println("contract address=", contractAddress)

	// fund the account
	node.StateDB.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftOver))
	testBalance := node.StateDB.GetBalance(testAddress)
	fmt.Println("after contract creation, testBalance=", testBalance, contractCode)

	// calling the contract
	method := abiObj.Methods["sendETH"]
	pm := vitalik
	input := append(method.ID, pm...)
	fmt.Println(hexutil.Encode(input))

	startTime := time.Now()
	fmt.Println("begin to exec contract")
	node.StateDB.SetCode(testAddress, contractCode)
	outputs, gasLeft, vmerr := node.Evm.Call(contractRef, testAddress, input, node.StateDB.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	endTime := time.Now()

	executionTime := endTime.Sub(startTime)
	fmt.Printf("function executed in %v nanoseconds\n", executionTime.Nanoseconds())

	tps := getTPS(startTime, endTime)
	fmt.Printf("Theoretical TPS is %v\n", tps)

	node.StateDB.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeft))
	testBalance = node.StateDB.GetBalance(testAddress)
	fmt.Println("after call contract, testBalance =", testBalance)

	for _, op := range method.Outputs {
		switch op.Type.String() {
		case "uint256":
			fmt.Printf("Output name=%s, value=%d\n", op.Name, big.NewInt(0).SetBytes(outputs))
		default:
			fmt.Println(op.Name, op.Type.String())
		}
	}
	fmt.Printf("Output %#v\n", hexutil.Encode(outputs))
}
