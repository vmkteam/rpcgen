package typescript

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

const (
	clientTS = "./testdata/catalogue_client.ts"
	classes  = "./testdata/catalogue_with_classes.ts"
	typeOnly = "./testdata/catalogue_with_types_only.ts"
)

var update = flag.Bool("update", false, "update .ts files")

func TestGenerateTypeScriptClient(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	cl := NewClient(rpc.SMD(), Settings{})

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate typescript client: %v", err)
	}

	if *update {
		var f *os.File
		f, err = os.Create(clientTS)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write(generated)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	testData, err := os.ReadFile(clientTS)
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	// cut first line with version from comparsion
	_, generatedBody, _ := bytes.Cut(generated, []byte{'\n'})
	_, testDataBody, _ := bytes.Cut(testData, []byte{'\n'})

	if !bytes.Equal(generatedBody, testDataBody) {
		t.Fatalf("bad generator output")
	}
}

func TestGenerateTypeScriptClasses(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	cl := NewClient(rpc.SMD(), Settings{WithClasses: true})

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate typescript client: %v", err)
	}

	if *update {
		var f *os.File
		f, err = os.Create(classes)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write(generated)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	testData, err := os.ReadFile(classes)
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	// cut first line with version from comparsion
	_, generatedBody, _ := bytes.Cut(generated, []byte{'\n'})
	_, testDataBody, _ := bytes.Cut(testData, []byte{'\n'})

	if !bytes.Equal(generatedBody, testDataBody) {
		t.Fatalf("bad generator output")
	}
}

func TestGenerateTypeScriptTypesOnly(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	cl := NewClient(rpc.SMD(), Settings{TypesOnly: true})

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate typescript client: %v", err)
	}

	if *update {
		var f *os.File
		f, err = os.Create(typeOnly)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write(generated)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	testData, err := os.ReadFile(typeOnly)
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	// cut first line with version from comparsion
	_, generatedBody, _ := bytes.Cut(generated, []byte{'\n'})
	_, testDataBody, _ := bytes.Cut(testData, []byte{'\n'})

	if !bytes.Equal(generatedBody, testDataBody) {
		t.Fatalf("bad generator output")
	}
}
