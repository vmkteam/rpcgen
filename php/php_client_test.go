package php

import (
	"bytes"
	"os"
	"testing"

	"github.com/vmkteam/rpcgen/v2/gen"
	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

func TestGeneratePHPClient(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	cl := NewClient(rpc.SMD(), "")

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate php client: %v", err)
	}

	testData, err := os.ReadFile("./testdata/RpcClient.php")
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	// cut version from comparsion
	generated = bytes.ReplaceAll(generated, []byte("v"+gen.DefaultGeneratorData().Version), []byte(""))
	testData = bytes.ReplaceAll(testData, []byte("v"+gen.DefaultGeneratorData().Version), []byte(""))

	if !bytes.Equal(generated, testData) {
		t.Fatalf("bad generator output")
	}
}
