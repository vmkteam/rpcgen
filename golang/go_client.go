package golang

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"

	"github.com/semrush/zenrpc/v2/smd"
)

const (
	definitionsPrefix = "#/definitions/"

	typeArray   = "array"
	typeInteger = "integer"
	typeNumber  = "number"
	typeObject  = "object"
	typeInt     = "int"
)

// Generator main package structure
type Generator struct {
	smd smd.Schema
}

// NewClient create Generator from zenrpc/v2 SMD.
func NewClient(schema smd.Schema) *Generator {
	return &Generator{smd: schema}
}

// Generate returns generated Go client.
func (g *Generator) Generate() ([]byte, error) {
	tmpl, err := template.New("").Funcs(templateFuncs).Parse(goTpl)
	if err != nil {
		return nil, err
	}

	// compile template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g); err != nil {
		return nil, err
	}

	//return buf.Bytes(), nil
	return format.Source(buf.Bytes())
}

// goMethod represent RPC method.
type goMethod struct {
	Params []goValue
	Result *goValue
}

// HasResult returns true if RPC method has return value. Used in template.
func (m goMethod) HasResult() bool {
	return m.Result != nil
}

// HasComplexResult returns true if RPC method result returns object or array. Used in template.
func (m goMethod) HasComplexResult() bool {
	return m.Result != nil && !m.SimpleResult()
}

// SimpleResult returns true if RPC method result is simple builtin type. Used in template.
func (m goMethod) SimpleResult() bool {
	return m.Result != nil && (m.Result.GoType() == typeInt || m.Result.GoType() == "string" || m.Result.GoType() == "bool")
}

// goModel represent object (param or result in RPC)
type goModel struct {
	Name   string
	Fields []goValue
}

// goValue represent object or simple type (param or result in RPC)
type goValue struct {
	Name            string
	Type            string
	ItemType        string
	ModelName       string
	ParentModelName string
	Comment         string
	Optional        bool

	Fields []goValue
}

// GoType convert schema type to Go type and return their literal
func (p goValue) GoType() string {
	switch p.Type {
	case typeInteger:
		return typeInt
	case typeNumber:
		return "float64"
	case typeObject:
		if p.ModelName != "" {
			return p.ModelName
		}

		return strings.Title(p.Name)
	case typeArray:
		return "[]" + goValue{Type: p.ItemType}.GoType()
	case "boolean":
		return "bool"
	}

	return p.Type
}

// GoItemType
func (p goValue) GoItemType() string {
	switch p.Type {
	case typeInteger:
		return typeInt
	case typeNumber:
		return "float64"
	case typeObject:
		return p.ModelName
	case typeArray:
		return goValue{Type: p.ItemType}.GoType()
	}

	return p.Type
}

// getNamespace extract namespace from string like "namespace.method"
func getNamespace(methodName string) string {
	return strings.Split(methodName, ".")[0]
}

// getMethodName extract method from string like "namespace.method"
func getMethodName(methodName string) string {
	const methodNameLen = 2

	arr := strings.Split(methodName, ".")
	if len(arr) != methodNameLen {
		return methodName
	}

	return arr[1]
}

// Models convert schema models to goModels. Used in template, should be exported.
func (g *Generator) Models() []goModel {
	return g.models()
}

// models convert smd1 v1 schema to goModels.
func (g *Generator) models() (res []goModel) {
	for serviceName, s := range g.smd.Services {
		for _, p := range s.Parameters {
			if isModel(p) {
				res = addModel(res, "", p.Name, p)
			}
		}

		if isModel(s.Returns) {
			arr := strings.Split(serviceName, ".")
			for i, v := range arr {
				arr[i] = strings.Title(v)
			}

			res = addModel(res, "", strings.Join(arr, "")+"Response", s.Returns)
		}
	}

	m := map[string]goModel{}
	for _, model := range res {
		m[model.Name] = model
	}

	res = res[0:0]
	for _, model := range m {
		res = append(res, model)
	}

	for k := range res {
		k := k
		sort.Slice(res[k].Fields, func(i, j int) bool { return res[k].Fields[i].Name < res[k].Fields[j].Name })
	}

	sort.Slice(res, func(i, j int) bool { return res[i].Name < res[j].Name })

	return res
}

