package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"StatusApp/internal/models"
)

func GetWaterTemperature() tea.Msg {
	location := os.Getenv("WATERTEMPERATURE_LOCATION_ID")
	url := fmt.Sprintf("https://www.yr.no/api/v0/locations/%s/nearestwatertemperatures", location)
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

	m := models.WaterTemperature{}
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	return models.WaterTempMsg{WaterTemp: m}
}
