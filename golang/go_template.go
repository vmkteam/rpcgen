package golang

// goTpl contains template for Go client
const goTpl = `// Code generated from jsonrpc schema by rpc_client_generator; DO NOT EDIT.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"sync/atomic"
	"strconv"

	"github.com/vmkteam/zenrpc/v2"
)

type Client struct {
	rpcClient *rpcClient

{{ range .NamespaceNames }}
	{{ title . }} *{{ title .}}{{ end }}}

func NewDefaultClient(endpoint string) *Client {
	return NewClient(endpoint, http.Header{}, &http.Client{})
}

func NewClient(endpoint string, header http.Header, httpClient *http.Client) *Client {
	c := &Client{
		rpcClient: newRPCClient(endpoint, header, httpClient),
	}
{{ range .NamespaceNames }}
	c.{{ title . }} = NewClient{{ title .}}(c.rpcClient){{ end }}

	return c
}

{{ range .Models }}
type {{ .Name }} struct {
	{{ range .Fields }}{{ if ne .Description "" }}// {{ .Description }}
	{{ end }}{{ title .Name }} {{ .GoType }} ` + "`json:\"{{ .Name }}{{if .Optional}},omitempty{{end}}\"`" + `
{{ end }}
}
{{ end }} 


{{ range .NamespaceNames }}
 {{ $lTitle := . }}{{ $namespace := . }}
type {{ title . }} struct {
	client *rpcClient
}

func NewClient{{ title . }} (client *rpcClient) *{{ title . }}  {
	return &{{ title . }}{
		client: client,
	}
}

{{ range $.NamespaceMethodNames .}} {{ $method := $.MethodByName $namespace . }}

{{ $method.CommentDescription }}
func (c *{{ title $lTitle}}) {{ title . }}(ctx context.Context, {{ range $method.Params }}{{ .Name }} {{ .GoType }}, {{ end }}) ( {{ if $method.HasResult }} res {{ $method.Returns.GoType }},  {{ else }}/* end  */ {{end}} err error) {
	_req := struct {
		{{ range $method.Params }}{{ title .Name }} {{ .GoType }}
{{ end }}
	} {
		{{ range $method.Params }}{{ title .Name }}: {{ .Name }}, {{ end }}
	}
	
	err = c.client.call(ctx, "{{ $namespace }}.{{ . }}", _req, {{ if $method.HasResult }} &res {{ else }} nil {{ end }})
	
	return
}
{{ end }}
{{ end }}

type rpcClient struct {
	endpoint string
	cl       *http.Client

	requestID uint64
	header    http.Header
}

func newRPCClient(endpoint string, header http.Header, httpClient *http.Client) *rpcClient {
	return &rpcClient{
		endpoint: endpoint,
		header:   header,
		cl:       httpClient,
	}
}

func (c *rpcClient) call(ctx context.Context, methodName string, request, result interface{}) error {
	// encode params
	bts, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("encode params: %w", err)
	}

	requestID := atomic.AddUint64(&c.requestID, 1)
	requestIDBts := json.RawMessage(strconv.Itoa(int(requestID)))

	req := zenrpc.Request{
		Version: zenrpc.Version,
		ID:      &requestIDBts,
		Method:  methodName,
		Params:  bts,
	}

	res, err := c.Exec(ctx, req)
	if err != nil {
		return err
	}

	if res.Error != nil {
		return res.Error
	}

	if res == nil || res.Result == nil {
		return nil
	}

	return json.Unmarshal(*res.Result, result)
}

// Exec makes http request to jsonrpc endpoint and returns json rpc response.
func (rc *rpcClient) Exec(ctx context.Context, rpcReq zenrpc.Request) (*zenrpc.Response, error) {
	c, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("json marshal call failed: %w", err)
	}

	buf := bytes.NewReader(c)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rc.endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header = rc.header
	req.Header.Add("Content-Type", "application/json")

	// Do request
	resp, err := rc.cl.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, fmt.Errorf("make request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response (%d)", resp.StatusCode)
	}

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body (%s) read failed: %w", bb, err)
	}

	var zresp zenrpc.Response
	if err = json.Unmarshal(bb, &zresp); err != nil {
		return nil, fmt.Errorf("json decode failed (%s): %w", bb, err)
	}

	return &zresp, nil
}
`
