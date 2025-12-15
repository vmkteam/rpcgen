package php

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

const rpcGenFilePath = "./testdata/RpcClient.php"

var update = flag.Bool("update", false, "update .php files")

func TestGeneratePHPClient(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	cl := NewClient(rpc.SMD(), "")

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate php client: %v", err)
	}

	if *update {
		var f *os.File
		f, err = os.Create(rpcGenFilePath)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write(generated)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	testData, err := os.ReadFile(rpcGenFilePath)
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	// cut first two lines with version from comparison
	generated = bytes.TrimPrefix(generated, []byte("<?php\n"))
	testData = bytes.TrimPrefix(testData, []byte("<?php\n"))

	_, generatedBody, _ := bytes.Cut(generated, []byte{'\n'})
	_, testDataBody, _ := bytes.Cut(testData, []byte{'\n'})

	if !bytes.Equal(generatedBody, testDataBody) {
		t.Fatalf("bad generator output")
	}
}
