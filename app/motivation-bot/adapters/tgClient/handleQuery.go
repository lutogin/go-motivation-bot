package tgClient

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/users/dto"
	"strings"
	"time"
)

type availableQueries struct {
	SetLang string `mapstructure:"setLang"`
}

func newAvailableQueries() availableQueries {
	return availableQueries{SetLang: "setLang"}
}

func saveLangForUser(t *TgClient, update *tgbotapi.Update, lang string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.userService.Create(ctx, usersDto.CreateUserDto{
		ChatId:    update.CallbackQuery.From.ID,
		FirstName: update.CallbackQuery.From.FirstName,
		LastName:  update.CallbackQuery.From.LastName,
		UserName:  update.CallbackQuery.From.UserName,
		Lang:      lang,
	})
}

func setLangHandler(t *TgClient, bot *tgbotapi.BotAPI, update *tgbotapi.Update, options string) {
	switch options {
	case "en":
		saveLangForUser(t, update, "en")

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Now, please send a time when you want to receive quotes. It should be 24h format time the next format '18:30'. Minutes must be a multiple of 30")
		_, err := bot.Send(msg)
		if err != nil {
			t.logger.Error(err)
		}

		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Bot is in test mode, so sending will be at the specified time but in UTC-0, please consider the offset in your timezone.")
		_, err = bot.Send(msg)
		if err != nil {
			t.logger.Error(err)
		}

	case "ru":
		saveLangForUser(t, update, "ru")

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Теперь, пожалуйста, укажите время, когда вы хотите получать цитаты. Это должно быть время в формате 24 часа, следующий формат '18:30' Минуты долджны быть кратные 30.")
		_, err := bot.Send(msg)
		if err != nil {
			t.logger.Error(err)
		}

		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Бот находится в тестовом режиме, поэтому отправка будет в указанное время, но по UTC-0, пожалуйста, учитывайте смещение по вашей таймзоне.")
		_, err = bot.Send(msg)
		if err != nil {
			t.logger.Error(err)
		}

	}
}

func (t *TgClient) HandleQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	queryArr := strings.Split(update.CallbackQuery.Data, ":")
	query := queryArr[0]
	option := queryArr[1]

	if query == "" || option == "" {
		t.SendMessage(update.Message.Chat.ID, "Wrong query.")
		return
	}

	availableQuery := newAvailableQueries()

	switch query {
	case availableQuery.SetLang:
		setLangHandler(t, bot, update, option)
	default:
		t.SendMessage(update.Message.Chat.ID, "Command not found.")
	}
}
