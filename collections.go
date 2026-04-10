/*
 * hlquery - Search beyond keywords.
 * http://www.hlquery.com
 *
 * Copyright (C) 2021-2026, Carlos F. Ferry <carlos.ferry@gmail.com>
 *
 * This file is part of hlquery Go API client, released under the BSD License version 3.
 */

package hlquery

import "fmt"

// Collections provides collection management operations
type Collections struct {
	client *Client
}

// List lists all collections
func (c *Collections) List(offset, limit int) (*Response, error) {
	path := fmt.Sprintf("/collections?offset=%d&limit=%d", offset, limit)
	return c.client.request("GET", path, nil)
}

// Get gets a collection by name
func (c *Collections) Get(name string) (*Response, error) {
	return c.client.request("GET", "/collections/"+name, nil)
}

// Create creates a new collection
func (c *Collections) Create(name string, schema map[string]interface{}) (*Response, error) {
	body := map[string]interface{}{
		"name":   name,
		"fields": schema["fields"],
	}
	return c.client.request("POST", "/collections", body)
}

// Delete deletes a collection
func (c *Collections) Delete(name string) (*Response, error) {
	return c.client.request("DELETE", "/collections/"+name, nil)
}

// Update updates a collection schema
func (c *Collections) Update(name string, schema map[string]interface{}) (*Response, error) {
	body := map[string]interface{}{
		"fields": schema["fields"],
	}
	return c.client.request("PUT", "/collections/"+name, body)
}

// GetFields gets formatted fields for a collection
func (c *Collections) GetFields(name string) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/fields", name)
	return c.client.request("GET", path, nil)
}
