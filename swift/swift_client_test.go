package swift

import (
	"os"
	"reflect"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

const rpcGenFilePath = "./testdata/rpc.generated.swift"
const protocolGenFilePath = "./testdata/protocol.generated.swift"
const multiProtocolGenFilePath = "./testdata/multiProtocol.generated.swift"

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
					"catalogue": testdata.CatalogueService{},
				},
				settings: Settings{},
			},
			outputFile: rpcGenFilePath,
		},
		{
			name: "generate protocol",
			fields: fields{
				servicesMap: map[string]zenrpc.Invoker{
					"catalogue": testdata.CatalogueService{},
				},
				settings: Settings{IsProtocol: true},
			},
			outputFile: protocolGenFilePath,
		},
		{
			name: "generate multi protocol",
			fields: fields{
				servicesMap: map[string]zenrpc.Invoker{
					"catalogue": testdata.CatalogueService{},
					"arith":     testdata.ArithService{},
				},
				settings: Settings{IsProtocol: true},
			},
			outputFile: multiProtocolGenFilePath,
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
			testData, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("open test data file: %v", err)
			}

			if !reflect.DeepEqual(got, testData) {
				t.Fatalf("bad generator output")
			}
		})
	}
}
