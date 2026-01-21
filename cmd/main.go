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
)

// --- Model ---

type mainModel struct {
	// Clients
	tsClient      tailscale.MachineGetter
	weatherClient weather.DataProvider
	truenasClient truenas.AppGetter

	// UI State
	width, height   int
	inMeeting       bool
	soonMeeting     bool
	outlined        bool
	tickCounter     int
	alternatingText bool

	// Data
	err              error
	lastUpdated      time.Time
	schedule         []schedule.Meeting
	tailscaleDevices tailscale.Devices
	tailscaleKey     time.Time
	truenasApps      []truenas.App
	weather          weather.Weather
	waterTemperature weather.WaterTemperatureInternal
}

func newModel() mainModel {
	// Load environment variables
	tsAPIKey := os.Getenv("TAILSCALE_API_KEY")
	tsTailnet := os.Getenv("TAILSCALE_TAILNET_ID")
	tsKeyID := os.Getenv("TAILSCALE_API_KEY_ID")
	weatherAPIKey := os.Getenv("WEATHERAPI_API_KEY")
	weatherLocation := os.Getenv("WEATHERAPI_LOCATION")
	waterLocationID := os.Getenv("WATERTEMPERATURE_LOCATION_ID")
	truenasURL := os.Getenv("TRUENAS_BASE_URL")
	truenasAPIKey := os.Getenv("TRUENAS_API_KEY")

	// Initialize clients
	tsClient := tailscale.NewClient(tsAPIKey, tsTailnet, tsKeyID)
	weatherClient := weather.NewClient(weatherAPIKey, weatherLocation, waterLocationID)
	truenasClient := truenas.NewClient(truenasURL, truenasAPIKey)

	return mainModel{
		tsClient:        tsClient,
		weatherClient:   weatherClient,
		truenasClient:   truenasClient,
		tickCounter:     60, // Start ready to fetch
		alternatingText: false,
	}
}

// --- Messages ---

type (
	tickMsg        time.Time
	errMsg         struct{ err error }
	fetchedDataMsg struct {
		tsDevices tailscale.Devices
		tsKey     time.Time
		tnApps    []truenas.App
		weather   weather.Weather
		waterTemp weather.WaterTemperature
		schedule  []schedule.Meeting
	}
)

func (e errMsg) Error() string { return e.err.Error() }

// --- Commands ---

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchData(m mainModel) tea.Cmd {
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
			tsDevices, err = m.tsClient.GetMachines(ctx)
			errs <- err
		}()
		go func() {
			tsKey, err = m.tsClient.GetKeyExpiry(ctx)
			errs <- err
		}()
		go func() {
			tnApps, err = m.truenasClient.GetApps(ctx)
			errs <- err
		}()
		go func() {
			weatherData, err = m.weatherClient.GetCurrentWeather(ctx)
			errs <- err
		}()
		go func() {
			waterTempData, err = m.weatherClient.GetWaterTemperature(ctx)
			errs <- err
		}()

		// Process results
		for i := 0; i < 5; i++ {
			if err := <-errs; err != nil {
				return errMsg{err} // Return on the first error
			}
		}

		// Schedule is loaded synchronously as it's a file read
		scheduleFile := os.Getenv("SCHEDULE_FILE_PATH")
		scheduleData, err = schedule.LoadSchedule(scheduleFile)
		if err != nil {
			return errMsg{err}
		}

		return fetchedDataMsg{
			tsDevices: tsDevices,
			tsKey:     tsKey,
			tnApps:    tnApps,
			weather:   weatherData,
			waterTemp: waterTempData,
			schedule:  scheduleData,
		}
	}
}

// --- Bubbletea Program ---

func (m mainModel) Init() tea.Cmd {
	// return tea.Batch(tickCmd(), fetchData(m))
	return tickCmd()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Also update legacy configs for now
		configs.WindowWidth = msg.Width
		configs.WindowHeight = msg.Height
		return m, nil

	case tickMsg:
		m.tickCounter++
		if m.tickCounter%5 == 0 {
			m.alternatingText = !m.alternatingText
		}
		if m.tickCounter >= 60 { // Fetch data every 60 seconds
			m.tickCounter = 0
			return m, fetchData(m)
		}
		return m, tickCmd()

	case fetchedDataMsg:
		m.err = nil
		m.lastUpdated = time.Now()
		m.tailscaleDevices = msg.tsDevices
		m.tailscaleKey = msg.tsKey
		m.truenasApps = msg.tnApps
		m.weather = msg.weather
		m.schedule = msg.schedule

		if len(msg.waterTemp.Embedded.NearestLocations) > 0 {
			loc := msg.waterTemp.Embedded.NearestLocations[0]
			m.waterTemperature = weather.WaterTemperatureInternal{
				Place:       loc.Location.Name,
				Temperature: loc.Temperature,
				LastUpdate:  loc.Time,
			}
		}
		return m, tickCmd()

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r": // Force refresh
			m.tickCounter = 60
			return m, nil
		}
	}
	return m, nil
}

// Implement local interfaces for the view functions
func (m mainModel) GetTruenasApps() []truenas.App          { return m.truenasApps }
func (m mainModel) GetTailscaleDevices() tailscale.Devices { return m.tailscaleDevices }
func (m mainModel) GetKeyExpiry() time.Time                { return m.tailscaleKey }
func (m mainModel) GetWeather() weather.Weather            { return m.weather }
func (m mainModel) GetWaterTemperature() weather.WaterTemperatureInternal {
	return m.waterTemperature
}
func (m mainModel) DisplayAlternatingText() bool { return m.alternatingText }

func (m mainModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Weather View (Simplified)
	iconPath := os.Getenv("WEATHER_ICON_PATH")
	weatherView := weather.View(m, iconPath)

	// Clock View
	clockView := clock.RenderClock(weatherView)

	// Tailscale View
	tsView := tailscale.View(m)

	// Top section
	top := lipgloss.JoinHorizontal(lipgloss.Left, tsView, clockView)

	// Schedule View
	scheduleView := schedule.View(m.schedule, &m.inMeeting, &m.soonMeeting, &m.outlined)

	mainContent := lipgloss.JoinVertical(lipgloss.Left, top, scheduleView)

	// Menu
	menuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Width(configs.ScheduleStyle.GetWidth()).
		AlignHorizontal(lipgloss.Center)
	menu := menuStyle.Render(
		fmt.Sprintf("q: quit | r: refresh | Last update: %s", m.lastUpdated.Format("15:04:05")),
	)

	final := lipgloss.JoinVertical(lipgloss.Top, mainContent, menu)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, final)
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
