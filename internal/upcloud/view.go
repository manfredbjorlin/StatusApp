package upcloud

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"StatusApp/configs"
)

const (
	runningText = "started"
)

func Status(servers []Server) string {
	result := make([]string, 0)

	for _, server := range servers {
		var icon string
		if server.State == runningText {
			icon = configs.RunningIcon
		} else {
			icon = configs.StoppedIcon
		}
		result = append(result, fmt.Sprintf("%s %s", icon, server.Title))
	}

	return strings.Join(result, " - ")
}

func View(servers []Server, alternatingText bool, accountInfo AccountInfo) string {
	accountInfo.Currency = strings.ReplaceAll(accountInfo.Currency, "EUR", "€")
	name := "Upcloud"
	nameFormatted := configs.BoldText.Foreground(configs.BrightGreen).Render(name)

	infoText := fmt.Sprintf(
		"Remaining Credits: %s %0.2f - Monthly usage: %s %0.2f",
		accountInfo.Currency,
		accountInfo.RemainingCredits,
		accountInfo.Currency,
		accountInfo.BillingSummary,
	)
	infoTextFormatted := fmt.Sprintf(
		"%s %s %0.2f - %s %s %0.2f",
		configs.BoldText.Render(
			"Remaining Credits:",
		),
		accountInfo.Currency,
		accountInfo.RemainingCredits,
		configs.BoldText.Render("Montly usage:"),
		accountInfo.Currency,
		accountInfo.BillingSummary,
	)

	spacing := configs.ApplicationWidth - len([]rune(name)) - len([]rune(infoText)) - 1
	spacingString := strings.Repeat(" ", spacing)

	result := lipgloss.NewStyle().
		Width(configs.ApplicationWidth).
		AlignHorizontal(lipgloss.Center).
		Render(nameFormatted + spacingString + infoTextFormatted)

	result += "\n"

	tableRows := make([][]string, 0)
	for _, server := range servers {

		status := server.State
		if status == runningText {
			status = configs.RunningIcon + " Running"
		} else {
			status = configs.StoppedIcon + status
		}

		memory, _ := strconv.Atoi(server.MemoryAmount)
		memory /= 1024

		tableRows = append(tableRows, []string{
			server.Title,
			status,
			server.CoreNumber,
			fmt.Sprintf("%v Gb", memory),
			server.Plan,
			server.Zone,
		})
	}
	headers := []string{
		"Hostname",
		"State",
		"Cores",
		"Memory",
		"Plan",
		"Zone",
	}

	result += table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(configs.NiceBlue)).
		Headers(headers...).
		Rows(tableRows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return configs.BoldText
			}
			return lipgloss.NewStyle()
		}).
		Width(configs.ApplicationWidth).
		String()

	result += "\n"

	return result
}
