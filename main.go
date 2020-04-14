package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"

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

		coords := strings.Split(update.Message.Text, " ")
		lat, err := strconv.ParseFloat(coords[0], 32)
		if err != nil {
			log.Println("Unable to parse latitude due to", err.Error())
			continue
		}
		lon, err := strconv.ParseFloat(coords[1], 32)
		if err != nil {
			log.Println("Unable to parse longitude due to", err.Error())
			continue
		}

		var responseMessage string
		weather, err := owm.OneCallByCoordinates(lat, lon)
		if err != nil {
			log.Printf("Failed to fetch weather data from OpenWeatherMap: %s", err.Error())
			responseMessage = "Не получилось загрузить прогноз погоды. Возможно, город введен некорректно. Попробуете еще раз?"
		} else {
			responseMessage = CreateCurrentWeatherReport(weather.Current)
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID, responseMessage)

		_, err = bot.Send(response)
		if err != nil {
			log.Printf("Failed to send response: %s", err.Error())
		}
	}

}
