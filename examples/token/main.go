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
	gm "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

var (
	testAddress  = common.HexToAddress("alice")
	toAddress    = common.HexToAddress("bob")
	amount       = big.NewInt(1)
	accountNonce = uint64(0)
	gasLimit     = uint64(100000)
	gasUsed      = uint64(1)
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
func getTPS(start time.Time, end time.Time) int64 {
	dur := end.Sub(start)
	sec, _ := time.ParseDuration("1s")

	return sec.Nanoseconds() / dur.Nanoseconds()
}

func main() {
	binFilePath := "./Token.bin"
	abiFilePath := "./Token.abi"
	data := loadBin(binFilePath)
	abiObj := loadAbi(abiFilePath)

	alice, err := testAddress.MarshalText()
	must(err)
	bob, err := toAddress.MarshalText()
	must(err)

	node := ec.NewNodeContext(gasLimit, gasUsed, testAddress, toAddress)
	fmt.Println("Alice Addr=", alice)
	fmt.Println("Bob Addr=", bob)

	// creating the contract
	fmt.Println("balance: ", node.StateDB.GetBalance(testAddress).Uint64())
	callerRef := vm.AccountRef(testAddress)

	_, contractAddress, gasLeftOver, vmerr := node.Evm.Create(callerRef, data, node.StateDB.GetBalance(testAddress).Uint64(), big.NewInt(0))
	fmt.Println("contract address=", contractAddress)
	must(vmerr)

	// deduct the gas
	node.StateDB.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftOver))
	testBalance := node.StateDB.GetBalance(testAddress)
	fmt.Println("after contract creation, testBalance=", testBalance)

	// MINT TRANSACTION

	method := abiObj.Methods["mint"]
	pm := gm.U256Bytes(big.NewInt(10))
	input := append(method.ID, pm...)
	//	fmt.Println(hexutil.Encode(input))

	// execute the transaction
	fmt.Println("begin to exec contract")
	//node.StateDB.SetCode(testAddress, contractCode)
	outputs, gasLeft, vmerr := node.Evm.Call(callerRef, contractAddress, input, node.StateDB.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	// after transaction cleanup
	node.StateDB.SetBalance(contractAddress, big.NewInt(0).SetUint64(gasLeft))
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

	// get the balance of the user
	method = abiObj.Methods["balanceOf"]
	input1 := append(method.ID, common.LeftPadBytes(testAddress[:], 32)...)
	outputs, gasLeft, vmerr = node.Evm.Call(callerRef, contractAddress, input1, node.StateDB.GetBalance(testAddress).Uint64(), big.NewInt(0))
	// deduct the gas
	node.StateDB.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftOver))
	testBalance = node.StateDB.GetBalance(testAddress)
	fmt.Println("after contract creation, testBalance=", testBalance)

	for _, op := range method.Outputs {
		switch op.Type.String() {
		case "uint256":
			fmt.Printf("Output name=%s, value=%d\n", op.Name, big.NewInt(0).SetBytes(outputs))
		default:
			fmt.Println(op.Name, op.Type.String(), hexutil.Encode((outputs)))
		}
	}

	// TRANSFER TRANSACTION

	// create the transaction data for transfer
	method = abiObj.Methods["transfer"]
	input2 := append(method.ID, common.LeftPadBytes(toAddress[:], 32)...)
	pm = gm.U256Bytes(big.NewInt(9))
	input2 = append(input2, pm...)

	startTime := time.Now()
	fmt.Println("begin to exec contract")

	// execute the transaction
	outputs, gasLeft, vmerr = node.Evm.Call(callerRef, contractAddress, input2, node.StateDB.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	// deduct the gas
	node.StateDB.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftOver))
	testBalance = node.StateDB.GetBalance(testAddress)
	fmt.Println("after contract creation, testBalance=", testBalance)

	endTime := time.Now()

	executionTime := endTime.Sub(startTime)
	fmt.Printf("function executed in %v nanoseconds\n", executionTime.Nanoseconds())

	tps := getTPS(startTime, endTime)
	fmt.Printf("Theoretical TPS is %v\n", tps)

	for _, op := range method.Outputs {
		res, err := abiObj.Unpack("transfer", outputs)
		must(err)

		fmt.Println("DECODED RESULT", res)

		switch op.Type.String() {
		case "uint256":
			fmt.Printf("Output name=%s, value=%d\n", op.Name, big.NewInt(0).SetBytes(outputs))
		default:
			fmt.Println(op.Name, op.Type.String(), hexutil.Encode(outputs))
		}
	}

	// get the balance of the user
	method = abiObj.Methods["balanceOf"]
	input3 := append(method.ID, common.LeftPadBytes(toAddress[:], 32)...)
	outputs, gasLeft, vmerr = node.Evm.Call(callerRef, contractAddress, input3, node.StateDB.GetBalance(testAddress).Uint64(), big.NewInt(0))
	// deduct the gas
	node.StateDB.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftOver))
	testBalance = node.StateDB.GetBalance(testAddress)
	fmt.Println("after contract creation, testBalance=", testBalance)

	// should be 9
	for _, op := range method.Outputs {
		switch op.Type.String() {
		case "uint256":
			fmt.Printf("Output name=%s, value=%d\n", op.Name, big.NewInt(0).SetBytes(outputs))
		default:
			fmt.Println(op.Name, op.Type.String(), hexutil.Encode((outputs)))
		}
	}

}

type ChainContext struct{}

func (cc ChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	fmt.Println("(cc ChainContext) GetHeader (hash common.Hash, number uint64)")
	return nil
}
