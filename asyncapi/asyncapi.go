// Package asyncapi generates AsyncAPI 3.0 documents for zenrpc servers that speak JSON-RPC 2.0
// over WebSocket: every SMD method becomes a send operation with a reply, and server to client
// notifications are described by a separate SMD schema, see Settings.Events.
package asyncapi

import (
	"encoding/json"

	"github.com/vmkteam/zenrpc/v2/smd"
)

// Generator main package structure.
type Generator struct {
	schema   smd.Schema
	settings Settings
}

// NewClient create Generator from zenrpc/v2 SMD.
func NewClient(schema smd.Schema, settings Settings) *Generator {
	return &Generator{schema: schema, settings: settings}
}

// Generate returns generated asyncapi schema.
func (g Generator) Generate() ([]byte, error) {
	doc, err := NewSchema(g.schema, g.settings)
	if err != nil {
		return nil, err
	}

	// two spaces: asyncapi tooling lints documents as yaml and warns on every tab.
	return json.MarshalIndent(doc, "", "  ")
}
