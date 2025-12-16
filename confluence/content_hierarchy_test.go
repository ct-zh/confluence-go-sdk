package confluence

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentService_GetChildPages(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/content/123/child/page", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Query().Get("start") != "0" {
			t.Errorf("Expected start=0, got %s", r.URL.Query().Get("start"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("Expected limit=10, got %s", r.URL.Query().Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"results":[{"id":"456","title":"Child Page"}],"size":1}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	opts := &GetChildPagesOptions{Start: 0, Limit: 10}
	results, _, err := client.Content.GetChildPages(context.Background(), "123", opts)

	if err != nil {
		t.Errorf("Content.GetChildPages returned error: %v", err)
	}
	if len(results.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results.Results))
	}
	if results.Results[0].Title != "Child Page" {
		t.Errorf("Expected Title 'Child Page', got '%s'", results.Results[0].Title)
	}
}

func TestContentService_GetAttachments(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/content/123/child/attachment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"results":[{"id":"att-1","title":"file.pdf","type":"attachment"}],"size":1}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	results, _, err := client.Content.GetAttachments(context.Background(), "123", nil)

	if err != nil {
		t.Errorf("Content.GetAttachments returned error: %v", err)
	}
	if len(results.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results.Results))
	}
	if results.Results[0].Title != "file.pdf" {
		t.Errorf("Expected Title 'file.pdf', got '%s'", results.Results[0].Title)
	}
}
