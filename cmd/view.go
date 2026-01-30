package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mbndr/figlet4go"

	"StatusApp/configs"
	"StatusApp/internal/clock"
	"StatusApp/internal/hosthatch"
	"StatusApp/internal/schedule"
	"StatusApp/internal/truenas"
	"StatusApp/internal/upcloud"
	"StatusApp/internal/weather"
)

func ViewInitialing(m BubbleTeaModel) string {
	return lipgloss.Place(
		m.WindowWidth,
		m.WindowHeight,
		lipgloss.Center,
		lipgloss.Center,
		configs.ClockStyle.Render("Initializing..."),
	)
}

func ViewServers(m BubbleTeaModel) string {
	ascii := figlet4go.NewAsciiRender()
	header, _ := ascii.Render("Servers")
	heading := lipgloss.NewStyle().Foreground(configs.HotPink).Render(header)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Top,
		heading,
		hosthatch.View(m.Data.HostHatchServers, m.AlternatingText),
		upcloud.View(
			m.Data.UpCloudServers,
			m.AlternatingText,
			m.Data.UpcloudAccountInfo,
		),
	)
	return mainContent
}

func ViewMain(m BubbleTeaModel, topLeft string) string {
	weatherView := weather.View(
		m.Data.Weather,
		m.Data.WaterTemperature,
		m.AlternatingText,
	)

	clockView := clock.RenderClock(weatherView)

	top := lipgloss.JoinHorizontal(lipgloss.Left, topLeft, clockView)

	scheduleView := schedule.View(m.Data.Schedule)

	statusStyle := lipgloss.NewStyle().
		Width(configs.ScheduleStyle.GetWidth()).
		AlignHorizontal(lipgloss.Center)
	statusLine := statusStyle.Render(
		truenas.Status(
			m.Data.TruenasApps,
		) + " - " + hosthatch.Status(
			m.Data.HostHatchServers,
		) + " - " + upcloud.Status(
			m.Data.UpCloudServers,
		),
	)

	return lipgloss.JoinVertical(lipgloss.Left, statusLine, top, scheduleView)
}
