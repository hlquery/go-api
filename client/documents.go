/*
 * hlquery - Search beyond keywords.
 * https://www.hlquery.com
 *
 * Copyright (C) 2021-2026, Carlos F. Ferry <carlos.ferry@gmail.com>
 *
 * This file is part of hlquery Go API client, released under the BSD License version 3.
 */

package hlquery

import "fmt"

// Documents provides document CRUD operations
type Documents struct {
	client *Client
}

// List lists documents in a collection
func (d *Documents) List(collection string, params map[string]interface{}) (*Response, error) {
	offset := 0
	limit := 10
	if o, ok := params["offset"].(int); ok {
		offset = o
	}
	if l, ok := params["limit"].(int); ok {
		limit = l
	}
	path := fmt.Sprintf("/collections/%s/documents?offset=%d&limit=%d", collection, offset, limit)
	return d.client.request("GET", path, nil)
}

// Get gets a document by ID
func (d *Documents) Get(collection, docID string) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	return d.client.request("GET", path, nil)
}

// Add adds a document to a collection
func (d *Documents) Add(collection string, document map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents", collection)
	return d.client.request("POST", path, document)
}

// Update updates a document
func (d *Documents) Update(collection, docID string, document map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	return d.client.request("PUT", path, document)
}

// Delete deletes a document
func (d *Documents) Delete(collection, docID string) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	return d.client.request("DELETE", path, nil)
}

// Import imports multiple documents (bulk import)
func (d *Documents) Import(collection string, documents []map[string]interface{}) (*Response, error) {
	body := map[string]interface{}{
		"documents": documents,
	}
	path := fmt.Sprintf("/collections/%s/documents/import", collection)
	return d.client.request("POST", path, body)
}
