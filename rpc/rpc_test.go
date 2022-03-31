package rpc

// import (
// 	"os"
// 	"testing"
// 	"tstorage"
// )

// func Test_RPCInsertRows(t *testing.T) {
// 	tmpDir, err := os.MkdirTemp("", "tstorage-test")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	jrpc, err := NewJSONRPC(3000, tstorage.WithDataPath(tmpDir))

// 	if err != nil {
// 		panic(err)
// 	}

// 	jrpc.Main()
// }
