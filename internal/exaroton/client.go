package exaroton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ExarotonClient interface {
	ListServers(ctx context.Context) ([]Server, error)
	ServerInfo(ctx context.Context, serverId string) (Server, error)
	ServerRam(ctx context.Context) (int, error)
	StartServer(ctx context.Context, serverId string) error
	StopServer(ctx context.Context, serverId string) error
	RestartServer(ctx context.Context, serverId string) error
	RemainingCredits(ctx context.Context) (float64, error)
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://api.exaroton.com/v1",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
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

func (client *Client) ListServers(ctx context.Context) ([]Server, error) {
	req, err := client.newHttpRequest(ctx, "GET", "/servers/", nil)
	if err != nil {
		return nil, err
	}
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	m := ApiResponse[[]Server]{}
	if err = json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m.Data, nil
}

func (client *Client) ServerInfo(ctx context.Context, serverId string) (Server, error) {
	return Server{}, nil
}
func (client *Client) ServerRam(ctx context.Context) (int, error)               { return 0, nil }
func (client *Client) StartServer(ctx context.Context, serverId string) error   { return nil }
func (client *Client) StopServer(ctx context.Context, serverId string) error    { return nil }
func (client *Client) RestartServer(ctx context.Context, serverId string) error { return nil }

func (client *Client) RemainingCredits(ctx context.Context) (float64, error) {
	req, err := client.newHttpRequest(ctx, "GET", "/account/", nil)
	if err != nil {
		return 0, err
	}
	res, err := client.httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	m := ApiResponse[Account]{}
	if err = json.NewDecoder(res.Body).Decode(&m); err != nil {
		return 0, err
	}

	return m.Data.Credits, nil
}
