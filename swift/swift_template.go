package swift

const client = `/// Code generated from jsonrpc schema by rpcgen v{{ .Version }}; DO NOT EDIT.

import Foundation

extension {{ .Class }}: RPCMethod {
    public var rpcMethod: String {
        switch self {
        case .batch(let requests): return requests.compactMap { $0.rpcMethod }.joined(separator: ",")
       	{{- range .Methods }}
        case .{{ .SafeName }}: return "{{ .Name }}"
       	{{- end }}
        }
    }
}

extension {{ .Class }}: RPCParameters {
    public var rpcParameters: [String: Any?]? {
        switch self {
        case .batch:
              return nil
{{- $methodsLen := len .Methods }}
{{- range $idx, $m := .Methods }}{{- $paramsLen := len .Parameters }}
        case {{ if gt $paramsLen 0 }}let {{ end }}.{{ .SafeName }}{{ if gt $paramsLen 0 }}({{ range $index, $item := .Parameters }}{{ .Name }}{{ if (notLast $index $paramsLen) }}, {{ end }}{{ end }}){{ end }}:
            return {{ if eq $paramsLen 0 }}nil{{ else }}[{{ range $index, $item := .Parameters }}"{{ .Name }}": {{ .Name }}{{ if or .IsArray .IsObject }}.any{{ end }}{{ if (notLast $index $paramsLen) }}, {{ end }}{{ end }}]{{ end }}
{{- if (notLast $idx $methodsLen) }}{{ print "\n" }}{{- end }}
{{- end }}
        }
    }
}

public enum {{ .Class }} {
    /// Make batch requests.
    case batch(requests: [{{ .Class }}])
{{ range .Methods }}{{- $paramsLen := len .Parameters }}
    {{- range .Description }}
    {{- if ne . "" }}
    /// {{ . }}
    {{- end }}
    {{- end }}
    {{- if ne .Returns.Type ""}}
    /// - Returns: {{ .Returns.Type }}{{ if .Returns.Optional }}?{{ end }}
    {{- end }}
    case {{ .SafeName }}{{ if gt $paramsLen 0 }}({{ range $index, $item := .Parameters }}{{ .Name }}: {{ .Type }}{{ if .Optional }}?{{ end }}{{ if (notLast $index $paramsLen) }}, {{ end }}{{ end }}){{ end }}
{{ end }}
}

{{ range .Models }}
{{- if .Description }}
/// {{ .Description }}
{{- end }}
public struct {{ .Name }}: Codable, Hashable {
    {{- $fieldsLen := len .Fields }}{{- range .Fields }}
    {{- if .Description }}
    /// {{ .Description }}
    {{- end }}
    {{- if and (not .Optional) (ne .DecodableDefault "") }}
    @DecodableDefault.{{ .DecodableDefault }}
    {{- end }}
    var {{ if .NeedEscaping }}{{ .SafeName }}{{ else }}{{ .Name }}{{ end }}: {{ .Type }}{{ if .Optional }}?{{ end }}
    {{- end }}
    init({{ range $index, $f := .Fields }}{{ .Name }}: {{ .Type }}{{ if .Optional }}? = nil{{ end }}{{ if (notLast $index $fieldsLen) }}, {{ end }}{{ end }}) {
        {{- range .Fields }}
        self.{{ .Name }} = {{ if .NeedEscaping }}{{ .SafeName }}{{ else }}{{ .Name }}{{ end }}
        {{- end }}
    }
}
{{ end }}
`

const protocolTemplate = `/// Code generated from jsonrpc schema by rpcgen v{{ .Version }}; DO NOT EDIT.

import Foundation
{{- range $service := .Protocols }}

protocol {{ $service.Class }}Networking {
{{- range $method := $service.Methods }}
    {{- range .Description }}
    {{- if ne . "" }}
    /// {{ . }}
    {{- end }}
    {{- end }}
    func {{ enumCaseName $method.Name }}({{- range $index, $item := $method.Parameters }}{{ $item.Name }}: {{ $item.Type }}{{ if $item.Optional }}?{{ end }}{{ if (notLast $index (len $method.Parameters)) }}, {{ end }}{{ end }}) async -> Result<{{ if $method.Returns.Type }}{{ $method.Returns.Type }}, {{ else }}{{ end }}RpcError>
{{- end }}
}

extension Networking: {{ $service.Class }}Networking {
{{- range $idx, $method := $service.Methods }}
    {{- range .Description }}
    {{- if ne . "" }}
    /// {{ . }}
    {{- end }}
    {{- end }}
    func {{ enumCaseName $method.Name }}(
        {{- range $index, $item := $method.Parameters -}}
            {{ $item.Name }}: {{ $item.Type }}{{ if $item.Optional }}? = nil{{ end }}{{ if (notLast $index (len $method.Parameters)) }},{{ end }}
        {{- end -}}
    ) async -> Result<{{ if $method.Returns.Type }}{{ $method.Returns.Type }}, {{ else }}{{ end }}RpcError> {
        await request(.{{ enumCaseName $method.Name }}{{ if $method.Parameters }}({{- range $index, $item := $method.Parameters }}{{ $item.Name }}: {{ $item.Name }}{{ if (notLast $index (len $method.Parameters)) }}, {{ end }}{{ end }}){{ else }}(){{ end }})
    }
{{- if (notLast $idx (len $service.Methods)) }}
{{ end }}
{{- end }}
}
{{ end }}`
