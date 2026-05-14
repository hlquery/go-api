/*
 * hlquery - Search beyond keywords.
 * SQL Examples
 *
 * Demonstrates collection-bound and top-level SQL operations.
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
	sqlAPI := client.SQLAPI()

	collection := "products"

	selectResults, err := sqlAPI.Search(
		collection,
		"SELECT id, title, price FROM products WHERE price > 100 ORDER BY price DESC LIMIT 3;",
		map[string]interface{}{
			"highlight": "false",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Collection SQL results:")
	selectJSON, _ := json.MarshalIndent(selectResults.Body, "", "    ")
	fmt.Println(string(selectJSON))

	showCollections, err := sqlAPI.Query("SHOW COLLECTIONS;")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nTop-level SQL results:")
	showJSON, _ := json.MarshalIndent(showCollections.Body, "", "    ")
	fmt.Println(string(showJSON))
}
