<div align="center">
  <img src="https://docs.hlquery.com/img/hlquery/2.png" alt="hlquery logo" width="200">
</div>

<div align="center">

**A clean, idiomatic Go client library for hlquery, designed with a familiar and intuitive API structure.**

[![Follow hlquery](https://img.shields.io/badge/Follow-%40hlquery-blue?logo=x&logoColor=white&labelColor=000000)](https://x.com/hlquery)
[![Go build](https://img.shields.io/badge/Go%20build-passing-brightgreen?logo=go&logoColor=white&labelColor=000000)](https://github.com/hlquery/go-api/actions/workflows/go-api.yml)
[![GitHub](https://img.shields.io/badge/GitHub-go--api-purple?logo=github&logoColor=white&labelColor=000000)](https://github.com/hlquery/go-api/)
[![hlquery](https://img.shields.io/badge/GitHub-hlquery-blue?logo=github&logoColor=white&labelColor=000000)](https://github.com/hlquery/hlquery/)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-a35a0f?logo=open-source-initiative&logoColor=white&labelColor=000000)](https://opensource.org/licenses/BSD-3-Clause)


</div>

### What is the hlquery Go API?

The hlquery Go API is the official Go client for [hlquery](https://github.com/hlquery/hlquery). It wraps the server's HTTP interface in a small, standard-library-friendly client that exposes collections, documents, search, and SQL helpers.

It is aimed at backend services, internal tools, and API servers that want hlquery integration without managing low-level HTTP details everywhere.

### Why use it?

Use it when you want a small, idiomatic Go surface with response helpers for status checks and parsed bodies, consistent auth handling, and direct access to both convenience methods and raw request execution.

### Install

```bash
$ go get github.com/hlquery/go-api/client
```

Or add it to `go.mod`:

```go
require github.com/hlquery/go-api v0.1.0
```

### Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"

    hlquery "github.com/hlquery/go-api/client"
)

func main() {
    baseURL := os.Getenv("HLQ_BASE_URL")
    if baseURL == "" {
        baseURL = os.Getenv("HLQUERY_BASE_URL")
    }
    if baseURL == "" {
        baseURL = "http://localhost:9200"
    }

    client := hlquery.NewClient(baseURL)

    health, err := client.Health()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %d\n", health.StatusCode)

    collections, err := client.ListCollections(0, 10)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(collections.Body)
}
```

### Auth

```go
client := hlquery.NewClient("http://localhost:9200", hlquery.ClientOptions{
    Token:      "your_token_here",
    AuthMethod: "bearer",
})

client.SetAuthToken("your_token_here", "bearer")
client.SetAuthToken("your_api_key_here", "api-key")
```

### SQL

```go
sqlAPI := client.SQLAPI()

rows, _ := sqlAPI.Query("SHOW COLLECTIONS;")
products, _ := sqlAPI.Search(
    "products",
    "SELECT id, title, price FROM products ORDER BY price DESC LIMIT 3;",
    nil,
)

fmt.Println(rows.Body)
fmt.Println(products.Body)
```

### Reduce Text Example

Use the raw request helper for custom module routes:

```go
moduleResponse, err := client.ExecuteRequest(
    "GET",
    "/modules/<name>/<route>?q=example",
    nil,
)
if err != nil {
    log.Fatal(err)
}

fmt.Println(moduleResponse.Body)
```

### Contributing

We welcome contributions from the community! All contributions must be released under the BSD 3-Clause license.

### How to Contribute

- Check existing [Go API issues](https://github.com/hlquery/go-api/issues) or create new ones
- Contribute Go client changes to [hlquery/go-api](https://github.com/hlquery/go-api)
- Contribute shared server/API changes to [hlquery/hlquery](https://github.com/hlquery/hlquery)
- Test and report bugs against the Go client
- Improve Go-specific documentation and examples

### Community

- [Documentation](https://docs.hlquery.com)
- [X (Twitter)](https://x.com/hlquery)
- [Go API GitHub](https://github.com/hlquery/go-api)
- [hlquery GitHub](https://github.com/hlquery/hlquery)

### License

The hlquery Go API is licensed under the [BSD 3-Clause License](https://opensource.org/licenses/BSD-3-Clause).
