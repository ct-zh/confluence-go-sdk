package confluence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/content/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"123","title":"Test Page"}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	content, _, err := client.Content.Get(context.Background(), "123", nil)

	if err != nil {
		t.Errorf("Content.Get returned error: %v", err)
	}
	if content.Title != "Test Page" {
		t.Errorf("Expected Title 'Test Page', got '%s'", content.Title)
	}
}

func TestContentService_Create(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/content", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var reqBody Content
		json.NewDecoder(r.Body).Decode(&reqBody)
		if reqBody.Title != "New Page" {
			t.Errorf("Expected Title 'New Page', got '%s'", reqBody.Title)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"1001","title":"New Page"}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	newContent := &Content{Title: "New Page", Type: "page"}
	content, _, err := client.Content.Create(context.Background(), newContent)

	if err != nil {
		t.Errorf("Content.Create returned error: %v", err)
	}
	if content.ID != "1001" {
		t.Errorf("Expected ID '1001', got '%s'", content.ID)
	}
}

func TestContentService_Update(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/content/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		var reqBody Content
		json.NewDecoder(r.Body).Decode(&reqBody)
		if reqBody.Version.Number != 2 {
			t.Errorf("Expected Version 2, got %d", reqBody.Version.Number)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"123","title":"Updated Page","version":{"number":2}}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	updateContent := &Content{
		Title:   "Updated Page",
		Version: &Version{Number: 2},
	}
	content, _, err := client.Content.Update(context.Background(), "123", updateContent)

	if err != nil {
		t.Errorf("Content.Update returned error: %v", err)
	}
	if content.Version.Number != 2 {
		t.Errorf("Expected Version 2, got %d", content.Version.Number)
	}
}

func TestContentService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/content/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.Content.Delete(context.Background(), "123")

	if err != nil {
		t.Errorf("Content.Delete returned error: %v", err)
	}
}
