package hlquery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(t *testing.T, check func(*testing.T, *http.Request, map[string]interface{})) *Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&body)
		}
		check(t, r, body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(server.Close)

	return NewClient(server.URL)
}

func TestSearchPerformUsesDocumentSearchRoute(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, _ map[string]interface{}) {
		if r.Method != "POST" {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/collections/products/documents/search" {
			t.Fatalf("path = %s, want /collections/products/documents/search", r.URL.Path)
		}
	})

	if _, err := client.Search().Perform("products", map[string]interface{}{"q": "laptop"}); err != nil {
		t.Fatal(err)
	}
}

func TestCollectionUpdateUsesServerRoute(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, _ map[string]interface{}) {
		if r.Method != "POST" {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/collections/products/update" {
			t.Fatalf("path = %s, want /collections/products/update", r.URL.Path)
		}
	})

	if _, err := client.Collections().Update("products", map[string]interface{}{"fields": []string{"title"}}); err != nil {
		t.Fatal(err)
	}
}

func TestCreateCollectionPreservesSchema(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, body map[string]interface{}) {
		if r.URL.Path != "/collections" {
			t.Fatalf("path = %s, want /collections", r.URL.Path)
		}
		if body["name"] != "products" {
			t.Fatalf("name = %v, want products", body["name"])
		}
		if _, ok := body["searchable_fields"]; !ok {
			t.Fatal("searchable_fields was dropped")
		}
	})

	_, err := client.Collections().Create("products", map[string]interface{}{
		"fields":            []map[string]interface{}{{"name": "title", "type": "string"}},
		"searchable_fields": []string{"title"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestResourceClientRoutes(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, _ map[string]interface{}) {
		if r.Method != "GET" {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/collections/products/synonyms" {
			t.Fatalf("path = %s, want /collections/products/synonyms", r.URL.Path)
		}
	})

	if _, err := client.Synonyms().List("products"); err != nil {
		t.Fatal(err)
	}
}
