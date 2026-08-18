package asyncapi

import (
	"bytes"
	"encoding/json"
)

// Version is a supported AsyncAPI specification version.
const Version = "3.0.0"

// WebSocketBindingVersion is a version of AsyncAPI WebSocket bindings.
const WebSocketBindingVersion = "0.1.0"

// Document is a root object of AsyncAPI document.
// See https://www.asyncapi.com/docs/reference/specification/v3.0.0.
type Document struct {
	AsyncAPI   string                `json:"asyncapi"`
	Info       Info                  `json:"info"`
	Servers    map[string]Server     `json:"servers,omitempty"`
	Channels   map[string]*Channel   `json:"channels"`
	Operations map[string]*Operation `json:"operations,omitempty"`
	Components *Components           `json:"components,omitempty"`
}

// Info provides metadata about the API.
type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

// Server is a message broker, a server or some other kind of computer program capable of sending
// and/or receiving data.
type Server struct {
	Host        string `json:"host"`
	Protocol    string `json:"protocol"`
	Description string `json:"description,omitempty"`
}

// Channel is an addressable component responsible for mediating the exchange of messages.
type Channel struct {
	Address     string              `json:"address"`
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Servers     []Reference         `json:"servers,omitempty"`
	Bindings    *ChannelBindings    `json:"bindings,omitempty"`
	Messages    map[string]*Message `json:"messages,omitempty"`
}

// ChannelBindings holds protocol-specific definitions for a channel.
type ChannelBindings struct {
	WebSocket *WebSocketBinding `json:"ws,omitempty"`
}

// WebSocketBinding describes a WebSocket handshake: the HTTP request that opens the connection.
type WebSocketBinding struct {
	Method         string  `json:"method,omitempty"`
	Query          *Schema `json:"query,omitempty"`
	Headers        *Schema `json:"headers,omitempty"`
	BindingVersion string  `json:"bindingVersion"`
}

// Message describes a single frame of the channel.
type Message struct {
	Name        string  `json:"name,omitempty"`
	Title       string  `json:"title,omitempty"`
	Summary     string  `json:"summary,omitempty"`
	Description string  `json:"description,omitempty"`
	ContentType string  `json:"contentType,omitempty"`
	Payload     *Schema `json:"payload,omitempty"`
}

// Operation describes a specific action the application performs on a channel.
type Operation struct {
	Action      string      `json:"action"`
	Channel     Reference   `json:"channel"`
	Title       string      `json:"title,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Messages    []Reference `json:"messages,omitempty"`
	Reply       *Reply      `json:"reply,omitempty"`
}

// Reply describes the message that is sent back in a request/reply interaction.
type Reply struct {
	Channel  Reference   `json:"channel"`
	Messages []Reference `json:"messages,omitempty"`
}

// Components holds reusable objects of the document.
type Components struct {
	Schemas Schemas `json:"schemas,omitempty"`
}

// Reference is a JSON Reference to another object of the document.
type Reference struct {
	Ref string `json:"$ref"`
}

// NewReference creates Reference to path.
func NewReference(path string) Reference {
	return Reference{Ref: path}
}

// Schema is a JSON Schema (draft-07) subset used by AsyncAPI payloads.
type Schema struct {
	Ref                  string    `json:"$ref,omitempty"`
	Type                 any       `json:"type,omitempty"` // string, or []string for nullable types
	Title                string    `json:"title,omitempty"`
	Description          string    `json:"description,omitempty"`
	Const                any       `json:"const,omitempty"`
	Properties           Schemas   `json:"properties,omitempty"`
	Required             []string  `json:"required,omitempty"`
	Items                *Schema   `json:"items,omitempty"`
	OneOf                []*Schema `json:"oneOf,omitempty"`
	Examples             []any     `json:"examples,omitempty"`
	AdditionalProperties *bool     `json:"additionalProperties,omitempty"`
}

// NamedSchema is a single entry of Schemas.
type NamedSchema struct {
	Name   string
	Schema *Schema
}

// Schemas is an ordered map of schemas: unlike map[string]*Schema it marshals in insertion order,
// so generated documents keep the field order of the original SMD.
type Schemas []NamedSchema

// Add appends schema with name to the list. Duplicates are not checked, see builder.addComponent.
func (s *Schemas) Add(name string, schema *Schema) {
	*s = append(*s, NamedSchema{Name: name, Schema: schema})
}

// Get returns schema by name.
func (s Schemas) Get(name string) (*Schema, bool) {
	for _, ns := range s {
		if ns.Name == name {
			return ns.Schema, true
		}
	}

	return nil, false
}

// MarshalJSON implements json.Marshaler and preserves insertion order.
func (s Schemas) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")
	for i, ns := range s {
		if i != 0 {
			buf.WriteString(",")
		}

		key, err := json.Marshal(ns.Name)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteString(":")

		val, err := json.Marshal(ns.Schema)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}
