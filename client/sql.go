/*
 * hlquery - Search beyond keywords.
 * https://www.hlquery.com
 *
 * Copyright (C) 2021-2026, Carlos F. Ferry <carlos.ferry@gmail.com>
 *
 * This file is part of hlquery Go API client, released under the BSD License version 3.
 */

package hlquery

// SQLAPI provides top-level and collection-bound SQL operations.
type SQLAPI struct {
	client *Client
}

// Query executes a top-level SQL query through GET /sql.
func (s *SQLAPI) Query(sql string, queryParams ...map[string]string) (*Response, error) {
	return s.client.SQL(sql, queryParams...)
}

// Exec executes a top-level SQL statement through POST /sql.
func (s *SQLAPI) Exec(sql string) (*Response, error) {
	return s.client.ExecSQL(sql)
}

// Search executes a collection-bound SQL SELECT through the search endpoint.
func (s *SQLAPI) Search(collection, sql string, params map[string]interface{}) (*Response, error) {
	return s.client.Search().SQL(collection, sql, params)
}
