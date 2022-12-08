package swift

import (
	"bytes"
	"os"
	"testing"

	"github.com/vmkteam/rpcgen/v2/gen"
	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

func TestGenerateSwiftClient(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	cl := NewClient(rpc.SMD(), Settings{})

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate swift client: %v", err)
	}

	testData, err := os.ReadFile("./testdata/rpc.generated.swift")
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
