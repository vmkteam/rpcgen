package golang

import (
	"reflect"
	"testing"

	"github.com/vmkteam/zenrpc/v2/smd"
)

func TestConvertJSONSchema(t *testing.T) {
	tc := []struct {
		in  smd.Schema
		out Schema
	}{
		{
			in: smd.Schema{
				Services: map[string]smd.Service{
					"simple.withoutArgsAndReturn": {},
				},
			},
			out: Schema{
				Namespaces: []Namespace{
					{
						Name: "simple",
						Methods: []Method{
							{
								Name: "withoutArgsAndReturn",
							},
						},
					},
				},
			},
		},
		// simple type argument
		{
			in: smd.Schema{
				Services: map[string]smd.Service{
					"simple.stringArg": {
						Parameters: []smd.JSONSchema{
							{
								Name: "test",
								Type: "string",
							},
							{
								Name: "testBool",
								Type: "boolean",
							},
						},
					},
				},
			},
			out: Schema{
				Namespaces: []Namespace{
					{
						Name: "simple",
						Methods: []Method{
							{
								Name: "stringArg",
								Params: []Value{
									{
										Name: "test",
										Type: "string",
									},
									{
										Name: "testBool",
										Type: "boolean",
									},
								},
							},
						},
					},
				},
			},
		},
		// simple type return
		{
			in: smd.Schema{
				Services: map[string]smd.Service{
					"simple.stringReturn": {
						Returns: smd.JSONSchema{
							Name: "test",
							Type: "string",
						},
					},
				},
			},
			out: Schema{
				Namespaces: []Namespace{
					{
						Name: "simple",
						Methods: []Method{
							{
								Name: "stringReturn",
								Returns: &Value{
									Name: "test",
									Type: "string",
								},
							},
						},
					},
				},
			},
		},
		// complex type with simple props param
		{
			in: smd.Schema{
				Services: map[string]smd.Service{
					"complex.objectSimplePropsParam": {
						Parameters: []smd.JSONSchema{
							{
								Name: "test",
								Type: "object",
								// Description is required for this type conversion
								Description: "ApiObject",
								Properties: smd.PropertyList{
									{
										Name: "test",
										Type: "string",
									},
								},
							},
						},
					},
				},
			},
			out: Schema{
				Namespaces: []Namespace{
					{
						Name: "complex",
						Methods: []Method{
							{
								Name: "objectSimplePropsParam",
								Params: []Value{
									{
										Name:      "test",
										Type:      "object",
										ModelName: "ApiObject",
									},
								},
								Models: []Model{
									{
										Name: "ApiObject",
										Fields: []Value{
											{
												Name: "test",
												Type: "string",
											},
										},
										IsParamModel: true,
										ParamName:    "test",
									},
								},
							},
						},
					},
				},
			},
		},
		// complex type with complex prop
		{
			in: smd.Schema{
				Services: map[string]smd.Service{
					"complex.objectComplexPropsParam": {
						Parameters: []smd.JSONSchema{
							{
								Name: "test",
								Type: "object",
								// Description is required for this type conversion
								Description: "ApiObject",
								Properties: smd.PropertyList{
									{
										Name: "test",
										Type: "object",
										Ref:  "#/definitions/ApiComplexType",
									},
								},
								Definitions: map[string]smd.Definition{
									"ApiComplexType": {
										Type: "object",
										Properties: smd.PropertyList{
											{
												Name: "test",
												Type: "string",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			out: Schema{
				Namespaces: []Namespace{
					{
						Name: "complex",
						Methods: []Method{
							{
								Name: "objectComplexPropsParam",
								Params: []Value{
									{
										Name:      "test",
										Type:      "object",
										ModelName: "ApiObject",
									},
								},
								Models: []Model{
									{
										Name: "ApiComplexType",
										Fields: []Value{
											{
												Name: "test",
												Type: "string",
											},
										},
									},
									{
										Name: "ApiObject",
										Fields: []Value{
											{
												Name:      "test",
												Type:      "object",
												ModelName: "ApiComplexType",
											},
										},
										IsParamModel: true,
										ParamName:    "test",
									},
								},
							},
						},
					},
				},
			},
		},
		// array with complex type
		{
			in: smd.Schema{
				Services: map[string]smd.Service{
					"complex.objectComplexArray": {
						Parameters: []smd.JSONSchema{
							{
								Name: "test",
								Type: "array",
								Items: map[string]string{
									"$ref": "#/definitions/ApiComplexType",
								},
								Definitions: map[string]smd.Definition{
									"ApiComplexType": {
										Type: "object",
										Properties: smd.PropertyList{
											{
												Name: "test",
												Type: "string",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			out: Schema{
				Namespaces: []Namespace{
					{
						Name: "complex",
						Methods: []Method{
							{
								Name: "objectComplexArray",
								Params: []Value{
									{
										Name:          "test",
										Type:          "array",
										ArrayItemType: "object",
										ModelName:     "ApiComplexType",
									},
								},
								Models: []Model{
									{
										Name: "ApiComplexType",
										Fields: []Value{
											{
												Name: "test",
												Type: "string",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// object param with array items
		{
			in: smd.Schema{
				Services: map[string]smd.Service{
					"complex.objectParamWithArray": {
						Parameters: []smd.JSONSchema{
							{
								Name: "test",
								Type: "object",
								Properties: smd.PropertyList{
									{
										Name: "items",
										Type: "array",
										Items: map[string]string{
											"$ref": "#/definitions/CartItem",
										},
									},
								},
								Definitions: map[string]smd.Definition{
									"CartItem": {
										Type: "object",
										Properties: smd.PropertyList{
											{
												Name: "count",
												Type: "integer",
											},
											{
												Name: "productId",
												Type: "integer",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			out: Schema{
				Namespaces: []Namespace{
					{
						Name: "complex",
						Methods: []Method{
							{
								Name: "objectParamWithArray",
								Params: []Value{
									{
										Name:      "test",
										Type:      "object",
										ModelName: "ComplexObjectParamWithArrayTestParam",
									},
								},
								Models: []Model{
									{
										Name: "CartItem",
										Fields: []Value{
											{
												Name: "count",
												Type: "integer",
											},
											{
												Name: "productId",
												Type: "integer",
											},
										},
									},
									{
										Name: "ComplexObjectParamWithArrayTestParam",
										Fields: []Value{
											{
												Name:          "items",
												Type:          "array",
												ArrayItemType: "object",
												ModelName:     "CartItem",
											},
										},
										IsParamModel: true,
										ParamName:    "test",
									},
								},
							},
						},
					},
				},
			},
		},
		// object param with no properties/definitions should result in "any"
		{
			in: smd.Schema{
				Services: map[string]smd.Service{
					"simple.anyParam": {
						Parameters: []smd.JSONSchema{
							{
								Name: "data",
								Type: "object",
							},
						},
						Returns: smd.JSONSchema{
							Type: smd.Object,
						},
					},
				},
			},
			out: Schema{
				Namespaces: []Namespace{
					{
						Name: "simple",
						Methods: []Method{
							{
								Name: "anyParam",
								Params: []Value{
									{
										Name:      "data",
										Type:      "object",
										ModelName: "",
									},
								},
								Returns: &Value{
									Type: "object",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tc {
		schema := NewSchema(tt.in)

		if !reflect.DeepEqual(schema, tt.out) {
			t.Errorf("bad conversion need=\n%+v\n%+v", tt.out, schema)
		}
	}
}
