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
	JsonRpc string `json:"jsonrpc"` // Version of the JSON RPC protocol, always set to 2.0
	Id      int    `json:"id"`      // Auto incrementing ID number for this request
	Error   []byte `json:"error"`   // Any error returned by the remote side
	Result  []byte `json:"result"`  // Whatever the remote side sends us in reply
}

func stringToRawMessage(s string) json.RawMessage {
	msg := json.RawMessage([]byte(fmt.Sprintf("{\"data\": \"%v\"}", s)))
	return msg
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

		// c.json(http.statusok, gin.h{
		// 	"message": "received!",
		// 	"request": req,
		// })

		resp := response{
			JsonRpc: "2.0",
			Id:      2,
			Error:   nil,
			Result:  []byte("hello"),
		}

		// Result:  json.RawMessage([]byte("0x3503de5f0c766c68f78a03a3b05036a5")),
		c.PureJSON(http.StatusOK, resp)

	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
