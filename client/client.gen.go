// Package client provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// Error defines model for Error.
type Error struct {
	Errors []string `json:"errors"`
}

// MessageResponse defines model for MessageResponse.
type MessageResponse struct {
	// Embedded struct due to allOf(#/components/schemas/NewMessage)
	NewMessage
	// Embedded fields due to inline allOf schema

	// timestamp full-date - RFC3339
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`

	// timestamp full-date - RFC3339
	UpdatedAt time.Time `json:"updated_at"`
}

// NewMessage defines model for NewMessage.
type NewMessage struct {
	Language     string `json:"language"`
	ProviderType string `json:"provider_type"`

	// Message content
	Template     string    `json:"template"`
	TemplateKeys *[]string `json:"template_keys,omitempty"`

	// will be used as part of the uniqunes of the message for example type could be greetings, driver arrived,...etc
	Type string `json:"type"`
}

// TriggerMessage defines model for TriggerMessage.
type TriggerMessage struct {

	// message type
	MessageType string   `json:"message_type"`
	UsersIds    []string `json:"users_ids"`
}

// AddMessageJSONBody defines parameters for AddMessage.
type AddMessageJSONBody NewMessage

// TriggerMessageJSONBody defines parameters for TriggerMessage.
type TriggerMessageJSONBody TriggerMessage

// AddMessageRequestBody defines body for AddMessage for application/json ContentType.
type AddMessageJSONRequestBody AddMessageJSONBody

// TriggerMessageRequestBody defines body for TriggerMessage for application/json ContentType.
type TriggerMessageJSONRequestBody TriggerMessageJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditor = fn
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// AddMessage request  with any body
	AddMessageWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	AddMessage(ctx context.Context, body AddMessageJSONRequestBody) (*http.Response, error)

	// TriggerMessage request  with any body
	TriggerMessageWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	TriggerMessage(ctx context.Context, body TriggerMessageJSONRequestBody) (*http.Response, error)

	// GetMessage request
	GetMessage(ctx context.Context, messageId string) (*http.Response, error)
}

func (c *Client) AddMessageWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewAddMessageRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) AddMessage(ctx context.Context, body AddMessageJSONRequestBody) (*http.Response, error) {
	req, err := NewAddMessageRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) TriggerMessageWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewTriggerMessageRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) TriggerMessage(ctx context.Context, body TriggerMessageJSONRequestBody) (*http.Response, error) {
	req, err := NewTriggerMessageRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetMessage(ctx context.Context, messageId string) (*http.Response, error) {
	req, err := NewGetMessageRequest(c.Server, messageId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewAddMessageRequest calls the generic AddMessage builder with application/json body
func NewAddMessageRequest(server string, body AddMessageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddMessageRequestWithBody(server, "application/json", bodyReader)
}

// NewAddMessageRequestWithBody generates requests for AddMessage with any type of body
func NewAddMessageRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/messages")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewTriggerMessageRequest calls the generic TriggerMessage builder with application/json body
func NewTriggerMessageRequest(server string, body TriggerMessageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewTriggerMessageRequestWithBody(server, "application/json", bodyReader)
}

// NewTriggerMessageRequestWithBody generates requests for TriggerMessage with any type of body
func NewTriggerMessageRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/messages/trigger")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewGetMessageRequest generates requests for GetMessage
func NewGetMessageRequest(server string, messageId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParam("simple", false, "message_id", messageId)
	if err != nil {
		return nil, err
	}

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/messages/%s", pathParam0)
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// AddMessage request  with any body
	AddMessageWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*AddMessageResponse, error)

	AddMessageWithResponse(ctx context.Context, body AddMessageJSONRequestBody) (*AddMessageResponse, error)

	// TriggerMessage request  with any body
	TriggerMessageWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*TriggerMessageResponse, error)

	TriggerMessageWithResponse(ctx context.Context, body TriggerMessageJSONRequestBody) (*TriggerMessageResponse, error)

	// GetMessage request
	GetMessageWithResponse(ctx context.Context, messageId string) (*GetMessageResponse, error)
}

type AddMessageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *MessageResponse
	JSON400      *Error
	JSON409      *Error
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r AddMessageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddMessageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type TriggerMessageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *MessageResponse
	JSON400      *Error
	JSON409      *Error
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r TriggerMessageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r TriggerMessageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetMessageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *MessageResponse
	JSON404      *Error
	JSON500      *Error
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetMessageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetMessageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// AddMessageWithBodyWithResponse request with arbitrary body returning *AddMessageResponse
func (c *ClientWithResponses) AddMessageWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*AddMessageResponse, error) {
	rsp, err := c.AddMessageWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseAddMessageResponse(rsp)
}

func (c *ClientWithResponses) AddMessageWithResponse(ctx context.Context, body AddMessageJSONRequestBody) (*AddMessageResponse, error) {
	rsp, err := c.AddMessage(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseAddMessageResponse(rsp)
}

// TriggerMessageWithBodyWithResponse request with arbitrary body returning *TriggerMessageResponse
func (c *ClientWithResponses) TriggerMessageWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*TriggerMessageResponse, error) {
	rsp, err := c.TriggerMessageWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseTriggerMessageResponse(rsp)
}

func (c *ClientWithResponses) TriggerMessageWithResponse(ctx context.Context, body TriggerMessageJSONRequestBody) (*TriggerMessageResponse, error) {
	rsp, err := c.TriggerMessage(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseTriggerMessageResponse(rsp)
}

// GetMessageWithResponse request returning *GetMessageResponse
func (c *ClientWithResponses) GetMessageWithResponse(ctx context.Context, messageId string) (*GetMessageResponse, error) {
	rsp, err := c.GetMessage(ctx, messageId)
	if err != nil {
		return nil, err
	}
	return ParseGetMessageResponse(rsp)
}

// ParseAddMessageResponse parses an HTTP response from a AddMessageWithResponse call
func ParseAddMessageResponse(rsp *http.Response) (*AddMessageResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &AddMessageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest MessageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseTriggerMessageResponse parses an HTTP response from a TriggerMessageWithResponse call
func ParseTriggerMessageResponse(rsp *http.Response) (*TriggerMessageResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &TriggerMessageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest MessageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetMessageResponse parses an HTTP response from a GetMessageWithResponse call
func ParseGetMessageResponse(rsp *http.Response) (*GetMessageResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetMessageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest MessageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
