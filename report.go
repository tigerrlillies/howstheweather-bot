package main

import (
	"fmt"

	"github.com/tigerrlillies/howstheweather/openweathermap"
)

// CreateCurrentWeatherReport creates weather report text, based on OpenWeatherMap CurrentWeatherData
func CreateCurrentWeatherReport(weather *openweathermap.Current) string {
	general := "Сейчас: " + createOverallWeatherReport(weather.Weather)

	wind := fmt.Sprintf("Скорость ветра: %v метра в секунду. ", weather.WindSpeed)
	temperature := fmt.Sprintf("Температура: %.1f°C, ощущается как %.1f°C. ",
		weather.Temp, weather.FeelsLike)

	humidity := fmt.Sprintf("Влажность воздуха: %d%%. ", weather.Humidity)

	uv := fmt.Sprintf("UV индекс: %.1f. ", weather.Uvi)

	return general + wind + temperature + humidity + uv
}

// CreateForecastReport creates weather report text,
func CreateForecastReport(weather *openweathermap.Daily) string {
	general := "Прогноз погоды на сегодня: " + createOverallWeatherReport(weather.Weather)

	morning := fmt.Sprintf("Утром ожидается %.1f°C, ощущается как %.1f°C. ",
		weather.Temperature.Morning, weather.FeelsLike.Morn)
	day := fmt.Sprintf("Днем будет %.1f°C, по ощущениям как %.1f°C. ",
		weather.Temperature.Day, weather.FeelsLike.Day)
	evening := fmt.Sprintf("Вечером ожидается %.1f°C, по ощущениям как %.1f°C, ",
		weather.Temperature.Eve, weather.FeelsLike.Eve)
	night := fmt.Sprintf("и, наконец, ночью %.1f°C и %.1f°C соответственно. ",
		weather.Temperature.Night, weather.FeelsLike.Night)

	temperature := morning + day + evening + night

	wind := fmt.Sprintf("Скорость ветра: %v м/с. ", weather.WindSpeed)

	humidity := fmt.Sprintf("Влажность воздуха: %d%%. ", weather.Humidity)

	uv := fmt.Sprintf("UV индекс: %.1f. ", weather.Uvi)

	return general + temperature + wind + humidity + uv
}

func createOverallWeatherReport(weather []*openweathermap.Weather) string {
	report := ""
	for _, w := range weather[:len(weather)-1] {
		report += fmt.Sprintf("%s, ", w.Description)
	}
	report += fmt.Sprintf("%s. ", weather[len(weather)-1].Description)

	return report
}
