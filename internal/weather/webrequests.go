package weather

import (
	"StatusApp/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func WeatherRequest() tea.Msg {
	apiKey := os.Getenv("WEATHERAPI_API_KEY")
	place := os.Getenv("WEATHERAPI_LOCATION")
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, place)
	c := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	res, err := c.Do(req)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	defer func(closer io.ReadCloser) {
		_ = closer.Close()
	}(res.Body)

	m := models.Weather{}
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	return models.WeatherMsg{Weather: m}
}
