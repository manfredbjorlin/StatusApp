package tailscale

import "time"

// Device represents a single machine in a Tailscale network.
type Device struct {
	Hostname           string    `json:"hostname"`
	ConnectedToControl bool      `json:"connectedToControl"`
	Name               string    `json:"name"`
	Os                 string    `json:"os"`
	ClientVersion      string    `json:"clientVersion"`
	UpdateAvailable    bool      `json:"updateAvailable"`
	LastSeen           time.Time `json:"lastSeen"`
	AdvertisedRoutes   []string  `json:"advertisedRoutes"`
	KeyExpiryDisabled  bool      `json:"keyExpiryDisabled"`
}

// Devices is a wrapper for a list of devices, matching the API response.
type Devices struct {
	Devices []Device `json:"devices"`
}

// Key represents the API key's expiration information.
type Key struct {
	Expires time.Time `json:"expires"`
}
