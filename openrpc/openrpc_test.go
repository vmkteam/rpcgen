package openrpc

import (
	"bytes"
	"os"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

func TestGenerateOpenRPCSchema(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})
	rpc.Register("phonebook", testdata.PhoneBook{})
	rpc.Register("arith", testdata.ArithService{})

	cl := NewClient(rpc.SMD(), "test", "")

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate openrpc client: %v", err)
	}

	testData, err := os.ReadFile("./testdata/openrpc.json")
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	if !bytes.Equal(generated, testData) {
		t.Fatalf("bad generator output")
	}

	// second run
	generated, err = cl.Generate()
	if err != nil {
		t.Fatalf("generate openrpc client: %v", err)
	}

	if !bytes.Equal(generated, testData) {
		t.Fatalf("bad generator output")
	}
}
