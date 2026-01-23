package main

import (
	"context"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"StatusApp/internal/schedule"
	"StatusApp/internal/tailscale"
	"StatusApp/internal/truenas"
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

		var tsDevices tailscale.Devices
		var tsKey time.Time
		var tnApps []truenas.App
		var weatherData *weather.WeatherForecastInternal
		var waterTempData weather.WaterTemperature
		var scheduleData []schedule.Meeting
		var err error

		// Run fetches in parallel
		errs := make(chan error, 5)

		go func() {
			tsDevices, err = m.Model.TailscaleClient.GetMachines(ctx)
			errs <- err
		}()
		go func() {
			tsKey, err = m.Model.TailscaleClient.GetKeyExpiry(ctx)
			errs <- err
		}()
		go func() {
			tnApps, err = m.Model.TruenasClient.GetApps(ctx)
			errs <- err
		}()
		go func() {
			weatherData, err = m.Model.WeatherClient.GetCurrentWeather(ctx)
			errs <- err
		}()
		go func() {
			waterTempData, err = m.Model.WeatherClient.GetWaterTemperature(ctx)
			errs <- err
		}()

		// Process results
		for range 5 {
			if err := <-errs; err != nil {
				return errorMsg{err} // Return on the first error
			}
		}

		// Schedule is loaded synchronously as it's a file read
		scheduleFile := os.Getenv("SCHEDULE_FILE_PATH")
		scheduleData, err = schedule.LoadSchedule(scheduleFile)
		if err != nil {
			return errorMsg{err}
		}

		return fetchedDataMsg{
			tailscaleDevices:   tsDevices,
			tailscaleKeyExpiry: tsKey,
			truenasApps:        tnApps,
			weather:            weatherData,
			waterTemperature:   waterTempData,
			schedule:           scheduleData,
		}
	}
}
