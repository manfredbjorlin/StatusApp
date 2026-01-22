package main

import (
	"fmt"
	"os"
	"strconv"
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

type BubbleTeaModel struct {
	CurrentScreen configs.Screen
	Model         core.MainModel
}

func newModel() BubbleTeaModel {
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

	return BubbleTeaModel{
		CurrentScreen: configs.ScreenMain,
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

// --- Bubbletea Program ---

func (m BubbleTeaModel) Init() tea.Cmd {
	return tickCmd()
}

func (m BubbleTeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Model.WindowWidth = msg.Width
		m.Model.WindowHeight = msg.Height
		return m, nil

	case tickMsg:
		return tickUpdate(m)
	case fetchedDataMsg:
		return mapFetchedData(m, msg)
	case errorMsg:
		m.Model.Error = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r": // Force refresh
			m.Model.TickCounter = configs.SecondsBetweenRefresh
		case "a":
			m.CurrentScreen = configs.ScreenApps
		case "m":
			m.CurrentScreen = configs.ScreenMain
		}
		number, err := strconv.Atoi(msg.String())
		if m.CurrentScreen == configs.ScreenApps && err == nil {
			_ = m.Model.TruenasClient.UpdateApp(m.Model.TruenasApps, number)
			m.Model.TickCounter = configs.SecondsBetweenRefresh
		}
	}
	return m, nil
}

func (m BubbleTeaModel) View() string {
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

	var topLeft string
	if m.CurrentScreen == configs.ScreenMain {
		topLeft = tailscale.View(
			m.Model.TailscaleDevices.Devices,
			m.Model.TailscaleKeyExpiry,
			m.Model.TruenasApps,
			m.Model.AlternatingText,
		)
	} else {
		topLeft = truenas.View(m.Model.TruenasApps)
	}

	// Top section
	top := lipgloss.JoinHorizontal(lipgloss.Left, topLeft, clockView)

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
			"q: quit | r: refresh | a: apps | m: main screen | Last update: %s",
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
