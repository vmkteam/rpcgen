package golang

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

func TestGenerateGoClient(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})
	rpc.Register("phonebook", testdata.PhoneBook{})
	rpc.Register("arith", testdata.ArithService{})

	cl := NewClient(rpc.SMD())

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate go client: %v", err)
	}

	testData, err := ioutil.ReadFile("./testdata/catalogue_client.go")
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	if !bytes.Equal(generated, testData) {
		t.Fatalf("bad generator output")
	}
}
