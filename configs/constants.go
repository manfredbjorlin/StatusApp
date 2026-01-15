package configs

import "github.com/charmbracelet/lipgloss"

const (
	BaseUrl    = "https://api.tailscale.com/api/v2"
	DateFormat = "2. January 2006"
)

var (
	CurrentScreen    = ScreenMain
	TailscaleWaits   = 0
	SoonMeeting      = false
	InMeeting        = false
	ScheduleOutlined = false
	TailscaleVersion = true
	TailscaleRenders = 0
	WindowHeight     = 0
	WindowWidth      = 0
	ColorPink        = "#ff07a9"
	ColorBlue        = "#07a9ff"
	ColorGreen       = "#a9ff07"
	ColorGrey        = "#c0c0c0"
	ColorWhite       = "#FFFFFF"
	// colorBg          = "#202020"
	ColorAltBg = "#282828"

	HotPink      = lipgloss.Color(ColorPink)
	NiceBlue     = lipgloss.Color(ColorBlue)
	BrightGreen  = lipgloss.Color(ColorGreen)
	StandardText = lipgloss.Color(ColorGrey)
	HeaderText   = lipgloss.Color(ColorWhite)
	AltBgColor   = lipgloss.Color(ColorAltBg)

	ClockStyle = lipgloss.NewStyle().
			Padding(1, 2, 1, 2).
			Margin(1, 1, 1, 1).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(NiceBlue)

	TailscaleStyle = lipgloss.NewStyle().
			Padding(1, 2, 0, 2).
			Margin(1, 1, 1, 1).
			Width(35)

	ScheduleStyle = lipgloss.NewStyle().
			Padding(1, 2, 1, 2).
			Margin(0, 1, 1, 1).
			Width(82).
			BorderStyle(lipgloss.RoundedBorder())

	BoldText = lipgloss.NewStyle().
			Bold(true).
			Foreground(HeaderText)
)

const (
	ScreenMain = iota
	ScreenApps
)

func SetBg(style lipgloss.Style, line int) lipgloss.Style {
	if line%2 == 0 {
		return style.Background(AltBgColor)
	} else {
		return style.UnsetBackground()
	}
}
