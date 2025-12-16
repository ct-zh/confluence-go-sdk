package confluence

import (
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentService_UploadAttachment(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/content/123/child/attachment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		// Verify Content-Type is multipart/form-data
		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(mediaType, "multipart/") {
			t.Errorf("Expected multipart content type, got %s", mediaType)
		}
		
		// Parse multipart body
		mr := multipart.NewReader(r.Body, params["boundary"])
		
		// Check file part
		p, err := mr.NextPart()
		if err != nil {
			t.Fatal(err)
		}
		if p.FormName() != "file" {
			t.Errorf("Expected form name 'file', got %s", p.FormName())
		}
		if p.FileName() != "test.txt" {
			t.Errorf("Expected filename 'test.txt', got %s", p.FileName())
		}
		
		// Check comment part
		p, err = mr.NextPart()
		if err != nil {
			t.Fatal(err)
		}
		if p.FormName() != "comment" {
			t.Errorf("Expected form name 'comment', got %s", p.FormName())
		}
		
		// Check X-Atlassian-Token header
		if r.Header.Get("X-Atlassian-Token") != "nocheck" {
			t.Errorf("Expected X-Atlassian-Token header to be 'nocheck', got %s", r.Header.Get("X-Atlassian-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"results":[{"id":"att-1","title":"test.txt","type":"attachment"}],"size":1}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	fileContent := strings.NewReader("hello world")
	results, _, err := client.Content.UploadAttachment(context.Background(), "123", "test.txt", fileContent, "test comment")

	if err != nil {
		t.Errorf("Content.UploadAttachment returned error: %v", err)
	}
	if len(results.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results.Results))
	}
	if results.Results[0].Title != "test.txt" {
		t.Errorf("Expected Title 'test.txt', got '%s'", results.Results[0].Title)
	}
}
