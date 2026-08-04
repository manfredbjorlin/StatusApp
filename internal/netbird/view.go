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
	alternatingText bool,
) string {
	var sb strings.Builder
	greenBold := lipgloss.NewStyle().Bold(true).Foreground(configs.BrightGreen)
	pinkBold := lipgloss.NewStyle().Bold(true).Foreground(configs.HotPink)

	fmt.Fprintf(
		&sb,
		"%s\n",
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
		name := caser.String(strings.Split(peer.Name, ".")[0])
		if peer.Hostname == hostname {
			name += " (this)"
		}
		spacing := 22 - len(name)
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
		}
		if !peer.InactivityExpirationEnabled {
			keyIcon = nameStyle.Render(
				" ",
			) + nameStyle.Foreground(configs.NiceBlue).
				Render("\ueb11")
			spacing -= 2
		}
		sb.WriteString(nameStyle.Render(
			" " + name,
		))
		sb.WriteString(exitIcon)
		sb.WriteString(keyIcon)
		sb.WriteString(nameStyle.Render(strings.Repeat(
			" ",
			spacing,
		)))

		if peer.Connected || !alternatingText {
			sb.WriteString(nameStyle.Render(" " + peer.Version))
			sb.WriteString("\n")
		} else {
			diffText := common.GetTimeDifferenceString(peer.LastSeen)
			sb.WriteString(nameStyle.Render(fmt.Sprintf(" \uf017 %4s", diffText)))
			sb.WriteString("\n")
		}
	}

	offlineDiff := time.Until(keyExpiry)
	diffText := common.GetTimeDifferenceString(keyExpiry)
	sb.WriteString("\nNetbird key expiry: ")
	keytext := fmt.Sprintf("%s%6s", "\uf017 ", diffText)
	if offlineDiff.Hours() < (24 * 4) {
		sb.WriteString(pinkBold.Render(keytext))
	} else {
		sb.WriteString(keytext)
	}

	res := configs.TailscaleStyle.Render(strings.TrimSuffix(sb.String(), "\n"))
	return res
}
