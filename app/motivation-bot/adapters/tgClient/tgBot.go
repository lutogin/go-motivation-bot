package tgClient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/config"
	"motivation-bot/integrations/forismatic"
	"motivation-bot/logging"
)

type Client struct {
	logger     *logging.Logger
	cfg        *config.Env
	forismatic *forismatic.Client
}

func NewClient(logger *logging.Logger, config *config.Env, forismatic *forismatic.Client) *Client {
	return &Client{
		logger:     logger,
		cfg:        config,
		forismatic: forismatic,
	}
}

func (c *Client) HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, client *forismatic.Client) {
	quote := client.GetQuote()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n\n<i>%s</i>", quote.QuoteText, quote.QuoteAuthor))
	msg.ParseMode = "html"

	bot.Send(msg)
}

func (c *Client) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
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
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
		bot.Send(msg)
	}

}

func (c *Client) HandleQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	data := update.CallbackQuery.Data
	switch data {
	case "setLang:en":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You chose EN")
		_, err := bot.Send(msg)
		if err != nil {
			c.logger.Error(err)
		}

	case "setLang:ru":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You chose RU")
		_, err := bot.Send(msg)
		if err != nil {
			c.logger.Error(err)
		}
	}

}

func (c *Client) Run() {
	bot, err := tgbotapi.NewBotAPI(c.cfg.TgApiKey)
	if err != nil {
		c.logger.Debugf("Bot connection was failed: %s \n", c.cfg.TgApiKey)
		c.logger.Panic(err)
	}

	bot.Debug = c.cfg.IsDebug

	c.logger.Infof("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// Specify allowed update types
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := bot.GetUpdatesChan(u)

	c.logger.Infoln("Listening messages.")

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				c.HandleCommand(bot, &update)
				continue
			} else {
				c.HandleMessage(bot, &update, c.forismatic)
			}
			continue
		} else if update.CallbackQuery != nil {
			// this is where we handle the callback query
			c.HandleQuery(bot, &update)
			continue
		}

	}
}
