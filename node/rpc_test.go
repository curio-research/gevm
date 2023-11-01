package node

import (
	"bytes"
	"encoding/json"
	"time"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
	gt "github.com/daweth/gevm/gevmtypes"
	"github.com/stretchr/testify/assert"
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
		Params:  []interface{}{"Hello", "World", "Go is fun!"},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Error marshaling data: %v", err)
	}

	req, _ := http.NewRequest("POST", "/rpc", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	a.Server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Additional assertions can go here
}

func TestRPCEthSend(t *testing.T) {
	// Implement the test for /rpc send endpoint
}
