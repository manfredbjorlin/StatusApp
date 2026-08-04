package main

import (
	"os"
	"time"

	"StatusApp/configs"
	"StatusApp/internal/exaroton"
	"StatusApp/internal/hosthatch"
	"StatusApp/internal/netbird"
	"StatusApp/internal/syncthing"
	"StatusApp/internal/tailscale"
	"StatusApp/internal/truenas"
	"StatusApp/internal/upcloud"
	"StatusApp/internal/weather"
	"StatusApp/pkg/core"
)

type BubbleTeaModel struct {
	TailscaleClient tailscale.MachineGetter
	NetBirdClient   netbird.MachineGetter
	WeatherClient   weather.DataProvider
	TruenasClient   truenas.AppGetter
	HostHatchClient hosthatch.HostHatchClient
	UpCloudClient   upcloud.UpCloudClient
	ExarotonClient  exaroton.ExarotonClient
	SyncthingClient syncthing.SyncthingClient

	CurrentScreen             configs.Screen
	WindowWidth, WindowHeight int
	TickCounter               int
	AlternatingText           bool

	Data core.MainModel
}

func newModel() BubbleTeaModel {
	return BubbleTeaModel{
		CurrentScreen: configs.ScreenMain,
		TailscaleClient: tailscale.NewClient(
			os.Getenv("TAILSCALE_API_KEY"),
			os.Getenv("TAILSCALE_TAILNET_ID"),
			os.Getenv("TAILSCALE_API_KEY_ID"),
		),
		NetBirdClient: netbird.NewClient(
			os.Getenv("NETBIRD_API_KEY"),
			os.Getenv("NETBIRD_USER_ID"),
			os.Getenv("NETBIRD_API_KEY_ID"),
		),
		WeatherClient: weather.NewClient(
			os.Getenv("WEATHERAPI_API_KEY"),
			os.Getenv("WEATHERAPI_LOCATION"),
			os.Getenv("WATERTEMPERATURE_LOCATION_ID"),
		),
		TruenasClient: truenas.NewClient(
			os.Getenv("TRUENAS_BASE_URL"),
			os.Getenv("TRUENAS_API_KEY"),
		),
		HostHatchClient: hosthatch.NewClient(os.Getenv("HOSTHATCH_API_KEY")),
		UpCloudClient:   upcloud.NewClient(os.Getenv("UPCLOUD_API_KEY")),
		ExarotonClient:  exaroton.NewClient(os.Getenv("EXAROTON_API_KEY")),
		SyncthingClient: syncthing.NewClient(os.Getenv("SYNCTHING_API_KEY")),
		TickCounter:     60, // Start ready to fetch
		AlternatingText: false,
		Data:            core.MainModel{},
	}
}

type (
	tickMsg        time.Time
	errorMsg       struct{ err error }
	fetchedDataMsg struct {
		time time.Time
		data core.MainModel
	}
)
