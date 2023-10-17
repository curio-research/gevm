package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// request is a JSON RPC request package assembled internally from the client
// method calls.
type request struct {
	JsonRpc string        `json:"jsonrpc"` // Version of the JSON RPC protocol, always set to 2.0
	Id      int           `json:"id"`      // Auto incrementing ID number for this request
	Method  string        `json:"method"`  // Remote procedure name to invoke on the server
	Params  []interface{} `json:"params"`  // List of parameters to pass through (keep types simple)
}

// response is a JSON RPC response package sent back from the API server.
type response struct {
	JsonRpc string          `json:"jsonrpc"` // Version of the JSON RPC protocol, always set to 2.0
	Id      int             `json:"id"`      // Auto incrementing ID number for this request
	Error   json.RawMessage `json:"error"`   // Any error returned by the remote side
	Result  json.RawMessage `json:"result"`  // Whatever the remote side sends us in reply
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"jsonrpc": "2.0", "id": 1, "result": "0x3503de5f0c766c68f78a03a3b05036a5"})
	})

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"jsonrpc": "2.0", "id": 1, "result": "0x3503de5f0c766c68f78a03a3b05036a5"})
	})

	r.POST("/rpc", func(c *gin.Context) {
		var req request
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		fmt.Println("printing the context", req)

		c.JSON(http.StatusOK, gin.H{
			"message": "Received!",
			"request": req,
		})

		//		c.JSON(http.StatusOK, gin.H{"jsonrpc": "2.0", "id": 1, "result": "0x3503de5f0c766c68f78a03a3b05036a5"})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
