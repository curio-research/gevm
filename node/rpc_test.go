package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	gt "github.com/daweth/gevm/gevmtypes"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var a *App

func TestMain(m *testing.M) {
	// Setup
	a = NewServer()
	go a.Server.Run(":8080") // Start server in a goroutine

	// Give the server a little time to start
	// Ideally, implement a wait mechanism to ensure the server is ready to handle requests
	time.Sleep(2 * time.Second)

	// Run tests
	exitVal := m.Run()

	// Teardown logic can go here if necessary

	// Exit with the value returned from m.Run()
	os.Exit(exitVal)
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	a.Server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Additional assertions can go here
}

func TestRPCEthCall(t *testing.T) {
	w := httptest.NewRecorder()

	data := gt.Request{
		JsonRpc: "2.0",
		Id:      9,
		Method:  "eth_call",
		Params:  []interface{}{"Hello", "World"},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Error marshaling data: %v", err)
	}

	req, _ := http.NewRequest("POST", "/rpc", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	a.Server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRPCEthSend(t *testing.T) {
	tx := types.NewTransaction(0, common.HexToAddress("bob"), big.NewInt(100), 1000000000, big.NewInt(20000), nil)

	// RLP encode the signed transaction
	rlpBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Fatalf("Failed to RLP encode transaction: %v", err)
	}
	rawTxHex := fmt.Sprintf("0x%x", rlpBytes)

	fmt.Printf("Raw TX: %s\n", rawTxHex)

	w := httptest.NewRecorder()

	data := gt.Request{
		JsonRpc: "2.0",
		Id:      9,
		Method:  "eth_sendRawTransaction",
		Params:  []interface{}{rawTxHex},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Error marshaling data: %v", err)
	}

	req, _ := http.NewRequest("POST", "/rpc", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	a.Server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	fmt.Println(w)
}

func TestTPCEthCallPrecompile(t *testing.T) {
	// value := big.NewInt(0) // this is a call transaction, the value is zero
	nonce := uint64(0)                  // account nonce
	gasLimit := uint64(1000000)         // gas limit for contract creation
	gasPrice := big.NewInt(10000000000) // Set an appropriate gas price
	value := big.NewInt(1)

	// Read the contract bytecode
	bytecodePath := "../examples/precompile/weather.bin"
	bytecode, err := ioutil.ReadFile(bytecodePath)
	if err != nil {
		log.Fatalf("unable to read contract bytecode: %v", err)
	}
	// Convert the byte slice to a string for use with common.FromHex
	bytecodeHex := string(bytecode)

	// **************************
	// CREATE THE CONTRACT VIA RPC
	// **************************

	// Construct the transaction
	tx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, common.FromHex(bytecodeHex))
	nonce++

	var buff bytes.Buffer
	err = rlp.Encode(&buff, tx)
	if err != nil {
		log.Fatalf("Failed to encode transaction: %v", err)
	}

	rawTxHex := fmt.Sprintf("0x%x", buff.Bytes())

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_sendRawTransaction",
		"params":  []interface{}{rawTxHex},
		"id":      1,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}

	req, _ := http.NewRequest("POST", "/rpc", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	a.Server.ServeHTTP(w, req)

	// Read the contract ABI
	abiPath := "../examples/precompile/weather.abi"
	contractABI, err := ioutil.ReadFile(abiPath)
	if err != nil {
		log.Fatalf("unable to read contract ABI: %v", err)
	}

	// **************************
	// CALL THE FUNCTION VIA RPC
	// **************************
	parsedABI, err := abi.JSON(bytes.NewReader(contractABI))
	if err != nil {
		log.Fatal(err)
	}

	// Packing the method call with its arguments
	data, err := parsedABI.Pack("getCurrentGameWeather")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(string(common.HexToAddress("bob")))

	tx = types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)

	// rawTxBytes := tx.GetRlp(0)
	// rawTxHex := fmt.Sprintf("0x%x", rawTxBytes) // This is the hex representation of the transaction

	fmt.Printf("Raw TX: %s\n", rawTxHex)

	data = gt.Request{
		JsonRpc: "2.0",
		Id:      9,
		Method:  "eth_call",
		Params:  []interface{}{rawTxHex},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Error marshaling data: %v", err)
	}

	req, _ = http.NewRequest("POST", "/rpc", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	a.Server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
