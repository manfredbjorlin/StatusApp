package tailscale

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(serverURL, apiKey, tailnet, keyID string) *Client {
	return &Client{
		baseURL:    serverURL,
		apiKey:     apiKey,
		tailnet:    tailnet,
		keyID:      keyID,
		httpClient: &http.Client{Timeout: 2 * time.Second},
	}
}

func TestGetMachines(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/tailnet/test-tailnet/devices" {
				t.Errorf("expected path '/tailnet/test-tailnet/devices', got '%s'", r.URL.Path)
			}
			user, _, ok := r.BasicAuth()
			if !ok || user != "test-key" {
				t.Errorf("expected basic auth username 'test-key', got '%s'", user)
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"devices": [{"name": "test-device"}]}`)
		}))
		defer server.Close()

		client := newTestClient(server.URL, "test-key", "test-tailnet", "")
		devices, err := client.GetMachines(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(devices.Devices) != 1 {
			t.Fatalf("expected 1 device, got %d", len(devices.Devices))
		}
		if devices.Devices[0].Name != "test-device" {
			t.Errorf("expected device name 'test-device', got '%s'", devices.Devices[0].Name)
		}
	})

	t.Run("API Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, `{"message": "internal server error"}`)
		}))
		defer server.Close()

		client := newTestClient(server.URL, "test-key", "test-tailnet", "")
		_, err := client.GetMachines(context.Background())
		if err == nil {
			t.Fatal("expected an error but got none")
		}
	})
}

func TestGetKeyExpiry(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expiryTime := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339Nano)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/tailnet/test-tailnet/keys/test-key-id" {
				t.Errorf("unexpected path for key expiry: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"expires": "%s"}`, expiryTime)
		}))
		defer server.Close()

		client := newTestClient(server.URL, "test-key", "test-tailnet", "test-key-id")
		expires, err := client.GetKeyExpiry(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		parsedTime, _ := time.Parse(time.RFC3339Nano, expiryTime)
		if !expires.Equal(parsedTime) {
			t.Errorf("expected expiry time %v, got %v", parsedTime, expires)
		}
	})
}
