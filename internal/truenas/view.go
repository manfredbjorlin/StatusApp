package truenas

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"StatusApp/configs"
)

func Status(apps []App) string {
	upToDate, toUpdate := GetAppStatus(apps)
	return fmt.Sprintf("Dodo: %s %d | %s %d", configs.OkIcon, upToDate, configs.FailIcon, toUpdate)
}

func View(apps []App) string {
	greenBold := lipgloss.NewStyle().Bold(true).Foreground(configs.BrightGreen)
	pinkBold := lipgloss.NewStyle().Bold(true).Foreground(configs.HotPink)
	ok := "\uf04b"
	var sb strings.Builder

	sb.WriteString(configs.BoldText.Render("\uf0c7 Apps with updates"))
	sb.WriteString("\n\n")

	var updateApps []string
	var longestString int
	for _, app := range apps {
		if app.UpgradeAvailable {
			appName := formatAppName(app.Name)
			if len(appName) > longestString {
				longestString = len(appName)
			}
		}
	}
	for _, app := range apps {
		if app.UpgradeAvailable {
			appName := formatAppName(app.Name)
			line := appName
			line += strings.Repeat(" ", longestString-len(appName)+2)
			if app.State == "RUNNING" {
				line += fmt.Sprintf("(%v) ", app.UpdateId)
				line += greenBold.Render(ok)

			} else {
				line += pinkBold.Render("Updating")
			}
			// var lineStyle lipgloss.Style
			// lineStyle = configs.SetBg(lineStyle)
			updateApps = append(updateApps, line)
		}
	}

	slices.Sort(updateApps)
	enumerator := func(_ list.Items, _ int) string {
		return pinkBold.Render("\uf00d")
	}

	toDisplay := updateApps
	if len(updateApps) > 9 {
		toDisplay = updateApps[:9]
	}
	mylist := list.New(toDisplay).
		Enumerator(enumerator).
		ItemStyleFunc(func(items list.Items, i int) lipgloss.Style {
			var lineStyle lipgloss.Style
			lineStyle = configs.SetBg(lineStyle, i)
			return lineStyle
		})
	sb.WriteString(mylist.String())

	return configs.TailscaleStyle.Render(sb.String())
}

func formatAppName(name string) string {
	titleCaser := cases.Title(language.English)
	output := titleCaser.String(strings.ReplaceAll(name, "-", " "))
	return output
}
