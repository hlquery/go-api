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
	if params == nil {
		params = map[string]interface{}{"offset": 0, "limit": 10}
	}
	path := fmt.Sprintf("/collections/%s/documents", encodePathPart(collection))
	return d.client.requestWithQueryValues("GET", path, nil, params)
}

// Get gets a document by ID
func (d *Documents) Get(collection, docID string) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", encodePathPart(collection), encodePathPart(docID))
	return d.client.request("GET", path, nil)
}

// Context gets contextual phrases for a document.
func (d *Documents) Context(collection, docID string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s/context", encodePathPart(collection), encodePathPart(docID))
	return d.client.requestWithQueryValues("GET", path, nil, params)
}

// Add adds a document to a collection
func (d *Documents) Add(collection string, document map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents", encodePathPart(collection))
	return d.client.request("POST", path, document)
}

// Update updates a document
func (d *Documents) Update(collection, docID string, document map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", encodePathPart(collection), encodePathPart(docID))
	return d.client.request("PUT", path, document)
}

// Delete deletes a document
func (d *Documents) Delete(collection, docID string) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", encodePathPart(collection), encodePathPart(docID))
	return d.client.request("DELETE", path, nil)
}

// DeleteByFilter deletes documents matching a filter expression.
func (d *Documents) DeleteByFilter(collection, filter string) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents", encodePathPart(collection))
	return d.client.requestWithQueryValues("DELETE", path, nil, map[string]interface{}{"filter_by": filter})
}

// UpdateByQuery updates documents matching a query/filter.
func (d *Documents) UpdateByQuery(collection string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/_update_by_query", encodePathPart(collection))
	return d.client.request("POST", path, params)
}

// DeleteByQuery deletes documents matching a query/filter.
func (d *Documents) DeleteByQuery(collection string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/_delete_by_query", encodePathPart(collection))
	return d.client.request("POST", path, params)
}

// FacetCounts computes facet counts for a collection.
func (d *Documents) FacetCounts(collection string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/facet_counts", encodePathPart(collection))
	return d.client.request("POST", path, params)
}

// Export exports documents from a collection.
func (d *Documents) Export(collection string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/export", encodePathPart(collection))
	return d.client.request("POST", path, params)
}

// Maybe returns suggestion candidates for a query in a collection.
func (d *Documents) Maybe(collection string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/documents/maybe", encodePathPart(collection))
	return d.client.request("POST", path, params)
}

// Import imports multiple documents (bulk import)
func (d *Documents) Import(collection string, documents []map[string]interface{}) (*Response, error) {
	body := map[string]interface{}{
		"documents": documents,
	}
	path := fmt.Sprintf("/collections/%s/documents/import", encodePathPart(collection))
	return d.client.request("POST", path, body)
}
