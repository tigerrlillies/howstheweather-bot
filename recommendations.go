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

func GetClothingRecommendations(weather *openweathermap.Current) string {
	categories := categorize(weather)

	recommendation := "Рекомендуем "
	for _, c := range categories[:len(categories)-1] {
		recommendation += recommend(c) + ", "
	}
	recommendation += recommend(categories[len(categories)-1]) + "."

	return recommendation
}

func categorize(weather *openweathermap.Current) []WeatherCategory {
	categories := make([]WeatherCategory, 0)

	t := weather.Temp
	if t < -20. {
		categories = append(categories, Freezing)
	} else if t < 0. {
		categories = append(categories, Cold)
	} else if t < 15. {
		categories = append(categories, Cool)
	} else if t < 25. {
		categories = append(categories, Warm)
	} else {
		categories = append(categories, Hot)
	}

	windSpeed := weather.WindSpeed
	if windSpeed > 10. {
		categories = append(categories, Windy)
	}

	if weather.Rain != nil && weather.Rain.OneHour > 0. {
		categories = append(categories, Rainy)
	}

	return categories
}

func recommend(c WeatherCategory) string {
	switch c {
	case Hot:
		return "одеться полегче - сейчас довольно жарко"
	case Warm:
		return "одеться как захочется, потому что будет тепло и приятно :)"
	case Cool:
		return "одеться потеплее, потому что будет прохладно"
	case Cold:
		return "одеться тепло"
	case Freezing:
		return "одеться максимально тепло, а еще лучше - никуда не выходить"
	case Rainy:
		return "не забыть зонтик"
	case Windy:
		return "взять ветровку или шарфик, потому что ветрено"
	default:
		return ""
	}
}
