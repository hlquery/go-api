/*
 * hlquery - Search beyond keywords.
 * Documents Examples
 *
 * Demonstrates document CRUD operations
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"

	hlquery "github.com/hlquery/go-api/client"
)

func main() {
	client := hlquery.NewClient("http://localhost:9200")

	collection := "collection"

	// List documents
	params := map[string]interface{}{
		"offset": 0,
		"limit":  10,
	}
	docs, err := client.Documents().List(collection, params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Documents:")
	docsJSON, _ := json.MarshalIndent(docs.Body, "", "    ")
	fmt.Println(string(docsJSON))

	// Get document
	doc, err := client.Documents().Get(collection, "doc_id")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nDocument:")
	docJSON, _ := json.MarshalIndent(doc.Body, "", "    ")
	fmt.Println(string(docJSON))

	// Add document
	newDoc := map[string]interface{}{
		"id":      "doc_1",
		"title":   "New Document",
		"content": "Document content",
	}
	addResult, err := client.Documents().Add(collection, newDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nAdd result:")
	addJSON, _ := json.MarshalIndent(addResult.Body, "", "    ")
	fmt.Println(string(addJSON))

	// Update document
	updatedDoc := map[string]interface{}{
		"title":   "Updated Document",
		"content": "Updated content",
	}
	updateResult, err := client.Documents().Update(collection, "doc_id", updatedDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nUpdate result:")
	updateJSON, _ := json.MarshalIndent(updateResult.Body, "", "    ")
	fmt.Println(string(updateJSON))

	// Delete document
	deleteResult, err := client.Documents().Delete(collection, "doc_id")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nDelete result:")
	deleteJSON, _ := json.MarshalIndent(deleteResult.Body, "", "    ")
	fmt.Println(string(deleteJSON))

	// Bulk import
	bulkDocs := []map[string]interface{}{
		{"id": "doc1", "title": "Doc 1"},
		{"id": "doc2", "title": "Doc 2"},
		{"id": "doc3", "title": "Doc 3"},
	}
	importResult, err := client.Documents().Import(collection, bulkDocs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nImport result:")
	importJSON, _ := json.MarshalIndent(importResult.Body, "", "    ")
	fmt.Println(string(importJSON))
}
