package core

import (
	"time"

	"StatusApp/internal/schedule"
	"StatusApp/internal/tailscale"
	"StatusApp/internal/truenas"
	"StatusApp/internal/weather"
)

type MainModel struct {
	// Clients
	TailscaleClient tailscale.MachineGetter
	WeatherClient   weather.DataProvider
	TruenasClient   truenas.AppGetter

	// UI State
	WindowWidth, WindowHeight int
	TickCounter               int
	AlternatingText           bool

	// Data
	Error              error
	LastUpdated        time.Time
	Schedule           []schedule.Meeting
	TailscaleDevices   tailscale.Devices
	TailscaleKeyExpiry time.Time
	TruenasApps        []truenas.App
	Weather            *weather.WeatherForecastInternal
	WaterTemperature   weather.WaterTemperatureInternal
	TruenasUpdateList  map[int]string
}
