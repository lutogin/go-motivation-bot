package tgClient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/integrations/forismatic"
)

func (t *TgClient) SendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := t.client.Send(msg)
	if err != nil {
		t.logger.Traceln(fmt.Sprintf("Error on sending message to %d", chatId))
		t.logger.Error(err)
	}
}

func (t *TgClient) SendMessageToAdmin(message string) {
	msg := tgbotapi.NewMessage(t.cfg.AdminChatId, message)
	t.client.Send(msg)
}

func (t *TgClient) SendQuote(chatId int64, quote forismatic.GetQuoteResponse) {
	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("%s\n\n<i>%s</i>", quote.QuoteText, quote.QuoteAuthor))
	msg.ParseMode = "html"

	_, err := t.client.Send(msg)
	if err != nil {
		t.logger.Traceln(fmt.Sprintf("Error on sending quote to %d", chatId))
		t.logger.Error(err)
	}
}
