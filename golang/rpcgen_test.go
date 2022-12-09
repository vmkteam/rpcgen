package golang

import (
	"bytes"
	"os"
	"testing"

	"github.com/vmkteam/rpcgen/v2/gen"
	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

func TestGenerateGoClient(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})
	rpc.Register("phonebook", testdata.PhoneBook{})
	rpc.Register("arith", testdata.ArithService{})

	cl := NewClient(rpc.SMD(), Settings{})

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate go client: %v", err)
	}

	testData, err := os.ReadFile("./testdata/catalogue_client.go")
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
