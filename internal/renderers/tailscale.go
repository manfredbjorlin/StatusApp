package renderers

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"StatusApp/configs"
	"StatusApp/internal/models"
	"StatusApp/internal/truenas"
)

func RenderTailscale(m models.Model) string {
	var sb strings.Builder
	greenBold := lipgloss.NewStyle().Bold(true).Foreground(configs.BrightGreen)
	pinkBold := lipgloss.NewStyle().Bold(true).Foreground(configs.HotPink)
	// sb.WriteString(
	// 	configs.BoldText.Render("Tailscale") + "\n\n",
	// )

	yes, no := truenas.GetAppStatus(m.TruenasApps)
	sb.WriteString(fmt.Sprintf("%-15s", "Dodo Apps:"))
	sb.WriteString(greenBold.Render("\uf00c"))
	sb.WriteString(fmt.Sprintf(" %d | ", yes))
	sb.WriteString(pinkBold.Render("\uf00d"))
	sb.WriteString(fmt.Sprintf(" %d\n\n", no))

	configs.TailscaleRenders++
	if configs.TailscaleRenders >= 5 {
		configs.TailscaleRenders = 0
		configs.TailscaleVersion = !configs.TailscaleVersion
	}
	for i, device := range m.Devices.Devices {
		deviceIcon := ""
		switch device.Os {
		case "linux":
			deviceIcon = "\uf17c"
		case "android":
			deviceIcon = "\ue70e"
		case "windows":
			deviceIcon = "\uf17a"
		}
		logoStyle := pinkBold
		if device.ConnectedToControl {
			logoStyle = greenBold
		}
		logoStyle = configs.SetBg(logoStyle, i)
		sb.WriteString(logoStyle.Render(deviceIcon))

		caser := cases.Title(language.BrazilianPortuguese)
		name := caser.String(strings.Split(device.Name, ".")[0])
		nameStyle := lipgloss.NewStyle()
		nameStyle = configs.SetBg(nameStyle, i)
		sb.WriteString(nameStyle.Render(fmt.Sprintf(" %-20s", name)))

		if device.ConnectedToControl || (!device.ConnectedToControl && configs.TailscaleVersion) {
			updateStyle := greenBold
			updateLogo := "\uf00c"

			if device.UpdateAvailable {
				updateLogo = "\uf00d"
				updateStyle = pinkBold
			}

			updateStyle = configs.SetBg(updateStyle, i)
			sb.WriteString(updateStyle.Render(updateLogo))

			shortVersion := strings.Split(device.ClientVersion, "-")[0]
			sb.WriteString(nameStyle.Render(" "+shortVersion) + "\n")
		} else {

			offlineDiff := time.Since(device.LastSeen)
			diffText := ""
			if offlineDiff.Hours() >= 24 {
				days := int(offlineDiff.Hours()) / 24
				diffText = fmt.Sprintf("%4d d", days)
			} else if offlineDiff.Hours() >= 1 {
				hours := int(offlineDiff.Hours())
				diffText = fmt.Sprintf("%4d H", hours)
			} else {
				minutes := int(offlineDiff.Minutes())
				diffText = fmt.Sprintf("%4d m", minutes)
			}
			sb.WriteString(nameStyle.Render(fmt.Sprintf("\uf017 %s", diffText)) + "\n")
		}
	}

	offlineDiff := time.Until(m.KeyExpiry)
	diffText := ""
	if offlineDiff.Hours() >= 24 {
		days := int(offlineDiff.Hours()) / 24
		diffText = fmt.Sprintf("%4d d", days)
	} else if offlineDiff.Hours() >= 1 {
		hours := int(offlineDiff.Hours())
		diffText = fmt.Sprintf("%4d H", hours)
	} else {
		minutes := int(offlineDiff.Minutes())
		diffText = fmt.Sprintf("%4d m", minutes)
	}
	sb.WriteString("\nTailscale key expiry: ")
	keytext := fmt.Sprintf("%0s", "\uf017 "+diffText)
	if offlineDiff.Hours() < (24 * 4) {
		sb.WriteString(pinkBold.Render(keytext))
	} else {
		sb.WriteString(keytext)
	}

	res := configs.TailscaleStyle.Render(strings.TrimSuffix(sb.String(), "\n"))
	return res
}
