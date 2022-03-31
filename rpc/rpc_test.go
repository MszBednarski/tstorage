package rpc

// import (
// 	"os"
// 	"testing"
// )

// func Test_RPCInsertRows(t *testing.T) {
// 	tmpDir, err := os.MkdirTemp("", "tstorage-test")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	opts := DefaultOptions()
// 	opts.dataPath = tmpDir

// 	jrpc, err := NewJSONRPC(opts)
// 	if err != nil {
// 		panic(err)
// 	}
// 	jrpc.Main()
// }
