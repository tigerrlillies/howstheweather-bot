package main

import (
	"log"
	"os"

	"github.com/Syfaro/telegram-bot-api"

	"github.com/tigerrlillies/howstheweather/openweathermap"
)

const (
	TelegramAPITokenEnv       = "TELEGRAM_API_TOKEN"
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
		command := update.Message.Command()
		text := update.Message.Text
		chatID := update.Message.Chat.ID

		log.Printf("Update received: command [%s], text [%s], chat ID [%d]",
			command, text, chatID)

		if command != "" {
			HandleCommand(command, chatID, bot)
			continue
		}

		if text != "" {
			HandleText(text, chatID, owm, bot)
		}
	}

}
