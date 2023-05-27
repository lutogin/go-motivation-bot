package tgClient

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func handleStartCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Define a row of buttons
	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("EN", "setLang:en"),
		tgbotapi.NewInlineKeyboardButtonData("RU", "setLang:ru"),
	)

	// Define the keyboard layout
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	// Define the message with the keyboard
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please choose a language of quotes:")
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}

func (t *TgClient) HandleCommand(update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		handleStartCommand(t.client, update)
		break
	case "changeSetup":
		handleStartCommand(t.client, update)
		break
	default:
		t.SendMessage(update.Message.Chat.ID, "I don't know that command.")
	}

}
