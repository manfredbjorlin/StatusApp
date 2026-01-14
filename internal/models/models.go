package models

import "time"

type Meeting struct {
	Time  time.Time
	End   time.Time
	Title string
	Room  string
	Rooms []string
}

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

type Weathercode struct {
	Code  int    `json:"code"`
	Day   string `json:"day"`
	Night string `json:"night"`
	Icon  int    `json:"icon"`
}

type Weather struct {
	Current Current `json:"current"`
}

type Current struct {
	Condition Condition `json:"condition"`
	Temp      float32   `json:"temp_c"`
	FeelsLike float32   `json:"feelslike_c"`
	IsDay     int       `json:"is_day"`
}

type Condition struct {
	Code int `json:"code"`
}

type Model struct {
	Devices   Devices
	KeyExpiry time.Time
	Misc      string
	Weather   Weather
}

type (
	ErrMsg       struct{ Err error }
	TailscaleMsg struct {
		Devices Devices
		Extra   string
	}
	WeatherMsg struct{ Weather Weather }
	TimeMsg    struct{ Time time.Time }
)
