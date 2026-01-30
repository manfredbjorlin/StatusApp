package tailscale

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MachineGetter defines the contract for fetching machines and key expiry.
type MachineGetter interface {
	GetMachines(ctx context.Context) (Devices, error)
	GetKeyExpiry(ctx context.Context) (time.Time, error)
}

// Client is a client for the Tailscale API.
type Client struct {
	baseURL    string
	apiKey     string
	tailnet    string
	keyID      string
	httpClient *http.Client
}

// NewClient creates a new, configured Tailscale client.
func NewClient(apiKey, tailnet, keyID string) *Client {
	return &Client{
		baseURL: "https://api.tailscale.com/api/v2",
		apiKey:  apiKey,
		tailnet: tailnet,
		keyID:   keyID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetMachines fetches all machines in the tailnet.
func (c *Client) GetMachines(ctx context.Context) (Devices, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/tailnet/%s/devices?fields=all", c.baseURL, c.tailnet),
		nil,
	)
	if err != nil {
		return Devices{}, fmt.Errorf("failed to create devices request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, "")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return Devices{}, fmt.Errorf("failed to execute devices request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return Devices{}, fmt.Errorf(
			"bad status code for devices: %d, body: %s",
			res.StatusCode,
			string(bodyBytes),
		)
	}

	m := Devices{}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return Devices{}, fmt.Errorf("failed to decode devices response: %w", err)
	}
	return m, nil
}

// GetKeyExpiry fetches the expiry date for the API key.
func (c *Client) GetKeyExpiry(ctx context.Context) (time.Time, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/tailnet/%s/keys/%s", c.baseURL, c.tailnet, c.keyID),
		nil,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to create key expiry request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, "")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to execute key expiry request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return time.Time{}, fmt.Errorf(
			"bad status code for key expiry: %d, body: %s",
			res.StatusCode,
			string(bodyBytes),
		)
	}

	k := Key{}
	if err := json.NewDecoder(res.Body).Decode(&k); err != nil {
		return time.Time{}, fmt.Errorf("failed to decode key expiry response: %w", err)
	}

	return k.Expires, nil
}
