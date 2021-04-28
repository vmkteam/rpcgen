package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vmkteam/zenrpc/v2/smd"
)

// NewSchema convert smd.Schema to Schema
func NewSchema(schema smd.Schema) Schema {
	var res Schema

	for _, namespaceName := range getNamspaceNames(schema) {
		res.Namespaces = append(res.Namespaces, Namespace{
			Name:    namespaceName,
			Methods: getNamespaceMethods(schema, namespaceName),
		})
	}

	sort.Slice(res.Namespaces, func(i, j int) bool {
		return res.Namespaces[i].Name < res.Namespaces[j].Name
	})

	return res
}

type Schema struct {
	Namespaces []Namespace
}

// NamespaceNames return all namespace names.
func (s Schema) NamespaceNames() (res []string) {
	for _, ns := range s.Namespaces {
		res = append(res, ns.Name)
	}

	return res
}

// NamespaceMethodNames return all methodNames by namespace name.
func (s Schema) NamespaceMethodNames(namespace string) (res []string) {
	ns := s.NamespaceByName(namespace)

	for _, method := range ns.Methods {
		res = append(res, method.Name)
	}

	sort.Strings(res)

	return
}

// Models return all models from all methods
func (s Schema) Models() (res []Model) {
	for _, ns := range s.Namespaces {
		for _, method := range ns.Methods {
			res = append(res, method.Models...)
		}
	}

	res = cleanModelList(res, "", "")

	// sort model fields
	for k := range res {
		k := k
		sort.Slice(res[k].Fields, func(i, j int) bool { return res[k].Fields[i].Name < res[k].Fields[j].Name })
	}

	// sort models
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res
}

func (s Schema) NamespaceByName(namespace string) Namespace {
	for _, ns := range s.Namespaces {
		if ns.Name == namespace {
			return ns
		}
	}

	return Namespace{}
}

func (s Schema) MethodByName(namespace, methodName string) Method {
	ns := s.NamespaceByName(namespace)

	for _, method := range ns.Methods {
		if method.Name == methodName {
			return method
		}
	}

	return Method{}
}

type Namespace struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name string

	// Params contains all method arguments
	Params []Value

	// RPC requires single one return value
	Returns *Value

	// Models contains all models which linked from Params or Return values by ModelName
	Models []Model

	// Docs fields

	// Description describes what the method does
	Description string

	// Typed errors
	Errors []Error
}

func (m Method) HasResult() bool {
	return m.Returns != nil
}

// CommentDescription add to head of all lines two slashes
func (m Method) CommentDescription() string {
	if len(m.Description) == 0 {
		return ""
	}

	return "// " + strings.ReplaceAll(m.Description, "\n", "\n// ")
}

type Error struct {
	Code    int
	Message string
}

type Model struct {
	Name        string
	Description string

	Fields []Value
}

type Value struct {
	Name        string
	Type        string
	Optional    bool
	Description string

	ArrayItemType string

	ModelName string
}

// GoType convert Value jsType to golang type
func (v *Value) GoType() string {
	if v == nil {
		return ""
	}

	if v.Type == smd.Array {
		if v.ModelName != "" {
			return fmt.Sprintf("[]%s", v.ModelName)
		}

		return fmt.Sprintf("[]%s", simpleGoType(v.ArrayItemType))
	} else if v.Type == smd.Object {
		return v.ModelName
	}

	return simpleGoType(v.Type)
}

func simpleGoType(jsType string) string {
	switch jsType {
	case "boolean":
		return "bool"
	case "number":
		return "float64"
	case "integer":
		return "int"
	default:
		return jsType
	}
}

// getNamspaceNames return all namspace names from schema.
func getNamspaceNames(schema smd.Schema) (res []string) {
	m := map[string]struct{}{}
	for name := range schema.Services {
		m[getNamespace(name)] = struct{}{}
	}

	for name := range m {
		res = append(res, name)
	}

	sort.Strings(res)

	return res
}

