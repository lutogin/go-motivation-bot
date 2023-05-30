package tgClient

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/users/dto"
	"strconv"
	"strings"
	"time"
)

type AvailableQueries struct {
	SetLang string `mapstructure:"setLang"`
	SetGMT  string `mapstructure:"setGMT"`
}

func newAvailableQueries() AvailableQueries {
	return AvailableQueries{SetLang: "setLang", SetGMT: "setGMT"}
}

func saveLangForUser(t *TgClient, update *tgbotapi.Update, lang string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.userService.Upsert(ctx, usersDto.UpdateUserDto{
		ChatId:    update.CallbackQuery.From.ID,
		FirstName: update.CallbackQuery.From.FirstName,
		LastName:  update.CallbackQuery.From.LastName,
		UserName:  update.CallbackQuery.From.UserName,
		Lang:      lang,
	})
}

func setGMT(t *TgClient, update *tgbotapi.Update, options string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gmtInt, _ := strconv.Atoi(options)

	t.userService.Update(ctx, usersDto.UpdateUserDto{
		ChatId:    update.CallbackQuery.From.ID,
		FirstName: update.CallbackQuery.From.FirstName,
		LastName:  update.CallbackQuery.From.LastName,
		UserName:  update.CallbackQuery.From.UserName,
		TimeZone:  gmtInt,
	})

	t.SendMessage(update.CallbackQuery.Message.Chat.ID, "Please, send a time when you want to receive quotes. \nIt should be in 24h format, as example '18:30'.\nMinutes must be a multiple of 30 (16:30, 17:00, 17:30 etc.)")
}

func sendMsgGMTElection(bot *tgbotapi.BotAPI, chatId int64) {
	var timezones = []string{"GMT-12", "GMT-11", "GMT-10", "GMT-9", "GMT-8",
		"GMT-7", "GMT-6", "GMT-5", "GMT-4", "GMT-3", "GMT-2", "GMT-1", "GMT-0",
		"GMT+1", "GMT+2", "GMT+3", "GMT+4", "GMT+5", "GMT+6", "GMT+7", "GMT+8",
		"GMT+9", "GMT+10", "GMT+11", "GMT+12"}

	msg := tgbotapi.NewMessage(chatId, "Please choose your timezone:")

	var row []tgbotapi.InlineKeyboardButton
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, tz := range timezones {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(tz, fmt.Sprintf("setGMT:%s", strings.Replace(tz, "GMT", "", 1))))
		if len(row) == 3 {
			keyboard = append(keyboard, row)
			row = nil
		}
	}
	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	bot.Send(msg)
}

func setLangHandler(t *TgClient, update *tgbotapi.Update, options string) {
	saveLangForUser(t, update, options)
	sendMsgGMTElection(t.client, update.CallbackQuery.Message.Chat.ID)
}

func (t *TgClient) HandleQuery(update *tgbotapi.Update) {
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
		setLangHandler(t, update, option)
	case availableQuery.SetGMT:
		setGMT(t, update, option)
	default:
		t.SendMessage(update.Message.Chat.ID, "Command not found.")
	}
}
