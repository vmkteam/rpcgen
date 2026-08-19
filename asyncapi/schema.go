package asyncapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/vmkteam/zenrpc/v2/smd"
)

const (
	contentTypeJSON = "application/json"
	componentPrefix = "#/components/schemas/"
	definitionsPref = "#/definitions/"
	errorSchemaName = "JsonRpcError"
)

// NewSchema converts SMD schema to AsyncAPI document. Every method of the schema becomes a send
// operation with a reply, every service of settings.Events becomes a receive operation.
func NewSchema(schema smd.Schema, settings Settings) (Document, error) {
	s := settings.withDefaults()

	b := &builder{
		settings:   s,
		channel:    NewReference("#/channels/" + s.Channel.Name),
		messages:   map[string]*Message{},
		operations: map[string]*Operation{},
		raw:        map[string][]byte{},
		visiting:   map[string]struct{}{},
		refs:       map[string]struct{}{},
	}

	// definitions are collected before conversion: an inline object may then be replaced by a
	// reference to a component of the same shape, see builder.jsonSchema
	b.collectSchemaDefinitions(schema)
	if s.Events != nil {
		b.collectSchemaDefinitions(*s.Events)

		for _, name := range sortedServices(*s.Events) {
			b.addEvent(name, s.Events.Services[name])
		}
	}

	var hasMethods bool
	for _, name := range sortedServices(schema) {
		if !s.MethodFilter(name) {
			continue
		}

		b.addMethod(name, schema.Services[name])
		hasMethods = true
	}

	if hasMethods {
		b.addComponent(errorSchemaName, jsonRPCErrorSchema())
	}

	if b.err != nil {
		return Document{}, b.err
	}

	if err := b.checkRefs(); err != nil {
		return Document{}, err
	}

	doc := Document{
		AsyncAPI: Version,
		Info: Info{
			Title:       s.Title,
			Version:     s.Version,
			Description: s.Description,
		},
		Servers: s.Servers,
		Channels: map[string]*Channel{s.Channel.Name: {
			Address:     s.Channel.Address,
			Title:       s.Channel.Title,
			Description: s.Channel.Description,
			Servers:     serverRefs(s.Servers),
			Bindings:    channelBindings(s.Channel),
			Messages:    b.messages,
		}},
		Operations: b.operations,
	}

	if used := b.usedComponents(doc); len(used) > 0 {
		doc.Components = &Components{Schemas: used}
	}

	return doc, nil
}

// usedComponents returns schemas the document actually references, sorted by name. Definitions are
// collected from the whole SMD, while MethodFilter may leave only a part of it in the document.
func (b *builder) usedComponents(doc Document) Schemas {
	roots := map[string]struct{}{}
	for _, channel := range doc.Channels {
		for _, message := range channel.Messages {
			schemaRefs(message.Payload, roots)
		}

		if channel.Bindings != nil && channel.Bindings.WebSocket != nil {
			schemaRefs(channel.Bindings.WebSocket.Query, roots)
			schemaRefs(channel.Bindings.WebSocket.Headers, roots)
		}
	}

	reachable := map[string]struct{}{}
	var visit func(string)
	visit = func(name string) {
		schema, ok := b.components.Get(name)
		if !ok {
			return
		}

		if _, seen := reachable[name]; seen {
			return
		}
		reachable[name] = struct{}{}

		refs := map[string]struct{}{}
		schemaRefs(schema, refs)

		for ref := range refs {
			visit(ref)
		}
	}

	for root := range roots {
		visit(root)
	}

	var used Schemas
	for _, ns := range b.components {
		if _, ok := reachable[ns.Name]; ok {
			used = append(used, ns)
		}
	}

	sort.Slice(used, func(i, j int) bool { return used[i].Name < used[j].Name })

	return used
}

// schemaRefs collects names of components the schema references.
func schemaRefs(schema *Schema, out map[string]struct{}) {
	if schema == nil {
		return
	}

	if name := strings.TrimPrefix(schema.Ref, componentPrefix); name != schema.Ref {
		out[name] = struct{}{}
	}

	for _, ns := range schema.Properties {
		schemaRefs(ns.Schema, out)
	}

	schemaRefs(schema.Items, out)

	for _, one := range schema.OneOf {
		schemaRefs(one, out)
	}
}

// builder accumulates the parts of a document while walking SMD schemas.
type builder struct {
	settings   Settings
	channel    Reference
	messages   map[string]*Message
	operations map[string]*Operation
	components Schemas

	raw      map[string][]byte   // marshalled components, to detect conflicting definitions
	visiting map[string]struct{} // definitions being converted right now, to break reference cycles
	refs     map[string]struct{} // referenced component names, to detect dangling refs
	err      error
}

func (b *builder) fail(err error) {
	if b.err == nil {
		b.err = err
	}
}

