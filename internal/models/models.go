package models

import "time"

type WeatherResponse struct {
	Location    string    `json:"location"`
	Temperature float64   `json:"temperature"`
	UVIndex     float64   `json:"uvindex"`
	Conditions  string    `json:"conditions"`
	Time        time.Time `json:"date"`
}

type WeatherData struct {
	Location string `json:"resolvedAddress"`
	Days     []struct {
		Temperature float64 `json:"temp"`
		UVIndex     float64 `json:"uvindex"`
		Conditions  string  `json:"conditions"`
	} `json:"days"`
}
