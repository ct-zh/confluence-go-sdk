package confluence

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSpaceService_EdgeCases 测试边缘情况
func TestSpaceService_EdgeCases(t *testing.T) {
	// Case 1: Special Characters in Space Key
	t.Run("SpecialCharsInKey", func(t *testing.T) {
		mux := http.NewServeMux()
		// 使用 Go 1.22+ 路由模式匹配
		mux.HandleFunc("GET /rest/api/space/{key}", func(w http.ResponseWriter, r *http.Request) {
			key := r.PathValue("key")
			// 期望 Server 接收到的是解压后的 key
			if key != "Space Key" {
				t.Errorf("Expected path param 'Space Key', got '%s'", key)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"key":"Space Key"}`)
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		client, _ := NewClient(server.URL)
		// 传入未转义的 Key，期望 Client 自动转义
		space, _, err := client.Space.Get(context.Background(), "Space Key", nil)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if space.Key != "Space Key" {
			t.Errorf("Expected key 'Space Key', got '%s'", space.Key)
		}
	})

	// Case 2: 404 Not Found
	t.Run("NotFound", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/rest/api/space/MISSING", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"statusCode":404,"message":"No space found"}`))
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		client, _ := NewClient(server.URL)
		_, _, err := client.Space.Get(context.Background(), "MISSING", nil)

		if err == nil {
			t.Error("Expected error for 404, got nil")
		}
	})

	// Case 3: Nil Options shouldn't panic
	t.Run("NilOptions", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/rest/api/space/TEST", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"key":"TEST"}`)
		})
		server := httptest.NewServer(mux)
		defer server.Close()

		client, _ := NewClient(server.URL)
		_, _, err := client.Space.Get(context.Background(), "TEST", nil)
		if err != nil {
			t.Errorf("Unexpected error with nil options: %v", err)
		}
	})
}
