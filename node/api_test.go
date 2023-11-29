package node

import (
	"bytes"
	"encoding/base64"
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

	txn := gt.Transaction{
		From:     common.HexToAddress("0x1").Hex(),
		To:       common.HexToAddress("0x2").Hex(),
		Gas:      1000000,
		GasPrice: 1000000000,
		Value:    1000000000000000000,
		Data:     "0x0",
	}

	rlpBytes, err := rlp.EncodeToBytes(txn)
	if err != nil {
		log.Fatalf("Failed to RLP encode transaction: %v", err)
	}

	// Convert RLP byte slice to hex string
	rawTxHex := hex.EncodeToString(rlpBytes)

	fmt.Println("Raw Transaction Hex:", rawTxHex)

	data := gt.Request{
		JsonRpc: "2.0",
		Id:      9,
		Method:  "eth_call",
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
}

func TestRPCUpsertAccount(t *testing.T) {
	// tx sending
	tx := types.NewTransaction(0, common.HexToAddress("UNKNOWN"), big.NewInt(100), 1000000000, big.NewInt(20000), nil)
	w := httptest.NewRecorder()
	data := gt.Request{
		JsonRpc: "2.0",
		Id:      9,
		Method:  "eth_sendRawTransaction",
		Params:  []interface{}{""},
	}

	// RLP encode the signed transaction
	rlpBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Fatalf("Failed to RLP encode transaction: %v", err)
	}
	rawTxHex := fmt.Sprintf("0x%x", rlpBytes)

	fmt.Printf("Raw TX: %s\n", rawTxHex)

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

func TestRPCGetBalance(t *testing.T) {
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
		Method:  "eth_getBalance",
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

	responseBody := w.Body.String()
	fmt.Println("Raw Response Body:", responseBody)

	// Attempt to unmarshal the response body
	var rmap map[string]interface{}
	err = json.Unmarshal([]byte(responseBody), &rmap)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check if the "result" key exists and assert it to a string.
	resultInterface, ok := rmap["result"]
	if !ok {
		log.Fatal("The key 'result' is not present in the response map.")
	}

	resultStr, ok := resultInterface.(string)
	if !ok {
		log.Fatal("The 'result' field is not of type string.")
	}

	// Decode the base64-encoded string to get the bytes.
	resultBytes, err := base64.StdEncoding.DecodeString(resultStr)
	if err != nil {
		log.Fatalf("Failed to decode base64 string: %v", err)
	}
	fmt.Println("result of getBalance:", BytesToUint64(resultBytes))

	// if resValue, ok := result["result"]; ok {
	// } else {
	// 	t.Fatal("Key 'result' not found in response JSON")
	// }
	// fmt.Println("result of getBalance", w)
}
func TestTPCEthCallPrecompile(t *testing.T) {
	// value := big.NewInt(0) // this is a call transaction, the value is zero
	nonce := uint64(0)                  // account nonce
	gasLimit := uint64(1000000)         // gas limit for contract creation
	gasPrice := big.NewInt(10000000000) // Set an appropriate gas price
	value := big.NewInt(1)
	contract := common.HexToAddress("precompile")

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

	rlpBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Fatalf("Failed to RLP encode transaction: %v", err)
	}
	rawTxHex := fmt.Sprintf("0x%x", rlpBytes)

	fmt.Printf("Raw TX: %s\n", rawTxHex)

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

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/rpc", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	assert.Equal(t, 200, w.Code)

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

	contractAddress := common.HexToAddress(contract.String())

	tx = types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)
	nonce++

	rlpBytes, err = rlp.EncodeToBytes(tx)
	if err != nil {
		log.Fatalf("Failed to RLP encode transaction: %v", err)
	}
	rawTxHex = fmt.Sprintf("0x%x", rlpBytes)

	fmt.Printf("Raw TX: %s\n", rawTxHex)

	data1 := gt.Request{
		JsonRpc: "2.0",
		Id:      9,
		Method:  "eth_call",
		Params:  []interface{}{rawTxHex},
	}

	jsonData, err := json.Marshal(data1)
	if err != nil {
		t.Fatalf("Error marshaling data: %v", err)
	}

	req, _ = http.NewRequest("POST", "/rpc", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	a.Server.ServeHTTP(w, req)

	fmt.Println("precompile", w.Result())
	assert.Equal(t, 200, w.Code)
}
