package weather

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type mockModel struct {
	weather   Weather
	waterTemp WaterTemperatureInternal
}

func (m mockModel) GetWeather() Weather {
	return m.weather
}

func (m mockModel) GetWaterTemperature() WaterTemperatureInternal {
	return m.waterTemp
}

func TestView(t *testing.T) {
	// Create a temporary file for weather icons
	icons := []Weathercode{
		{Code: 1000, Day: "☀️", Night: "🌙"},
	}
	tmpfile, err := os.CreateTemp("", "weather-icons-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	content, _ := json.Marshal(icons)
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	model := mockModel{
		weather: Weather{
			Current: Current{
				Condition: Condition{Code: 1000},
				IsDay:     1,
				Temp:      15.5,
			},
		},
		waterTemp: WaterTemperatureInternal{
			Temperature: 12.1,
		},
	}

	output := View(model, tmpfile.Name())

	if !strings.Contains(output, "☀️") {
		t.Error("expected day icon")
	}
	if !strings.Contains(output, "15.5°C") {
		t.Error("expected air temperature")
	}
	if !strings.Contains(output, "12.1°C") {
		t.Error("expected water temperature")
	}
}

func TestGetWeatherIcon_fileNotFound(t *testing.T) {
	_, err := getWeatherIcon(1000, 1, "non-existent-file.json")
	if err == nil {
		t.Fatal("expected an error for a missing icon file, but got none")
	}
}
