package openrpc

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/iancoleman/orderedmap"
	openrpc "github.com/vmkteam/meta-schema"
	"github.com/vmkteam/zenrpc/v2/smd"
)

func NewSchema(schema smd.Schema, title, baseurl string) openrpc.OpenrpcDocument {
	orpc := openrpc.OpenrpcEnum0

	bs, _ := json.Marshal(schema)
	ver := openrpc.InfoObjectVersion(fmt.Sprintf("v0.0.0-%x", md5.Sum(bs)))

	name := openrpc.InfoObjectProperties(title)
	desc := openrpc.InfoObjectDescription(schema.Description)
	surl := openrpc.ServerObjectUrl(fmt.Sprintf("%s%s", baseurl, schema.Target))
	sname := openrpc.ServerObjectName(sanitizeHost(baseurl))

	doc := openrpc.OpenrpcDocument{
		Openrpc: &orpc,
		Info: &openrpc.InfoObject{
			Title:       &name,
			Description: &desc,
			Version:     &ver,
		},
		Servers: &openrpc.Servers{{
			Name: &sname,
			Url:  &surl,
		}},
		Methods:    newMethods(schema),
		Components: newComponents(schema),
	}

	return doc
}

func newMethods(schema smd.Schema) *openrpc.Methods {
	methods := openrpc.Methods{}

	for serviceName, service := range schema.Services {
		parts := strings.Split(serviceName, ".")

		name := openrpc.MethodObjectName(serviceName)
		tag := openrpc.TagObjectName(parts[0])

		method := openrpc.MethodObject{
			Name:           &name,
			Tags:           &openrpc.MethodObjectTags{{TagObject: &openrpc.TagObject{Name: &tag}}},
			ParamStructure: newParamStruct(service),
			Params:         newParams(serviceName, service),
			Result:         newResult(serviceName, service),
			Errors:         newErrors(service),
		}

		if service.Description != "" {
			desc := openrpc.MethodObjectSummary(service.Description)
			method.Summary = &desc
		}

		methods = append(methods, openrpc.MethodOrReference{MethodObject: &method})
	}

	sort.Slice(methods, func(i, j int) bool {
		return string(*methods[i].MethodObject.Name) < string(*methods[j].MethodObject.Name)
	})

	return &methods
}

func newParamStruct(service smd.Service) *openrpc.MethodObjectParamStructure {
	for _, param := range service.Parameters {
		if param.Name == "" {
			param := openrpc.MethodObjectParamStructureEnum0
			return &param
		} else {
			param := openrpc.MethodObjectParamStructureEnum1
			return &param
		}
	}

	return nil
}

func newParams(serviceName string, service smd.Service) *openrpc.MethodObjectParams {
	conts := openrpc.MethodObjectParams{}
	for i, param := range service.Parameters {
		name := openrpc.ContentDescriptorObjectName(param.Name)
		if param.Name == "" {
			name = openrpc.ContentDescriptorObjectName(fmt.Sprintf("%s%d", serviceName, i))
		}

		cont := openrpc.ContentDescriptorObject{
			Name:   &name,
			Schema: newJSONSchema(serviceName, param),
		}

		if param.Description != "" {
			desc := openrpc.ContentDescriptorObjectSummary(param.Description)
			cont.Summary = &desc
		}

		if !param.Optional {
			req := openrpc.ContentDescriptorObjectRequired(true)
			cont.Required = &req
		}

		conts = append(conts, openrpc.ContentDescriptorOrReference{ContentDescriptorObject: &cont})
	}

	return &conts
}

func newResult(serviceName string, service smd.Service) *openrpc.MethodObjectResult {
	switch service.Returns.Type {
	case "":
		return &openrpc.MethodObjectResult{ReferenceObject: &openrpc.ReferenceObject{Ref: cdRefName(respName("null"))}}
	case smd.Array:
		if itemType, ok := service.Returns.Items["type"]; ok {
			return &openrpc.MethodObjectResult{ReferenceObject: &openrpc.ReferenceObject{Ref: cdRefName(arrayRespName(itemType))}}
		}
		fallthrough
	case smd.Object:
		name := openrpc.ContentDescriptorObjectName(objName(respName(serviceName)))
		var desc *openrpc.ContentDescriptorObjectSummary
		if service.Returns.Description != "" {
			d := openrpc.ContentDescriptorObjectSummary(service.Returns.Description)
			desc = &d
		}

		return &openrpc.MethodObjectResult{ContentDescriptorObject: &openrpc.ContentDescriptorObject{
			Name:    &name,
			Summary: desc,
			Schema:  newJSONSchema(respName(serviceName), service.Returns),
		}}
	default:
		return &openrpc.MethodObjectResult{ReferenceObject: &openrpc.ReferenceObject{Ref: cdRefName(respName(service.Returns.Type))}}
	}
}

