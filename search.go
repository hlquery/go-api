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

// Search provides search operations
type Search struct {
	client *Client
}

// Perform performs a search query
func (s *Search) Perform(collection string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/search", collection)
	return s.client.request("POST", path, params)
}

// Multi performs multiple searches
func (s *Search) Multi(searches []map[string]interface{}) (*Response, error) {
	body := map[string]interface{}{
		"searches": searches,
	}
	return s.client.request("POST", "/multi_search", body)
}

// Vector performs a vector search
func (s *Search) Vector(collection string, params map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/collections/%s/vector_search", collection)
	return s.client.request("POST", path, params)
}
