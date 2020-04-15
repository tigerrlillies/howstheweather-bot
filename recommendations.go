package main

import "github.com/tigerrlillies/howstheweather/openweathermap"

type WeatherCategory string

const (
	// temperature
	Hot      WeatherCategory = "hot"
	Warm     WeatherCategory = "warm"
	Cool     WeatherCategory = "cool"
	Cold     WeatherCategory = "cold"
	Freezing WeatherCategory = "freezing"

	// precipitation
	Rainy WeatherCategory = "rainy"

	// wind
	Windy WeatherCategory = "windy"
)

func GetClothingRecommendations(weather *openweathermap.Daily) string {
	categories := categorize(weather)

	recommendation := "Рекомендуем "
	for _, c := range categories[:len(categories)-1] {
		recommendation += recommend(c) + ", "
	}
	recommendation += "а еще " + recommend(categories[len(categories)-1]) + "."

	return recommendation
}

func categorize(weather *openweathermap.Daily) []WeatherCategory {
	categories := make([]WeatherCategory, 0)

	averageTemp := weather.Temperature.Min + (weather.Temperature.Max-weather.Temperature.Min)/2
	if averageTemp < -20 {
		categories = append(categories, Freezing)
	} else if averageTemp < 0 {
		categories = append(categories, Cold)
	} else if averageTemp < 15 {
		categories = append(categories, Cool)
	} else if averageTemp < 25 {
		categories = append(categories, Warm)
	} else {
		categories = append(categories, Hot)
	}

	windSpeed := weather.WindSpeed
	if windSpeed > 5 {
		categories = append(categories, Windy)
	}

	if weather.Rain > 0 {
		categories = append(categories, Rainy)
	}

	return categories
}

func recommend(c WeatherCategory) string {
	switch c {
	case Hot:
		return "одеться полегче — сейчас довольно жарко"
	case Warm:
		return "одеться как захочется, потому что температура будет комфортной"
	case Cool:
		return "одеться потеплее, потому что будет прохладно"
	case Cold:
		return "одеться тепло"
	case Freezing:
		return "одеться максимально тепло, а еще лучше — никуда не выходить"
	case Rainy:
		return "не забыть зонтик"
	case Windy:
		return "взять ветровку или шарфик, потому что ветрено"
	default:
		return ""
	}
}