// addModelProperty convert prop to goModel and add it to allModels. Return allModels.
func addModelProperty(allModels []goModel, parent, name string, prop smd.Property) []goModel {
	if !isModelProperty(prop) {
		return allModels
	}

	return addModel(allModels, parent, strings.Title(name), smd.JSONSchema{
		Name:        strings.Title(name),
		Type:        prop.Type,
		Description: prop.Description,
		Properties:  prop.Definitions[name].Properties,
		Definitions: prop.Definitions,
		Items:       prop.Items,
	})
}

// addModel convert smdModel to goModel and add it to allModels. Return allModels.
func addModel(allModels []goModel, parentName, modelName string, smdModel smd.JSONSchema) []goModel {
	if !isModel(smdModel) {
		return allModels
	}

	// check duplicates
	for _, v := range allModels {
		if v.Name == strings.Title(modelName) {
			return allModels
		}
	}

	for name, def := range smdModel.Definitions {
		allModels = addModel(allModels, modelName, strings.Title(name), smd.JSONSchema{
			Name:        name,
			Type:        name,
			Properties:  def.Properties,
			Definitions: smdModel.Definitions,
			Items:       smdModel.Items,
		})
	}

	switch smdModel.Type {
	case typeObject:
		smdModel.Name = modelName
		smdModel.Description = modelName
		allModels = append(allModels, objectToModel(smdModel))
		// add recursive fields
		for name, field := range smdModel.Properties {
			if isModelProperty(field) {
				if field.Ref != "" {
					name = strings.TrimPrefix(field.Ref, definitionsPrefix)
				} else {
					name = fmt.Sprintf("%s%s", modelName, strings.Title(name))
				}

				field.Definitions = smdModel.Definitions
				allModels = addModelProperty(allModels, parentName, name, field)
			}
		}
	case typeArray:
		modelName := strings.TrimPrefix(smdModel.Items["$ref"], definitionsPrefix)
		def := smdModel.Definitions[modelName]
		allModels = append(allModels, definitionToModel(modelName, def))
		// add recursive fields
		for name, field := range def.Properties {
			if isModelProperty(field) {
				if field.Ref != "" {
					name = strings.TrimPrefix(field.Ref, definitionsPrefix)
				} else {
					name = fmt.Sprintf("%s%s", modelName, strings.Title(name))
				}

				allModels = addModelProperty(allModels, parentName, name, field)
			}
		}
	}

	return allModels
}

// objectToModel convert smdModel to goModel
func objectToModel(smdModel smd.JSONSchema) (res goModel) {
	if smdModel.Description != "" && !strings.Contains(smdModel.Description, " ") {
		res.Name = strings.Title(smdModel.Description)
	} else {
		res.Name = strings.Title(smdModel.Name)
	}

	for fieldName, fieldData := range smdModel.Properties {
		res.Fields = append(res.Fields, newGoParamFromProperty(res.Name, fieldName, fieldData))
	}

	sort.Slice(res.Fields, func(i, j int) bool { return res.Fields[i].Name < res.Fields[j].Name })

	return res
}

// definitionToModel convert smdDefinition to goModel
func definitionToModel(name string, smdDefinition smd.Definition) (res goModel) {
	res.Name = name
	for fieldName, fieldData := range smdDefinition.Properties {
		res.Fields = append(res.Fields, newGoParamFromProperty(name, fieldName, fieldData))
	}

	return res
}

// isModelProperty return true if smdProperty is object
func isModelProperty(smdProperty smd.Property) bool {
	return smdProperty.Type == typeObject || smdProperty.Type == typeArray && smdProperty.Items["$ref"] != ""
}

