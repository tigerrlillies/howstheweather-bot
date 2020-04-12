package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/briandowns/openweathermap"
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

	weather, err := openweathermap.NewCurrent("C", "ru", os.Getenv(OpenWeatherMapAPITokenEnv))
	if err != nil {
		log.Fatalf("Unable to connect to OpenWeatherAPI: %s", err.Error())
	}

	log.Printf("Started listening for updates")
	for update := range updateChan {
		log.Printf("Update received: command %s, text %s, chat ID %d",
			update.Message.Command(), update.Message.Text, update.Message.Chat.ID)

		var responseMessage string
		err := weather.CurrentByName(update.Message.Text)
		if err != nil {
			log.Printf("Failed to fetch weather data from OpenWeatherMap: %s", err.Error())
			responseMessage = "Не получилось загрузить прогноз погоды. Возможно, город введен некорректно. Попробуете еще раз?"
		} else {
			responseMessage = fmt.Sprintf("Ожидается такая вот погодка: %s", weather.Weather[0].Description)
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID, responseMessage)

		_, err = bot.Send(response)
		if err != nil {
			log.Printf("Failed to send response: %s", err.Error())
		}
	}

}
