/*
 * hlquery - Search beyond keywords.
 * https://www.hlquery.com
 *
 * Copyright (C) 2021-2026, Carlos F. Ferry <carlos.ferry@gmail.com>
 *
 * This file is part of hlquery Go API client, released under the BSD License version 3.
 * You are free to redistribute and/or modify this software
 * under the terms of the BSD License.
 * For more details, please visit: https://docs.hlquery.com
 */

package hlquery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// Client represents the hlquery API client
type Client struct {
	baseURL          string
	httpClient       *http.Client
	token            string
	authMethod       string // "bearer" or "api-key"
	maxResponseBytes int64
	initErr          error
}

// ClientOptions represents optional client configuration
type ClientOptions struct {
	Token            string
	AuthMethod       string // "bearer" or "api-key"
	Timeout          time.Duration
	HTTPClient       *http.Client
	MaxResponseBytes int64
}

const defaultMaxResponseBytes int64 = 64 << 20

// NewClient creates a new hlquery API client
func NewClient(baseURL string, options ...ClientOptions) *Client {
	client, err := NewClientWithError(baseURL, options...)
	if err != nil {
		client = newClientWithOptions(baseURL, options...)
		client.initErr = err
	}
	return client
}

// NewClientWithError creates a new client and validates the base URL.
func NewClientWithError(baseURL string, options ...ClientOptions) (*Client, error) {
	client := newClientWithOptions(baseURL, options...)
	if err := validateBaseURL(client.baseURL); err != nil {
		return client, err
	}
	return client, nil
}

func newClientWithOptions(baseURL string, options ...ClientOptions) *Client {
	client := &Client{
		baseURL:          strings.TrimRight(baseURL, "/"),
		httpClient:       &http.Client{Timeout: 30 * time.Second},
		authMethod:       "bearer",
		maxResponseBytes: defaultMaxResponseBytes,
	}

	if len(options) > 0 {
		opt := options[0]
		if opt.HTTPClient != nil {
			client.httpClient = opt.HTTPClient
		}
		if opt.Token != "" {
			client.token = opt.Token
		}
		if opt.AuthMethod != "" {
			client.authMethod = opt.AuthMethod
		}
		if opt.Timeout > 0 {
			client.httpClient.Timeout = opt.Timeout
		}
		if opt.MaxResponseBytes > 0 {
			client.maxResponseBytes = opt.MaxResponseBytes
		}
	}

	return client
}

func validateBaseURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid base URL: scheme must be http or https")
	}
	if parsed.Host == "" {
		return fmt.Errorf("invalid base URL: host is required")
	}
	return nil
}

// SetAuthToken sets the authentication token
func (c *Client) SetAuthToken(token string, authMethod string) {
	c.token = token
	if authMethod != "" {
		c.authMethod = authMethod
	}
}

// Response represents an API response
type Response struct {
	StatusCode int
	Body       map[string]interface{}
	Data       interface{}
	RawBody    []byte
	Headers    http.Header
}

