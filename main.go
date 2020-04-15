package main

import (
	"log"
	"os"

	"github.com/Syfaro/telegram-bot-api"

	"github.com/tigerrlillies/howstheweather/openweathermap"
)

const (
	TelegramAPITokenEnv = "TELEGRAM_API_TOKEN"
	OpenWeatherMapAPITokenEnv = "OPENWEATHERMAP_API_TOKEN"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(TelegramAPITokenEnv))
	if err != nil {
		log.Fatalf("Unable to initialize the bot: %s", err.Error())
	}

	log.Printf("Authorized as %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChan, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatalf("Failed to initialize update channel: %s", err.Error())
	}

	owm := openweathermap.New(os.Getenv(OpenWeatherMapAPITokenEnv))

	log.Printf("Started listening for updates")
	for update := range updateChan {
		log.Printf("Update received: command %s, text %s, chat ID %d",
			update.Message.Command(), update.Message.Text, update.Message.Chat.ID)

		city := update.Message.Text

		var response string

		lat, lon, err := owm.GetCityCoordinates(city)
		if err != nil {
			log.Println("Failed to get city coordinates: ", err.Error())

			switch err {
			case openweathermap.ErrCityNotFound:
				response = "Такой город не найден"
			default:
				response = "Не получилось загрузить прогноз погоды. Попробуете еще раз?"
			}

		} else {
			weather, err := owm.OneCallByCoordinates(lat, lon)
			if err != nil {
				log.Printf("Failed to fetch weather data from OpenWeatherMap: %s", err.Error())
				response = "Не получилось загрузить прогноз погоды. Попробуете еще раз?"
			} else {
				response = CreateForecastReport(weather.Daily[0]) + "\n" + GetClothingRecommendations(weather.Current)
			}
		}

		_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response))
		if err != nil {
			log.Printf("Failed to send response: %s", err.Error())
		}
	}

}
