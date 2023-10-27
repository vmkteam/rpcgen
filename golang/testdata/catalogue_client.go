// Code generated from jsonrpc schema by rpcgen v2.4.2; DO NOT EDIT.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vmkteam/zenrpc/v2"
)

var (
	// Always import time package. Generated models can contain time.Time fields.
	_ time.Time
)

type Client struct {
	rpcClient *rpcClient

	Arith     *svcArith
	Catalogue *svcCatalogue
	Phonebook *svcPhonebook
}

func NewDefaultClient(endpoint string) *Client {
	return NewClient(endpoint, http.Header{}, &http.Client{})
}

func NewClient(endpoint string, header http.Header, httpClient *http.Client) *Client {
	c := &Client{
		rpcClient: newRPCClient(endpoint, header, httpClient),
	}

	c.Arith = newClientArith(c.rpcClient)
	c.Catalogue = newClientCatalogue(c.rpcClient)
	c.Phonebook = newClientPhonebook(c.rpcClient)

	return c
}

type Address struct {
	City   string `json:"City"`
	Street string `json:"Street"`
}

type Campaign struct {
	Groups []Group `json:"groups"`
	ID     int     `json:"id"`
}

type Group struct {
	Child  *Group   `json:"child,omitempty"`
	Groups []Group  `json:"groups"`
	ID     int      `json:"id"`
	Nodes  []Group  `json:"nodes"`
	Sub    SubGroup `json:"sub"`
	Title  string   `json:"title"`
}

type ModelPoint struct {
	ConnectedObject ObjectsAbstractObject `json:"ConnectedObject"`
	// coordinate
	X int `json:"X"`
	// coordinate
	Y int `json:"Y"`
}

type ObjectsAbstractObject struct {
}