// APIError describes a non-2xx API response.
type APIError struct {
	StatusCode int
	Message    string
	Response   *Response
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return fmt.Sprintf("hlquery API error: HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("hlquery API error: HTTP %d", e.StatusCode)
}

// IsSuccess returns true if the response status code is 2xx
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if the response status code is 4xx or 5xx
func (r *Response) IsError() bool {
	return r.StatusCode >= 400
}

// GetError returns the error message from the response
func (r *Response) GetError() string {
	if r.Body != nil {
		if err, ok := r.Body["error"].(string); ok {
			return err
		}
		if msg, ok := r.Body["message"].(string); ok {
			return msg
		}
	}
	return fmt.Sprintf("HTTP %d", r.StatusCode)
}

// GetStatusCode returns the HTTP status code.
func (r *Response) GetStatusCode() int {
	return r.StatusCode
}

// GetBody returns the decoded JSON object body, when the response body is an object.
func (r *Response) GetBody() map[string]interface{} {
	return r.Body
}

// request performs an HTTP request
func (c *Client) request(method, path string, body interface{}) (*Response, error) {
	return c.requestWithQuery(method, path, body, nil)
}

// requestWithContext performs an HTTP request with caller-controlled cancellation.
func (c *Client) requestWithContext(ctx context.Context, method, path string, body interface{}) (*Response, error) {
	return c.requestWithQueryValuesContext(ctx, method, path, body, nil)
}

// requestWithQuery performs an HTTP request with optional query parameters.
func (c *Client) requestWithQuery(method, path string, body interface{}, queryParams map[string]string) (*Response, error) {
	params := map[string]interface{}{}
	for key, value := range queryParams {
		params[key] = value
	}
	return c.requestWithQueryValues(method, path, body, params)
}

// requestWithQueryValues performs an HTTP request with flexible query parameters.
func (c *Client) requestWithQueryValues(method, path string, body interface{}, queryParams map[string]interface{}) (*Response, error) {
	return c.requestWithQueryValuesContext(context.Background(), method, path, body, queryParams)
}

func (c *Client) requestWithQueryValuesContext(ctx context.Context, method, path string, body interface{}, queryParams map[string]interface{}) (*Response, error) {
	if c.initErr != nil {
		return nil, c.initErr
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	fullURL := strings.TrimRight(c.baseURL, "/") + path
	if len(queryParams) > 0 {
		values := url.Values{}
		for key, value := range queryParams {
			addQueryValue(values, key, value)
		}
		fullURL += "?" + values.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		if c.authMethod == "api-key" {
			req.Header.Set("X-API-Key", c.token)
		} else {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	reader := io.Reader(resp.Body)
	if c.maxResponseBytes > 0 {
		reader = io.LimitReader(resp.Body, c.maxResponseBytes+1)
	}

	respBody, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if c.maxResponseBytes > 0 && int64(len(respBody)) > c.maxResponseBytes {
		return nil, fmt.Errorf("response body exceeds %d bytes", c.maxResponseBytes)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		RawBody:    respBody,
		Headers:    resp.Header,
	}

	// Try to parse JSON response
	var jsonBody interface{}
	if err := json.Unmarshal(respBody, &jsonBody); err == nil {
		response.Data = jsonBody
		if objectBody, ok := jsonBody.(map[string]interface{}); ok {
			response.Body = objectBody
		}
	}

	if !response.IsSuccess() {
		return response, &APIError{
			StatusCode: response.StatusCode,
			Message:    response.GetError(),
			Response:   response,
		}
	}

	return response, nil
}

func addQueryValue(values url.Values, key string, value interface{}) {
	if value == nil {
		return
	}

	switch typed := value.(type) {
	case string:
		values.Set(key, typed)
	case []string:
		for _, item := range typed {
			values.Add(key, item)
		}
	case []interface{}:
		for _, item := range typed {
			addQueryValue(values, key, item)
		}
	case bool:
		if typed {
			values.Set(key, "true")
		} else {
			values.Set(key, "false")
		}
	default:
		kind := reflect.TypeOf(value).Kind()
		if kind == reflect.Map || kind == reflect.Slice || kind == reflect.Array {
			if encoded, err := json.Marshal(value); err == nil {
				values.Set(key, string(encoded))
				return
			}
		}
		values.Set(key, fmt.Sprint(typed))
	}
}

func encodePathPart(value string) string {
	return url.PathEscape(value)
}

func copyMap(input map[string]interface{}) map[string]interface{} {
	output := map[string]interface{}{}
	for key, value := range input {
		output[key] = value
	}
	return output
}

func mapFromJSONStruct(input interface{}) map[string]interface{} {
	encoded, err := json.Marshal(input)
	if err != nil {
		return map[string]interface{}{}
	}
	var output map[string]interface{}
	if err := json.Unmarshal(encoded, &output); err != nil {
		return map[string]interface{}{}
	}
	return output
}

func requireNonEmpty(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s must be a non-empty string", field)
	}
	return nil
}

// CollectionSchema is a typed helper for common collection creation calls.
type CollectionSchema struct {
	Name                string                   `json:"name,omitempty"`
	Fields              []map[string]interface{} `json:"fields,omitempty"`
	SearchableFields    []string                 `json:"searchable_fields,omitempty"`
	DefaultSortingField string                   `json:"default_sorting_field,omitempty"`
}

// SearchParams is a typed helper for common search calls.
type SearchParams struct {
	Query    string `json:"q,omitempty"`
	QueryBy  string `json:"query_by,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	FilterBy string `json:"filter_by,omitempty"`
	SortBy   string `json:"sort_by,omitempty"`
}

// Document is a typed alias for JSON document payloads.
type Document map[string]interface{}

// Health checks the server health status
func (c *Client) Health() (*Response, error) {
	return c.request("GET", "/health", nil)
}

// Status gets server status.
func (c *Client) Status() (*Response, error) {
	return c.request("GET", "/status", nil)
}

// Ready checks server readiness.
func (c *Client) Ready() (*Response, error) {
	return c.request("GET", "/ready", nil)
}

// Ping checks basic server reachability.
func (c *Client) Ping() (*Response, error) {
	return c.request("GET", "/ping", nil)
}

// Stats gets server statistics
func (c *Client) Stats() (*Response, error) {
	return c.request("GET", "/stats", nil)
}

// Info gets server information
func (c *Client) Info() (*Response, error) {
	return c.request("GET", "/", nil)
}

// Etc gets protocol and runtime metadata.
func (c *Client) Etc() (*Response, error) {
	return c.request("GET", "/etc", nil)
}

// Links lists configured cluster links
func (c *Client) Links() (*Response, error) {
	return c.request("GET", "/links", nil)
}

// LinksPing pings configured cluster links
func (c *Client) LinksPing() (*Response, error) {
	return c.request("GET", "/links/ping", nil)
}

// SQL executes a top-level SQL query through GET /sql.
func (c *Client) SQL(sql string, queryParams ...map[string]string) (*Response, error) {
	if sql == "" {
		return nil, fmt.Errorf("SQL query must be a non-empty string")
	}

	params := map[string]string{"sql": sql}
	if len(queryParams) > 0 {
		for key, value := range queryParams[0] {
			params[key] = value
		}
	}

	return c.requestWithQuery("GET", "/sql", nil, params)
}

// ExecSQL executes a top-level SQL statement through POST /sql.
func (c *Client) ExecSQL(sql string) (*Response, error) {
	if sql == "" {
		return nil, fmt.Errorf("SQL query must be a non-empty string")
	}

	return c.request("POST", "/sql", map[string]interface{}{"exec": sql})
}

// LinksConnect adds a cluster link (in-memory only)
func (c *Client) LinksConnect(endpointOrHost string, port ...int) (*Response, error) {
	if err := requireNonEmpty(endpointOrHost, "endpoint or host"); err != nil {
		return nil, err
	}
	body := map[string]interface{}{}
	if len(port) > 0 {
		body["host"] = endpointOrHost
		body["port"] = port[0]
	} else {
		body["endpoint"] = endpointOrHost
	}
	return c.request("POST", "/links/connect", body)
}

// LinksDisconnect removes a cluster link (in-memory only)
func (c *Client) LinksDisconnect(endpointOrHost string, port ...int) (*Response, error) {
	if err := requireNonEmpty(endpointOrHost, "endpoint or host"); err != nil {
		return nil, err
	}
	body := map[string]interface{}{}
	if len(port) > 0 {
		body["host"] = endpointOrHost
		body["port"] = port[0]
	} else {
		body["endpoint"] = endpointOrHost
	}
	return c.request("POST", "/links/disconnect", body)
}

// Collections returns a Collections API client
func (c *Client) Collections() *Collections {
	return &Collections{client: c}
}

// Documents returns a Documents API client
func (c *Client) Documents() *Documents {
	return &Documents{client: c}
}

// Search returns a Search API client
func (c *Client) Search() *Search {
	return &Search{client: c}
}

// SQLAPI returns a SQL API client
func (c *Client) SQLAPI() *SQLAPI {
	return &SQLAPI{client: c}
}

// System returns a System API client.
func (c *Client) System() *System {
	return &System{client: c}
}

// SAMAPI returns a SAM API client.
func (c *Client) SAMAPI() *SAM {
	return &SAM{client: c}
}

// SAM returns a SAM API client.
func (c *Client) SAM() *SAM {
	return c.SAMAPI()
}

// Aliases returns an Aliases API client.
func (c *Client) Aliases() *Aliases {
	return &Aliases{client: c}
}

// Synonyms returns a Synonyms API client.
func (c *Client) Synonyms() *Synonyms {
	return &Synonyms{client: c}
}

// Stopwords returns a Stopwords API client.
func (c *Client) Stopwords() *Stopwords {
	return &Stopwords{client: c}
}

// Overrides returns an Overrides API client.
func (c *Client) Overrides() *Overrides {
	return &Overrides{client: c}
}

// Keys returns an API keys client.
func (c *Client) Keys() *Keys {
	return &Keys{client: c}
}

// Users returns a Users API client.
func (c *Client) Users() *Users {
	return &Users{client: c}
}

// Modules returns a Modules API client.
func (c *Client) Modules() *Modules {
	return &Modules{client: c}
}

// Convenience methods

// ListCollections lists all collections
func (c *Client) ListCollections(offset, limit int) (*Response, error) {
	return c.requestWithQueryValues("GET", "/collections", nil, map[string]interface{}{"offset": offset, "limit": limit})
}

// ListCollectionsDistributed lists collections across all configured nodes
func (c *Client) ListCollectionsDistributed() (*Response, error) {
	return c.request("GET", "/collections/distributed", nil)
}

// GetCollection gets a collection by name
func (c *Client) GetCollection(name string) (*Response, error) {
	if err := requireNonEmpty(name, "collection name"); err != nil {
		return nil, err
	}
	return c.request("GET", "/collections/"+encodePathPart(name), nil)
}

// CreateCollection creates a new collection
func (c *Client) CreateCollection(name string, schema map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(name, "collection name"); err != nil {
		return nil, err
	}
	body := copyMap(schema)
	body["name"] = name
	return c.request("POST", "/collections", body)
}

// CreateCollectionSchema creates a new collection using a typed schema payload.
func (c *Client) CreateCollectionSchema(schema CollectionSchema) (*Response, error) {
	return c.Collections().CreateSchema(schema)
}

// DeleteCollection deletes a collection
func (c *Client) DeleteCollection(name string) (*Response, error) {
	if err := requireNonEmpty(name, "collection name"); err != nil {
		return nil, err
	}
	return c.request("DELETE", "/collections/"+encodePathPart(name), nil)
}

// UpdateCollection updates a collection schema.
func (c *Client) UpdateCollection(name string, schema map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(name, "collection name"); err != nil {
		return nil, err
	}
	return c.request("POST", "/collections/"+encodePathPart(name)+"/update", schema)
}

// GetCollectionLanguage gets language metadata for a collection.
func (c *Client) GetCollectionLanguage(name string) (*Response, error) {
	if err := requireNonEmpty(name, "collection name"); err != nil {
		return nil, err
	}
	return c.request("GET", "/collections/"+encodePathPart(name)+"/lang", nil)
}

// ListDocuments lists documents in a collection
func (c *Client) ListDocuments(collection string, offset, limit int) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/collections/%s/documents", encodePathPart(collection))
	return c.requestWithQueryValues("GET", path, nil, map[string]interface{}{"offset": offset, "limit": limit})
}

// GetDocument gets a document by ID
func (c *Client) GetDocument(collection, docID string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(docID, "document id"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/collections/%s/documents/%s", encodePathPart(collection), encodePathPart(docID))
	return c.request("GET", path, nil)
}

// AddDocument adds a document to a collection
func (c *Client) AddDocument(collection string, document map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/collections/%s/documents", encodePathPart(collection))
	return c.request("POST", path, document)
}

// UpdateDocument updates a document
func (c *Client) UpdateDocument(collection, docID string, document map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(docID, "document id"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/collections/%s/documents/%s", encodePathPart(collection), encodePathPart(docID))
	return c.request("PUT", path, document)
}

// DeleteDocument deletes a document
func (c *Client) DeleteDocument(collection, docID string) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	if err := requireNonEmpty(docID, "document id"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/collections/%s/documents/%s", encodePathPart(collection), encodePathPart(docID))
	return c.request("DELETE", path, nil)
}

// ImportDocuments imports multiple documents
func (c *Client) ImportDocuments(collection string, documents []map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"documents": documents,
	}
	path := fmt.Sprintf("/collections/%s/documents/import", encodePathPart(collection))
	return c.request("POST", path, body)
}

// Search performs a search query
func (c *Client) SearchDocuments(collection string, params map[string]interface{}) (*Response, error) {
	if err := requireNonEmpty(collection, "collection name"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/collections/%s/documents/search", encodePathPart(collection))
	return c.request("POST", path, params)
}

// SearchDocumentsTyped performs a search using typed parameters.
func (c *Client) SearchDocumentsTyped(collection string, params SearchParams) (*Response, error) {
	return c.Search().PerformTyped(collection, params)
}

// SQLSearch executes a collection-bound SQL SELECT through the search endpoint.
func (c *Client) SQLSearch(collection, sql string, params ...map[string]interface{}) (*Response, error) {
	if len(params) > 0 {
		return c.Search().SQL(collection, sql, params[0])
	}
	return c.Search().SQL(collection, sql, nil)
}

// ExecuteRequest performs an arbitrary HTTP request
func (c *Client) ExecuteRequest(method, path string, body interface{}) (*Response, error) {
	return c.request(method, path, body)
}

// ExecuteRequestWithContext performs an arbitrary HTTP request with caller-controlled cancellation.
func (c *Client) ExecuteRequestWithContext(ctx context.Context, method, path string, body interface{}) (*Response, error) {
	return c.requestWithContext(ctx, method, path, body)
}

// ExecuteRequestWithQuery performs an arbitrary HTTP request with query parameters.
func (c *Client) ExecuteRequestWithQuery(method, path string, body interface{}, queryParams map[string]interface{}) (*Response, error) {
	return c.requestWithQueryValues(method, path, body, queryParams)
}

// ExecuteRequestWithQueryContext performs an arbitrary HTTP request with query params and caller-controlled cancellation.
func (c *Client) ExecuteRequestWithQueryContext(ctx context.Context, method, path string, body interface{}, queryParams map[string]interface{}) (*Response, error) {
	return c.requestWithQueryValuesContext(ctx, method, path, body, queryParams)
}

// Flush flushes all data from the database
func (c *Client) Flush() (*Response, error) {
	return c.request("POST", "/flush", nil)
}
