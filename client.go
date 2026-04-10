/*
 * hlquery - Search beyond keywords.
 * http://www.hlquery.com
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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the hlquery API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	authMethod string // "bearer" or "api-key"
}

// ClientOptions represents optional client configuration
type ClientOptions struct {
	Token      string
	AuthMethod string // "bearer" or "api-key"
	Timeout    time.Duration
}

// NewClient creates a new hlquery API client
func NewClient(baseURL string, options ...ClientOptions) *Client {
	client := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		authMethod: "bearer",
	}

	if len(options) > 0 {
		opt := options[0]
		if opt.Token != "" {
			client.token = opt.Token
		}
		if opt.AuthMethod != "" {
			client.authMethod = opt.AuthMethod
		}
		if opt.Timeout > 0 {
			client.httpClient.Timeout = opt.Timeout
		}
	}

	return client
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
	RawBody    []byte
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

// request performs an HTTP request
func (c *Client) request(method, path string, body interface{}) (*Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		RawBody:    respBody,
	}

	// Try to parse JSON response
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(respBody, &jsonBody); err == nil {
		response.Body = jsonBody
	}

	return response, nil
}

// Health checks the server health status
func (c *Client) Health() (*Response, error) {
	return c.request("GET", "/health", nil)
}

// Stats gets server statistics
func (c *Client) Stats() (*Response, error) {
	return c.request("GET", "/stats", nil)
}

// Info gets server information
func (c *Client) Info() (*Response, error) {
	return c.request("GET", "/", nil)
}

// Links lists configured cluster links
func (c *Client) Links() (*Response, error) {
	return c.request("GET", "/links", nil)
}

// LinksPing pings configured cluster links
func (c *Client) LinksPing() (*Response, error) {
	return c.request("GET", "/links/ping", nil)
}

// LinksConnect adds a cluster link (in-memory only)
func (c *Client) LinksConnect(endpointOrHost string, port ...int) (*Response, error) {
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

// Convenience methods

// ListCollections lists all collections
func (c *Client) ListCollections(offset, limit int) (*Response, error) {
	path := fmt.Sprintf("/collections?offset=%d&limit=%d", offset, limit)
	return c.request("GET", path, nil)
}

// ListCollectionsDistributed lists collections across all configured nodes
func (c *Client) ListCollectionsDistributed() (*Response, error) {
	return c.request("GET", "/collections/distributed", nil)
}

// GetCollection gets a collection by name
func (c *Client) GetCollection(name string) (*Response, error) {
	return c.request("GET", "/collections/"+name, nil)
}

// CreateCollection creates a new collection
func (c *Client) CreateCollection(name string, schema map[string]interface{}) (*Response, error) {
	body := map[string]interface{}{
		"name":   name,
		"fields": schema["fields"],
	}
	return c.request("POST", "/collections", body)
}

// DeleteCollection deletes a collection
func (c *Client) DeleteCollection(name string) (*Response, error) {
	return c.request("DELETE", "/collections/"+name, nil)
}

// ListDocuments lists documents in a collection
func (c *Client) ListDocuments(collection string, offset, limit int) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents?offset=%d&limit=%d", collection, offset, limit)
	return c.request("GET", path, nil)
}

// GetDocument gets a document by ID
func (c *Client) GetDocument(collection, docID string) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	return c.request("GET", path, nil)
}

// AddDocument adds a document to a collection
func (c *Client) AddDocument(collection string, document map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents", collection)
	return c.request("POST", path, document)
}

// UpdateDocument updates a document
func (c *Client) UpdateDocument(collection, docID string, document map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	return c.request("PUT", path, document)
}

// DeleteDocument deletes a document
func (c *Client) DeleteDocument(collection, docID string) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	return c.request("DELETE", path, nil)
}

// ImportDocuments imports multiple documents
func (c *Client) ImportDocuments(collection string, documents []map[string]interface{}) (*Response, error) {
	body := map[string]interface{}{
		"documents": documents,
	}
	path := fmt.Sprintf("/collections/%s/documents/import", collection)
	return c.request("POST", path, body)
}

// Search performs a search query
func (c *Client) SearchDocuments(collection string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/search", collection)
	return c.request("POST", path, params)
}

// ExecuteRequest performs an arbitrary HTTP request
func (c *Client) ExecuteRequest(method, path string, body interface{}) (*Response, error) {
	return c.request(method, path, body)
}

// Flush flushes all data from the database
func (c *Client) Flush() (*Response, error) {
	return c.request("POST", "/flush", nil)
}
