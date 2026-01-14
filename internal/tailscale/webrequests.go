package tailscale

import (
	"StatusApp/configs"
	"StatusApp/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TailscaleRequest() tea.Msg {
	tailnet := os.Getenv("TAILSCALE_TAILNET_ID")
	c := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tailnet/%s/devices", configs.BaseUrl, tailnet), nil)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	req.SetBasicAuth(os.Getenv("TAILSCALE_API_KEY"), "")
	res, err := c.Do(req)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	defer func(closer io.ReadCloser) {
		_ = closer.Close()
	}(res.Body)

	m := models.Devices{}
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	return models.TailscaleMsg{Devices: m, Extra: time.Now().Format("15:04:05")}
}

func GetKeyExpiry() tea.Msg {
	tailnet := os.Getenv("TAILSCALE_TAILNET_ID")
	apiKey := os.Getenv("TAILSCALE_API_KEY")
	keyId := os.Getenv("TAILSCALE_API_KEY_ID")

	c := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/tailnet/%s/keys/%s", configs.BaseUrl, tailnet, keyId),
		nil,
	)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	req.SetBasicAuth(apiKey, "")
	res, err := c.Do(req)
	if err != nil {
		return models.ErrMsg{Err: err}
	}
	defer func(closer io.ReadCloser) {
		_ = closer.Close()
	}(res.Body)

	type Key struct {
		Expires time.Time `json:"expires"`
	}

	k := Key{}
	err = json.NewDecoder(res.Body).Decode(&k)
	if err != nil {
		return models.ErrMsg{Err: err}
	}

	return models.TimeMsg{Time: k.Expires}
}

