package php

import (
	"bytes"
	"os"
	"testing"

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

	// cut first two lines with version from comparsion
	generated = bytes.TrimPrefix(generated, []byte("<?php\n"))
	testData = bytes.TrimPrefix(testData, []byte("<?php\n"))

	_, generatedBody, _ := bytes.Cut(generated, []byte{'\n'})
	_, testDataBody, _ := bytes.Cut(testData, []byte{'\n'})

	if !bytes.Equal(generatedBody, testDataBody) {
		t.Fatalf("bad generator output")
	}
}
