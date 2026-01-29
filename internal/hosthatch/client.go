package hosthatch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
		path,
		&buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create devices request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+client.apiKey)

	return req, nil
}

func (client *Client) ListServers(ctx context.Context) ([]Server, error) {
	req, err := client.newHttpRequest(ctx, "GET", fmt.Sprintf("%s/v1/servers", client.baseURL), nil)
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

func (client *Client) BootServer(ctx context.Context, serverId int) error {
	return client.bootServer(ctx, serverId, "boot")
}

func (client *Client) RebootServer(ctx context.Context, serverId int) error {
	return client.bootServer(ctx, serverId, "reboot")
}

func (client *Client) ShutdownServer(ctx context.Context, serverId int) error {
	return client.bootServer(ctx, serverId, "shutdown")
}

func (client *Client) bootServer(ctx context.Context, serverId int, action string) error {
	req, err := client.newHttpRequest(
		ctx,
		"POST",
		fmt.Sprintf("%s/v1/servers/%v/%s", client.baseURL, serverId, action),
		nil,
	)
	if err != nil {
		return err
	}
	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	m := ApiResult[string]{}

	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return err
	} else if m.Result != "success" {
		return fmt.Errorf("error rebooting: %s", m.Result)
	}
	return nil
}
