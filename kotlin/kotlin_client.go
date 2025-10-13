package kotlin

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/vmkteam/rpcgen/v2/gen"
	"github.com/vmkteam/zenrpc/v2/smd"
)

const (
	BasePackageAPI = "api"
	Bool           = "Boolean"
	Int            = "Int"
	Double         = "Double"
	Float          = "Float"
	String         = "String"
	Long           = "Long"
	Timestamp      = "ZonedDateTime"
	List           = "List"

	DefaultBoolFalse = "false"
	DefaultString    = "\"\""
	DefaultInteger   = "0"
	DefaultDouble    = "0.0"
	DefaultFloat     = "0f"
	DefaultTimestamp = "ZonedDateTime.now()"
	DefaultList      = "emptyList()"
)

var (
	id                  = "id"
	linebreakRegex      = regexp.MustCompile("[\r\n]+")
	kotlinDefaultValues = map[string]string{
		String:    DefaultString,
		Bool:      DefaultBoolFalse,
		Float:     DefaultFloat,
		Double:    DefaultDouble,
		Timestamp: DefaultTimestamp,
		Int:       DefaultInteger,
		Long:      DefaultInteger,
	}
	doubleSuffixes = map[string]struct{}{
		"Lat":       {},
		"Lon":       {},
		"Latitude":  {},
		"Longitude": {},
	}
)

type templateData struct {
	gen.GeneratorData
	Methods    []Method
	Models     []Model
	PackageAPI string
}

type Method struct {
	Name        string
	SafeName    string
	Description []string
	Errors      []Parameter
	Parameters  []Parameter
	Returns     Parameter
}

type Model struct {
	Name        string
	Description string
	IsInitial   bool
	Fields      []Parameter
}

type Parameter struct {
	Name         string
	Description  string
	Type         string
	BaseType     string
	ReturnType   string
	Optional     bool
	DefaultValue string
	IsObject     bool
	Properties   []Parameter
}

type Settings struct {
	Class      string
	IsProtocol bool
	PackageAPI string
}

type Generator struct {
	schema   smd.Schema
	settings Settings
}

func NewClient(schema smd.Schema, settings Settings) *Generator {
	return &Generator{schema: schema, settings: settings}
}

// Generate return kotlin generated code by template
func (g *Generator) Generate() ([]byte, error) {
	data := g.fillTemplateData()

	return g.executeTemplate(data)
}

// fillTemplateData prepare template data
func (g *Generator) fillTemplateData() templateData {
	if g.settings.PackageAPI == "" {
		g.settings.PackageAPI = BasePackageAPI
	}

	data := templateData{GeneratorData: gen.DefaultGeneratorData(), PackageAPI: g.settings.PackageAPI}

	modelsMap := make(map[string]Model)
	servicesMap := make(map[string][]Method)

	for serviceName, service := range g.schema.Services {
		var (
			params []Parameter
		)
		desc := linebreakRegex.ReplaceAllString(service.Description, "\n")
		method := Method{
			Name:     serviceName,
			Errors:   g.prepareErrs(service.Errors),
			SafeName: strings.ReplaceAll(serviceName, ".", ""),
			Returns:  g.prepareParameter(service.Returns),
		}
		if desc != "" {
			method.Description = strings.Split(desc, "\n")
		}
		if len(service.Returns.Definitions) > 0 {
		}
		g.defToModelMap(modelsMap, service.Returns.Definitions)
		paramToModelMap(modelsMap, method.Returns)
		for _, param := range service.Parameters {
			p := g.prepareParameter(param)
			g.defToModelMap(modelsMap, param.Definitions)
			paramToModelMap(modelsMap, p)
			params = append(params, p)
		}
		method.Parameters = params
		data.Methods = append(data.Methods, method)
		if g.settings.IsProtocol {
			parts := strings.SplitN(method.Name, ".", 2)
			if len(parts) != 2 {
				continue
			}
			servicesMap[parts[0]] = append(servicesMap[parts[0]], method)
		}
	}

	for _, v := range modelsMap {
		data.Models = append(data.Models, v)
	}

	initialModels := make(map[string]struct{})
	for _, m := range data.Methods {
		if len(m.Parameters) > 0 {
			for _, param := range m.Parameters {
				if param.IsObject {
					initialModels[param.ReturnType] = struct{}{}
				}
			}
		}
	}

	for i := range data.Models {
		if _, ok := initialModels[data.Models[i].Name]; ok {
			data.Models[i].IsInitial = true
		}
	}

	g.sortTemplateData(&data)

	return data
}

func (g *Generator) sortTemplateData(data *templateData) {
	sort.Slice(data.Methods, func(i, j int) bool {
		return data.Methods[i].Name < data.Methods[j].Name
	})

	sort.Slice(data.Models, func(i, j int) bool {
		return data.Models[i].Name < data.Models[j].Name
	})

	for idx := range data.Models {
		sort.Slice(data.Models[idx].Fields, func(i, j int) bool {
			return data.Models[idx].Fields[i].Name < data.Models[idx].Fields[j].Name
		})
	}
}

