package models

import "time"

type WaterTemperatureInternal struct {
	Place       string
	Temperature float64
	LastUpdate  time.Time
}

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
