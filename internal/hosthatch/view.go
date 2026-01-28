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

func View(servers []Server, alternatingText bool) string {
	spacing := (configs.ScheduleStyle.GetWidth() - 9) / 2
	result := configs.BoldText.Foreground(configs.BrightGreen).
		Render(strings.Repeat(" ", spacing) + "HostHatch")
	result += "\n"

	ok := "\uf04b"

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
		if status == "Active" {
			status = lipgloss.NewStyle().Foreground(configs.BrightGreen).Render(ok) + " Running"
		} else {
			status = lipgloss.NewStyle().Foreground(configs.HotPink).Render(status)
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
