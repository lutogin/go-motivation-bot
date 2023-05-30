package tgClient

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	usersDto "motivation-bot/users/dto"
	"time"
)

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

func handleDeleteSubscriptionCommand(t *TgClient, update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.userService.DeleteByChatId(ctx, usersDto.DeleteUserByChatIdDto{
		ChatId: update.Message.From.ID,
	})
	t.SendMessage(update.Message.Chat.ID, "You've been unsubscribed from the notifications. \nHave a nice day!")
}

func (t *TgClient) HandleCommand(update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		handleStartCommand(t.client, update)
		break
	case "setup":
		handleStartCommand(t.client, update)
		break
	case "remove":
		handleDeleteSubscriptionCommand(t, update)
		break
	default:
		t.SendMessage(update.Message.Chat.ID, "I don't know that command.")
	}

}
