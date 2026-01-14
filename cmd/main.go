package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"

	"StatusApp/configs"
	"StatusApp/internal/models"
	"StatusApp/internal/renderers"
	"StatusApp/internal/tailscale"
	"StatusApp/internal/weather"
)

type tickMsg time.Time

type mainModel struct {
	models.Model
}

func (m mainModel) Init() tea.Cmd {
	return tickCmd()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Could not load .env")
	}
	m := mainModel{}
	m.Devices.Devices = []models.Device{
		{
			Hostname:           "None",
			ConnectedToControl: true,
			Name:               "None",
			Os:                 "none",
		},
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		configs.WindowHeight = msg.Height
		configs.WindowWidth = msg.Width
		return m, nil
	case tickMsg:
		if configs.TailscaleWaits == 0 {
			devs := tailscale.TailscaleRequest()
			m.Devices = devs.(models.TailscaleMsg).Devices
			m.Misc = devs.(models.TailscaleMsg).Extra
			weather := weather.WeatherRequest()
			m.Weather = weather.(models.WeatherMsg).Weather
			keyExpiry := tailscale.GetKeyExpiry()
			m.KeyExpiry = keyExpiry.(models.TimeMsg).Time
			configs.TailscaleWaits = 60
		}
		configs.TailscaleWaits = configs.TailscaleWaits - 1
		return m, tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		return m, nil
	}
}

func (m mainModel) View() string {
	result := ""
	clock := renderers.RenderClock(m.Model)
	tailscale := renderers.RenderTailscale(m.Model)
	result = lipgloss.JoinHorizontal(
		lipgloss.Left,
		tailscale,
		clock,
	)
	schedule := renderers.RenderSchedule()
	result = lipgloss.JoinVertical(
		lipgloss.Top,
		result,
		schedule,
	)

	centered := lipgloss.Place(
		configs.WindowWidth,
		configs.WindowHeight,
		lipgloss.Center,
		lipgloss.Center,
		result,
	)

	return centered
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
