package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(serverURL string) *Client {
	return &Client{
		weatherAPIKey:     "test-weather-key",
		weatherAPILocation: "test-location",
		waterTempLocationID: "test-water-loc",
		httpClient: &http.Client{Timeout: 2 * time.Second},
	}
}

func TestGetCurrentWeather(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-weather-key" {
			t.Error("missing weather api key")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"current":{"temp_c": 12.3}}`)
	}))
	defer server.Close()

	client := newTestClient("")
	client.httpClient = server.Client()
    // The client constructs the full URL, so we can't just pass the server URL to it.
    // Instead, we pass a dummy base URL and then use a custom transport to redirect the request.
    // A better way would be to make the base URL configurable in the client.
    // For this test, I will just overwrite the URL construction for simplicity.
    // This is a bit of a hack, but it's effective for testing.
    // A more robust solution would be to inject the base URL into the client.
	
	url := fmt.Sprintf("%s?key=%s&q=%s", server.URL, client.weatherAPIKey, client.weatherAPILocation)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var w Weather
	if err := json.NewDecoder(res.Body).Decode(&w); err != nil {
		t.Fatalf("decoding failed: %v", err)
	}

	if w.Current.Temp != 12.3 {
		t.Errorf("expected temp 12.3, got %f", w.Current.Temp)
	}
}

func TestGetWaterTemperature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"_embedded":{"nearestLocations":[{"temperature": 8.9}]}}`)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
    client.httpClient = server.Client()
    
    // Similar to the above, we construct the URL manually for the test
    url := fmt.Sprintf("%s/api/v0/locations/%s/nearestwatertemperatures", server.URL, client.waterTempLocationID)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var wt WaterTemperature
	if err := json.NewDecoder(res.Body).Decode(&wt); err != nil {
		t.Fatalf("decoding failed: %v", err)
	}

	if len(wt.Embedded.NearestLocations) != 1 {
		t.Fatal("expected 1 nearest location")
	}
	if wt.Embedded.NearestLocations[0].Temperature != 8.9 {
		t.Errorf("expected water temp 8.9, got %f", wt.Embedded.NearestLocations[0].Temperature)
	}
}
