package netbird

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"StatusApp/configs"
	"StatusApp/internal/common"
)

func View(
	peers []Peer,
	keyExpiry time.Time,
	latestVersion string,
	alternatingText bool,
) string {
	var sb strings.Builder
	greenBold := lipgloss.NewStyle().Bold(true).Foreground(configs.BrightGreen)
	pinkBold := lipgloss.NewStyle().Bold(true).Foreground(configs.HotPink)
	blueBold := lipgloss.NewStyle().Bold(true).Foreground(configs.NiceBlue)

	fmt.Fprintf(
		&sb,
		"\n%s\n",
		strings.Repeat(" ", 6)+configs.BoldText.Render("\uef08 The Bird Valley"),
	)

	sb.WriteString("\n")

	hostname, _ := os.Hostname()
	for i, peer := range peers {
		deviceIcon := "\uf17c"
		if strings.HasPrefix(peer.Os, "Android") {
			deviceIcon = "\ue70e"
		} else if strings.HasPrefix(peer.Os, "Windows") {
			deviceIcon = "\uf17a"
		}
		logoStyle := pinkBold
		if peer.Connected {
			logoStyle = greenBold
		}
		logoStyle = configs.SetBg(logoStyle, i)
		sb.WriteString(logoStyle.Render(deviceIcon))

		caser := cases.Title(language.BrazilianPortuguese)
		name := caser.String(peer.Name)
		indicator := ""
		indicatorPadding := 0
		if peer.Hostname == hostname {
			indicator = " \uf256"
			indicatorPadding = 2
		}
		spacing := 21 - len(name) - indicatorPadding
		nameStyle := lipgloss.NewStyle()
		nameStyle = configs.SetBg(nameStyle, i)
		exitIcon := ""
		keyIcon := ""
		exitNode := false
		for _, group := range peer.Groups {
			if group.Name == "Exit Nodes" {
				exitNode = true
				break
			}
		}
		if exitNode {
			exitIcon = nameStyle.Render(
				" ",
			) + nameStyle.Foreground(configs.NiceBlue).
				Render("\uea6e")
			spacing -= 2

			if len(peer.CountryCode) > 0 {
				location := " " + peer.CityName
				spacing -= len(location)
				exitIcon += nameStyle.Render(location)
			}
		}
		if peer.InactivityExpirationEnabled {
			keyIcon = nameStyle.Render(
				" ",
			) + nameStyle.Foreground(configs.NiceBlue).
				Render("\ue641")
			spacing -= 2
		}
		sb.WriteString(nameStyle.Render(
			" " + name,
		))
		sb.WriteString(nameStyle.Foreground(configs.BrightGreen).Render(indicator))
		sb.WriteString(exitIcon)
		sb.WriteString(keyIcon)
		sb.WriteString(nameStyle.Render(strings.Repeat(
			" ",
			spacing,
		)))

		if peer.Connected || !alternatingText {
			updateStyle := greenBold
			updateLogo := "\uf00c"

			if strings.HasPrefix(peer.Os, "Android") {
				updateLogo = "\uf128"
				updateStyle = blueBold
			} else if peer.Version != latestVersion {
				updateLogo = "\uf00d"
				updateStyle = pinkBold
			}

			sb.WriteString(updateStyle.Render(updateLogo))
			sb.WriteString(nameStyle.Render(" " + peer.Version))
			sb.WriteString("\n")
		} else {
			diffText := common.GetTimeDifferenceString(peer.LastSeen)
			sb.WriteString(nameStyle.Render(fmt.Sprintf("\uf017 %6s", diffText)))
			sb.WriteString("\n")
		}
	}

	offlineDiff := time.Until(keyExpiry)
	diffText := common.GetTimeDifferenceString(keyExpiry)
	sb.WriteString("\nNetbird key expiry:    ")
	keytext := fmt.Sprintf("%s%6s", "\uf017 ", diffText)
	if offlineDiff.Hours() < (24 * 4) {
		sb.WriteString(pinkBold.Render(keytext))
	} else {
		sb.WriteString(keytext)
	}

	res := configs.TailscaleStyle.Render(strings.TrimSuffix(sb.String(), "\n"))
	return res
}
