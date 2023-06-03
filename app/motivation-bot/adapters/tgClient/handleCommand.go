package tgClient

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/localization"
	usersDto "motivation-bot/users/dto"
	"time"
)

func handleStartCommand(t *TgClient, update *tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	lang := update.Message.From.LanguageCode
	// Greeting
	t.SendMessage(chatId, localization.Tr("It's a simple bot for keep you more motivated. \nGlory to Ukraine!!!", lang))

	// Define a row of buttons
	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("EN", "setLang:en"),
		tgbotapi.NewInlineKeyboardButtonData("RU", "setLang:ru"),
	)

	// Define the keyboard layout
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	// Define the message with the keyboard
	msg := tgbotapi.NewMessage(chatId, localization.Tr("Please, choose a language of quotes:", lang))
	msg.ReplyMarkup = keyboard

	t.client.Send(msg)
}

func handleStopCommand(t *TgClient, update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.userService.DeleteByChatId(ctx, usersDto.DeleteUserByChatIdDto{
		ChatId: update.Message.From.ID,
	})
	t.SendMessage(update.Message.Chat.ID, localization.Tr("You've been unsubscribed from the notifications. \nHave a nice day!", update.Message.From.LanguageCode))
}

func (t *TgClient) HandleCommand(update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		handleStartCommand(t, update)
		break
	case "setup":
		handleStartCommand(t, update)
		break
	case "stop":
		handleStopCommand(t, update)
		break
	default:
		t.SendMessage(update.Message.Chat.ID, localization.Tr("I don't know that command.", update.Message.From.LanguageCode))
	}

}
