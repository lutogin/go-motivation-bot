package tgClient

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (t *TgClient) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		// Define a row of buttons
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("EN", "setLang:en"),
			tgbotapi.NewInlineKeyboardButtonData("RU", "setLang:ru"),
		)

		// Define the keyboard layout
		keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

		// Define the message with the keyboard
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a language:")
		msg.ReplyMarkup = keyboard

		bot.Send(msg)

		break
	case "setup":
		// Define a row of buttons
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("EN", "setLang:en"),
			tgbotapi.NewInlineKeyboardButtonData("RU", "setLang:ru"),
		)

		// Define the keyboard layout
		keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

		// Define the message with the keyboard
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a language:")
		msg.ReplyMarkup = keyboard

		bot.Send(msg)

		break
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
		bot.Send(msg)
	}

}
