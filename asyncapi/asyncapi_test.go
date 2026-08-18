package asyncapi

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/vmkteam/zenrpc/v2"
	"github.com/vmkteam/zenrpc/v2/smd"
	"github.com/vmkteam/zenrpc/v2/testdata"
)

const asyncAPIFile = "./testdata/asyncapi.json"

var update = flag.Bool("update", false, "update golden files")

func testSettings(events smd.Schema) Settings {
	return Settings{
		Title:       "test",
		Description: "Test websocket api.",
		Servers: map[string]Server{
			"local": {Host: "localhost:8080", Protocol: "ws", Description: "local server"},
		},
		Channel: ChannelSettings{
			Name:        "main",
			Address:     "/ws",
			Description: "Test channel.",
			Query: &Schema{
				Type:       smd.Object,
				Required:   []string{"id"},
				Properties: Schemas{{Name: "id", Schema: &Schema{Type: smd.Integer}}},
			},
			Headers: &Schema{
				Type:       smd.Object,
				Properties: Schemas{{Name: "Authorization", Schema: &Schema{Type: smd.String}}},
			},
		},
		Events: &events,
	}
}

func TestGenerateAsyncAPISchema(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})
	rpc.Register("arith", testdata.ArithService{})

	events := zenrpc.NewServer(zenrpc.Options{})
	events.Register("ws", testdata.PrintService{})

	cl := NewClient(rpc.SMD(), testSettings(events.SMD()))

	generated, err := cl.Generate()
	if err != nil {
		t.Fatalf("generate asyncapi schema: %v", err)
	}

	if *update {
		if err := os.WriteFile(asyncAPIFile, generated, 0o600); err != nil {
			t.Fatalf("write golden file: %v", err)
		}
	}

	testData, err := os.ReadFile(asyncAPIFile)
	if err != nil {
		t.Fatalf("open test data file: %v", err)
	}

	if !bytes.Equal(generated, testData) {
		t.Fatalf("bad generator output")
	}

	// second run
	generated, err = cl.Generate()
	if err != nil {
		t.Fatalf("generate asyncapi schema: %v", err)
	}

	if !bytes.Equal(generated, testData) {
		t.Fatalf("bad generator output")
	}
}

// TestEventMethodNames checks that notifications are named as they are sent on the wire, not as
// they are named in Go.
func TestEventMethodNames(t *testing.T) {
	events := zenrpc.NewServer(zenrpc.Options{})
	events.Register("ws", testdata.PrintService{})

	doc, err := NewSchema(smd.Schema{}, testSettings(events.SMD()))
	if err != nil {
		t.Fatalf("generate asyncapi schema: %v", err)
	}

	msg, ok := doc.Channels["main"].Messages["ws_printOptional"]
	if !ok {
		t.Fatalf("message ws_printOptional not found")
	}

	if msg.Name != "ws.printOptional" {
		t.Fatalf("got message name %q, want ws.printOptional", msg.Name)
	}

	if _, ok := doc.Operations["receive_ws_printOptional"]; !ok {
		t.Fatalf("operation receive_ws_printOptional not found")
	}
}

// TestOptionalIsNullable checks that an optional parameter is both absent in required and allowed
// to be null: in zenrpc it is a Go pointer serialized as null.
func TestOptionalIsNullable(t *testing.T) {
	events := zenrpc.NewServer(zenrpc.Options{})
	events.Register("ws", testdata.PrintService{})

	doc, err := NewSchema(smd.Schema{}, testSettings(events.SMD()))
	if err != nil {
		t.Fatalf("generate asyncapi schema: %v", err)
	}

	params, ok := doc.Components.Schemas.Get("WsPrintOptionalParams")
	if !ok {
		t.Fatalf("schema WsPrintOptionalParams not found")
	}

	if len(params.Required) != 0 {
		t.Fatalf("got required %v, want none", params.Required)
	}

	prop, ok := params.Properties.Get("s")
	if !ok {
		t.Fatalf("property s not found")
	}

	data, err := json.Marshal(prop.Type)
	if err != nil {
		t.Fatalf("marshal type: %v", err)
	}

	if string(data) != `["string","null"]` {
		t.Fatalf("got type %s, want [\"string\",\"null\"]", data)
	}
}

