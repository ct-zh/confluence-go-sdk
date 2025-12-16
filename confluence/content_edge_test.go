package confluence

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestContentService_EdgeCases(t *testing.T) {
	t.Run("ContextCancelled", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // Wait longer than context
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client, _ := NewClient(server.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, _, err := client.Content.GetChildPages(ctx, "123", nil)
		if err == nil {
			t.Error("Expected context cancellation error, got nil")
		}
	})

	t.Run("ServerError_500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL)
		_, resp, err := client.Content.GetAttachments(context.Background(), "123", nil)

		if err == nil {
			t.Error("Expected error for 500 response, got nil")
		}
		if resp == nil || resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected 500 status code, got %v", resp)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"invalid": json`)) // Broken JSON
		}))
		defer server.Close()

		client, _ := NewClient(server.URL)
		_, _, err := client.Content.GetChildPages(context.Background(), "123", nil)
		if err == nil {
			t.Error("Expected JSON decode error, got nil")
		}
	})

	t.Run("URLEscaping_SpecialChars", func(t *testing.T) {
		// Test that special characters in contentID are properly escaped
		specialID := "Test Page/Title?&"

		mux := http.NewServeMux()
		// Go 1.22+ routing matches path
		mux.HandleFunc("GET /rest/api/content/{id}/child/page", func(w http.ResponseWriter, r *http.Request) {
			// Check if the URL sent by client was actually escaped
			// r.URL.Path contains the path.
			// Note: The httptest server might decode it before it reaches handler?
			// Let's check RequestURI which is raw.
			// url.PathEscape escapes '?' to '%3F', '/' to '%2F', but '&' is valid in path so it remains '&'
			expectedPart := "Test%20Page%2FTitle%3F&"
			if !strings.Contains(r.RequestURI, expectedPart) {
				t.Errorf("Expected encoded URI part '%s', got %s", expectedPart, r.RequestURI)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"results":[],"size":0}`))
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		client, _ := NewClient(server.URL)
		_, _, err := client.Content.GetChildPages(context.Background(), specialID, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
