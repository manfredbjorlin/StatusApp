package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"

	"StatusApp/configs"
	"StatusApp/internal/clock"
	"StatusApp/internal/schedule"
	"StatusApp/internal/tailscale"
	"StatusApp/internal/truenas"
	"StatusApp/internal/weather"
	"StatusApp/pkg/core"
)

// --- Model ---

type MainModel struct {
	Model core.MainModel
}

func newModel() MainModel {
	// Load environment variables
	tailscaleAPIKey := os.Getenv("TAILSCALE_API_KEY")
	tailscaleTailnet := os.Getenv("TAILSCALE_TAILNET_ID")
	tailscaleKeyID := os.Getenv("TAILSCALE_API_KEY_ID")
	weatherAPIKey := os.Getenv("WEATHERAPI_API_KEY")
	weatherLocation := os.Getenv("WEATHERAPI_LOCATION")
	waterLocationID := os.Getenv("WATERTEMPERATURE_LOCATION_ID")
	truenasURL := os.Getenv("TRUENAS_BASE_URL")
	truenasAPIKey := os.Getenv("TRUENAS_API_KEY")

	// Initialize clients
	tailscaleClient := tailscale.NewClient(tailscaleAPIKey, tailscaleTailnet, tailscaleKeyID)
	weatherClient := weather.NewClient(weatherAPIKey, weatherLocation, waterLocationID)
	truenasClient := truenas.NewClient(truenasURL, truenasAPIKey)

	return MainModel{
		Model: core.MainModel{
			TailscaleClient: tailscaleClient,
			WeatherClient:   weatherClient,
			TruenasClient:   truenasClient,
			TickCounter:     60, // Start ready to fetch
			AlternatingText: false,
		},
	}
}

// --- Messages ---

type (
	tickMsg        time.Time
	errorMsg       struct{ err error }
	fetchedDataMsg struct {
		tailscaleDevices   tailscale.Devices
		tailscaleKeyExpiry time.Time
		truenasApps        []truenas.App
		weather            weather.Weather
		waterTemperature   weather.WaterTemperature
		schedule           []schedule.Meeting
	}
)

// --- Commands ---

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchData(m MainModel) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var tsDevices tailscale.Devices
		var tsKey time.Time
		var tnApps []truenas.App
		var weatherData weather.Weather
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

// --- Bubbletea Program ---

func (m MainModel) Init() tea.Cmd {
	return tickCmd()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Model.WindowWidth = msg.Width
		m.Model.WindowHeight = msg.Height
		return m, nil

	case tickMsg:
		m.Model.TickCounter++
		if m.Model.TickCounter%configs.SecondsBetweenAlternatingText == 0 {
			m.Model.AlternatingText = !m.Model.AlternatingText
		}
		if m.Model.TickCounter >= configs.SecondsBetweenRefresh {
			m.Model.TickCounter = 0
			return m, tea.Batch(fetchData(m), tickCmd())
		}
		return m, tickCmd()

	case fetchedDataMsg:
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

	case errorMsg:
		m.Model.Error = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r": // Force refresh
			m.Model.TickCounter = configs.SecondsBetweenRefresh
			return m, nil
		}
	}
	return m, nil
}

func (m MainModel) View() string {
	if m.Model.WindowWidth == 0 || m.Model.WindowHeight == 0 {
		return "Initializing..."
	}

	// Weather View (Simplified)
	iconPath := os.Getenv("WEATHER_ICON_PATH")
	weatherView := weather.View(
		m.Model.Weather,
		m.Model.WaterTemperature,
		m.Model.AlternatingText,
		iconPath,
	)

	// Clock View
	clockView := clock.RenderClock(weatherView)

	// Tailscale View
	tailscaleView := tailscale.View(
		m.Model.TailscaleDevices.Devices,
		m.Model.TailscaleKeyExpiry,
		m.Model.TruenasApps,
		m.Model.AlternatingText,
	)

	// Top section
	top := lipgloss.JoinHorizontal(lipgloss.Left, tailscaleView, clockView)

	// Schedule View
	scheduleView := schedule.View(m.Model.Schedule)

	mainContent := lipgloss.JoinVertical(lipgloss.Left, top, scheduleView)

	// Menu
	menuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(configs.ColorDimGrey)).
		Width(configs.ScheduleStyle.GetWidth()).
		AlignHorizontal(lipgloss.Center)
	menu := menuStyle.Render(
		fmt.Sprintf(
			"q: quit | r: refresh | Last update: %s",
			m.Model.LastUpdated.Format("15:04:05"),
		),
	)

	final := lipgloss.JoinVertical(lipgloss.Top, mainContent, menu)

	return lipgloss.Place(
		m.Model.WindowWidth,
		m.Model.WindowHeight,
		lipgloss.Center,
		lipgloss.Center,
		final,
	)
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Could not load .env file")
		os.Exit(1)
	}

	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
