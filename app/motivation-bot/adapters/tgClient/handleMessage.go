package tgClient

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/common"
	"motivation-bot/localization"
	"motivation-bot/users/dto"
	"strconv"
	"strings"
	"time"
)

func (t *TgClient) HandleMessage(update *tgbotapi.Update) {
	incomeMsg := update.Message.Text
	senderChatId := update.Message.Chat.ID
	lang := update.Message.From.LanguageCode

	if !common.Is24TimeFormat(update.Message.Text) {
		t.SendMessage(senderChatId, localization.Tr("Wrong format. \nExample: '12:30', '13:00', '00:30'", lang))
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
	err = t.userService.Update(ctxUpd, usersDto.UpdateUserDto{UserName: update.Message.From.UserName, AlertingTime: common.GetTimeFromDate(nowInUserTZ)})
	if err != nil {
		t.SendMessage(senderChatId, localization.Tr("Seem like you are not registered. Try to use /start command.", lang))
		return
	}

	if time.Now().After(nowInUserTZ) {
		t.SendMessage(senderChatId, localization.Tr("Bot set this time for you.", lang))
		t.SendMessage(senderChatId, localization.Tr("There is one quote for you for today:", lang))

		quote := t.forismatic.GetQuote(user.Lang)
		t.SendQuote(senderChatId, quote)
	} else {
		t.SendMessage(senderChatId, localization.Tr("Bot set this time for you.", lang))
		t.SendMessage(senderChatId, fmt.Sprintf("%s%s.", localization.Tr("You will receive a quote later today, at ", lang), incomeMsg))
	}
}
