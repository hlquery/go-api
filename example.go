/*
 * hlquery - Search beyond keywords.
 * https://www.hlquery.com
 *
 * Copyright (C) 2021-2026, Carlos F. Ferry <carlos.ferry@gmail.com>
 *
 * This file is part of hlquery Go API client, released under the BSD License version 3.
 */

package main

import (
	"encoding/json"
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

	fmt.Printf("Health Status: %d\n", health.StatusCode)
	healthJSON, _ := json.MarshalIndent(health.Body, "", "    ")
	fmt.Printf("Health Body: %s\n", healthJSON)

	collections, err := client.ListCollections(0, 10)
	if err != nil {
		log.Fatal(err)
	}

	collectionsJSON, _ := json.MarshalIndent(collections.Body, "", "    ")
	fmt.Printf("Collections: %s\n", collectionsJSON)
}
