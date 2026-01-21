package truenas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AppGetter defines the contract for fetching TrueNAS apps.
type AppGetter interface {
	GetApps(ctx context.Context) ([]App, error)
}

// Client is a client for the TrueNAS API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new, configured TrueNAS client.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetApps fetches all applications from the TrueNAS API.
func (c *Client) GetApps(ctx context.Context) ([]App, error) {
	url := fmt.Sprintf("%s/api/v2.0/app", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create truenas apps request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute truenas apps request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("bad status code for truenas apps: %d, body: %s", res.StatusCode, string(bodyBytes))
	}

	var m []App
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode truenas apps response: %w", err)
	}
	return m, nil
}
