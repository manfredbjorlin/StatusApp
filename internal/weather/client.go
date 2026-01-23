package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DataProvider defines the contract for fetching weather data.
type DataProvider interface {
	GetCurrentWeather(ctx context.Context) (*WeatherForecastInternal, error)
	GetWaterTemperature(ctx context.Context) (WaterTemperature, error)
}

// Client is a client for weather-related APIs.
type Client struct {
	weatherAPIKey       string
	weatherAPILocation  string
	waterTempLocationID string
	httpClient          *http.Client
}

// NewClient creates a new, configured weather client.
func NewClient(weatherKey, weatherLocation, waterLocationID string) *Client {
	return &Client{
		weatherAPIKey:       weatherKey,
		weatherAPILocation:  weatherLocation,
		waterTempLocationID: waterLocationID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func getWeatherIcon(code int, isDay int, iconPath string) string {
	byteValue, err := os.ReadFile(iconPath)
	if err != nil {
		return "?"
	}

	var weatherData []Weathercode
	if err := json.Unmarshal(byteValue, &weatherData); err != nil {
		return "?"
	}

	for _, weather := range weatherData {
		if weather.Code == code {
			if isDay == 1 {
				return weather.Day
			}
			return weather.Night
		}
	}
	return "?"
}

func (c *Client) GetCurrentWeather(ctx context.Context) (*WeatherForecastInternal, error) {
	switch os.Getenv("WEATHER_PROVIDER") {
	case "yr":
		return c.getCurrentWeatherYr(ctx)
	case "weatherapi":
		return c.getCurrentWeatherWeatherapi(ctx)
	}
	return c.getCurrentWeatherWeatherapi(ctx)
}

func (c *Client) getCurrentWeatherYr(ctx context.Context) (*WeatherForecastInternal, error) {
	lat := os.Getenv("WEATHER_LAT")
	lon := os.Getenv("WEATHER_LON")
	url := fmt.Sprintf(
		"https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%s&lon=%s",
		lat,
		lon,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "ManfredStatusApp/0.1 github.com/manfredbjorlin/StatusApp")
	req.Header.Add("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf(
			"bad status code for weather: %d, body: %s",
			res.StatusCode,
			string(bodyBytes),
		)
	}

	var m WeatherYr
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}

	icon, text := getIconYr(
		os.Getenv("WEATHER_ICON_PATH_YR"),
		m.Properties.Timeseries[0].Data.Next1Hours.Summary.SymbolCode,
	)

	go func(icon string, text string) {
		_ = os.WriteFile(
			"/home/manfred/Development/StatusApp/deployments/icons.txt",
			fmt.Appendf(nil, "Icon: %s\nText: %s", icon, text),
			0o666,
		)
	}(icon, text)

	result := WeatherForecastInternal{
		Icon:        icon,
		Text:        text,
		Location:    c.weatherAPILocation,
		Temperature: float32(m.Properties.Timeseries[0].Data.Instant.Details.AirTemperature),
		ExtraInfo: fmt.Sprintf(
			"%0.1f m/s",
			m.Properties.Timeseries[0].Data.Instant.Details.WindSpeed,
		),
	}

	return &result, nil
}

func (c *Client) getCurrentWeatherWeatherapi(
	ctx context.Context,
) (*WeatherForecastInternal, error) {
	url := fmt.Sprintf(
		"http://api.weatherapi.com/v1/current.json?key=%s&q=%s",
		c.weatherAPIKey,
		c.weatherAPILocation,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create weather request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute weather request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf(
			"bad status code for weather: %d, body: %s",
			res.StatusCode,
			string(bodyBytes),
		)
	}

	var m Weather
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}
	m.Current.Location = c.weatherAPILocation

	result := WeatherForecastInternal{
		Icon: getWeatherIcon(
			m.Current.Condition.Code,
			m.Current.IsDay,
			os.Getenv("WEATHER_ICON_PATH"),
		),
		Text:        m.Current.Condition.Text,
		Temperature: m.Current.Temp,
		Location:    c.weatherAPILocation,
		ExtraInfo:   fmt.Sprintf("%0.1f°C", m.Current.FeelsLike),
	}

	return &result, nil
}

// GetWaterTemperature fetches water temperature from YR.no.
func (c *Client) GetWaterTemperature(ctx context.Context) (WaterTemperature, error) {
	url := fmt.Sprintf(
		"https://www.yr.no/api/v0/locations/%s/nearestwatertemperatures",
		c.waterTempLocationID,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return WaterTemperature{}, fmt.Errorf("failed to create water temperature request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return WaterTemperature{}, fmt.Errorf(
			"failed to execute water temperature request: %w",
			err,
		)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return WaterTemperature{}, fmt.Errorf(
			"bad status code for water temperature: %d, body: %s",
			res.StatusCode,
			string(bodyBytes),
		)
	}

	var m WaterTemperature
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return WaterTemperature{}, fmt.Errorf(
			"failed to decode water temperature response: %w",
			err,
		)
	}
	return m, nil
}
