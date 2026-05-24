/*
 * hlquery - Search beyond keywords.
 * https://www.hlquery.com
 *
 * Copyright (C) 2021-2026, Carlos F. Ferry <carlos.ferry@gmail.com>
 *
 * This file is part of hlquery Go API client, released under the BSD License version 3.
 */

package hlquery

// Collections provides collection management operations
type Collections struct {
	client *Client
}

// List lists all collections
func (c *Collections) List(offset, limit int) (*Response, error) {
	return c.client.requestWithQueryValues("GET", "/collections", nil, map[string]interface{}{"offset": offset, "limit": limit})
}

// ListDistributed lists collections across configured nodes.
func (c *Collections) ListDistributed() (*Response, error) {
	return c.client.request("GET", "/collections/distributed", nil)
}

// Get gets a collection by name
func (c *Collections) Get(name string) (*Response, error) {
	return c.client.request("GET", "/collections/"+encodePathPart(name), nil)
}

// Create creates a new collection
func (c *Collections) Create(name string, schema map[string]interface{}) (*Response, error) {
	body := copyMap(schema)
	body["name"] = name
	return c.client.request("POST", "/collections", body)
}

// Delete deletes a collection
func (c *Collections) Delete(name string) (*Response, error) {
	return c.client.request("DELETE", "/collections/"+encodePathPart(name), nil)
}

// Update updates a collection schema
func (c *Collections) Update(name string, schema map[string]interface{}) (*Response, error) {
	return c.client.request("POST", "/collections/"+encodePathPart(name)+"/update", schema)
}

// GetLanguage gets language metadata for a collection.
func (c *Collections) GetLanguage(name string) (*Response, error) {
	return c.client.request("GET", "/collections/"+encodePathPart(name)+"/lang", nil)
}

// GetFields gets formatted fields for a collection
func (c *Collections) GetFields(name string) (*Response, error) {
	return c.Get(name)
}
