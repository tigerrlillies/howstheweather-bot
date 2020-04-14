package main

import (
	"fmt"

	"github.com/tigerrlillies/howstheweather/openweathermap"
)

// CreateCurrentWeatherReport creates weather report text, based on OpenWeatherMap CurrentWeatherData
func CreateCurrentWeatherReport(weather *openweathermap.Current) string {
	general := "Ожидается "
	for _, w := range weather.Weather[:len(weather.Weather)-1] {
		general += fmt.Sprintf("%s, ", w.Description)
	}
	general += fmt.Sprintf("%s. ", weather.Weather[len(weather.Weather)-1].Description)

	wind := fmt.Sprintf("Скорость ветра: %v метра в секунду. ", weather.WindSpeed)
	temperature := fmt.Sprintf("Температура: %.1f градусов по Цельсию, ощущается как %.1f. ", weather.Temp, weather.FeelsLike)

	humidity := fmt.Sprintf("Влажность воздуха: %d процентов. ", weather.Humidity)

	uv := fmt.Sprintf("UV индекс: %.1f. ", weather.Uvi)

	return general + wind + temperature + humidity + uv
}

// CreateForecastReport creates weather report text,
func CreateForecastReport(weather *openweathermap.Hourly) string {
	return ""
}