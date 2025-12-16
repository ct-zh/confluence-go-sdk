package confluence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpaceService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/space/DS", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":123,"key":"DS","name":"Dev Space"}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	space, _, err := client.Space.Get(context.Background(), "DS", nil)

	if err != nil {
		t.Errorf("Space.Get returned error: %v", err)
	}
	if space.Name != "Dev Space" {
		t.Errorf("Expected Name 'Dev Space', got '%s'", space.Name)
	}
}

func TestSpaceService_Create(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/space", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var reqBody Space
		json.NewDecoder(r.Body).Decode(&reqBody)
		if reqBody.Key != "NEW" {
			t.Errorf("Expected Key 'NEW', got '%s'", reqBody.Key)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":1001,"key":"NEW","name":"New Space"}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	newSpace := &Space{Key: "NEW", Name: "New Space"}
	space, _, err := client.Space.Create(context.Background(), newSpace, nil)

	if err != nil {
		t.Errorf("Space.Create returned error: %v", err)
	}
	if space.ID != 1001 {
		t.Errorf("Expected ID 1001, got %d", space.ID)
	}
}

func TestSpaceService_Update(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/space/DS", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		var reqBody Space
		json.NewDecoder(r.Body).Decode(&reqBody)
		if reqBody.Name != "Updated Space" {
			t.Errorf("Expected Name 'Updated Space', got '%s'", reqBody.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":123,"key":"DS","name":"Updated Space"}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	updateSpace := &Space{Name: "Updated Space"}
	space, _, err := client.Space.Update(context.Background(), "DS", updateSpace)

	if err != nil {
		t.Errorf("Space.Update returned error: %v", err)
	}
	if space.Name != "Updated Space" {
		t.Errorf("Expected Name 'Updated Space', got '%s'", space.Name)
	}
}

func TestSpaceService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/space/DS", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.Space.Delete(context.Background(), "DS")

	if err != nil {
		t.Errorf("Space.Delete returned error: %v", err)
	}
}
