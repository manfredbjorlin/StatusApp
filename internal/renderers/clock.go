package renderers

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mbndr/figlet4go"

	"StatusApp/configs"
	"StatusApp/internal/models"
)

func RenderClock(m models.Model) string {
	currentTime := time.Now()
	digits := make([]int, 4)
	hour := currentTime.Hour()
	minute := currentTime.Minute()

	digits[0] = hour / 10
	digits[1] = hour % 10
	digits[2] = minute / 10
	digits[3] = minute % 10

	clock := ""
	lipglosspink := lipgloss.NewStyle().
		Bold(true).
		Foreground(configs.HotPink).
		Width(40).
		Align(lipgloss.Center)

	ascii := figlet4go.NewAsciiRender()
	// text, _ := os.ReadFile("mobius.txt")
	// clock = string(text)
	opts := figlet4go.NewRenderOptions()
	opts.FontName = "big"
	_ = ascii.LoadFont("/home/manfred/Development/StatusApp/")
	clock, _ = ascii.Render(currentTime.Format("15:04"))

	clock = lipglosspink.Render(clock)

	weatherIcon := getWeatherIcon(m.Weather.Current.Condition.Code, m.Weather.Current.IsDay)

	withText := lipgloss.JoinVertical(
		lipgloss.Center,
		clock,
		configs.BoldText.Render(time.Now().Format(configs.DateFormat)),
		configs.BoldText.Render(
			fmt.Sprintf(
				"%s  %v\ue33eC (%v\ue33eC)",
				weatherIcon,
				m.Weather.Current.Temp,
				m.Weather.Current.FeelsLike,
			),
		),
	)

	return configs.ClockStyle.Render(withText)
}
