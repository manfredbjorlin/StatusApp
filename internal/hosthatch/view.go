package hosthatch

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"StatusApp/configs"
)

const (
	runningText = "Active"
)

func Status(servers []Server) string {
	result := make([]string, 0)

	caser := cases.Title(language.BrazilianPortuguese)

	for _, server := range servers {
		var icon string
		if server.State == runningText {
			icon = configs.RunningIcon
		} else {
			icon = configs.StoppedIcon
		}
		result = append(result, fmt.Sprintf("%s %s", icon, caser.String(server.Hostname)))
	}

	return strings.Join(result, " - ")
}

func View(servers []Server, alternatingText bool) string {
	result := configs.BoldText.Foreground(configs.BrightGreen).
		Width(configs.ScheduleStyle.GetWidth()).
		AlignHorizontal(lipgloss.Left).
		Render(" HostHatch")
	result += "\n"

	tableRows := make([][]string, 0)
	for _, server := range servers {

		caser := cases.Title(language.BrazilianPortuguese)
		name := caser.String(server.Hostname)
		billing := fmt.Sprintf(
			"%s %v.00/%s",
			strings.ReplaceAll(server.Billing.Currency, "USD", "$"),
			server.Billing.RecurringCost,
			server.Billing.BillingCycle,
		)
		if alternatingText {
			billing = server.Billing.NextDue + strings.Repeat(
				" ",
				len(billing)-len(server.Billing.NextDue),
			)
		}

		status := server.State
		if status == runningText {
			status = configs.RunningIcon + " Running"
		} else {
			status = configs.StoppedIcon + status
		}
		tableRows = append(tableRows, []string{
			name,
			status,
			server.Product.Location,
			server.Product.Name,
			server.Product.Image,
			billing,
		})
	}
	headers := []string{
		"Hostname",
		"State",
		"Location",
		"Product",
		"Image",
		"Billing",
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
		Width(configs.ScheduleStyle.GetWidth()).
		String()

	result += "\n"
	return result
}
