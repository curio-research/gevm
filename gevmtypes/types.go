package gevmtypes

// request is a JSON RPC request package assembled internally from the client
// method calls.
type Request struct {
	JsonRpc string        `json:"jsonrpc"` // Version of the JSON RPC protocol, always set to 2.0
	Id      int           `json:"id"`      // Auto incrementing ID number for this request
	Method  string        `json:"method"`  // Remote procedure name to invoke on the server
	Params  []interface{} `json:"params"`  // List of parameters to pass through (keep types simple)
}

// response is a JSON RPC response package sent back from the API server.
type Response struct {
	JsonRpc string `json:"jsonrpc"` // Version of the JSON RPC protocol, always set to 2.0
	Id      int    `json:"id"`      // Auto incrementing ID number for this request
	Error   []byte `json:"error"`   // Any error returned by the remote side
	Result  []byte `json:"result"`  // Whatever the remote side sends us in reply
	GasLeft uint64 `json:"gasLeft"` // Gas left over from the transaction
}

// transaction is the data payload from the caller
type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      uint64 `json:"gas"`
	GasPrice uint64 `json:"gasPrice"`
	Value    uint64 `json:"value"`
	Data     string `json:"data"`
}

// weather is a type sent when changing / getting weather
type Weather struct {
	Weather int `json:"weather"`
}

// keep track of request ids
type Ids struct {
	EthCall              int
	EthSendRawTransaction int
	EthSend              int
}

type BlockNumber string

type Address string