// isModel return true if smdValue is object
func isModel(smdValue smd.JSONSchema) bool {
	if smdValue.Type == typeArray && smdValue.Items["$ref"] != "" {
		return true
	}

	if smdValue.Type == typeObject {
		return true
	}

	return false
}

// newGoParamFromProperty convert prop to goValue
func newGoParamFromProperty(parentName, name string, prop smd.Property) goValue {
	if prop.Type == typeObject && prop.Ref == "" {
		name = fmt.Sprintf("%s%s", parentName, strings.Title(name))
	}
	p := smd.JSONSchema{
		Name:        name,
		Type:        prop.Type,
		Definitions: prop.Definitions,
		Items:       prop.Items,
	}

	return newGoParam(p)
}

// newGoParam convert smdValue to goValue
func newGoParam(smdValue smd.JSONSchema) goValue {
	var (
		t   = smdValue.Type
		res goValue
	)

	res.Name = smdValue.Name

	if smdValue.Type == typeObject && smdValue.Description != "" {
		t = smdValue.Description

		for fieldName, property := range smdValue.Properties {
			if property.Type == typeObject && property.Ref == "" {
				fieldName = fmt.Sprintf("%s%s", smdValue.Name, strings.Title(fieldName))
			}
			val := newGoParamFromProperty(res.Name, fieldName, property)
			val.ParentModelName = smdValue.Name
			res.Fields = append(res.Fields, val)
		}

		sort.Slice(res.Fields, func(i, j int) bool { return res.Fields[i].Name < res.Fields[j].Name })
	}

	itemType := ""
	if smdValue.Type == typeArray {
		itemType = smdValue.Items["type"]
		if ref, ok := smdValue.Items["$ref"]; ok {
			itemType = strings.TrimPrefix(ref, definitionsPrefix)
		}
	}

	res.Type = t
	res.ItemType = itemType

	return res
}

// newGoParams convert slice of smd1 objects to slice goValues
func newGoParams(params []smd.JSONSchema) (res []goValue) {
	for _, p := range params {
		res = append(res, newGoParam(p))
	}

	return
}

// Method return method by namespace and methodName.
func (g *Generator) Method(namespace, methodName string) goMethod {
	for name, srv := range g.smd.Services {
		if getNamespace(name) == namespace && getMethodName(name) == methodName {
			var result *goValue

			if srv.Returns.Type != "" {
				if srv.Returns.Type == typeObject && srv.Returns.Name == "" {
					srv.Returns.Name = fmt.Sprintf("%s%sResponse", strings.Title(namespace), strings.Title(methodName))
				}
				r := newGoParam(srv.Returns)
				result = &r
			}

			return goMethod{
				Params: newGoParams(srv.Parameters),
				Result: result,
			}
		}
	}

	return goMethod{}
}

// NamespaceNames return all namespace names.
func (g *Generator) NamespaceNames() (res []string) {
	m := map[string]struct{}{}
	for name := range g.smd.Services {
		m[getNamespace(name)] = struct{}{}
	}

	for name := range m {
		res = append(res, name)
	}

	sort.Strings(res)

	return
}

// NamespaceMethodNames return all methodNames by namespace name.
func (g *Generator) NamespaceMethodNames(namespace string) (res []string) {
	for name := range g.smd.Services {
		if getNamespace(name) != namespace {
			continue
		}

		res = append(res, getMethodName(name))
	}

	sort.Strings(res)

	return
}

var templateFuncs = template.FuncMap{
	"title": func(s string) string {
		if strings.EqualFold(s, "id") {
			return "ID"
		}

		if strings.HasSuffix(s, "Id") {
			s = strings.TrimSuffix(s, "Id")
			s += "ID"

			return strings.Title(s)
		}

		return strings.Title(s)
	},
}
