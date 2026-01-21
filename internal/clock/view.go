package clock

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mbndr/figlet4go"

	"StatusApp/configs"
)

// RenderClock creates a string with the ASCII art clock and status information.
func RenderClock(weatherPart string) string {
	currentTime := time.Now()

	// Render ASCII Clock
	ascii := figlet4go.NewAsciiRender()
	// opts := figlet4go.NewRenderOptions()
	// opts.FontName = "big"
	// _ = ascii.LoadFont(fontPath)
	clockStr, _ := ascii.Render(currentTime.Format("15:04"))

	lipglossPink := lipgloss.NewStyle().
		Bold(true).
		Foreground(configs.HotPink).
		Width(40).
		Align(lipgloss.Center)

	clock := lipglossPink.Render(clockStr)

	// Prepare status line (error or weather)

	// Join components vertically
	withText := lipgloss.JoinVertical(
		lipgloss.Center,
		clock,
		configs.BoldText.Render(time.Now().Format(configs.DateFormat)),
		configs.BoldText.Render(weatherPart),
	)

	return configs.ClockStyle.Render(withText)
}
