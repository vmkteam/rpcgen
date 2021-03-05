package golang

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/semrush/zenrpc/v2"
	"github.com/semrush/zenrpc/v2/testdata"
)

func TestGenerateCatalogue(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	cl := NewClient(rpc.SMD())

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate go cloent: %v", err)
	}

	testData, err := ioutil.ReadFile("./testdata/catalogue_client.go")
	if err != nil {
		t.Fatalf("open test data file")
	}

	if !bytes.Equal(generated, testData) {
		t.Fatalf("bad generator output")
	}
}
