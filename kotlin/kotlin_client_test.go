package kotlin

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/vmkteam/rpcgen/v2/kotlin/testdata"
	"github.com/vmkteam/zenrpc/v2"
)

const testDataPath = "./testdata/rpc.generated.kt"
const testDataProtocolPath = "./testdata/protocol.generated.kt"

var update = flag.Bool("update", false, "update .kt files")

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
			outputFile: testDataPath,
		},
		{
			name: "generate multi protocol",
			fields: fields{
				servicesMap: map[string]zenrpc.Invoker{
					"arith": testdata.ArithService{},
				},
				settings: Settings{
					IsProtocol: true,
					Imports:    []string{"api.debug.model.*", "api.TransportOption", "api.Transport"},
				},
			},
			outputFile: testDataProtocolPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rpc := zenrpc.NewServer(zenrpc.Options{})
			for sn, ss := range tt.fields.servicesMap {
				rpc.Register(sn, ss)
			}

			cl := NewClient(rpc.SMD(), tt.fields.settings)
			got, err := cl.Generate()
			if err != nil {
				t.Fatalf("generate client: %v", err)
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

func Test_kotlinTypeID(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     bool
	}{
		{
			name:     "id",
			typeName: "id",
			want:     true,
		},
		{
			name:     "name",
			typeName: "name",
			want:     false,
		},
		{
			name:     "nameIds",
			typeName: "nameIds",
			want:     true,
		},
		{
			name:     "nameIDs",
			typeName: "nameIDs",
			want:     true,
		},
		{
			name:     "ID",
			typeName: "ID",
			want:     true,
		},
		{
			name:     "nameId",
			typeName: "nameId",
			want:     true,
		},
		{
			name:     "nameID",
			typeName: "nameID",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kotlinTypeID(tt.typeName); got != tt.want {
				t.Errorf("kotlinTypeID() = %v, want %v", got, tt.want)
			}
		})
	}
}
