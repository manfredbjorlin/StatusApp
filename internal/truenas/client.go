package truenas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

// AppGetter defines the contract for fetching TrueNAS apps.
type AppGetter interface {
	GetApps(ctx context.Context) ([]App, error)
	UpdateApp(apps []App, updateId int) error
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
		return nil, fmt.Errorf(
			"bad status code for truenas apps: %d, body: %s",
			res.StatusCode,
			string(bodyBytes),
		)
	}

	var m []App
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode truenas apps response: %w", err)
	}
	return assignUpdateId(m), nil
}

func assignUpdateId(apps []App) []App {
	var appList []string
	for _, app := range apps {
		if app.UpgradeAvailable {
			appList = append(appList, app.Name)
		}
	}
	slices.Sort(appList)
	var resultApps []App
	for _, app := range apps {
		if slices.Contains(appList, app.Name) {
			a := app
			a.UpdateId = slices.Index(appList, app.Name) + 1
			resultApps = append(resultApps, a)
		} else {
			resultApps = append(resultApps, app)
		}
	}
	return resultApps
}

func (c *Client) UpdateApp(apps []App, updateId int) error {
	for _, app := range apps {
		if app.UpgradeAvailable && app.UpdateId == updateId && app.State == "RUNNING" {
			return c.updateApp(context.Background(), app.Name)
		}
	}
	return fmt.Errorf("could not find app to upgrade")
}

func (c *Client) updateApp(ctx context.Context, appId string) error {
	url := fmt.Sprintf("%s/api/v2.0/app/upgrade", c.baseURL)
	jsonData := fmt.Sprintf("{ \"app_name\": \"%s\", \"options\": {} }", appId)
	reader := bytes.NewReader([]byte(jsonData))
	req, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v", res.StatusCode)
	}
	return nil
}
