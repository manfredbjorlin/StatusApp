package tailscale

import (
	"strings"
	"testing"
	"time"

	"StatusApp/internal/truenas"
)

// mockModel implements the modelWithData interface for testing the view.
type mockModel struct {
	apps    []truenas.App
	devices Devices
	expiry  time.Time
}

func (m mockModel) GetTruenasApps() []truenas.App {
	return m.apps
}

func (m mockModel) GetTailscaleDevices() Devices {
	return m.devices
}

func (m mockModel) GetKeyExpiry() time.Time {
	return m.expiry
}

func TestView(t *testing.T) {
	t.Run("Renders devices and expiry", func(t *testing.T) {
		model := mockModel{
			apps: []truenas.App{
				{Name: "app1", State: "ACTIVE"},
				{Name: "app2", State: "STOPPED"},
			},
			devices: Devices{
				Devices: []Device{
					{Name: "linux-device.example.com", Os: "linux", ConnectedToControl: true, ClientVersion: "1.2.3-v1"},
					{Name: "win-device.example.com", Os: "windows", ConnectedToControl: false, LastSeen: time.Now().Add(-2 * time.Hour)},
				},
			},
			expiry: time.Now().Add(10 * 24 * time.Hour),
		}

		output := View(model)

		if !strings.Contains(output, "Dodo Apps:") {
			t.Error("output does not contain Truenas app status")
		}
		if !strings.Contains(output, "Linux device") {
			t.Error("output does not contain linux device name")
		}
		if !strings.Contains(output, "Win device") {
			t.Error("output does not contain windows device name")
		}
		if !strings.Contains(output, "2 H") { // Check for offline time
			t.Error("output does not contain offline time for windows device")
		}
		if !strings.Contains(output, "Tailscale key expiry:") {
			t.Error("output does not contain key expiry information")
		}
	})

	t.Run("Handles model mismatch", func(t *testing.T) {
		type wrongModel struct{}
		output := View(wrongModel{})
		if !strings.Contains(output, "Error: Could not render Tailscale view") {
			t.Error("view did not return an error for a mismatched model")
		}
	})
}
