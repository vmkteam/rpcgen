package openrpc

import (
	"encoding/json"

	openrpc "github.com/open-rpc/meta-schema"
	"github.com/vmkteam/zenrpc/v2/smd"
)

// Generator main package structure
type Generator struct {
	schema openrpc.OpenrpcDocument
}

// NewClient create Generator from zenrpc/v2 SMD.
func NewClient(schema smd.Schema, title string) *Generator {
	return &Generator{schema: NewSchema(schema, title)}
}

// Generate returns generated openrpc schema.
func (g Generator) Generate() ([]byte, error) {
	return json.MarshalIndent(g.schema, "", "	")
}