// addEvent adds a server to client notification.
func (b *builder) addEvent(name string, service smd.Service) {
	method := b.settings.EventMethodName(name)
	pas, k := pascal(method), key(method)

	var params *Schema
	if p := b.paramsSchema(service, pas+"Params"); p != nil {
		b.addComponent(pas+"Params", p)
		params = &Schema{Ref: b.componentRef(pas + "Params")}
	}

	summary, description := splitDescription(service.Description)

	b.messages[k] = &Message{
		Name:        method,
		Title:       method,
		Summary:     summary,
		Description: description,
		ContentType: contentTypeJSON,
		Payload:     notificationPayload(pas+"Frame", method, params),
	}

	b.operations["receive_"+k] = &Operation{
		Action:   "receive",
		Channel:  b.channel,
		Title:    method,
		Summary:  summary,
		Messages: []Reference{NewReference(b.messagePath(k))},
	}
}

// addMethod adds a client to server request and its reply.
func (b *builder) addMethod(name string, service smd.Service) {
	pas, k := pascal(name), key(name)
	respKey := k + "_response"

	summary, _ := splitDescription(service.Description)

	b.messages[k] = &Message{
		Name:        name,
		Title:       name + " request",
		Description: cleanDescription(service.Description),
		ContentType: contentTypeJSON,
		Payload:     requestPayload(pas+"Request", name, b.paramsSchema(service, pas+"Params")),
	}

	b.messages[respKey] = &Message{
		Name:        name + ".response",
		Title:       name + " response",
		Description: errorsTable(service.Errors),
		ContentType: contentTypeJSON,
		Payload:     responsePayload(pas+"Response", b.resultSchema(service, pas), b.componentRef(errorSchemaName)),
	}

	b.operations["send_"+k] = &Operation{
		Action:   "send",
		Channel:  b.channel,
		Title:    name,
		Summary:  summary,
		Messages: []Reference{NewReference(b.messagePath(k))},
		Reply:    &Reply{Channel: b.channel, Messages: []Reference{NewReference(b.messagePath(respKey))}},
	}
}

// paramsSchema builds params object from service parameters. Returns nil for a method without them.
func (b *builder) paramsSchema(service smd.Service, title string) *Schema {
	if len(service.Parameters) == 0 {
		return nil
	}

	sch := &Schema{Type: smd.Object, Title: title}
	var required []string

	for _, param := range service.Parameters {
		prop := b.jsonSchema(param, title+pascal(param.Name))
		if param.Optional {
			prop = nullable(prop)
		} else {
			required = append(required, param.Name)
		}

		sch.Properties.Add(param.Name, prop)
	}
	sch.Required = required

	return sch
}

// resultSchema builds the result field of a reply.
func (b *builder) resultSchema(service smd.Service, pascalName string) *Schema {
	if service.Returns.Type == "" {
		return &Schema{Type: "null", Description: "the method returns nothing"}
	}

	sch := b.jsonSchema(service.Returns, pascalName+"Result")
	if service.Returns.Optional {
		sch = nullable(sch)
	}

	return sch
}

// jsonSchema converts SMD schema to JSON Schema, hoisting definitions to components.
func (b *builder) jsonSchema(schema smd.JSONSchema, fallbackTitle string) *Schema {
	b.collectDefinitions(schema.Definitions)

	var sch *Schema
	switch schema.Type {
	case smd.Object:
		// a named object without properties is described in full somewhere else in the schema
		if len(schema.Properties) == 0 && b.hasComponent(schema.TypeName) {
			return &Schema{Ref: b.componentRef(schema.TypeName)}
		}

		title := schema.TypeName
		if title == "" {
			title = fallbackTitle
		}

		sch = b.propertiesSchema(schema.Properties, title)
		if ref := b.componentOf(schema.TypeName, sch); ref != nil {
			return ref
		}
	case smd.Array:
		sch = &Schema{Type: smd.Array, Items: b.itemsSchema(schema.Items)}
	default:
		sch = &Schema{Type: schema.Type}
	}

	if d := prose(schema.Description, schema.Type); d != "" {
		sch.Description = d
	}

	return sch
}

// propertiesSchema converts a flat SMD property list to an object schema.
func (b *builder) propertiesSchema(props smd.PropertyList, title string) *Schema {
	sch := &Schema{Type: smd.Object, Title: title}
	var required []string

	for _, prop := range props {
		b.collectDefinitions(prop.Definitions)

		p := b.propertySchema(prop)
		if prop.Optional {
			p = nullable(p)
		} else {
			required = append(required, prop.Name)
		}

		sch.Properties.Add(prop.Name, p)
	}
	sch.Required = required

	return sch
}

