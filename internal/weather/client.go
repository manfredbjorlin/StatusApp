package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DataProvider defines the contract for fetching weather data.
type DataProvider interface {
	GetCurrentWeather(ctx context.Context) (Weather, error)
	GetWaterTemperature(ctx context.Context) (WaterTemperature, error)
}

// Client is a client for weather-related APIs.
type Client struct {
	weatherAPIKey    string
	weatherAPILocation string
	waterTempLocationID string
	httpClient       *http.Client
}

// NewClient creates a new, configured weather client.
func NewClient(weatherKey, weatherLocation, waterLocationID string) *Client {
	return &Client{
		weatherAPIKey:    weatherKey,
		weatherAPILocation: weatherLocation,
		waterTempLocationID: waterLocationID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetCurrentWeather fetches current weather from WeatherAPI.com.
func (c *Client) GetCurrentWeather(ctx context.Context) (Weather, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", c.weatherAPIKey, c.weatherAPILocation)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Weather{}, fmt.Errorf("failed to create weather request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Weather{}, fmt.Errorf("failed to execute weather request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return Weather{}, fmt.Errorf("bad status code for weather: %d, body: %s", res.StatusCode, string(bodyBytes))
	}

	var m Weather
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return Weather{}, fmt.Errorf("failed to decode weather response: %w", err)
	}
	return m, nil
}

// GetWaterTemperature fetches water temperature from YR.no.
func (c *Client) GetWaterTemperature(ctx context.Context) (WaterTemperature, error) {
	url := fmt.Sprintf("https://www.yr.no/api/v0/locations/%s/nearestwatertemperatures", c.waterTempLocationID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return WaterTemperature{}, fmt.Errorf("failed to create water temperature request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return WaterTemperature{}, fmt.Errorf("failed to execute water temperature request: %w", err)
	}
	defer res.Body.Close()
	
	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return WaterTemperature{}, fmt.Errorf("bad status code for water temperature: %d, body: %s", res.StatusCode, string(bodyBytes))
	}

	var m WaterTemperature
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return WaterTemperature{}, fmt.Errorf("failed to decode water temperature response: %w", err)
	}
	return m, nil
}
