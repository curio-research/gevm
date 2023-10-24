package node

import (
	"encoding/json"
	"fmt"
	"net/http"

	cvm "github.com/daweth/gevm/core"
	gt "github.com/daweth/gevm/gevmtypes"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server *gin.Engine
	Node   cvm.NodeCtx
}

func NewServer() *App {
	app := &App{
		Server: gin.Default(),
		Node:   cvm.Default(),
	}

	app.Server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"jsonrpc": "2.0", "id": 1, "result": "0x3503de5f0c766c68f78a03a3b05036a5"})
	})

	// app.Server.POST("/addAccount", func(c *gin.Context)) {

	// }

	app.Server.POST("/rpc", func(c *gin.Context) {
		var req gt.Request
		var resp gt.Response
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		fmt.Println("printing the request", req.Method)

		switch m := req.Method; m {

		case "eth_call":
			fmt.Println("eth call")
			resp = app.handleEthCall(req)
		case "eth_send":
			fmt.Println("eth send")
		case "eth_sendRawTransaction":
			fmt.Println("eth send raw txn")
		case "eth_getBalance":
			fmt.Println("get bal")
		case "eth_getTransactionCount":
			fmt.Println("get transaction count")
		case "eth_getCode":
			fmt.Println("handling get code")
		default:
			resp = gt.Response{
				JsonRpc: "2.0",
				Id:      2,
				Error:   nil,
				Result:  []byte("hello"),
			}
		}
		
		// Result:  json.RawMessage([]byte("0x3503de5f0c766c68f78a03a3b05036a5")),
		c.PureJSON(http.StatusOK, resp)
	})

	return app
}

func (app *App) handleEthCall(r gt.Request) gt.Response {
	p := r.Params

	var tx gt.Transaction
	if err := json.Unmarshal([]byte(p[0].(string)), &tx); err != nil {
		// Handle the error
	}
	// tx := p[0].BindJSON(Transaction)
	//bl := p[1].BindJSON(BlockNumber)

	// handle transaction here using fn from
	o, g := app.Node.HandleTransaction(tx)

	return gt.Response{
		JsonRpc: "2.0",
		Id:      9,
		Error:   []byte(""),
		Result:  o,
		GasLeft: g,
	}
}

