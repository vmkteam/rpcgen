package swift

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/vmkteam/rpcgen/v2/swift/testdata"
	"github.com/vmkteam/zenrpc/v2"
)

const rpcGenFilePath = "./testdata/rpc.generated.swift"
const protocolGenFilePath = "./testdata/protocol.generated.swift"

var update = flag.Bool("update", false, "update .swift files")

func TestGenerator_Generate(t *testing.T) {
	type fields struct {
		settings    Settings
		servicesMap map[string]zenrpc.Invoker
	}
	tests := []struct {
		name       string
		fields     fields
		outputFile string
		wantErr    bool
	}{
		{
			name: "generate rpc",
			fields: fields{
				servicesMap: map[string]zenrpc.Invoker{
					"arith": testdata.ArithService{},
				},
				settings: Settings{},
			},
			outputFile: rpcGenFilePath,
		},
		{
			name: "generate multi protocol",
			fields: fields{
				servicesMap: map[string]zenrpc.Invoker{
					"arith": testdata.ArithService{},
				},
				settings: Settings{IsProtocol: true},
			},
			outputFile: protocolGenFilePath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rpc := zenrpc.NewServer(zenrpc.Options{})
			for sn, service := range tt.fields.servicesMap {
				rpc.Register(sn, service)
			}

			cl := NewClient(rpc.SMD(), tt.fields.settings)
			got, err := cl.Generate()
			if err != nil {
				t.Fatalf("generate swift client: %v", err)
			}

			if *update {
				var f *os.File
				f, err = os.Create(tt.outputFile)
				if err != nil {
					t.Fatal(err)
				}
				_, err = f.Write(got)
				if err != nil {
					t.Fatal(err)
				}
				return
			}

			testData, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("open test data file: %v", err)
			}

			_, generatedBody, _ := bytes.Cut(got, []byte{'\n'})
			_, testDataBody, _ := bytes.Cut(testData, []byte{'\n'})

			if !bytes.Equal(generatedBody, testDataBody) {
				t.Fatalf("bad generator output")
			}
		})
	}
}
