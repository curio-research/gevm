package main

import (
	"fmt"

	server "github.com/daweth/gevm/node"
)

func main() {
	s := server.NewServer()
	fmt.Println("Start the server on port 8080")
	s.Server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
