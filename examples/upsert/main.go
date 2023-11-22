package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	// "github.com/daweth/gevm/vm"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	testAddress      = common.HexToAddress("alice")
	toAddress        = common.HexToAddress("bob")
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
	alice, err := testAddress.MarshalText()
	must(err)
	vitalik, err := uncreatedAddress.MarshalText()
	must(err)

	// node := ec.NewNodeContext(gasLimit, gasUsed, testAddress, toAddress)
	fmt.Println("Alice Addr=", alice)
	fmt.Println("Vitalik Addr=", vitalik)

	

}
