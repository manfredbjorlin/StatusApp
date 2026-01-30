package exaroton

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"StatusApp/configs"
)

var serverStatus = map[int]string{
	0:  "Offline",
	1:  "Online",
	2:  "Starting",
	3:  "Stopping",
	4:  "Restarting",
	5:  "Saving",
	6:  "Loading",
	7:  "Crashed",
	8:  "Pending",
	9:  "Transferring",
	10: "Preparing",
}

func Status(servers []Server) string {
	result := make([]string, 0)

	excludedServers := make([]string, 0)
	serversToExclude := os.Getenv("EXAROTON_EXCLUDE_SERVERS")
	if len(serversToExclude) > 0 {
		excludedServers = strings.Split(serversToExclude, ",")
	}
	for _, server := range servers {
		if len(excludedServers) > 0 && slices.Contains(excludedServers, server.Name) {
			continue
		}
		var icon string
		name := server.Name
		if server.Status == 1 {
			icon = configs.RunningIcon
			name += fmt.Sprintf(" (%d)", server.Players.Count)
		} else {
			icon = configs.StoppedIcon
		}
		result = append(result, fmt.Sprintf("%s %s", icon, name))
	}

	return strings.Join(result, " - ")
}

func View(servers []Server, alternatingText bool, remainingCredit float64) string {
	var result strings.Builder

	name := "Exaroton"
	nameFormatted := configs.BoldText.Foreground(configs.BrightGreen).Render(name)

	infoText := fmt.Sprintf("Remaining Credits: %0.2f", remainingCredit)
	infoTextFormatted := fmt.Sprintf(
		"%s %0.2f",
		configs.BoldText.Render("Remaining Credits: "),
		remainingCredit,
	)
	spacing := configs.ScheduleStyle.GetWidth() - len([]rune(name)) - len([]rune(infoText)) - 4
	spacingString := strings.Repeat(" ", spacing)

	result.WriteString(lipgloss.NewStyle().
		Width(configs.ScheduleStyle.GetWidth()).
		AlignHorizontal(lipgloss.Center).
		Render(nameFormatted + spacingString + infoTextFormatted))

	result.WriteString("\n")

	tableRows := make([][]string, 0)
	for _, server := range servers {

		status := serverStatus[server.Status]
		if server.Status == 1 {
			status = configs.RunningIcon + " " + status
		} else {
			status = configs.StoppedIcon + " " + status
		}

		players := "<None>"
		if server.Players.Count > 0 {
			players = strings.Join(server.Players.List, ", ")
		}

		tableRows = append(tableRows, []string{
			server.Name,
			status,
			server.Address,
			server.Software.Name,
			server.Software.Version,
			players,
		})
	}
	headers := []string{
		"Name",
		"Status",
		"Address",
		"Software",
		"Version",
		"Players",
	}

	result.WriteString(table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(configs.NiceBlue)).
		BorderRow(true).
		Headers(headers...).
		Rows(tableRows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return configs.BoldText
			}
			return lipgloss.NewStyle()
		}).
		Width(configs.ScheduleStyle.GetWidth()).
		String())

	result.WriteString("\n")

	return result.String()
}
