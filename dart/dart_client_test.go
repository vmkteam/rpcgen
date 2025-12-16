package dart

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

const rpcGenFilePath = "./testdata/client.dart"

var update = flag.Bool("update", false, "update .dart files")

func TestGenerateDartClient(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	cl := NewClient(rpc.SMD(), Settings{Part: "client"})

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate dart client: %v", err)
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

	testData, err := os.ReadFile("./testdata/client.dart")
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