func newErrors(service smd.Service) *openrpc.MethodObjectErrors {
	var errors openrpc.MethodObjectErrors
	for c, m := range service.Errors {
		code := openrpc.ErrorObjectCode(int64(c))
		mess := openrpc.ErrorObjectMessage(m)

		errors = append(errors, openrpc.ErrorOrReference{
			ErrorObject: &openrpc.ErrorObject{
				Code:    &code,
				Message: &mess,
			},
		})
	}

	sort.Slice(errors, func(i, j int) bool {
		return int64(*errors[i].ErrorObject.Code) > int64(*errors[i].ErrorObject.Code)
	})

	if len(errors) == 0 {
		return nil
	}

	return &errors
}

func newComponents(schema smd.Schema) *openrpc.Components {
	components := openrpc.SchemaComponents{}
	descriptors := openrpc.ContentDescriptorComponents{}

	for serviceName, service := range schema.Services {
		for _, param := range service.Parameters {
			parseComponentsFromSchema(serviceName, param, components)
		}

		parseDescriptorsFromSchema(respName(serviceName), service.Returns, components, descriptors)
	}

	return &openrpc.Components{Schemas: &components, ContentDescriptors: &descriptors}
}

// recursive hell
func parseComponentsFromSchema(serviceName string, schema smd.JSONSchema, components openrpc.SchemaComponents) {
	sch := newJSONSchema(serviceName, schema)
	if sch.JSONSchemaObject.Ref != nil && len(schema.Properties) > 0 {
		base := refBase(sch.JSONSchemaObject.Ref)

		if _, ok := components[base]; !ok {
			components[base] = newPropertiesFromList(schema.Properties, components)
		}
	}

	if len(schema.Definitions) > 0 {
		parseComponentsFromDefinitions(schema.Definitions, components)
	}
}

func parseDescriptorsFromSchema(serviceName string, schema smd.JSONSchema, components openrpc.SchemaComponents, descriptors openrpc.ContentDescriptorComponents) {
	sch := newJSONSchema(serviceName, schema)

	var (
		descriptor *openrpc.ContentDescriptorObject
		component  *openrpc.JSONSchema
		base       string
	)

	switch schema.Type {
	case "":
		descriptor = newSimpleResponse("null")
		base = string(*descriptor.Name)
	case smd.Array:
		if itemType, ok := schema.Items["type"]; ok {
			descriptor = newSimpleArrayResponse(itemType)
			base = string(*descriptor.Name)
		}
	case smd.Object:
		if sch.JSONSchemaObject.Ref != nil && len(schema.Properties) > 0 {
			component = newPropertiesFromList(schema.Properties, components)
			base = refBase(sch.JSONSchemaObject.Ref)
		}
	default:
		descriptor = newSimpleResponse(schema.Type)
		base = string(*descriptor.Name)
	}

	if descriptor != nil {
		if _, ok := descriptors[base]; !ok {
			descriptors[base] = descriptor
		}
	}

	if component != nil {
		if _, ok := components[base]; !ok {
			components[base] = component
		}
	}

	if len(schema.Definitions) > 0 {
		parseComponentsFromDefinitions(schema.Definitions, components)
	}
}

func parseComponentsFromDefinitions(definitions map[string]smd.Definition, components openrpc.SchemaComponents) {
	for n, definition := range definitions {
		name := objName(n)
		if _, ok := components[name]; !ok {
			components[name] = newPropertiesFromList(definition.Properties, components)
		}
	}
}

func newPropertiesFromList(props smd.PropertyList, components openrpc.SchemaComponents) *openrpc.JSONSchema {
	required := openrpc.StringArray{}
	result := orderedmap.New()

	for _, prop := range props {
		if len(prop.Definitions) > 0 {
			parseComponentsFromDefinitions(prop.Definitions, components)
		}

		if !prop.Optional {
			required = append(required, openrpc.StringDoaGddGA(prop.Name))
		}

		var desc *openrpc.Description
		if prop.Description != "" {
			d := openrpc.Description(prop.Description)
			desc = &d
		}

		switch prop.Type {
		case smd.Object:
			result.Set(prop.Name, openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
				Ref:         refName(prop.Ref),
				Description: desc,
			}})
		case smd.Array:
			items := &openrpc.JSONSchemaObject{}
			if prop.Items["$ref"] != "" {
				items = &openrpc.JSONSchemaObject{
					Ref: refName(prop.Items["$ref"]),
				}
			} else {
				itemsType := openrpc.SimpleTypes(prop.Items["type"])
				items = &openrpc.JSONSchemaObject{
					Type: &openrpc.Type{SimpleTypes: &itemsType},
				}
			}

			typ := openrpc.SimpleTypes(prop.Type)
			result.Set(prop.Name, openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
				Description: desc,
				Type:        &openrpc.Type{SimpleTypes: &typ},
				Items:       &openrpc.Items{JSONSchema: &openrpc.JSONSchema{JSONSchemaObject: items}},
			}})

		default:
			typ := openrpc.SimpleTypes(prop.Type)
			result.Set(prop.Name, openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
				Description: desc,
				Type:        &openrpc.Type{SimpleTypes: &typ},
			}})
		}
	}

	return &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
		Properties: result,
		Required:   &required,
	}}
}

