/*
 * hlquery - Search beyond keywords.
 * Search Examples
 *
 * Demonstrates search operations
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

	collection := "collection"

	// Simple search
	params := map[string]interface{}{
		"q":        "search query",
		"query_by": "title,content",
		"limit":    10,
	}
	results, err := client.Search().Perform(collection, params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Search results:")
	resultsJSON, _ := json.MarshalIndent(results.Body, "", "    ")
	fmt.Println(string(resultsJSON))

	// Search with filters
	filterParams := map[string]interface{}{
		"q":         "laptop",
		"query_by":  "title,description",
		"filter_by": "price:>1000",
		"sort_by":   "price:asc",
		"limit":     10,
	}
	filterResults, err := client.Search().Perform(collection, filterParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nFiltered search results:")
	filterJSON, _ := json.MarshalIndent(filterResults.Body, "", "    ")
	fmt.Println(string(filterJSON))

	fieldParams := map[string]interface{}{
		"q":        "title:laptop",
		"query_by": "title,content",
		"limit":    10,
	}
	fieldResults, err := client.Search().Perform(collection, fieldParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nField search results:")
	fieldJSON, _ := json.MarshalIndent(fieldResults.Body, "", "    ")
	fmt.Println(string(fieldJSON))

	orParams := map[string]interface{}{
		"q":        "title:laptop OR title:notebook",
		"query_by": "title,content",
		"limit":    10,
	}
	orResults, err := client.Search().Perform(collection, orParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nBoolean OR search results:")
	orJSON, _ := json.MarshalIndent(orResults.Body, "", "    ")
	fmt.Println(string(orJSON))

	notParams := map[string]interface{}{
		"q":        "title:laptop NOT title:refurbished",
		"query_by": "title,content",
		"limit":    10,
	}
	notResults, err := client.Search().Perform(collection, notParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nBoolean NOT search results:")
	notJSON, _ := json.MarshalIndent(notResults.Body, "", "    ")
	fmt.Println(string(notJSON))

	phraseParams := map[string]interface{}{
		"q":        "\"wireless keyboard\"",
		"query_by": "title",
		"limit":    10,
	}
	phraseResults, err := client.Search().Perform(collection, phraseParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nPhrase search results:")
	phraseJSON, _ := json.MarshalIndent(phraseResults.Body, "", "    ")
	fmt.Println(string(phraseJSON))

	wildcardParams := map[string]interface{}{
		"q":        "laptop*",
		"query_by": "title,content",
		"limit":    10,
	}
	wildcardResults, err := client.Search().Perform(collection, wildcardParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nWildcard search results:")
	wildcardJSON, _ := json.MarshalIndent(wildcardResults.Body, "", "    ")
	fmt.Println(string(wildcardJSON))

	// Multi-search
	searches := []map[string]interface{}{
		{"collection": "products", "q": "laptop", "query_by": "title"},
		{"collection": "articles", "q": "laptop", "query_by": "content"},
	}
	multiResults, err := client.Search().Multi(searches)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nMulti-search results:")
	multiJSON, _ := json.MarshalIndent(multiResults.Body, "", "    ")
	fmt.Println(string(multiJSON))

	// Vector search
	vectorParams := map[string]interface{}{
		"vector": []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		"limit":  10,
	}
	vectorResults, err := client.Search().Vector(collection, vectorParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nVector search results:")
	vectorJSON, _ := json.MarshalIndent(vectorResults.Body, "", "    ")
	fmt.Println(string(vectorJSON))
}
