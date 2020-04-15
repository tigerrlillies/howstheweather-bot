package main

import (
	"log"

	"github.com/Syfaro/telegram-bot-api"

	"github.com/tigerrlillies/howstheweather/openweathermap"
)

const (
	StartCommand = "start"

	CityNotFoundErrorMessage        = "Такой город не найден. Попробуете еще раз?"
	CommandNotFoundErrorMessage     = "Такая команда мне пока неизвестна."
	UnableToLoadWeatherErrorMessage = "Не получилось загрузить прогноз погоды. Попробуете еще раз?"

	WelcomeReply = "Введите название любого города, и я расскажу вам прогноз погоды на сегодня."
)

func HandleCommand(command string, chatID int64, bot *tgbotapi.BotAPI) {
	switch command {
	case StartCommand:
		_, err := bot.Send(tgbotapi.NewMessage(chatID, WelcomeReply))
		if err != nil {
			log.Printf("Failed to send response: %s", err.Error())
		}
	default:
		_, err := bot.Send(tgbotapi.NewMessage(chatID, CommandNotFoundErrorMessage))
		if err != nil {
			log.Printf("Failed to send response: %s", err.Error())
		}
	}
}

func HandleText(text string, chatID int64, owm *openweathermap.OpenWeatherMap, bot *tgbotapi.BotAPI) {
	city := text

	var response string

	lat, lon, err := owm.GetCityCoordinates(city)
	if err != nil {
		log.Println("Failed to get city coordinates: ", err.Error())

		switch err {
		case openweathermap.ErrCityNotFound:
			response = CityNotFoundErrorMessage
		default:
			response = UnableToLoadWeatherErrorMessage
		}

	} else {
		weather, err := owm.OneCallByCoordinates(lat, lon)
		if err != nil {
			log.Printf("Failed to fetch weather data from OpenWeatherMap: %s", err.Error())
			response = UnableToLoadWeatherErrorMessage
		} else {
			response = CreateForecastReport(weather.Daily[0]) + "\n" + GetClothingRecommendations(weather.Daily[0])
		}
	}

	_, err = bot.Send(tgbotapi.NewMessage(chatID, response))
	if err != nil {
		log.Printf("Failed to send response: %s", err.Error())
	}
}
