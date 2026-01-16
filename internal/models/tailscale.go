package models

import "time"

type Device struct {
	Hostname           string    `json:"hostname"`
	ConnectedToControl bool      `json:"connectedToControl"`
	Name               string    `json:"name"`
	Os                 string    `json:"os"`
	ClientVersion      string    `json:"clientVersion"`
	UpdateAvailable    bool      `json:"updateAvailable"`
	LastSeen           time.Time `json:"lastSeen"`
}

type Devices struct {
	Devices []Device `json:"devices"`
}
