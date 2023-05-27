package tgClient

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/common"
	"motivation-bot/users/dto"
	"strconv"
	"strings"
	"time"
)

func (t *TgClient) HandleMessage(update *tgbotapi.Update) {
	incomeMsg := update.Message.Text
	senderChatId := update.Message.Chat.ID

	if !common.Is24TimeFormat(update.Message.Text) {
		t.SendMessage(senderChatId, "Wrong format. \nExample: '12:30', '13:00', '00:30'")
		return
	}

	incomeDateArr := strings.Split(incomeMsg, ":")
	hours, _ := strconv.Atoi(incomeDateArr[0])
	minutes, _ := strconv.Atoi(incomeDateArr[1])

	ctxGetUser, cancelGetUser := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelGetUser()
	user, err := t.userService.GetByChatId(ctxGetUser, usersDto.GetUserByChatIdDto{ChatId: senderChatId})
	fmt.Println(user)

	now := time.Now()
	nowInUserTZ := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		hours,
		minutes,
		0,
		0,
		now.Location(),
	).Add(time.Duration(user.TimeZone*-1) * time.Hour)

	ctxUpd, cancelUpd := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelUpd()
	err = t.userService.Update(ctxUpd, usersDto.UpdateUserDto{UserName: update.Message.From.UserName, AlertingTime: usersDto.AlertingTime{Hours: nowInUserTZ.Hour(), Minutes: nowInUserTZ.Minute()}})
	if err != nil {
		t.SendMessage(senderChatId, "Seem like you are not registered. True to user /start command.")
		return
	}

	if time.Now().After(nowInUserTZ) {
		t.SendMessage(senderChatId, "Bot set this time for you.")
		t.SendMessage(senderChatId, "There is one quote for you for today:")

		quote := t.forismatic.GetQuote("ru")
		t.SendQuote(senderChatId, quote)
	} else {
		t.SendMessage(senderChatId, "Bot set this time for you.")
		t.SendMessage(senderChatId, fmt.Sprintf("You will receive a quote later today, at %s.", incomeMsg))
	}
}
