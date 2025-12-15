package golang

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

const rpcGenFilePath = "./testdata/catalogue_client.go.test"

var update = flag.Bool("update", false, "update .go files")

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

	// cut first line with version from comparison
	_, generatedBody, _ := bytes.Cut(generated, []byte{'\n'})
	_, testDataBody, _ := bytes.Cut(testData, []byte{'\n'})

	if !bytes.Equal(generatedBody, testDataBody) {
		t.Fatalf("bad generator output")
	}
}
