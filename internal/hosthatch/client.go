package hosthatch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type HostHatchClient interface {
	ListServers(ctx context.Context) ([]Server, error)
	ServerInfo(ctx context.Context, serverId int) (Server, error)
	ServerStatus(ctx context.Context, serverId int) (string, error)
	BootServer(ctx context.Context, serverId int) error
	RebootServer(ctx context.Context, serverId int) error
	ShutdownServer(ctx context.Context, serverId int) error
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://cloud.hosthatch.com/api",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (client *Client) newHttpRequest(ctx context.Context, path string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		path,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create devices request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+client.apiKey)

	return req, nil
}

func (client *Client) ListServers(ctx context.Context) ([]Server, error) {
	req, err := client.newHttpRequest(ctx, fmt.Sprintf("%s/v1/servers", client.baseURL))
	if err != nil {
		return nil, err
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	m := ApiResult[Servers]{}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode devices response: %w", err)
	}

	jsonData, _ := json.Marshal(m)
	_ = os.WriteFile(
		"/home/manfred/Development/StatusApp/deployments/hosthatch.json",
		jsonData,
		0o666,
	)

	return m.Data.Servers, nil
}

func (client *Client) ServerInfo(ctx context.Context, serverId int) (Server, error) {
	return Server{}, nil
}

func (client *Client) ServerStatus(
	ctx context.Context,
	serverId int,
) (string, error) {
	return "", nil
}
func (client *Client) BootServer(ctx context.Context, serverId int) error     { return nil }
func (client *Client) RebootServer(ctx context.Context, serverId int) error   { return nil }
func (client *Client) ShutdownServer(ctx context.Context, serverId int) error { return nil }
