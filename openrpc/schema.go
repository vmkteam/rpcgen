package openrpc

import (
	openrpc "github.com/open-rpc/meta-schema"
	"github.com/vmkteam/zenrpc/v2/smd"
)

func NewSchema(schema smd.Schema) openrpc.OpenrpcDocument {
	orpc := openrpc.OpenrpcEnum0
	desc := openrpc.InfoObjectDescription(schema.Description)
	surl := openrpc.ServerObjectUrl(schema.Target)

	doc := openrpc.OpenrpcDocument{
		Openrpc: &orpc,
		Info: &openrpc.InfoObject{
			Description: &desc,
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
		name := openrpc.MethodObjectName(n)
		method := openrpc.MethodObject{
			Name:           &name,
			Params:         newParams(service),
			ParamStructure: newParamStruct(service),
			Result:         newResult(service),
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

func newParams(service smd.Service) *openrpc.MethodObjectParams {
	return nil
}

func newParamStruct(service smd.Service) *openrpc.MethodObjectParamStructure {
	return nil
}

func newResult(service smd.Service) *openrpc.MethodObjectResult {
	return nil
}

func newErrors(service smd.Service) *openrpc.MethodObjectErrors {
	return nil
}

func newComponents(schema smd.Schema) *openrpc.Components {
	return nil
}