type Person struct {
	// Addresses Could be nil or len() == 0.
	Addresses []Address `json:"Addresses"`
	BirthDate string    `json:"BirthDate"`
	// Deleted is flag for
	Deleted   bool   `json:"Deleted"`
	FirstName string `json:"FirstName"`
	// ID is Unique Identifier for person
	ID       int      `json:"ID"`
	LastName string   `json:"LastName"`
	Mobile   []string `json:"Mobile"`
	// Phone is main phone
	Phone     string   `json:"Phone"`
	WorkPhone *string  `json:"WorkPhone,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type PersonSearch struct {
	ByAddress *Address `json:"ByAddress,omitempty"`
	// ByName is filter for searching person by first name or last name.
	ByName  *string `json:"ByName,omitempty"`
	ByPhone string  `json:"ByPhone"`
	ByType  *string `json:"ByType,omitempty"`
}

type Quotient struct {
	// Quo docs
	Quo int `json:"Quo"`
	// Rem docs
	Rem int `json:"rem"`
}

type SubGroup struct {
	ID    int     `json:"id"`
	Nodes []Group `json:"nodes"`
	Title string  `json:"title"`
}

type svcArith struct {
	client *rpcClient
}

func newClientArith(client *rpcClient) *svcArith {
	return &svcArith{
		client: client,
	}
}

var (
	ErrArithCheckError500 = zenrpc.NewError(500, fmt.Errorf("test error"))
)

// CheckError throws error is isErr true.
func (c *svcArith) CheckError(ctx context.Context, isErr bool) (err error) {
	_req := struct {
		IsErr bool
	}{
		IsErr: isErr,
	}

	err = c.client.call(ctx, "arith.CheckError", _req, nil)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 500 {
			err = ErrArithCheckError500
		}
	}

	return
}

var (
	ErrArithCheckZenRPCError500 = zenrpc.NewError(500, fmt.Errorf("test error"))
)

// CheckZenRPCError throws zenrpc error is isErr true.
func (c *svcArith) CheckZenRPCError(ctx context.Context, isErr bool) (err error) {
	_req := struct {
		IsErr bool
	}{
		IsErr: isErr,
	}

	err = c.client.call(ctx, "arith.CheckZenRPCError", _req, nil)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 500 {
			err = ErrArithCheckZenRPCError500
		}
	}

	return
}

var (
	ErrArithDivide401 = zenrpc.NewError(401, fmt.Errorf("we do not serve 1"))
)

// Divide divides two numbers.
func (c *svcArith) Divide(ctx context.Context, a int, b int) (res *Quotient, err error) {
	_req := struct {
		A int
		B int
	}{
		A: a, B: b,
	}

	err = c.client.call(ctx, "arith.Divide", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 401 {
			err = ErrArithDivide401
		}
	}

	return
}

func (c *svcArith) DoSomething(ctx context.Context) (err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "arith.DoSomething", _req, nil)

	return
}

func (c *svcArith) DoSomethingWithPoint(ctx context.Context, p ModelPoint) (res ModelPoint, err error) {
	_req := struct {
		P ModelPoint
	}{
		P: p,
	}

	err = c.client.call(ctx, "arith.DoSomethingWithPoint", _req, &res)

	return
}

func (c *svcArith) GetPoints(ctx context.Context) (res []ModelPoint, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "arith.GetPoints", _req, &res)

	return
}

// Multiply multiples two digits and returns result.
func (c *svcArith) Multiply(ctx context.Context, a int, b int) (res int, err error) {
	_req := struct {
		A int
		B int
	}{
		A: a, B: b,
	}

	err = c.client.call(ctx, "arith.Multiply", _req, &res)

	return
}

// Pi returns math.Pi.
func (c *svcArith) Pi(ctx context.Context) (res float64, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "arith.Pi", _req, &res)

	return
}

func (c *svcArith) Positive(ctx context.Context) (res bool, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "arith.Positive", _req, &res)

	return
}

// Pow returns x**y, the base-x exponential of y. If Exp is not set then default value is 2.
func (c *svcArith) Pow(ctx context.Context, base float64, exp *float64) (res float64, err error) {
	_req := struct {
		Base float64
		Exp  *float64
	}{
		Base: base, Exp: exp,
	}

	err = c.client.call(ctx, "arith.Pow", _req, &res)

	return
}

// Sum sums two digits and returns error with error code as result and IP from context.
func (c *svcArith) Sum(ctx context.Context, a int, b int) (res bool, err error) {
	_req := struct {
		A int
		B int
	}{
		A: a, B: b,
	}

	err = c.client.call(ctx, "arith.Sum", _req, &res)

	return
}

// SumArray returns sum all items from array
func (c *svcArith) SumArray(ctx context.Context, array []float64) (res float64, err error) {
	_req := struct {
		Array []float64
	}{
		Array: array,
	}

	err = c.client.call(ctx, "arith.SumArray", _req, &res)

	return
}

type svcCatalogue struct {
	client *rpcClient
}

func newClientCatalogue(client *rpcClient) *svcCatalogue {
	return &svcCatalogue{
		client: client,
	}
}

func (c *svcCatalogue) First(ctx context.Context, groups []Group) (res bool, err error) {
	_req := struct {
		Groups []Group
	}{
		Groups: groups,
	}

	err = c.client.call(ctx, "catalogue.First", _req, &res)

	return
}

func (c *svcCatalogue) Second(ctx context.Context, campaigns []Campaign) (res bool, err error) {
	_req := struct {
		Campaigns []Campaign
	}{
		Campaigns: campaigns,
	}

	err = c.client.call(ctx, "catalogue.Second", _req, &res)

	return
}

func (c *svcCatalogue) Third(ctx context.Context) (res Campaign, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "catalogue.Third", _req, &res)

	return
}

type svcPhonebook struct {
	client *rpcClient
}

func newClientPhonebook(client *rpcClient) *svcPhonebook {
	return &svcPhonebook{
		client: client,
	}
}

var (
	ErrPhonebookByID404 = zenrpc.NewError(404, fmt.Errorf("person was not found"))
)

// ById returns Person from DB.
func (c *svcPhonebook) ByID(ctx context.Context, id int) (res *Person, err error) {
	_req := struct {
		ID int
	}{
		ID: id,
	}

	err = c.client.call(ctx, "phonebook.ById", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrPhonebookByID404
		}
	}

	return
}

// Delete marks person as deleted.
func (c *svcPhonebook) Delete(ctx context.Context, id int) (res bool, err error) {
	_req := struct {
		ID int
	}{
		ID: id,
	}

	err = c.client.call(ctx, "phonebook.Delete", _req, &res)

	return
}

// Get returns all people from DB.
func (c *svcPhonebook) Get(ctx context.Context, search PersonSearch, page *int, count *int) (res []Person, err error) {
	_req := struct {
		Search PersonSearch
		Page   *int
		Count  *int
	}{
		Search: search, Page: page, Count: count,
	}

	err = c.client.call(ctx, "phonebook.Get", _req, &res)

	return
}

// Removes deletes person from DB.
func (c *svcPhonebook) Remove(ctx context.Context, id int) (res bool, err error) {
	_req := struct {
		ID int
	}{
		ID: id,
	}

	err = c.client.call(ctx, "phonebook.Remove", _req, &res)

	return
}

var (
	ErrPhonebookSave400 = zenrpc.NewError(400, fmt.Errorf("invalid request"))
	ErrPhonebookSave401 = zenrpc.NewError(401, fmt.Errorf("use replace=true"))
)

// Save saves person to DB.
func (c *svcPhonebook) Save(ctx context.Context, p Person, replace *bool) (res int, err error) {
	_req := struct {
		P       Person
		Replace *bool
	}{
		P: p, Replace: replace,
	}

	err = c.client.call(ctx, "phonebook.Save", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 400 {
			err = ErrPhonebookSave400
		}
		if v.Code == 401 {
			err = ErrPhonebookSave401
		}
	}

	return
}

// ValidateSearch returns given search as result.
func (c *svcPhonebook) ValidateSearch(ctx context.Context, search *PersonSearch) (res *PersonSearch, err error) {
	_req := struct {
		Search *PersonSearch
	}{
		Search: search,
	}

	err = c.client.call(ctx, "phonebook.ValidateSearch", _req, &res)

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

func (rc *rpcClient) call(ctx context.Context, methodName string, request, result interface{}) error {
	// encode params
	bts, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("encode params: %w", err)
	}

	requestID := atomic.AddUint64(&rc.requestID, 1)
	requestIDBts := json.RawMessage(strconv.Itoa(int(requestID)))

	req := zenrpc.Request{
		Version: zenrpc.Version,
		ID:      &requestIDBts,
		Method:  methodName,
		Params:  bts,
	}

	res, err := rc.Exec(ctx, req)
	if err != nil {
		return err
	}

	if res == nil {
		return nil
	}

	if res.Error != nil {
		return res.Error
	}

	if res.Result == nil {
		return nil
	}

	if result == nil {
		return nil
	}

	return json.Unmarshal(*res.Result, result)
}

// Exec makes http request to jsonrpc endpoint and returns json rpc response.
func (rc *rpcClient) Exec(ctx context.Context, rpcReq zenrpc.Request) (*zenrpc.Response, error) {
	if n, ok := ctx.Value("JSONRPC2-Notification").(bool); ok && n {
		rpcReq.ID = nil
	}

	c, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("json marshal call failed: %w", err)
	}

	buf := bytes.NewReader(c)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rc.endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header = rc.header.Clone()
	req.Header.Add("Content-Type", "application/json")

	if xRequestID, ok := ctx.Value("X-Request-Id").(string); ok && req.Header.Get("X-Request-Id") == "" && xRequestID != "" {
		req.Header.Add("X-Request-Id", xRequestID)
	}

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

	var zresp zenrpc.Response
	if rpcReq.ID == nil {
		return &zresp, nil
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body (%s) read failed: %w", bb, err)
	}

	if err = json.Unmarshal(bb, &zresp); err != nil {
		return nil, fmt.Errorf("json decode failed (%s): %w", bb, err)
	}

	return &zresp, nil
}