func (g *Generator) executeTemplate(data templateData) ([]byte, error) {
	t := model
	if g.settings.IsProtocol {
		t = protocolTemplate
	}

	gen.TemplateFuncs["hasDescriptions"] = hasDescriptions

	tmpl, err := template.New("kotlin_client").Funcs(gen.TemplateFuncs).Parse(t)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// propertiesToParams convert smd.PropertyList to []Parameter
func (g *Generator) propertiesToParams(list smd.PropertyList) []Parameter {
	parameters := make([]Parameter, 0, len(list))
	for _, prop := range list {
		p := Parameter{
			Name:        prop.Name,
			Optional:    prop.Optional,
			Description: prop.Description,
		}

		pType := kotlinType(prop.Type, p.Name)
		if prop.Type == smd.Object && prop.Ref != "" {
			pType = strings.TrimPrefix(prop.Ref, gen.DefinitionsPrefix)
			p.IsObject = true
		}

		if prop.Type == smd.Array {
			pType = arrayType(prop.Items, false)
		}

		p.Type = pType

		p.DefaultValue = kotlinDefault(pType)

		parameters = append(parameters, p)
	}
	return parameters
}

// prepareErrs convert errs map to []Parameter
func (g *Generator) prepareErrs(errs map[int]string) []Parameter {
	pp := make([]Parameter, 0, len(errs))
	for code, err := range errs {
		pp = append(pp, Parameter{
			Name:        strconv.Itoa(code),
			Description: err,
		})
	}

	return pp
}

// prepareParameter create Parameter from smd.JSONSchema
func (g *Generator) prepareParameter(in smd.JSONSchema) Parameter {
	out := Parameter{
		Name:        in.Name,
		Description: in.Description,
		BaseType:    kotlinType(in.Type),
		Optional:    in.Optional,
		Properties:  g.propertiesToParams(in.Properties),
	}

	pType := out.BaseType
	out.ReturnType = pType
	if in.Type == smd.Object {
		typeName := in.TypeName
		if typeName == "" && in.Description != "" && smd.IsSMDTypeName(in.Description, in.Type) {
			typeName = in.Description
		}

		if typeName != "" {
			pType = typeName
			out.ReturnType = pType

		}
		out.IsObject = true
	}
	out.Type = pType

	if in.Type == smd.Array {
		out.Type = arrayType(in.Items, false)
		out.ReturnType = arrayType(in.Items, true)
	}

	return out
}

// defToModelMap convert smd.Definition to Model and add to models map
func (g *Generator) defToModelMap(modelMap map[string]Model, definitions map[string]smd.Definition) {
	for name, def := range definitions {
		modelMap[name] = Model{Name: name, Fields: g.propertiesToParams(def.Properties)}
	}
}

// paramToModelMap add Parameter to model map if parameter is object
func paramToModelMap(modelsMap map[string]Model, p Parameter) {
	if p.IsObject && p.BaseType != p.Type {
		modelsMap[p.Type] = Model{
			Name:   p.Type,
			Fields: p.Properties,
		}
	}
}

// kotlinType convert smd types to kotlin types
func kotlinType(smdType string, propNames ...string) string {
	var name string
	if len(propNames) > 0 {
		name = propNames[0]
	}

	switch smdType {
	case smd.String:
		if kotlinTypeTimestamp(name) {
			return Timestamp
		}
		return String
	case smd.Boolean:
		return Bool
	case smd.Float:
		if kotlinTypeDouble(name) {
			return Double
		}
		return Float
	case smd.Integer:
		if kotlinTypeID(name) {
			return Long
		}
		return Int
	}
	return smdType
}

// kotlinDefault convert smd types to kotlin default types
func kotlinDefault(smdType string) string {
	if val, ok := kotlinDefaultValues[smdType]; ok {
		return val
	}

	if strings.Contains(smdType, List) {
		return DefaultList
	}

	return smdType
}

// kotlinTypeDouble check if property need set type Double
func kotlinTypeDouble(name string) bool {
	if _, ok := doubleSuffixes[name]; ok {
		return true
	}

	for suffix := range doubleSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}

	return false
}

// kotlinTypeID check if property is ID set type Long
func kotlinTypeID(name string) bool {
	return name == id || strings.HasSuffix(name, strings.ToUpper(id)) || strings.HasSuffix(name, "Id")
}

// kotlinTypeTimestamp check if is time property set Timestamp
func kotlinTypeTimestamp(name string) bool {
	return strings.HasSuffix(name, "edAt")
}

func arrayType(items map[string]string, isReturnType bool) string {
	var subType string
	if scalar, ok := items["type"]; ok {
		subType = kotlinType(scalar)
	}
	if ref, ok := items["$ref"]; ok {
		subType = strings.TrimPrefix(ref, gen.DefinitionsPrefix)
	}

	if isReturnType {
		return subType
	}

	return fmt.Sprintf("List<%s>", subType)
}

// hasDescriptions check description
func hasDescriptions(m Method) bool {
	return m.Returns.Description != "" || len(m.Description) > 0 || len(m.Parameters) > 0
}