// getNamespaceMethods return all namespace methods
func getNamespaceMethods(schema smd.Schema, namespace string) (res []Method) {
	for name, service := range schema.Services {
		if getNamespace(name) == namespace {
			res = append(res, newMethod(service, getNamespace(name), getMethodName(name)))
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res
}

// getNamespace extract namespace from string like "namespace.method"
func getNamespace(methodName string) string {
	return strings.Split(methodName, ".")[0]
}

// newMethod convert smd.Service method to Method
func newMethod(service smd.Service, namespace, methodName string) Method {
	var args []Value
	for _, arg := range service.Parameters {
		args = append(args, newValue(arg))
	}

	var methodErrors []Error
	for code, message := range service.Errors {
		methodErrors = append(methodErrors, Error{
			Code:    code,
			Message: message,
		})
	}

	// sort for prevent smd's errors map shuffle
	sort.SliceStable(methodErrors, func(i, j int) bool {
		return methodErrors[i].Code < methodErrors[j].Code
	})

	method := Method{
		Name:        methodName,
		Description: service.Description,
		Params:      args,
		Models:      newMethodModels(service, namespace, methodName),
		Errors:      methodErrors,
	}

	if service.Returns.Type != "" {
		method.Returns = newValuePointer(service.Returns)

		// fix return model name
		if method.Returns.Type == smd.Object && method.Returns.ModelName == "" {
			method.Returns.ModelName = fmt.Sprintf("%s%sResponse", titleFirstLetter(namespace), titleFirstLetter(methodName))
		}
	}

	return method
}

// newValue convert smd.JSONSchema to value
func newValue(in smd.JSONSchema) Value {
	value := Value{
		Name:     in.Name,
		Type:     in.Type,
		Optional: in.Optional,
	}
	// in is simple builtin type

	// in is array
	if in.Type == smd.Array {
		if in.Items["type"] != "" { // simple type
			value.ArrayItemType = in.Items["type"]
		} else { // complex type
			value.ArrayItemType = "object"
			value.ModelName = strings.TrimPrefix(in.Items["$ref"], definitionsPrefix)
		}
	}

	// in is object
	if in.Type == smd.Object {
		value.ModelName = in.Description
	} else {
		value.Description = in.Description
	}

	return value
}

func newValuePointer(in smd.JSONSchema) *Value {
	v := newValue(in)

	return &v
}

func newValueFromProp(in smd.Property) Value {
	value := Value{
		Name:     in.Name,
		Type:     in.Type,
		Optional: in.Optional,
	}

	// in is an object
	if in.Type == smd.Object {
		value.ModelName = in.Description

		if in.Ref != "" {
			value.ModelName = strings.TrimPrefix(in.Ref, definitionsPrefix)
		}
	} else {
		value.Description = in.Description
	}

	// in is an array
	if in.Type == smd.Array {
		if in.Items["type"] != "" { // simple type
			value.ArrayItemType = in.Items["type"]
		} else { // complex type
			value.ArrayItemType = "object"
			value.ModelName = strings.TrimPrefix(in.Items["$ref"], definitionsPrefix)
		}
	}

	return value
}

// newMethodModels scrape all method models
//
// Model will be created from:
// - method has argument or return with type `array` with $ref to definition. Name by definition name.
// - method has argument or return with type `object`. Name if model gather from Description field.
// - method return object has property with type `object` and $ref to definition.
// - definition has property with type `object` and $ref to definition.
//
// NOTE: definitions limited to smd.JSONSchema level, but in our Schema
// the definition will be raised to a higher level in the scope of the method.
func newMethodModels(method smd.Service, namespace, methodName string) (res []Model) {
	// get models from params
	for _, arg := range method.Parameters {
		res = append(res, newModels(arg)...)
	}

	// get models from return
	res = append(res, newModels(method.Returns)...)

	return cleanModelList(res, namespace, methodName)
}

// cleanModelList remove duplicates, add NamespaceMethod prefix for each model.
func cleanModelList(models []Model, namespace, methodName string) (res []Model) {
	var dup = map[string]struct{}{}

	for _, model := range models {
		// fix empty model name
		// is always response because argument required named params.
		if model.Name == "" {
			model.Name = fmt.Sprintf("%s%sResponse", titleFirstLetter(namespace), titleFirstLetter(methodName))
		}

		// filter out duplicates
		if _, has := dup[model.Name]; has {
			continue
		}

		dup[model.Name] = struct{}{}
		res = append(res, model)
	}

	return res
}

// newModels convert
func newModels(in smd.JSONSchema) (res []Model) {
	// definitions
	for name, def := range in.Definitions {
		res = append(res, convertDefinitionToModel(def, name))
	}

	// object args and return
	if in.Type == smd.Object {
		var values []Value

		for _, prop := range in.Properties {
			values = append(values, newValueFromProp(prop))
		}

		name := in.Name
		if in.Description != "" {
			name = in.Description
		}

		res = append(res, Model{
			Name:   name,
			Fields: values,
		})
	}

	return res
}

// convertDefinitionToModel convert smd definition to Object.
// Definition is always an object.
func convertDefinitionToModel(def smd.Definition, name string) Model {
	model := Model{
		Name: name,
	}

	for _, property := range def.Properties {
		model.Fields = append(model.Fields, newValueFromProp(property))
	}

	return model
}
