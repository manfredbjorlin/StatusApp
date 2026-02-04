package syncthing

import "time"

type Connections struct {
	Connections map[string]Connection `json:"connections"`
}

type Connection struct {
	Connected     bool      `json:"connected"`
	ClientVersion string    `json:"clientVersion"`
	Paused        bool      `json:"paused"`
	LastSync      time.Time `json:"at"`
	Device        Device
	InBytesTotal  int `json:"inBytesTotal"`
	OutBytesTotal int `json:"outBytesTotal"`
}

type Device struct {
	DeviceID string `json:"deviceID"`
	Name     string `json:"name"`
}
