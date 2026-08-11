package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type MachineGetter interface {
	GetLatestVersion(ctx context.Context) (string, error)
	GetMachines(ctx context.Context) ([]Peer, error)
	GetKeyExpiry(ctx context.Context) (time.Time, error)
}

type Client struct {
	baseURL    string
	apiKey     string
	userID     string
	keyID      string
	httpClient *http.Client
}

func NewClient(apiKey, userID, keyID string) *Client {
	return &Client{
		baseURL: "https://netbird.manfred.no/api",
		apiKey:  apiKey,
		userID:  userID,
		keyID:   keyID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetLatestVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.github.com/repos/netbirdio/netbird/releases/latest",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create version request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute version request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf(
			"bad status code for devices: %d, body: %s",
			res.StatusCode,
			string(bodyBytes),
		)
	}

	m := Version{}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return "", fmt.Errorf("failed to decode Version response: %w", err)
	} else {
		return strings.TrimPrefix(m.Name, "v"), nil
	}
}

func (c *Client) GetMachines(ctx context.Context) ([]Peer, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/peers", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create peers request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute peers request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf(
			"bad status code for devices: %d, body: %s",
			res.StatusCode,
			string(bodyBytes),
		)
	}

	m := make([]Peer, 0)
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode peers response: %w", err)
	}
	return m, nil
}

func (c *Client) GetKeyExpiry(ctx context.Context) (time.Time, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/users/%s/tokens/%s", c.baseURL, c.userID, c.keyID),
		nil,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to create key expiry request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
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