// TestSharedTypeIsReferenced checks that a type described in full both in a method and in an event
// lands in components once and is referenced from both places.
func TestSharedTypeIsReferenced(t *testing.T) {
	item := smd.JSONSchema{
		Name: "item", Type: smd.Object, TypeName: "Item",
		Properties:  smd.PropertyList{{Name: "id", Type: smd.Integer}},
		Definitions: map[string]smd.Definition{"Item": {Type: smd.Object, Properties: smd.PropertyList{{Name: "id", Type: smd.Integer}}}},
	}

	events := smd.Schema{Services: map[string]smd.Service{
		"ws.ItemChanged": {Parameters: []smd.JSONSchema{item}},
	}}
	schema := smd.Schema{Services: map[string]smd.Service{
		"a.SaveItem": {Parameters: []smd.JSONSchema{item}},
	}}

	doc, err := NewSchema(schema, Settings{Events: &events})
	if err != nil {
		t.Fatalf("generate asyncapi schema: %v", err)
	}

	params, ok := doc.Components.Schemas.Get("WsItemChangedParams")
	if !ok {
		t.Fatalf("schema WsItemChangedParams not found")
	}

	prop, ok := params.Properties.Get("item")
	if !ok {
		t.Fatalf("property item not found")
	}

	if prop.Ref != "#/components/schemas/Item" {
		t.Fatalf("got ref %q, want #/components/schemas/Item", prop.Ref)
	}
}

// TestUnusedComponentsAreDropped checks that filtering methods out also drops the schemas that
// only they used: definitions are collected from the whole SMD.
func TestUnusedComponentsAreDropped(t *testing.T) {
	rpc := zenrpc.NewServer(zenrpc.Options{})
	rpc.Register("catalogue", testdata.CatalogueService{})

	events := zenrpc.NewServer(zenrpc.Options{})
	events.Register("ws", testdata.PrintService{})

	full, err := NewSchema(rpc.SMD(), testSettings(events.SMD()))
	if err != nil {
		t.Fatalf("generate asyncapi schema: %v", err)
	}

	settings := testSettings(events.SMD())
	settings.MethodFilter = func(string) bool { return false }

	frames, err := NewSchema(rpc.SMD(), settings)
	if err != nil {
		t.Fatalf("generate asyncapi schema: %v", err)
	}

	if _, ok := full.Components.Schemas.Get("Campaign"); !ok {
		t.Fatalf("schema Campaign not found in the full document")
	}

	if frames.Components != nil {
		if _, ok := frames.Components.Schemas.Get("Campaign"); ok {
			t.Fatalf("schema Campaign is kept while nothing references it")
		}
	}
}

// TestConflictingDefinitions checks that the same type name with a different shape fails the
// generation instead of silently ending up in the document once.
func TestConflictingDefinitions(t *testing.T) {
	def := func(propName string) map[string]smd.Definition {
		return map[string]smd.Definition{"Item": {
			Type:       smd.Object,
			Properties: smd.PropertyList{{Name: propName, Type: smd.String}},
		}}
	}

	schema := smd.Schema{Services: map[string]smd.Service{
		"a.First": {Parameters: []smd.JSONSchema{{
			Name: "item", Type: smd.Object, Definitions: def("one"),
			Properties: smd.PropertyList{{Name: "item", Type: smd.Object, Ref: "#/definitions/Item"}},
		}}},
		"a.Second": {Parameters: []smd.JSONSchema{{
			Name: "item", Type: smd.Object, Definitions: def("two"),
			Properties: smd.PropertyList{{Name: "item", Type: smd.Object, Ref: "#/definitions/Item"}},
		}}},
	}}

	_, err := NewSchema(schema, Settings{})
	if err == nil {
		t.Fatalf("got no error, want a conflict")
	}

	if !strings.Contains(err.Error(), `conflicting definitions of "Item"`) {
		t.Fatalf("got error %q, want a conflict of Item", err)
	}
}

// TestMissingRef checks that a reference to a type that is not described in SMD fails the
// generation instead of producing a document with a dangling ref.
func TestMissingRef(t *testing.T) {
	schema := smd.Schema{Services: map[string]smd.Service{
		"a.First": {Parameters: []smd.JSONSchema{{
			Name: "item", Type: smd.Object,
			Properties: smd.PropertyList{{Name: "item", Type: smd.Object, Ref: "#/definitions/Unknown"}},
		}}},
	}}

	_, err := NewSchema(schema, Settings{})
	if err == nil {
		t.Fatalf("got no error, want a missing ref")
	}

	if !strings.Contains(err.Error(), "Unknown") {
		t.Fatalf("got error %q, want a missing ref of Unknown", err)
	}
}

func TestLowerCamelMethod(t *testing.T) {
	for in, want := range map[string]string{
		"ws.NewMessage": "ws.newMessage",
		"NewMessage":    "newMessage",
		"ws.":           "ws.",
		"":              "",
	} {
		if got := LowerCamelMethod(in); got != want {
			t.Fatalf("LowerCamelMethod(%q) = %q, want %q", in, got, want)
		}
	}
}
