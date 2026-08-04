package core

import (
	"time"

	"StatusApp/internal/exaroton"
	"StatusApp/internal/hosthatch"
	"StatusApp/internal/netbird"
	"StatusApp/internal/schedule"
	"StatusApp/internal/syncthing"
	"StatusApp/internal/tailscale"
	"StatusApp/internal/truenas"
	"StatusApp/internal/upcloud"
	"StatusApp/internal/weather"
)

type MainModel struct {
	Error                error
	LastUpdated          time.Time
	Schedule             []schedule.Meeting
	TailscaleDevices     tailscale.Devices
	TailscaleKeyExpiry   time.Time
	NetBirdPeers         []netbird.Peer
	NetBirdKeyExpiry     time.Time
	TruenasApps          []truenas.App
	Weather              *weather.WeatherForecastInternal
	WaterTemperature     weather.WaterTemperatureInternal
	TruenasUpdateList    map[int]string
	HostHatchServers     []hosthatch.Server
	UpCloudServers       []upcloud.Server
	UpcloudAccountInfo   upcloud.AccountInfo
	ExarotonServers      []exaroton.Server
	ExarotonCreditLeft   float64
	SyncthingConnections []syncthing.Connection
}
