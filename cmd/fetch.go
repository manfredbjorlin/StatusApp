package main

import (
	"context"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"StatusApp/internal/schedule"
	"StatusApp/internal/upcloud"
	"StatusApp/internal/weather"
)

func mapFetchedData(m BubbleTeaModel, msg fetchedDataMsg) (tea.Model, tea.Cmd) {
	m.Model.Error = nil
	m.Model.LastUpdated = time.Now()
	m.Model.TailscaleDevices = msg.tailscaleDevices
	m.Model.TailscaleKeyExpiry = msg.tailscaleKeyExpiry
	m.Model.TruenasApps = msg.truenasApps
	m.Model.Weather = msg.weather
	m.Model.Schedule = msg.schedule
	m.Model.HostHatchServers = msg.hosthatchServers
	m.Model.UpCloudServers = msg.upcloudServers
	m.Model.UpcloudAccountInfo = msg.upcloudAccountInfo

	if len(msg.waterTemperature.Embedded.NearestLocations) > 0 {
		loc := msg.waterTemperature.Embedded.NearestLocations[0]
		m.Model.WaterTemperature = weather.WaterTemperatureInternal{
			Place:       loc.Location.Name,
			Temperature: loc.Temperature,
			LastUpdate:  loc.Time,
		}
	}
	return m, nil
}

func fetchData(m BubbleTeaModel) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var result fetchedDataMsg
		var err error

		type updateFuncs func()
		var errs chan error

		updateFunctions := []updateFuncs{
			func() {
				result.tailscaleDevices, err = m.Model.TailscaleClient.GetMachines(ctx)
				errs <- err
			},
			func() {
				result.tailscaleKeyExpiry, err = m.Model.TailscaleClient.GetKeyExpiry(ctx)
				errs <- err
			}, func() {
				result.truenasApps, err = m.Model.TruenasClient.GetApps(ctx)
				errs <- err
			}, func() {
				result.weather, err = m.Model.WeatherClient.GetCurrentWeather(ctx)
				errs <- err
			}, func() {
				result.waterTemperature, err = m.Model.WeatherClient.GetWaterTemperature(ctx)
				errs <- err
			}, func() {
				result.hosthatchServers, err = m.Model.HostHatchClient.ListServers(ctx)
				errs <- err
			}, func() {
				result.upcloudServers, err = m.Model.UpCludClient.ListServers(ctx)
				errs <- err
			}, func() {
				_, remaining, err := m.Model.UpCludClient.RemainingCredits(
					ctx,
				)
				if err != nil {
					errs <- err
					return
				}
				currency, usage, err := m.Model.UpCludClient.BillingSummary(
					ctx,
					time.Now().Year(),
					int(time.Now().Month()),
				)
				errs <- err
				result.upcloudAccountInfo = upcloud.AccountInfo{
					RemainingCredits: remaining,
					Currency:         currency,
					BillingSummary:   usage,
				}
			},
		}

		errs = make(chan error, len(updateFunctions))
		for _, updFunc := range updateFunctions {
			go updFunc()
		}

		for range len(updateFunctions) {
			if err := <-errs; err != nil {
				return errorMsg{err} // Return on the first error
			}
		}

		// Schedule is loaded synchronously as it's a file read
		scheduleFile := os.Getenv("SCHEDULE_FILE_PATH")
		result.schedule, err = schedule.LoadSchedule(scheduleFile)
		if err != nil {
			return errorMsg{err}
		}

		return result
	}
}
