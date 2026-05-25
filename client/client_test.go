package hlquery

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

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

func TestNewClientWithErrorValidatesBaseURL(t *testing.T) {
	if _, err := NewClientWithError("://bad"); err == nil {
		t.Fatal("expected invalid base URL error")
	}

	client := NewClient("://bad")
	if _, err := client.Health(); err == nil {
		t.Fatal("expected stored invalid base URL error")
	}
}

func TestRequestHeadersAndAuth(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, _ map[string]interface{}) {
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("Accept = %q, want application/json", got)
		}
		if got := r.Header.Get("Content-Type"); got != "" {
			t.Fatalf("Content-Type on body-less request = %q, want empty", got)
		}
		if got := r.Header.Get("X-API-Key"); got != "secret" {
			t.Fatalf("X-API-Key = %q, want secret", got)
		}
	})
	client.SetAuthToken("secret", "api-key")

	if _, err := client.Health(); err != nil {
		t.Fatal(err)
	}
}

func TestBodyRequestSetsContentType(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, _ map[string]interface{}) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}
	})

	if _, err := client.CreateCollection("products", map[string]interface{}{}); err != nil {
		t.Fatal(err)
	}
}

func TestCustomHTTPClientIsUsed(t *testing.T) {
	called := false
	client := NewClient("http://example.test", ClientOptions{
		HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			called = true
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
				Request:    r,
			}, nil
		})},
	})

	if _, err := client.Health(); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("custom HTTP client transport was not used")
	}
}

func TestQueryEncodingAndPathEscaping(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, _ map[string]interface{}) {
		if r.URL.EscapedPath() != "/collections/a+b%2Fc/documents/search" {
			t.Fatalf("path = %s, want escaped path preserving plus and slash", r.URL.EscapedPath())
		}
		if got := r.URL.Query().Get("q"); got != "red shoes" {
			t.Fatalf("q = %q, want red shoes", got)
		}
		if got := r.URL.Query()["tag"]; len(got) != 2 || got[0] != "a" || got[1] != "b" {
			t.Fatalf("tag query = %#v, want [a b]", got)
		}
	})

	if _, err := client.Search().PerformGET("a+b/c", map[string]interface{}{"q": "red shoes", "tag": []string{"a", "b"}}); err != nil {
		t.Fatal(err)
	}
}

func TestAPIErrorReturnedForNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	}))
	t.Cleanup(server.Close)

	resp, err := NewClient(server.URL).Health()
	if err == nil {
		t.Fatal("expected API error")
	}
	if resp == nil || resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("response = %#v, want status 400", resp)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if !strings.Contains(apiErr.Error(), "bad request") {
		t.Fatalf("API error = %q, want message", apiErr.Error())
	}
}

func TestResponseBodyLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, strings.Repeat("x", 8))
	}))
	t.Cleanup(server.Close)

	client := NewClient(server.URL, ClientOptions{MaxResponseBytes: 4})
	if _, err := client.Health(); err == nil {
		t.Fatal("expected response size error")
	}
}

func TestExecuteRequestWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := NewClient("http://127.0.0.1")
	if _, err := client.ExecuteRequestWithContext(ctx, "GET", "/health", nil); err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestNilParamsDoNotPanic(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, _ map[string]interface{}) {
		if r.URL.Query().Get("q") != "hello" {
			t.Fatalf("q = %q, want hello", r.URL.Query().Get("q"))
		}
	})

	if _, err := client.SAM().Search("", "hello", nil); err != nil {
		t.Fatal(err)
	}
}

func TestEmptyRouteParametersReturnError(t *testing.T) {
	client := NewClient("http://127.0.0.1")
	if _, err := client.Documents().Get("", "doc"); err == nil {
		t.Fatal("expected empty collection error")
	}
	if _, err := client.Documents().Get("products", ""); err == nil {
		t.Fatal("expected empty document id error")
	}
}

func TestTypedHelpers(t *testing.T) {
	client := testClient(t, func(t *testing.T, r *http.Request, body map[string]interface{}) {
		switch r.URL.Path {
		case "/collections":
			if body["name"] != "products" {
				t.Fatalf("name = %v, want products", body["name"])
			}
		case "/collections/products/documents/search":
			if body["q"] != "laptop" {
				t.Fatalf("q = %v, want laptop", body["q"])
			}
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	})

	if _, err := client.Collections().CreateSchema(CollectionSchema{Name: "products"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Search().PerformTyped("products", SearchParams{Query: "laptop"}); err != nil {
		t.Fatal(err)
	}
}
