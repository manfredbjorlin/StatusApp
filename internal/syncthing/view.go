package syncthing

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/mbndr/figlet4go"

	"StatusApp/configs"
)

func Status(connections []Connection) string {
	total := len(connections)
	connected := total
	for _, conn := range connections {
		if !conn.Connected {
			connected--
		}
	}

	return fmt.Sprintf(
		"Syncthing: %s %d | %s %d",
		configs.OkIcon,
		connected,
		configs.FailIcon,
		total-connected,
	)
}

func View(connections []Connection) string {
	ascii := figlet4go.NewAsciiRender()
	header, _ := ascii.Render("Syncthing")
	result := lipgloss.NewStyle().Foreground(configs.HotPink).Render(header)

	result += "\n"

	tableRows := make([][]string, 0)
	for _, conn := range connections {
		connected := configs.FailIcon + " Not connected"
		lastSync := ""
		if conn.Connected {
			connected = configs.OkIcon + " Connected"
			lastSync = conn.LastSync.Format("15:04:05")
		}
		row := []string{
			conn.Device.Name,
			connected,
			lastSync,
			conn.ClientVersion,
			fmt.Sprint(conn.InBytesTotal),
			fmt.Sprint(conn.OutBytesTotal),
		}
		tableRows = append(tableRows, row)
	}
	headers := []string{
		"Name",
		"Status",
		"Last Sync",
		"Version",
		"Bytes in",
		"Bytes Out",
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

	return result + "\n"
}
