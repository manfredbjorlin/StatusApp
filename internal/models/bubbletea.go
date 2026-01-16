package models

import "time"

type Model struct {
	Devices     Devices
	KeyExpiry   time.Time
	Misc        string
	Weather     Weather
	WaterTemp   WaterTemperatureInternal
	TruenasApps []TruenasApp
}

type (
	ErrMsg       struct{ Err error }
	TailscaleMsg struct {
		Devices Devices
		Extra   string
	}
	WeatherMsg   struct{ Weather Weather }
	WaterTempMsg struct{ WaterTemp WaterTemperature }
	TimeMsg      struct{ Time time.Time }
	TruenasMsg   struct{ Apps []TruenasApp }
)
