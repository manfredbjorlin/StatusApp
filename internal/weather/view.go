package weather

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"StatusApp/internal/common"
)

// getWeatherIcon finds the correct weather icon based on the weather code and time of day.
// It reads from a JSON file specified by the WEATHER_ICON_PATH environment variable.
func getWeatherIcon(code int, isDay int, iconPath string) (string, error) {
	byteValue, err := os.ReadFile(iconPath)
	if err != nil {
		return "", fmt.Errorf("failed to read weather icon file: %w", err)
	}

	var weatherData []Weathercode
	if err := json.Unmarshal(byteValue, &weatherData); err != nil {
		return "", fmt.Errorf("failed to unmarshal weather icons: %w", err)
	}

	for _, weather := range weatherData {
		if weather.Code == code {
			if isDay == 1 {
				return weather.Day, nil
			}
			return weather.Night, nil
		}
	}
	return "", fmt.Errorf("no icon found for code %d", code)
}

// View renders the weather information.
func View(
	weather Weather,
	waterTemp WaterTemperatureInternal,
	alternatingText bool,
	iconPath string,
) string {
	var sb strings.Builder
	var style lipgloss.Style

	icon, err := getWeatherIcon(weather.Current.Condition.Code, weather.Current.IsDay, iconPath)
	if err != nil {
		icon = "?"
	}

	if alternatingText {
		sb.WriteString(
			style.Render(
				fmt.Sprintf(
					"%-2s %5.1f°C (%5.1f°C)",
					icon,
					weather.Current.Temp,
					weather.Current.FeelsLike,
				),
			),
		)
	} else {
		lastUpdateString := common.GetTimeDifferenceString(waterTemp.LastUpdate)
		sb.WriteString(style.Render(fmt.Sprintf(" %-2s %s %2.1f°C (%s)", "\uef30", waterTemp.Place, waterTemp.Temperature, lastUpdateString)))
	}

	res := sb.String()
	return res
}
