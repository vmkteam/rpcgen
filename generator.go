package rpcgen

import (
	"fmt"
	"github.com/dizzyfool/rpcgen/v2/openrpc"
	"net/http"

	"github.com/dizzyfool/rpcgen/v2/golang"
	"github.com/dizzyfool/rpcgen/v2/php"
	"github.com/dizzyfool/rpcgen/v2/typescript"

	smd1 "github.com/vmkteam/zenrpc/smd"
	"github.com/vmkteam/zenrpc/v2/smd"
)

type RPCGen struct {
	schema smd.Schema
}

type Generator interface {
	Generate() ([]byte, error)
}

func (g RPCGen) GoClient() Generator {
	return golang.NewClient(g.schema)
}

func (g RPCGen) PHPClient(phpNamespace string) Generator {
	return php.NewClient(g.schema, phpNamespace)
}

func (g RPCGen) TSClient(typeMapper typescript.TypeMapper) Generator {
	return typescript.NewClient(g.schema, typeMapper)
}

func (g RPCGen) OpenRPC() Generator {
	return openrpc.NewClient(g.schema)
}

// FromSMD create Generator from smd schema
func FromSMD(schema smd.Schema) *RPCGen {
	return &RPCGen{schema: schema}
}

// FromSMDv1 create Generator from smd v1 schema
func FromSMDv1(schema smd1.Schema) *RPCGen {
	return FromSMD(smdv1ToSMD(schema))
}

// HTTP handlers

// Handler create http handler with fn
func Handler(gen Generator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bts, err := gen.Generate()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("error: %s", err)))
			return
		}

		_, _ = w.Write(bts)
	}
}
