package clock

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mbndr/figlet4go"

	"StatusApp/configs"
)

// Is it the best thing to have the clock view render the error message? Maybe not, but it makes it clean...
func RenderClock(weatherPart string, err error) string {
	currentTime := time.Now()

	ascii := figlet4go.NewAsciiRender()
	var clockStr string
	if err == nil {
		clockStr, _ = ascii.Render(currentTime.Format("15:04"))
	} else {
		clockStr, _ = ascii.Render("Error")
		weatherPart = err.Error()
	}

	lipglossPink := lipgloss.NewStyle().
		Bold(true).
		Foreground(configs.HotPink).
		Width(40).
		Align(lipgloss.Center)

	clock := lipglossPink.Render(clockStr)

	withText := lipgloss.JoinVertical(
		lipgloss.Center,
		clock,
		configs.BoldText.Render(time.Now().Format(configs.DateFormat)),
		configs.BoldText.Render(weatherPart),
	)

	return configs.ClockStyle.Render(withText)
}
