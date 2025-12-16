package confluence

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchService_Search(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		cql := r.URL.Query().Get("cql")
		if cql != "type=page" {
			t.Errorf("Expected cql=type=page, got %s", cql)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `
		{
			"results": [
				{
					"id": "1",
					"title": "Result 1",
					"type": "page"
				},
				{
					"id": "2",
					"title": "Result 2",
					"type": "page"
				}
			],
			"start": 0,
			"limit": 25,
			"size": 2
		}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	opts := &SearchOptions{
		CQL: "type=page",
	}
	results, _, err := client.Search.Search(context.Background(), opts)

	if err != nil {
		t.Fatalf("Search.Search returned error: %v", err)
	}
	if len(results.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results.Results))
	}
	if results.Results[0].Title != "Result 1" {
		t.Errorf("Expected first result title 'Result 1', got '%s'", results.Results[0].Title)
	}
}
