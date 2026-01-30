package core

import (
	"time"

	"StatusApp/internal/hosthatch"
	"StatusApp/internal/schedule"
	"StatusApp/internal/tailscale"
	"StatusApp/internal/truenas"
	"StatusApp/internal/upcloud"
	"StatusApp/internal/weather"
)

type MainModel struct {
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
	HostHatchServers   []hosthatch.Server
	UpCloudServers     []upcloud.Server
	UpcloudAccountInfo upcloud.AccountInfo
}
