package weather

import "time"

// --- WeatherAPI Structures ---

// Weathercode maps a condition code to its description and icon.
type Weathercode struct {
	Code  int    `json:"code"`
	Day   string `json:"day"`
	Night string `json:"night"`
	Icon  int    `json:"icon"`
}

// Weather represents the top-level response from the WeatherAPI.
type Weather struct {
	Current Current `json:"current"`
}

// Current holds the current weather conditions.
type Current struct {
	Condition Condition `json:"condition"`
	Temp      float32   `json:"temp_c"`
	FeelsLike float32   `json:"feelslike_c"`
	IsDay     int       `json:"is_day"`
}

// Condition holds the specific weather condition code.
type Condition struct {
	Code int `json:"code"`
}

// --- Water Temperature (YR.no) Structures ---

// WaterTemperatureInternal is a simplified struct for application use.
type WaterTemperatureInternal struct {
	Place       string
	Temperature float64
	LastUpdate  time.Time
}

// WaterTemperature represents the full, complex response from the YR.no API.
type WaterTemperature struct {
	Embedded struct {
		NearestLocations []struct {
			Location struct {
				Category struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"category"`
				ID       string `json:"id"`
				Name     string `json:"name"`
				Position struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"position"`
				Elevation    int `json:"elevation"`
				CoastalPoint struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"coastalPoint"`
				TimeZone string `json:"timeZone"`
				URLPath  string `json:"urlPath"`
				Country  struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"country"`
				Region struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"region"`
				Subregion struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"subregion"`
				IsInOcean bool `json:"isInOcean"`
			} `json:"location"`
			ID                   int       `json:"id"`
			Temperature          float64   `json:"temperature"`
			Time                 time.Time `json:"time"`
			Source               int       `json:"source"`
			DistanceFromLocation int       `json:"distanceFromLocation"`
			SourceDisplayName    string    `json:"sourceDisplayName,omitempty"`
		} `json:"nearestLocations"`
	} `json:"_embedded"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Parent struct {
			Href string `json:"href"`
		} `json:"parent"`
	} `json:"_links"`
}
