package tailscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"StatusApp/configs"
	"StatusApp/internal/common"
	"StatusApp/internal/truenas"
)

// View renders the Tailscale status information.
// Note: This function has dependencies on other parts of the model (like truenas).
// A future refactoring could be to pass only the necessary data to this function
// to make it more self-contained.
func View(m interface{}) string {
	// This is a temporary solution to handle the model dependency.
	// The ideal approach would be to define a clear model/viewmodel in main.go
	// and pass only the required data.
	type modelWithData interface {
		GetTruenasApps() []truenas.App
		GetTailscaleDevices() Devices
		GetKeyExpiry() time.Time
		DisplayAlternatingText() bool
	}

	appModel, ok := m.(modelWithData)
	if !ok {
		return "Error: Could not render Tailscale view due to model mismatch"
	}

	var sb strings.Builder
	greenBold := lipgloss.NewStyle().Bold(true).Foreground(configs.BrightGreen)

	pinkBold := lipgloss.NewStyle().Bold(true).Foreground(configs.HotPink)

	yes, no := truenas.GetAppStatus(appModel.GetTruenasApps())
	fmt.Fprintf(&sb, "% -15s", "Dodo Apps:")
	sb.WriteString(greenBold.Render("\uf00c"))
	fmt.Fprintf(&sb, " %d | ", yes)
	sb.WriteString(pinkBold.Render("\uf00d"))
	fmt.Fprintf(&sb, " %d\n\n", no)

	for i, device := range appModel.GetTailscaleDevices().Devices {
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
		sb.WriteString(nameStyle.Render(fmt.Sprintf(" % -20s", name)))

		if device.ConnectedToControl || !appModel.DisplayAlternatingText() {
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

			diffText := common.GetTimeDifferenceString(device.LastSeen)
			sb.WriteString(nameStyle.Render(fmt.Sprintf("\uf017 %6s", diffText)) + "\n")
		}
	}

	keyExpiry := appModel.GetKeyExpiry()
	offlineDiff := time.Until(keyExpiry)
	diffText := common.GetTimeDifferenceString(keyExpiry)
	sb.WriteString("\nTailscale key expiry: ")
	keytext := fmt.Sprintf("%s%6s", "\uf017 ", diffText)
	if offlineDiff.Hours() < (24 * 4) {
		sb.WriteString(pinkBold.Render(keytext))
	} else {
		sb.WriteString(keytext)
	}

	res := configs.TailscaleStyle.Render(strings.TrimSuffix(sb.String(), "\n"))
	return res
}
