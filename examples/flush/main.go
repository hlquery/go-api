/*
 * hlquery - Search beyond keywords.
 * Flush Example
 *
 * Demonstrates the flush operation:
 * 1. Create a fake collection
 * 2. Create a fake document
 * 3. Check collection count
 * 4. Flush all data
 * 5. Re-check collection count (should be 0)
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hlquery/go-api"
)

func main() {
	client := hlquery.NewClient("http://localhost:9200")

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("FLUSH EXAMPLE")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Step 1: Create a fake collection
	fmt.Println("Step 1: Creating a fake collection...")
	collection_name := fmt.Sprintf("flush_test_collection_%d", time.Now().Unix())

	schema := map[string]interface{}{
		"fields": []map[string]interface{}{
			{"name": "title", "type": "string"},
			{"name": "content", "type": "string"},
			{"name": "value", "type": "int"},
		},
	}

	create_result, err := client.Collections().Create(collection_name, schema)
	if err != nil {
		log.Fatal(err)
	}
	if create_result.IsSuccess() {
		fmt.Printf("  ✓ Collection '%s' created successfully\n", collection_name)
	} else {
		fmt.Printf("  ✗ Failed to create collection: %d\n", create_result.StatusCode)
		errorJSON, _ := json.MarshalIndent(create_result.Body, "", "  ")
		fmt.Printf("  Error: %s\n", string(errorJSON))
		return
	}

	fmt.Println()

	// Step 2: Create a fake document
	fmt.Println("Step 2: Creating a fake document...")
	doc := map[string]interface{}{
		"id":      fmt.Sprintf("flush_test_doc_%d", time.Now().Unix()),
		"title":   "Flush Test Document",
		"content": "This is a test document for flush example",
		"value":   42,
	}

	add_result, err := client.Documents().Add(collection_name, doc)
	if err != nil {
		log.Fatal(err)
	}
	if add_result.IsSuccess() {
		fmt.Printf("  ✓ Document '%s' added successfully\n", doc["id"])
	} else {
		fmt.Printf("  ✗ Failed to add document: %d\n", add_result.StatusCode)
		errorJSON, _ := json.MarshalIndent(add_result.Body, "", "  ")
		fmt.Printf("  Error: %s\n", string(errorJSON))
	}

	fmt.Println()

	// Step 3: Check collection count before flush
	fmt.Println("Step 3: Checking collection count before flush...")
	collections_before, err := client.ListCollections(0, 1000)
	if err != nil {
		log.Fatal(err)
	}
	count_before := 0
	if collections_before.IsSuccess() {
		body := collections_before.Body
		if collections_list, ok := body["collections"].([]interface{}); ok {
			count_before = len(collections_list)
			fmt.Printf("  Collections before flush: %d\n", count_before)
			if count_before == 0 {
				fmt.Println("  Warning: No collections found before flush")
			}
		}
	} else {
		fmt.Printf("  Failed to list collections: %d\n", collections_before.StatusCode)
	}

	fmt.Println()

	// Step 4: Flush all data
	fmt.Println("Step 4: Flushing all data...")
	flush_result, err := client.Flush()
	if err != nil {
		log.Fatal(err)
	}
	if flush_result.IsSuccess() {
		body := flush_result.Body
		collections_deleted := 0
		if cd, ok := body["collections_deleted"].(float64); ok {
			collections_deleted = int(cd)
		}
		fmt.Println("  Flush completed successfully")
		fmt.Printf("  Collections deleted: %d\n", collections_deleted)
		message := "N/A"
		if msg, ok := body["message"].(string); ok {
			message = msg
		}
		fmt.Printf("  Message: %s\n", message)
	} else {
		fmt.Printf("  Flush failed: %d\n", flush_result.StatusCode)
		errorJSON, _ := json.MarshalIndent(flush_result.Body, "", "  ")
		fmt.Printf("  Error: %s\n", string(errorJSON))
		return
	}

	fmt.Println()

	// Step 5: Re-check collection count after flush
	fmt.Println("Step 5: Checking collection count after flush...")
	collections_after, err := client.ListCollections(0, 1000)
	if err != nil {
		log.Fatal(err)
	}
	count_after := -1
	if collections_after.IsSuccess() {
		body := collections_after.Body
		if collections_list, ok := body["collections"].([]interface{}); ok {
			count_after = len(collections_list)
			fmt.Printf("  Collections after flush: %d\n", count_after)

			if count_after == 0 {
				fmt.Println("  SUCCESS: All collections have been flushed")
			} else {
				fmt.Printf("  Warning: Expected 0 collections, but found %d\n", count_after)
			}
		}
	} else {
		fmt.Printf("  Failed to list collections: %d\n", collections_after.StatusCode)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("FLUSH EXAMPLE COMPLETED")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("Summary:")
	fmt.Printf("  Collections before flush: %d\n", count_before)
	fmt.Printf("  Collections after flush: %d\n", count_after)
	fmt.Println()
}
