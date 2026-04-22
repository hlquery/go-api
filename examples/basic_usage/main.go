/*
 * hlquery - Search beyond keywords.
 * Basic Usage Examples
 *
 * Demonstrates basic operations with the hlquery Go client
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hlquery/go-api"
)

func main() {
	// Initialize client
	baseURL := os.Getenv("HLQ_BASE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("HLQUERY_BASE_URL")
	}
	if baseURL == "" {
		baseURL = "http://localhost:9200"
	}
	client := hlquery.NewClient(baseURL)

	// Health check
	health, err := client.Health()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Health Status: %d\n", health.StatusCode)
	healthJSON, _ := json.MarshalIndent(health.Body, "", "    ")
	fmt.Printf("Health Body: %s\n", healthJSON)

	// List collections
	collections, err := client.ListCollections(0, 10)
	if err != nil {
		log.Fatal(err)
	}
	if collections.IsSuccess() {
		body := collections.Body
		count := 0
		if cols, ok := body["collections"].([]interface{}); ok {
			count = len(cols)
		}
		fmt.Printf("Found %d collections\n", count)
	}

	// With authentication
	authenticatedClient := hlquery.NewClient(baseURL, hlquery.ClientOptions{
		Token:      "your_token_here",
		AuthMethod: "bearer",
	})

	// Or set token dynamically
	client.SetAuthToken("your_token_here", "bearer")

	// Use authenticated client
	_, _ = authenticatedClient.Health()
}
