package asyncapi

import (
	"strings"
	"unicode"

	"github.com/vmkteam/zenrpc/v2/smd"
)

const (
	defaultTitle       = "json-rpc 2.0 websocket api"
	defaultVersion     = "1.0.0"
	defaultChannelName = "main"
	defaultAddress     = "/"
)

// Settings holds everything that can not be taken from SMD: transport description and the list of
// server to client events.
type Settings struct {
	// Title of the API. Defaults to "json-rpc 2.0 websocket api".
	Title string

	// Version of the API. Defaults to "1.0.0".
	Version string

	// Description of the API in CommonMark. A good place for transport notes: handshake,
	// authorization, connection limits.
	Description string

	// Servers of the API, keyed by server name.
	Servers map[string]Server

	// Channel describes the websocket endpoint itself.
	Channel ChannelSettings

	// Events is an optional schema of server to client notifications. Every service of this schema
	// becomes a receive operation: a JSON-RPC notification (a request without id) with method name
	// built by EventMethodName and params built from service parameters.
	//
	// Such a schema is usually produced by a documentation-only zenrpc server: register a service
	// with a method per event and take its SMD without exposing the service to clients.
	Events *smd.Schema

	// EventMethodName maps an SMD service name to the method name on the wire. Defaults to
	// LowerCamelMethod: "ws.NewMessage" becomes "ws.newMessage".
	EventMethodName func(string) string

	// MethodFilter reports whether an SMD method should be exposed as a send operation. Nil filter
	// keeps every method of the schema.
	MethodFilter func(string) bool
}

// ChannelSettings describes a websocket channel and its handshake.
type ChannelSettings struct {
	// Name of the channel, used as a key in the channels object. Defaults to "main".
	Name string

	// Address is a path of the websocket endpoint, e.g. "/v1/rpc/ws". Defaults to "/".
	Address string

	// Title of the channel. Defaults to Address.
	Title string

	// Description of the channel in CommonMark.
	Description string

	// Query is a schema of the query string of the handshake request.
	Query *Schema

	// Headers is a schema of the headers of the handshake request.
	Headers *Schema
}

// withDefaults returns settings with empty values replaced by defaults.
func (s Settings) withDefaults() Settings {
	if s.Title == "" {
		s.Title = defaultTitle
	}
	if s.Version == "" {
		s.Version = defaultVersion
	}
	if s.Channel.Name == "" {
		s.Channel.Name = defaultChannelName
	}
	if s.Channel.Address == "" {
		s.Channel.Address = defaultAddress
	}
	if s.Channel.Title == "" {
		s.Channel.Title = s.Channel.Address
	}
	if s.EventMethodName == nil {
		s.EventMethodName = LowerCamelMethod
	}
	if s.MethodFilter == nil {
		s.MethodFilter = func(string) bool { return true }
	}

	return s
}

// LowerCamelMethod lowers the first letter of the method part of an SMD service name: SMD keeps Go
// method names ("ws.NewMessage") while zenrpc dispatches case-insensitively and hand-written
// notifications usually use lower camel case ("ws.newMessage").
func LowerCamelMethod(name string) string {
	idx := strings.LastIndex(name, ".")
	if idx == -1 || idx == len(name)-1 {
		return lowerFirst(name)
	}

	return name[:idx+1] + lowerFirst(name[idx+1:])
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}