func newJSONSchema(serviceName string, schema smd.JSONSchema) *openrpc.JSONSchema {
	var desc *openrpc.Description
	if schema.Description != "" {
		d := openrpc.Description(schema.Description)
		desc = &d
	}

	switch schema.Type {
	case smd.Object:
		var ref *openrpc.Ref
		if isObjName(schema.Description) {
			ref = refName(schema.Description)
			desc = nil
		} else {
			ref = refName(objName(serviceName, schema.Name))
		}

		return &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
			Ref:         ref,
			Description: desc,
		}}
	case smd.Array:
		items := &openrpc.JSONSchemaObject{}
		if schema.Items["$ref"] != "" {
			items = &openrpc.JSONSchemaObject{
				Ref: refName(schema.Items["$ref"]),
			}
		} else {
			itemsType := openrpc.SimpleTypes(schema.Items["type"])
			items = &openrpc.JSONSchemaObject{
				Type: &openrpc.Type{SimpleTypes: &itemsType},
			}
		}

		typ := openrpc.SimpleTypes(schema.Type)
		return &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
			Type:        &openrpc.Type{SimpleTypes: &typ},
			Items:       &openrpc.Items{JSONSchema: &openrpc.JSONSchema{JSONSchemaObject: items}},
			Description: desc,
		}}

	default:
		typ := openrpc.SimpleTypes(schema.Type)
		return &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
			Type:        &openrpc.Type{SimpleTypes: &typ},
			Description: desc,
		}}
	}
}

func refName(name string) *openrpc.Ref {
	if strings.HasPrefix(name, "#/definitions/") {
		name = strings.Replace(name, "#/definitions/", "", 1)
	}

	ref := openrpc.Ref(fmt.Sprintf("#/components/schemas/%s", objName(name)))
	return &ref
}

func cdRefName(name string) *openrpc.Ref {
	ref := openrpc.Ref(fmt.Sprintf("#/components/contentDescriptors/%s", objName(name)))
	return &ref
}

func refBase(ref *openrpc.Ref) string {
	return path.Base(string(*ref))
}

func isObjName(name string) bool {
	return name != "" && objName(name) == name
}

func objName(names ...string) string {
	buf := strings.Builder{}
	for i, name := range names {
		if i == 0 {
			buf.WriteString(name)
		} else {
			buf.WriteString(strings.Title(name))
		}
	}

	return strings.Title(regexp.MustCompile(`[^a-zA-z1-9_]`).ReplaceAllString(buf.String(), ""))
}

func respName(name string) string {
	return strings.Title(fmt.Sprintf("%sResponse", name))
}

func arrayRespName(name string) string {
	return strings.Title(fmt.Sprintf("%sArrayResponse", name))
}

func newSimpleArrayResponse(typeName string) *openrpc.ContentDescriptorObject {
	name := openrpc.ContentDescriptorObjectName(arrayRespName(typeName))
	desc := openrpc.ContentDescriptorObjectSummary(fmt.Sprintf("%s array response", typeName))
	typ := openrpc.SimpleTypes(typeName)

	return &openrpc.ContentDescriptorObject{
		Name:    &name,
		Summary: &desc,
		Schema: &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
			Items: &openrpc.Items{JSONSchema: &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
				Type: &openrpc.Type{SimpleTypes: &typ},
			}}},
		}},
	}
}

func newSimpleResponse(typeName string) *openrpc.ContentDescriptorObject {
	name := openrpc.ContentDescriptorObjectName(respName(typeName))
	desc := openrpc.ContentDescriptorObjectSummary(fmt.Sprintf("%s response", typeName))
	typ := openrpc.SimpleTypes(typeName)

	return &openrpc.ContentDescriptorObject{
		Name:    &name,
		Summary: &desc,
		Schema: &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
			Type: &openrpc.Type{SimpleTypes: &typ},
		}},
	}
}

func sanitizeHost(baseurl string) string {
	u, err := url.Parse(baseurl)
	if err != nil {
		return objName(baseurl)
	}

	return u.Host
}
