package models

type Weathercode struct {
	Code  int    `json:"code"`
	Day   string `json:"day"`
	Night string `json:"night"`
	Icon  int    `json:"icon"`
}

type Weather struct {
	Current Current `json:"current"`
}

type Current struct {
	Condition Condition `json:"condition"`
	Temp      float32   `json:"temp_c"`
	FeelsLike float32   `json:"feelslike_c"`
	IsDay     int       `json:"is_day"`
}

type Condition struct {
	Code int `json:"code"`
}
