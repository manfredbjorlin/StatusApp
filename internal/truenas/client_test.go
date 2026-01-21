package truenas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetApps(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-key" {
			t.Errorf("expected auth header 'Bearer test-key', got '%s'", authHeader)
		}
		if r.URL.Path != "/api/v2.0/app" {
			t.Errorf("expected path '/api/v2.0/app', got '%s'", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"name": "test-app", "state": "ACTIVE"}]`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	apps, err := client.GetApps(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}
	if apps[0].Name != "test-app" {
		t.Errorf("expected app name 'test-app', got '%s'", apps[0].Name)
	}
}
