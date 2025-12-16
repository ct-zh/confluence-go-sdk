package confluence

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// faultyReader simulates an error during reading
type faultyReader struct {
	data string
	read int
}

func (r *faultyReader) Read(p []byte) (n int, err error) {
	if r.read >= len(r.data) {
		return 0, errors.New("simulated read error")
	}
	n = copy(p, r.data[r.read:])
	r.read += n
	return n, nil
}

func TestContentService_UploadAttachment_EdgeCases(t *testing.T) {
	// Setup a mock server that expects multipart data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Consume body to trigger read
		_, err := io.Copy(io.Discard, r.Body)
		if err != nil {
			// This is expected when client disconnects or pipe fails
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[]}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)

	t.Run("ReadError", func(t *testing.T) {
		// Test that an error in the reader is propagated
		reader := &faultyReader{data: "some data"}
		_, _, err := client.Content.UploadAttachment(context.Background(), "123", "faulty.txt", reader, "")

		if err == nil {
			t.Fatal("Expected error due to faulty reader, got nil")
		}
		if !strings.Contains(err.Error(), "simulated read error") {
			t.Errorf("Expected error containing 'simulated read error', got: %v", err)
		}
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		// Test context cancellation during upload
		// We use a pipe to block the read until we cancel context
		pr, pw := io.Pipe()
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
			pw.Close()
		}()

		_, _, err := client.Content.UploadAttachment(ctx, "123", "cancel.txt", pr, "")

		if err == nil {
			t.Fatal("Expected error due to context cancellation, got nil")
		}
		// The error might be "context canceled" or wrapped
		if !strings.Contains(err.Error(), "context canceled") && !strings.Contains(err.Error(), "pipe closed") {
			t.Logf("Got error: %v (Acceptable if related to cancel/close)", err)
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		// Test uploading an empty file
		reader := strings.NewReader("")
		_, resp, err := client.Content.UploadAttachment(context.Background(), "123", "empty.txt", reader, "")
		if err != nil {
			t.Fatalf("Unexpected error for empty file: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
