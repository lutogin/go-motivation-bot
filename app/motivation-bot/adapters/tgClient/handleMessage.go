package tgClient

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/common"
	"motivation-bot/integrations/forismatic"
	"motivation-bot/users/dto"
	"strconv"
	"strings"
	"time"
)

func (t *TgClient) HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, client *forismatic.Client) {
	incomeMsg := update.Message.Text
	senderChatId := update.Message.Chat.ID

	if !common.Is24TimeFormat(update.Message.Text) {
		msg := tgbotapi.NewMessage(senderChatId, "Wrong format. \nExample: '12:30', '13:00', '00:30'")
		_, err := bot.Send(msg)
		if err != nil {
			t.logger.Error(err)
		}

		return
	}

	incomeDateArr := strings.Split(incomeMsg, ":")
	hours, _ := strconv.Atoi(incomeDateArr[0])
	minutes, _ := strconv.Atoi(incomeDateArr[1])

	now := time.Now()
	setTime := time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0, now.Location())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := t.userService.Update(ctx, usersDto.UpdateUserDto{UserName: update.Message.From.UserName, AlertingTime: usersDto.AlertingTime{Hours: hours, Minutes: minutes}})
	if err != nil {
		t.SendMessage(senderChatId, "Seem like you are not registered. True to user /start command.")
		return
	}

	if time.Now().After(setTime) {
		t.SendMessage(senderChatId, "Bot set this time for you.")
		t.SendMessage(senderChatId, "There is one quote for you for today:")

		quote := client.GetQuote("ru")
		t.SendQuote(senderChatId, quote)
	} else {
		t.SendMessage(senderChatId, "Bot set this time for you.")
		t.SendMessage(senderChatId, fmt.Sprintf("You will receive a quote later today, at %s.", incomeMsg))
	}
}
