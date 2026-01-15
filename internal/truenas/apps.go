package truenas

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

func GetApps() tea.Msg {
	apiKey := os.Getenv("TRUENAS_API_KEY")
	baseUrl := os.Getenv("TRUENAS_BASE_URL")
	c := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2.0/app", baseUrl), nil)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	res, err := c.Do(req)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	defer func(closer io.ReadCloser) {
		_ = closer.Close()
	}(res.Body)

	m := make([]models.TruenasApp, 0)
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return models.ErrMsg{Err: err}
	}

	return models.TruenasMsg{Apps: m}
}

func GetAppStatus(apps []models.TruenasApp) (int, int) {
	upToDate := 0
	toUpgrade := 0

	for _, app := range apps {
		if app.UpgradeAvailable {
			toUpgrade++
		} else {
			upToDate++
		}
	}

	return upToDate, toUpgrade
}
