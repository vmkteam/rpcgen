package golang

import (
	"bytes"
	"go/format"
	"strings"
	"text/template"
	"unicode"

	"github.com/vmkteam/zenrpc/v2/smd"
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
	schema Schema
}

// NewClient create Generator from zenrpc/v2 SMD.
func NewClient(schema smd.Schema) *Generator {
	return &Generator{schema: NewSchema(schema)}
}

// Generate returns generated Go client.
func (g *Generator) Generate() ([]byte, error) {
	tmpl, err := template.New("").Funcs(templateFuncs).Parse(goTpl)
	if err != nil {
		return nil, err
	}

	// compile template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.schema); err != nil {
		return nil, err
	}

	//return buf.Bytes(), nil
	return format.Source(buf.Bytes())
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

// titleFirstLetter upper case first letter of str
func titleFirstLetter(str string) string {
	ss := []rune(str)
	for i, v := range ss {
		return string(unicode.ToUpper(v)) + string(ss[i+1:])
	}

	return ""
}
