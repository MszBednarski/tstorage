package main

import (
	"log"
	"tstorage/rpc"
)

func main() {
	jrpc, err := rpc.NewJSONRPC(rpc.DefaultOptions())
	if err != nil {
		panic(err)
	}
	log.Fatal(jrpc.Main())
}
