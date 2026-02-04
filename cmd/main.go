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
	"StatusApp/internal/syncthing"
	"StatusApp/internal/tailscale"
	"StatusApp/internal/truenas"
)

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m BubbleTeaModel) Init() tea.Cmd {
	return tickCmd()
}

func (m BubbleTeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		return m, nil

	case tickMsg:
		return tickUpdate(&m)
	case fetchedDataMsg:
		m.Data = msg.data
		m.Data.LastUpdated = msg.time
		m.Data.Error = nil
		return m, nil
	case errorMsg:
		m.Data.Error = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r": // Force refresh
			m.TickCounter = configs.SecondsBetweenRefresh
		case "a":
			m.CurrentScreen = configs.ScreenApps
		case "m":
			m.CurrentScreen = configs.ScreenMain
		case "s":
			m.CurrentScreen = configs.ScreenServers
		case "y":
			m.CurrentScreen = configs.ScreenSyncthing
		}
		number, err := strconv.Atoi(msg.String())
		if m.CurrentScreen == configs.ScreenApps && err == nil {
			_ = m.TruenasClient.UpdateApp(m.Data.TruenasApps, number)
			updatedModel := make([]truenas.App, 0)
			for _, _app := range m.Data.TruenasApps {
				if _app.UpdateId == number {
					_app.State = "STOPPED"
				}
				updatedModel = append(updatedModel, _app)
			}
			m.Data.TruenasApps = updatedModel
			m.TickCounter = configs.SecondsBetweenRefresh / 2
		}
	}
	return m, nil
}

func (m BubbleTeaModel) View() string {
	// if m.Data.Error != nil {
	// 	return lipgloss.Place(
	// 		m.WindowWidth,
	// 		m.WindowHeight,
	// 		lipgloss.Center,
	// 		lipgloss.Center,
	// 		fmt.Sprintf("Error: %s", m.Data.Error.Error()),
	// 	)
	// }
	if m.WindowWidth == 0 || m.WindowHeight == 0 || m.Data.Schedule == nil {
		return ViewInitialing(m)
	}
	var mainContent string
	var menuItemOrder []string

	switch m.CurrentScreen {
	case configs.ScreenServers:
		mainContent = ViewServers(m)
		menuItemOrder = []string{"q", "r", "m", "a", "y"}
	case configs.ScreenMain:
		topLeft := tailscale.View(
			m.Data.TailscaleDevices.Devices,
			m.Data.TailscaleKeyExpiry,
			m.AlternatingText,
		)
		mainContent = ViewMain(m, topLeft)
		menuItemOrder = []string{"q", "r", "a", "s", "y"}
	case configs.ScreenApps:
		topLeft := truenas.View(m.Data.TruenasApps)
		mainContent = ViewMain(m, topLeft)
		menuItemOrder = []string{"q", "r", "m", "s", "y"}
	case configs.ScreenSyncthing:
		mainContent = syncthing.View(m.Data.SyncthingConnections)
		menuItemOrder = []string{"q", "r", "m", "s"}
	}

	menuItems := map[string]string{
		"q": "quit",
		"r": "refresh",
		"a": "apps",
		"s": "servers",
		"m": "main screen",
		"e": "exclude",
		"y": "syncthing",
	}

	menu := ViewMenu(menuItems, menuItemOrder, m.Data.LastUpdated)

	final := lipgloss.JoinVertical(lipgloss.Top, mainContent, menu)

	return lipgloss.Place(
		m.WindowWidth,
		m.WindowHeight,
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
