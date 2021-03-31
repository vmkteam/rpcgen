// Code generated from jsonrpc schema by rpc_client_generator; DO NOT EDIT.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/vmkteam/zenrpc/v2"
)

type Client struct {
	rpcClient *rpcClient

	Catalogue *Catalogue
}

func NewDefaultClient(endpoint string) *Client {
	return NewClient(endpoint, http.Header{}, &http.Client{})
}

func NewClient(endpoint string, header http.Header, httpClient *http.Client) *Client {
	c := &Client{
		rpcClient: newRPCClient(endpoint, header, httpClient),
	}

	c.Catalogue = NewClientCatalogue(c.rpcClient)

	return c
}

type Campaign struct {
	Group []Group `json:"group"`
	ID    int     `json:"id"`
}

type CatalogueThirdResponse struct {
	Group []Group `json:"group"`
	ID    int     `json:"id"`
}

type Group struct {
}

type SubGroup struct {
}

type Catalogue struct {
	client *rpcClient
}

func NewClientCatalogue(client *rpcClient) *Catalogue {
	return &Catalogue{
		client: client,
	}
}

func (c *Catalogue) First(ctx context.Context, groups []Group) (res bool, err error) {
	_req := struct {
		Groups []Group
	}{
		Groups: groups,
	}
	err = c.client.call(ctx, "catalogue.First", _req, &res)
	return
}

func (c *Catalogue) Second(ctx context.Context, campaigns []Campaign) (res bool, err error) {
	_req := struct {
		Campaigns []Campaign
	}{
		Campaigns: campaigns,
	}
	err = c.client.call(ctx, "catalogue.Second", _req, &res)
	return
}

func (c *Catalogue) Third(ctx context.Context) (res CatalogueThirdResponse, err error) {
	_req := struct {
	}{}
	err = c.client.call(ctx, "catalogue.Third", _req, &res)
	return
}

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
