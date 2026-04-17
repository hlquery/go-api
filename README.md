<div align="center">
  <img src="https://docs.hlquery.com/img/hlquery/2.png" alt="hlquery logo" width="200">
</div>

<div align="center">

**A clean, idiomatic Go client library for hlquery, designed with a familiar and intuitive API structure.**

[![Follow hlquery](https://img.shields.io/badge/Follow-%40hlquery-blue?logo=x&logoColor=white)](https://x.com/hlquery)
[![Commit Activity](https://img.shields.io/github/commit-activity/m/hlquery/hlquery)](https://github.com/hlquery/go-api/pulse)
[![hlquery](https://img.shields.io/badge/GitHub-hlquery-181717?logo=github&logoColor=white)](https://github.com/hlquery/hlquery/stargazers)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

</div>

# hlquery Go API Client

## Features

- **No External Dependencies**: Uses only Go standard library (except uuid for examples)
- **Type-safe Responses**: Response objects with helper methods
- **Consistent API Design**: Familiar structure for Go developers
- **Authentication Support**: Bearer token and X-API-Key authentication
- **Comprehensive Operations**: Collections, Documents, and Search APIs
- **Error Handling**: Clear error messages and status codes

## Installation

```bash
go get github.com/hlquery/go-api
```

Or add to your `go.mod`:

```go
require github.com/hlquery/go-api v0.1.0
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/hlquery/go-api"
)

func main() {
    // Initialize client
    client := hlquery.NewClient("http://localhost:9200")
    
    // Health check
    health, err := client.Health()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %d\n", health.StatusCode)
    
    // List collections
    collections, err := client.ListCollections(0, 10)
    if err != nil {
        log.Fatal(err)
    }
    if collections.IsSuccess() {
        body := collections.Body
        fmt.Printf("Found collections: %v\n", body["collections"])
    }
}
```

### With Authentication

```go
// Method 1: Set token in constructor
client := hlquery.NewClient("http://localhost:9200", hlquery.ClientOptions{
    Token:      "your_token_here",
    AuthMethod: "bearer", // or "api-key"
})

// Method 2: Set token dynamically
client := hlquery.NewClient("http://localhost:9200")
client.SetAuthToken("your_token_here", "bearer")

// Method 3: Use X-API-Key
client.SetAuthToken("your_token_here", "api-key")
```

### Reduce Text Example

If the `ai_search` module is enabled, you can call it with `ExecuteRequest` and a query string:

```go
package main

import (
    "fmt"
    "log"
    "net/url"

    hlquery "github.com/hlquery/go-api"
)

func main() {
    client := hlquery.NewClient("http://localhost:9200")

    path := "/modules/ai_search/talk?q=" +
        url.QueryEscape("summarize onboarding guide in docs") +
        "&run=true"

    summary, err := client.ExecuteRequest("GET", path, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(summary.Body)
}
```

## API Reference

### Client Initialization

```go
client := hlquery.NewClient(baseURL string, options ...ClientOptions)
```

**Parameters**:
- `baseURL` (string, required): Base URL of hlquery server (e.g., `"http://localhost:9200"`)
- `options` (ClientOptions, optional): Client options
  - `Token` (string): Authentication token
  - `AuthMethod` (string): Authentication method (`"bearer"` or `"api-key"`)
  - `Timeout` (time.Duration): Request timeout

### Health & Status

#### `Health()`

Check server health status.

```go
health, err := client.Health()
if err != nil {
    log.Fatal(err)
}
if health.IsSuccess() {
    body := health.Body
    fmt.Printf("Status: %v\n", body["status"])
}
```

#### `Stats()`

Get server statistics.

```go
stats, err := client.Stats()
```

#### `Info()`

Get server information.

```go
info, err := client.Info()
```

### Collections API

#### Using Collections API Object

```go
collections := client.Collections()

// List collections
result, err := collections.List(0, 10)

// Get collection
result, err := collections.Get("my_collection")

// Create collection
schema := map[string]interface{}{
    "fields": []map[string]interface{}{
        {"name": "title", "type": "string"},
        {"name": "price", "type": "float"},
    },
}
result, err := collections.Create("new_collection", schema)

// Delete collection
result, err := collections.Delete("collection_name")

// Update collection
result, err := collections.Update("collection_name", schema)

// Get formatted fields
result, err := collections.GetFields("my_collection")
```

#### Direct Methods

```go
// List collections
collections, err := client.ListCollections(offset, limit)

// Get collection
collection, err := client.GetCollection("collection_name")

// Create collection
result, err := client.CreateCollection("collection_name", schema)

// Delete collection
result, err := client.DeleteCollection("collection_name")
```

### Documents API

#### Using Documents API Object

```go
documents := client.Documents()

// List documents
params := map[string]interface{}{
    "offset": 0,
    "limit":  10,
}
result, err := documents.List("collection_name", params)

// Get document
result, err := documents.Get("collection_name", "document_id")

// Add document
doc := map[string]interface{}{
    "id":    "doc1",
    "title": "Example",
    "price": 99.99,
}
result, err := documents.Add("collection_name", doc)

// Update document
result, err := documents.Update("collection_name", "doc_id", doc)

// Delete document
result, err := documents.Delete("collection_name", "doc_id")

// Import documents (bulk)
docs := []map[string]interface{}{doc1, doc2, doc3}
result, err := documents.Import("collection_name", docs)
```

#### Direct Methods

```go
// List documents
docs, err := client.ListDocuments("collection_name", offset, limit)

// Get document
doc, err := client.GetDocument("collection_name", "doc_id")

// Add document
result, err := client.AddDocument("collection_name", doc)

// Update document
result, err := client.UpdateDocument("collection_name", "doc_id", doc)

// Delete document
result, err := client.DeleteDocument("collection_name", "doc_id")

// Import documents
result, err := client.ImportDocuments("collection_name", docs)
```

### Search API

#### Using Search API Object

```go
search := client.Search()

// Perform search
params := map[string]interface{}{
    "q":         "laptop",
    "query_by":  "title,description",
    "filter_by": "price:>1000",
    "sort_by":   "price:asc",
    "limit":     10,
}
result, err := search.Perform("collection_name", params)

// Multi-search
searches := []map[string]interface{}{
    {"collection": "products", "q": "laptop"},
    {"collection": "articles", "q": "laptop"},
}
result, err := search.Multi(searches)

// Vector search
params := map[string]interface{}{
    "vector": []float64{0.1, 0.2, 0.3},
    "limit":  10,
}
result, err := search.Vector("collection_name", params)
```

#### Direct Method

```go
// Search
params := map[string]interface{}{
    "q":        "laptop",
    "query_by": "title",
}
results, err := client.SearchDocuments("collection_name", params)
```

### Ranking helpers

`ComputeRankSignal` and `AttachRankSort` live in the root package so you can recompute `rank_signal` every time new hits arrive and let the API sort by that field.

```go
params := map[string]interface{}{"q": "guide"}
signal := hlquery.ComputeRankSignal(float64(popularity), float64(hitLog), nil)
params["rank_signal"] = signal
hlquery.AttachRankSort(params, "rank_signal", "desc")
results, err := client.SearchDocuments("collection_name", params)
```

## Response Objects

All API methods return a `Response` object with helper methods:

```go
response, err := client.Health()

// Get HTTP status code
statusCode := response.StatusCode

// Get response body
body := response.Body

// Check if successful
if response.IsSuccess() {
    // Handle success
}

// Check if error
if response.IsError() {
    error := response.GetError()
    fmt.Printf("Error: %s\n", error)
}
```

## Error Handling

The client returns errors for network issues and HTTP errors:

```go
result, err := client.CreateCollection("collection", schema)
if err != nil {
    // Handle network/request errors
    log.Fatal(err)
}

if result.IsError() {
    // Handle HTTP errors (4xx, 5xx)
    error := result.GetError()
    fmt.Printf("API Error: %s\n", error)
}
```

## Examples

See the `examples/` directory for complete examples:

- `basic_usage.go` - Basic usage examples
- `collections.go` - Collection management
- `documents.go` - Document CRUD operations
- `search.go` - Search operations

Run examples:

```bash
cd examples
go run basic_usage.go
go run collections.go
go run documents.go
```

## Requirements

- Go >= 1.21
