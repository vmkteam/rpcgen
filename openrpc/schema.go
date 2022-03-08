package openrpc

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	openrpc "github.com/open-rpc/meta-schema"
	"github.com/vmkteam/zenrpc/v2/smd"
	"strings"
)

func NewSchema(schema smd.Schema, title string) openrpc.OpenrpcDocument {
	orpc := openrpc.OpenrpcEnum0

	bs, _ := json.Marshal(schema)
	ver := openrpc.InfoObjectVersion(fmt.Sprintf("v0.0.0-%x", md5.Sum(bs)))

	name := openrpc.InfoObjectProperties(title)
	desc := openrpc.InfoObjectDescription(schema.Description)
	surl := openrpc.ServerObjectUrl(schema.Target)

	doc := openrpc.OpenrpcDocument{
		Openrpc: &orpc,
		Info: &openrpc.InfoObject{
			Title:       &name,
			Description: &desc,
			Version:     &ver,
		},
		Servers: &openrpc.Servers{{
			Url: &surl,
		}},
		Methods:    newMethods(schema),
		Components: newComponents(schema),
	}

	return doc
}

func newMethods(schema smd.Schema) *openrpc.Methods {
	methods := openrpc.Methods{}

	for n, service := range schema.Services {
		parts := strings.Split(n, ".")

		name := openrpc.MethodObjectName(n)
		tag := openrpc.TagObjectName(parts[0])

		method := openrpc.MethodObject{
			Name:           &name,
			Tags:           &openrpc.MethodObjectTags{{TagObject: &openrpc.TagObject{Name: &tag}}},
			ParamStructure: newParamStruct(service),
			Params:         newParams(n, service),
			Result:         newResult(n, service),
			Errors:         newErrors(service),
		}

		if service.Description != "" {
			desc := openrpc.MethodObjectDescription(service.Description)
			method.Description = &desc
		}

		methods = append(methods, openrpc.MethodOrReference{MethodObject: &method})
	}

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
			Schema: newJSONSchema(param, false),
		}

		if param.Description != "" {
			desc := openrpc.ContentDescriptorObjectDescription(param.Description)
			cont.Description = &desc
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
	name := openrpc.ContentDescriptorObjectName(fmt.Sprintf("%sResult", serviceName))
	return &openrpc.MethodObjectResult{ContentDescriptorObject: &openrpc.ContentDescriptorObject{
		Name:   &name,
		Schema: newJSONSchema(service.Returns, false),
	}}
}

func newErrors(service smd.Service) *openrpc.MethodObjectErrors {
	return nil
}

func newComponents(schema smd.Schema) *openrpc.Components {
	return nil
}

func newJSONSchema(schema smd.JSONSchema, withProps bool) *openrpc.JSONSchema {
	switch schema.Type {
	case "object":
		if withProps {

		} else {
			return &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
				Ref: refName(schema.Description),
			}}
		}
	case "array":
		items := &openrpc.JSONSchemaObject{}
		if schema.Items["$ref"] != "" {
			ref := openrpc.Ref(schema.Items["$ref"])
			items = &openrpc.JSONSchemaObject{
				Ref: &ref,
			}
		} else {
			itemsType := openrpc.SimpleTypes(schema.Items["type"])
			items = &openrpc.JSONSchemaObject{
				Type: &openrpc.Type{SimpleTypes: &itemsType},
			}
		}

		typ := openrpc.SimpleTypes(schema.Type)
		return &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
			Type:  &openrpc.Type{SimpleTypes: &typ},
			Items: &openrpc.Items{JSONSchema: &openrpc.JSONSchema{JSONSchemaObject: items}},
		}}

	default:
		typ := openrpc.SimpleTypes(schema.Type)
		return &openrpc.JSONSchema{JSONSchemaObject: &openrpc.JSONSchemaObject{
			Type: &openrpc.Type{SimpleTypes: &typ},
		}}
	}

	return nil
}

func refName(name string) *openrpc.Ref {
	ref := openrpc.Ref(fmt.Sprintf("#/components/schemas/%s", name))
	return &ref
}