func (b *builder) propertySchema(prop smd.Property) *Schema {
	switch {
	case prop.Ref != "":
		// json schema draft-07 ignores keywords next to $ref, so a description here would be lost
		return &Schema{Ref: b.refName(prop.Ref)}
	case prop.Type == smd.Array:
		return &Schema{Type: smd.Array, Description: prose(prop.Description, prop.Type), Items: b.itemsSchema(prop.Items)}
	default:
		return &Schema{Type: prop.Type, Description: prose(prop.Description, prop.Type)}
	}
}

func (b *builder) itemsSchema(items map[string]string) *Schema {
	if ref := items["$ref"]; ref != "" {
		return &Schema{Ref: b.refName(ref)}
	}

	if typ := items["type"]; typ != "" {
		return &Schema{Type: typ}
	}

	return &Schema{}
}

// collectSchemaDefinitions hoists definitions of every service of the schema to components.
func (b *builder) collectSchemaDefinitions(schema smd.Schema) {
	for _, name := range sortedServices(schema) {
		service := schema.Services[name]

		for _, param := range service.Parameters {
			b.collectDefinitions(param.Definitions)
		}
		b.collectDefinitions(service.Returns.Definitions)
	}
}

// componentOf returns a reference to an already collected component of the same shape. SMD spells
// a type out in full at every usage, so without this the same type is both a component (from a
// definition of another method) and an inline copy. Types with the same name but a different shape
// stay inline: zenrpc allows them and a wrong shared schema is worse than a repeated one.
func (b *builder) componentOf(typeName string, schema *Schema) *Schema {
	if typeName == "" {
		return nil
	}

	component, ok := b.raw[typeName]
	if !ok {
		return nil
	}

	probe := *schema
	probe.Title = ""

	data, err := json.Marshal(&probe)
	if err != nil || !bytes.Equal(component, data) {
		return nil
	}

	return &Schema{Ref: b.componentRef(typeName)}
}

// collectDefinitions hoists SMD definitions to components.schemas.
func (b *builder) collectDefinitions(defs map[string]smd.Definition) {
	names := make([]string, 0, len(defs))
	for name := range defs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if _, ok := b.visiting[name]; ok {
			continue
		}

		b.visiting[name] = struct{}{}
		b.addComponent(name, b.propertiesSchema(defs[name].Properties, ""))
		delete(b.visiting, name)
	}
}

// addComponent registers a component schema. The same name with a different shape is an error:
// silently keeping the first one would produce a document that lies about half of its usages.
func (b *builder) addComponent(name string, schema *Schema) {
	data, err := json.Marshal(schema)
	if err != nil {
		b.fail(fmt.Errorf("marshal schema %s: %w", name, err))
		return
	}

	if prev, ok := b.raw[name]; ok {
		if !bytes.Equal(prev, data) {
			b.fail(fmt.Errorf("asyncapi: conflicting definitions of %q", name))
		}

		return
	}

	b.raw[name] = data
	b.components.Add(name, schema)
}

func (b *builder) hasComponent(name string) bool {
	if name == "" {
		return false
	}

	_, ok := b.raw[name]

	return ok
}

func (b *builder) componentRef(name string) string {
	b.refs[name] = struct{}{}

	return componentPrefix + name
}

func (b *builder) refName(ref string) string {
	return b.componentRef(strings.TrimPrefix(ref, definitionsPref))
}

func (b *builder) messagePath(key string) string {
	return b.channel.Ref + "/messages/" + key
}

// checkRefs reports references to components that were never added.
func (b *builder) checkRefs() error {
	var missing []string
	for name := range b.refs {
		if _, ok := b.raw[name]; !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) == 0 {
		return nil
	}
	sort.Strings(missing)

	return fmt.Errorf("asyncapi: missing schemas for refs: %s", strings.Join(missing, ", "))
}

// notificationPayload builds a JSON-RPC 2.0 notification frame: a request without id.
func notificationPayload(title, method string, params *Schema) *Schema {
	payload := &Schema{Type: smd.Object, Title: title}
	payload.Properties.Add("jsonrpc", &Schema{Type: smd.String, Const: "2.0"})
	payload.Properties.Add("method", &Schema{Type: smd.String, Const: method})
	payload.Required = []string{"jsonrpc", "method"}

	if params != nil {
		payload.Properties.Add("params", params)
		payload.Required = append(payload.Required, "params")
	}
	payload.AdditionalProperties = falseValue()

	return payload
}

// requestPayload builds a JSON-RPC 2.0 request frame.
func requestPayload(title, method string, params *Schema) *Schema {
	payload := &Schema{Type: smd.Object, Title: title}
	payload.Properties.Add("jsonrpc", &Schema{Type: smd.String, Const: "2.0"})
	payload.Properties.Add("id", idSchema())
	payload.Properties.Add("method", &Schema{Type: smd.String, Const: method})
	payload.Required = []string{"jsonrpc", "id", "method"}

	if params != nil {
		payload.Properties.Add("params", params)
		payload.Required = append(payload.Required, "params")
	}
	payload.AdditionalProperties = falseValue()

	return payload
}

