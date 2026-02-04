package syncthing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type SyncthingClient interface {
	ListConnections(ctx context.Context) ([]Connection, error)
	ListDevices(ctx context.Context) ([]Device, error)
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "http://syncthing.manfred.no",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (client *Client) newHttpRequest(
	ctx context.Context,
	verb string,
	path string,
	body any,
) (*http.Request, error) {
	var buffer bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buffer).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		strings.ToUpper(verb),
		client.baseURL+path,
		&buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create devices request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+client.apiKey)

	return req, nil
}

func (client *Client) ListConnections(ctx context.Context) ([]Connection, error) {
	req, err := client.newHttpRequest(ctx, "GET", "/rest/system/connections", nil)
	if err != nil {
		return nil, err
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	m := Connections{}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode connections response: %w", err)
	}

	devices, err := client.getDevices(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Connection, 0)

	for key, connection := range m.Connections {
		connection.Device = Device{
			DeviceID: key,
			Name:     devices[key],
		}
		result = append(result, connection)
	}

	return result, nil
}

func (client *Client) getDevices(ctx context.Context) (map[string]string, error) {
	devices, err := client.ListDevices(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, 0)
	for _, device := range devices {
		result[device.DeviceID] = device.Name
	}

	return result, nil
}

func (client *Client) ListDevices(ctx context.Context) ([]Device, error) {
	req, err := client.newHttpRequest(ctx, "GET", "/rest/config/devices", nil)
	if err != nil {
		return nil, err
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	m := []Device{}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode devices response: %w", err)
	}

	return m, nil
}
