package upcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UpCloudClient interface {
	ListServers(ctx context.Context) ([]Server, error)
	ServerInfo(ctx context.Context, uuid uuid.UUID) (Server, error)
	StartServer(ctx context.Context, uuid uuid.UUID) error
	StopServer(ctx context.Context, uuid uuid.UUID) error
	RestartServer(ctx context.Context, uuid uuid.UUID) error
	RemainingCredits(ctx context.Context) (string, float64, error)
	BillingSummary(ctx context.Context, year int, month int) (string, float64, error)
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://api.upcloud.com",
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
	req, err := client.newHttpRequest(ctx, "GET", fmt.Sprintf("%s/1.3/server", client.baseURL), nil)
	if err != nil {
		return nil, err
	}
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	m := Servers{}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode devices response: %w", err)
	}

	return m.Data.Server, nil
}

func (client *Client) ServerInfo(ctx context.Context, uuid uuid.UUID) (Server, error) {
	return Server{}, nil
}
func (client *Client) StartServer(ctx context.Context, uuid uuid.UUID) error   { return nil }
func (client *Client) StopServer(ctx context.Context, uuid uuid.UUID) error    { return nil }
func (client *Client) RestartServer(ctx context.Context, uuid uuid.UUID) error { return nil }
func (client *Client) RemainingCredits(ctx context.Context) (string, float64, error) {
	req, err := client.newHttpRequest(
		ctx,
		"GET",
		fmt.Sprintf("%s/1.3/account", client.baseURL),
		nil,
	)
	if err != nil {
		return "", 0.0, nil
	}
	res, err := client.httpClient.Do(req)
	if err != nil {
		return "", 0.0, err
	}

	m := Account{}
	if err = json.NewDecoder(res.Body).Decode(&m); err != nil {
		return "", 0.0, err
	}

	return "EUR", m.Data.Credits / 100, nil
}

func (client *Client) BillingSummary(
	ctx context.Context,
	year int,
	month int,
) (string, float64, error) {
	req, err := client.newHttpRequest(
		ctx,
		"GET",
		fmt.Sprintf("%s/1.3/account/billing/summary/%v-%02d", client.baseURL, year, month),
		nil,
	)
	if err != nil {
		return "", 0.0, nil
	}
	res, err := client.httpClient.Do(req)
	if err != nil {
		return "", 0.0, err
	}

	m := Billing{}
	if err = json.NewDecoder(res.Body).Decode(&m); err != nil {
		return "", 0.0, err
	}

	return m.Currency, m.TotalAmount, nil
}
