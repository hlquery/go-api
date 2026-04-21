/*
 * hlquery - Search beyond keywords.
 * Collections Examples
 *
 * Demonstrates collection management operations
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hlquery/go-api"
)

func main() {
	client := hlquery.NewClient("http://localhost:9200")

	// List collections
	collections, err := client.Collections().List(0, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Collections:")
	collectionsJSON, _ := json.MarshalIndent(collections.Body, "", "    ")
	fmt.Println(string(collectionsJSON))

	// Get collection details
	if collections.IsSuccess() {
		body := collections.Body
		if collectionsList, ok := body["collections"].([]interface{}); ok && len(collectionsList) > 0 {
			var collectionName string
			if first, ok := collectionsList[0].(map[string]interface{}); ok {
				collectionName = first["name"].(string)
			} else if first, ok := collectionsList[0].(string); ok {
				collectionName = first
			}

			if collectionName != "" {
				// Get collection
				collection, err := client.Collections().Get(collectionName)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("\nCollection details:")
				collectionJSON, _ := json.MarshalIndent(collection.Body, "", "    ")
				fmt.Println(string(collectionJSON))

				// Get formatted fields
				fields, err := client.Collections().GetFields(collectionName)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("\nCollection fields:")
				fieldsJSON, _ := json.MarshalIndent(fields.Body, "", "    ")
				fmt.Println(string(fieldsJSON))
			}
		}
	}

	// Create collection
	schema := map[string]interface{}{
		"fields": []map[string]interface{}{
			{"name": "title", "type": "string"},
			{"name": "content", "type": "string"},
			{"name": "embedding", "type": "float[]"},
		},
	}
	createResult, err := client.Collections().Create("new_collection", schema)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nCreate result:")
	createJSON, _ := json.MarshalIndent(createResult.Body, "", "    ")
	fmt.Println(string(createJSON))

	// Delete collection
	deleteResult, err := client.Collections().Delete("collection_name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nDelete result:")
	deleteJSON, _ := json.MarshalIndent(deleteResult.Body, "", "    ")
	fmt.Println(string(deleteJSON))
}
