package tgClient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleTooLongMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Message to long. Please write to @lutogin for fix that. 👺")
	bot.Send(msg)
}

func handleCommands(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		msg.Text = "👋 Welcome! To receive messages with a customized font, just send me the text you want to customizing. Supports only LATIN symbols 🇺🇸.\n🇺🇦Glory to Ukraine!"
	case "help":
		msg.Text = "I am a simple bot for customization text. Just try to send me any message. Supports only LATIN symbols 🇺🇸."
	default:
		msg.Text = "👨‍💻 I don't know that command. Just write me a message containing what you do want to transform."
	}

	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}
