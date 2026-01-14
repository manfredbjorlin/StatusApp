package renderers

import (
	"encoding/json"
	"log"
	"os"

	"StatusApp/internal/models"
)

func getWeatherIcon(code int, isDay int) string {
	filePath := os.Getenv("WEATHER_ICON_PATH")
	byteValue, _ := os.ReadFile(filePath)

	var weatherData []models.Weathercode

	err := json.Unmarshal(byteValue, &weatherData)
	if err != nil {
		log.Fatal(err)
	}

	for _, weather := range weatherData {
		if weather.Code == code {
			if isDay == 1 {
				return weather.Day
			}
			return weather.Night
		}
	}
	return ""
}
