package rpcgen

import (
	"fmt"
	"net/http"

	"github.com/vmkteam/rpcgen/golang"
	"github.com/vmkteam/rpcgen/typescript"

	smd1 "github.com/semrush/zenrpc/smd"
	"github.com/semrush/zenrpc/v2/smd"
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

func (g RPCGen) TSClient(typeMapper typescript.TypeMapper) Generator {
	return typescript.NewClient(g.schema, typeMapper)
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