// responsePayload builds a JSON-RPC 2.0 response frame.
func responsePayload(title string, result *Schema, errorRef string) *Schema {
	payload := &Schema{
		Type:        smd.Object,
		Title:       title,
		Description: "Exactly one of `result` and `error` is present in the frame.",
	}
	payload.Properties.Add("jsonrpc", &Schema{Type: smd.String, Const: "2.0"})
	payload.Properties.Add("id", idSchema())
	payload.Properties.Add("result", result)
	payload.Properties.Add("error", &Schema{Ref: errorRef})
	payload.Required = []string{"jsonrpc", "id"}
	payload.AdditionalProperties = falseValue()

	return payload
}

func idSchema() *Schema {
	return &Schema{
		Type:        []string{smd.Integer, smd.String},
		Description: "request id; the reply arrives in a separate frame with the same id",
	}
}

func jsonRPCErrorSchema() *Schema {
	sch := &Schema{
		Type:        smd.Object,
		Title:       errorSchemaName,
		Description: "JSON-RPC 2.0 error. Method specific codes are listed in the reply description.",
	}
	sch.Properties.Add("code", &Schema{Type: smd.Integer})
	sch.Properties.Add("message", &Schema{Type: smd.String})
	sch.Properties.Add("data", &Schema{Description: "optional error details"})
	sch.Required = []string{"code", "message"}

	return sch
}

// nullable allows null in schema: an optional value in zenrpc is a Go pointer which is serialized
// as null instead of being omitted.
func nullable(schema *Schema) *Schema {
	if schema.Ref != "" || len(schema.OneOf) > 0 {
		return &Schema{OneOf: []*Schema{schema, {Type: "null"}}}
	}

	if typ, ok := schema.Type.(string); ok && typ != "" {
		out := *schema
		out.Type = []string{typ, "null"}

		return &out
	}

	return schema
}

// errorsTable renders SMD errors as a markdown table.
func errorsTable(errors map[int]string) string {
	if len(errors) == 0 {
		return ""
	}

	codes := make([]int, 0, len(errors))
	for code := range errors {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	var sb strings.Builder
	sb.WriteString("Method errors:\n\n| code | message |\n|---|---|")
	for _, code := range codes {
		fmt.Fprintf(&sb, "\n| `%d` | %s |", code, errors[code])
	}

	return sb.String()
}

func channelBindings(channel ChannelSettings) *ChannelBindings {
	if channel.Query == nil && channel.Headers == nil {
		return nil
	}

	return &ChannelBindings{WebSocket: &WebSocketBinding{
		Method:         "GET",
		Query:          channel.Query,
		Headers:        channel.Headers,
		BindingVersion: WebSocketBindingVersion,
	}}
}

func serverRefs(servers map[string]Server) []Reference {
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)

	refs := make([]Reference, 0, len(names))
	for _, name := range names {
		refs = append(refs, NewReference("#/servers/"+name))
	}

	return refs
}

func sortedServices(schema smd.Schema) []string {
	names := make([]string, 0, len(schema.Services))
	for name := range schema.Services {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

// splitDescription returns the first line of a doc comment as a summary and the whole comment as a
// description when there is more than one line.
func splitDescription(description string) (summary, full string) {
	description = cleanDescription(description)
	if description == "" {
		return "", ""
	}

	first, _, multiline := strings.Cut(description, "\n")
	if !multiline {
		return description, ""
	}

	return strings.TrimSpace(first), description
}

// cleanDescription drops linter directives that go doc comments carry into SMD.
func cleanDescription(description string) string {
	lines := strings.Split(description, "\n")

	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "nolint") {
			continue
		}

		kept = append(kept, line)
	}

	return strings.TrimSpace(strings.Join(kept, "\n"))
}

// prose drops descriptions that zenrpc uses to carry a type name instead of human text.
func prose(description, typ string) string {
	if smd.IsSMDTypeName(description, typ) {
		return ""
	}

	return description
}

// key converts an SMD method name to an AsyncAPI object key: only letters, digits, _ and - allowed.
func key(name string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			return r
		}

		return '_'
	}, name)
}

// pascal converts an SMD method name to a schema title: "ws.newMessage" becomes "WsNewMessage".
func pascal(name string) string {
	var sb strings.Builder
	upper := true

	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			upper = true
			continue
		}

		if upper {
			sb.WriteRune(unicode.ToUpper(r))
			upper = false
			continue
		}

		sb.WriteRune(r)
	}

	return sb.String()
}

func falseValue() *bool {
	v := false

	return &v
}
