package main

import (
	"context"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	logger "StatusApp/internal"
	"StatusApp/internal/exaroton"
	"StatusApp/internal/hosthatch"
	"StatusApp/internal/schedule"
	"StatusApp/internal/upcloud"
	"StatusApp/internal/weather"
	"StatusApp/pkg/core"
)

func fetchData(m *BubbleTeaModel) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var err error
		type updateFuncs func()
		var errs chan error

		result := core.MainModel{}

		updateFunctions := []updateFuncs{
			func() {
				result.TailscaleDevices, err = m.TailscaleClient.GetMachines(ctx)
				errs <- err
			},
			func() {
				result.TailscaleKeyExpiry, err = m.TailscaleClient.GetKeyExpiry(ctx)
				errs <- err
			},
			func() {
				result.NetBirdPeers, err = m.NetBirdClient.GetMachines(ctx)
				errs <- err
			}, func() {
				result.NetBirdKeyExpiry, err = m.NetBirdClient.GetKeyExpiry(ctx)
				errs <- err
			}, func() {
				result.TruenasApps, err = m.TruenasClient.GetApps(ctx)
				errs <- err
			}, func() {
				result.Weather, err = m.WeatherClient.GetCurrentWeather(ctx)
				errs <- err
			}, func() {
				waterTemp, err := m.WeatherClient.GetWaterTemperature(ctx)
				if len(waterTemp.Embedded.NearestLocations) > 0 {
					loc := waterTemp.Embedded.NearestLocations[0]
					result.WaterTemperature = weather.WaterTemperatureInternal{
						Place:       loc.Location.Name,
						Temperature: loc.Temperature,
						LastUpdate:  loc.Time,
					}
				}
				errs <- err
			}, func() {
				result.HostHatchServers, err = m.HostHatchClient.ListServers(ctx)
				if err != nil {
					result.HostHatchServers = []hosthatch.Server{
						{Hostname: "HostHatch: N/A"},
					}
				}
				errs <- err
			}, func() {
				result.UpCloudServers, err = m.UpCloudClient.ListServers(ctx)
				errs <- err
			}, func() {
				result.ExarotonServers, err = m.ExarotonClient.ListServers(ctx)
				if err != nil {
					result.ExarotonServers = []exaroton.Server{
						{Name: "Exaroton: N/A", Status: 0},
					}
					errs <- err
				}
				result.ExarotonCreditLeft, err = m.ExarotonClient.RemainingCredits(ctx)
				errs <- err
			}, func() {
				result.SyncthingConnections, err = m.SyncthingClient.ListConnections(ctx)
				errs <- err
			}, func() {
				_, remaining, err := m.UpCloudClient.RemainingCredits(
					ctx,
				)
				if err != nil {
					errs <- err
					return
				}
				currency, usage, err := m.UpCloudClient.BillingSummary(
					ctx,
					time.Now().Year(),
					int(time.Now().Month()),
				)
				errs <- err
				result.UpcloudAccountInfo = upcloud.AccountInfo{
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
				logger.LogError(err.Error())
			}
		}

		// Schedule is loaded synchronously as it's a file read
		scheduleFile := os.Getenv("SCHEDULE_FILE_PATH")
		result.Schedule, err = schedule.LoadSchedule(scheduleFile)
		if err != nil {
			return errorMsg{err}
		}

		for {
			select {
			case <-time.After(20 * time.Second):
				cancel()
			case <-ctx.Done():
				return fetchedDataMsg{
					time: time.Now(),
					data: result,
				}
			}
		}
	}
}
