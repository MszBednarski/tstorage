package main

import (
	"log"
	"time"
	"tstorage"
	"tstorage/rpc"
)

const (
	DAY       = 24 * time.Hour
	TWO_YEARS = 24 * 365 * 2 * time.Hour
)

func main() {
	jrpc, err := rpc.NewJSONRPC(3000,
		tstorage.WithDataPath("./data"),
		tstorage.WithTimestampPrecision("s"),
		tstorage.WithRetention(TWO_YEARS),
		tstorage.WithPartitionDuration(DAY),
	)
	if err != nil {
		panic(err)
	}
	log.Fatal(jrpc.Main())
}
