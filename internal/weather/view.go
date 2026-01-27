package weather

import (
	"fmt"
	"strings"

	"StatusApp/configs"
	"StatusApp/internal/common"
)

// getWeatherIcon finds the correct weather icon based on the weather code and time of day.
// It reads from a JSON file specified by the WEATHER_ICON_PATH environment variable.

// View renders the weather information.
func View(
	weather *WeatherForecastInternal,
	waterTemp WaterTemperatureInternal,
	alternatingText bool,
) string {
	if weather == nil {
		return ""
	}
	var sb strings.Builder

	noStyle := configs.BoldText.UnsetBold().UnsetForeground()

	if !alternatingText {
		sb.WriteString(
			configs.BoldText.Render(fmt.Sprintf("%-2s%s ", weather.Icon, weather.Text)),
		)
		sb.WriteString(noStyle.Render(weather.Location))
		sb.WriteString(configs.BoldText.Render(fmt.Sprintf(" %0.1f°C ", weather.Temperature)))
		sb.WriteString(noStyle.Render(fmt.Sprintf("(%s)", weather.ExtraInfo)))
	} else {
		lastUpdateString := common.GetTimeDifferenceString(waterTemp.LastUpdate)
		sb.WriteString(configs.BoldText.Render(fmt.Sprintf("%-2s ", "\uef30")))
		sb.WriteString(noStyle.Render(waterTemp.Place))
		sb.WriteString(configs.BoldText.Render(fmt.Sprintf(" %2.1f°C", waterTemp.Temperature)))
		sb.WriteString(noStyle.Render(fmt.Sprintf(" (%s)", lastUpdateString)))
	}

	res := sb.String()
	return res
}
