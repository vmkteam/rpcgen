package openrpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/smd"
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
		t.Fatalf("generate go client: %v", err)
	}

	testData, err := ioutil.ReadFile("./testdata/openrpc.json")
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	if !bytes.Equal(generated, testData) {
		t.Fatalf("bad generator output")
	}
}

func TestGenerateOpenRPCClientAPISRV(t *testing.T) {
	bs, err := ioutil.ReadFile("./testdata/api.json")

	schema := smd.Schema{}
	if err := json.Unmarshal(bs, &schema); err != nil {
		t.Fatalf("read smd: %v", err)
	}

	cl := NewClient(schema, "test", "")

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate go client: %v", err)
	}

	err = ioutil.WriteFile("./testdata/api_openrpc.json", generated, os.ModePerm)
	if err != nil {
		t.Fatalf("save openrpc: %v", err)
	}
}
