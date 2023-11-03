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
	Server  *gin.Engine
	Node    cvm.NodeCtx
	Weather gt.Weather
	Count   *gt.Ids
}

func NewServer() *App {
	app := &App{
		Server:  gin.Default(),
		Node:    cvm.Default(),
		Weather: gt.Weather{},
		Count:   &gt.Ids{},
	}

	// simple sanity check
	app.Server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"jsonrpc": "2.0", "id": 1, "result": "0x3503de5f0c766c68f78a03a3b05036a5"})
	})

	// since we don't have signatures they need to be manually created in db
	// app.Server.POST("/addAccount", func(c *gin.Context)) {
	// }

	app.Server.POST("/setWeather", func(c *gin.Context) {
		fmt.Println("printing the request body", c.Request.Body)
		var w gt.Weather

		if err := c.BindJSON(&w); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		app.handleSetWeather(w)
		fmt.Println("printing the new weather value", w.Weather)
		c.JSON(http.StatusOK, gin.H{"status": "Success"})
	})

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
			resp = app.handleEthSend(req)
		case "eth_sendRawTransaction":
			fmt.Println("eth send raw txn")
			resp = app.handleEthSendRawTransaction(req)
		case "eth_getBalance":
			fmt.Println("get bal")
			resp = app.handleEthGetBalance(req)
		case "eth_getTransactionCount":
			fmt.Println("get transaction count")
			// resp = app.handleGetTxCount(req)
		case "eth_getCode":
			fmt.Println("handling get code")
			// resp = app.handleGetCode(req)
		default:
			resp = gt.Response{
				JsonRpc: "2.0",
				Id:      2,
				Error:   nil,
				Result:  []byte(""),
			}
		}

		c.PureJSON(http.StatusOK, resp)
	})

	return app
}

func (app *App) handleEthCall(r gt.Request) gt.Response {
	p := r.Params

	var tx gt.Transaction
	if err := json.Unmarshal([]byte(p[0].(string)), &tx); err != nil {
		// handle the error
	}
	// tx := p[0].BindJSON(Transaction)
	//bl := p[1].BindJSON(BlockNumber)

	o, g := app.Node.HandleTransaction(tx)

	currId := app.Count.EthCall
	app.Count.EthCall++

	return gt.Response{
		JsonRpc: "2.0",
		Id:      currId,
		Error:   []byte(""),
		Result:  o,
		GasLeft: g,
	}
}

func (app *App) handleEthSend(r gt.Request) gt.Response {
	p := r.Params

	var tx gt.Transaction
	if err := json.Unmarshal([]byte(p[0].(string)), &tx); err != nil {
		// handle the error
	}

	// check that no transaction data exists

	o, g := app.Node.HandleTransaction(tx)

	currId := app.Count.EthSend
	app.Count.EthSend++

	return gt.Response{
		JsonRpc: "2.0",
		Id:      currId,
		Error:   []byte(""),
		Result:  o,
		GasLeft: g,
	}
}

func (app *App) handleEthSendRawTransaction(r gt.Request) gt.Response {
	p := r.Params

	var tx gt.Transaction
	if err := json.Unmarshal([]byte(p[0].(string)), &tx); err != nil {
		// handle the error
	}

	o, g := app.Node.HandleTransaction(tx)

	currId := app.Count.EthSend
	app.Count.EthSendRawTransaction++

	return gt.Response{
		JsonRpc: "2.0",
		Id:      currId,
		Error:   []byte(""),
		Result:  o,
		GasLeft: g,
	}
}

func (app *App) handleEthGetBalance(r gt.Request) gt.Response {
	p := r.Params

	var tx gt.Transaction
	if err := json.Unmarshal([]byte(p[0].(string)), &tx); err != nil {
		// handle the error
	}

	o := app.Node.HandleGetBalance(tx)

	currId := app.Count.EthSend
	app.Count.EthSendRawTransaction++

	return gt.Response{
		JsonRpc: "2.0",
		Id:      currId,
		Error:   []byte(""),
		Result:  Uint64ToBytes(o),
		GasLeft: 0,
	}

}

func (app *App) handleSetWeather(r gt.Weather) {
	w := &app.Weather
	w.Weather = r.Weather
}
