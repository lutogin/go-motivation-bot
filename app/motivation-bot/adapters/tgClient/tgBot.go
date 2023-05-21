package tgClient

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/common"
	"motivation-bot/config"
	"motivation-bot/integrations/forismatic"
	"motivation-bot/logging"
	"motivation-bot/users"
	userDto "motivation-bot/users/dto"
	"strconv"
	"strings"
	"time"
)

type TgClient struct {
	logger      *logging.Logger
	cfg         *config.Env
	forismatic  *forismatic.Client
	client      *tgbotapi.BotAPI
	userService *users.Service
}

func NewClient(logger *logging.Logger, config *config.Env, forismatic *forismatic.Client, userService *users.Service) *TgClient {
	bot, err := tgbotapi.NewBotAPI(config.TgApiKey)
	if err != nil {
		logger.Debugf("Bot connection was failed: %s \n", config.TgApiKey)
		logger.Panic(err)
	}

	bot.Debug = config.IsDebug

	logger.Infof("Authorized on account %s\n", bot.Self.UserName)

	return &TgClient{
		logger:      logger,
		cfg:         config,
		forismatic:  forismatic,
		client:      bot,
		userService: userService,
	}
}

func (t *TgClient) HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, client *forismatic.Client) {
	incomeMsg := update.Message.Text
	senderChatId := update.Message.Chat.ID

	if !common.Is24TimeFormat(update.Message.Text) {
		msg := tgbotapi.NewMessage(senderChatId, "Wrong format.")
		_, err := bot.Send(msg)
		if err != nil {
			t.logger.Error(err)
		}
	}

	fmt.Println("Is24TimeFormat")

	incomeDateArr := strings.Split(incomeMsg, ":")
	incomeHr, _ := strconv.Atoi(incomeDateArr[0])
	incomeMin, _ := strconv.Atoi(incomeDateArr[1])

	fmt.Println(incomeDateArr)
	fmt.Println(incomeHr)
	fmt.Println(incomeMin)
	fmt.Println(incomeMin)

	now := time.Now()
	setTime := time.Date(now.Year(), now.Month(), now.Day(), incomeHr, incomeMin, 0, 0, now.Location())

	if time.Now().After(setTime) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		t.userService.Update(ctx, userDto.UpdateUserDto{UserName: update.Message.From.UserName, AlertingDate: setTime.Add(24 * time.Hour)})

		t.SendMessage(senderChatId, "Bot set this time for you.")
		t.SendMessage(senderChatId, "There is one quote for you for today:")

		quote := client.GetQuote("ru")
		t.SendQuote(senderChatId, quote)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		t.userService.Update(ctx, userDto.UpdateUserDto{UserName: update.Message.From.UserName, AlertingDate: setTime})

		t.SendMessage(senderChatId, "Bot set this time for you.")
		t.SendMessage(senderChatId, fmt.Sprintf("You will receive a quote later today, at %s.", incomeMsg))
	}
}

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
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
		bot.Send(msg)
	}

}

func (t *TgClient) HandleQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	data := update.CallbackQuery.Data
	switch data {
	case "setLang:en":
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		t.userService.Create(ctx, userDto.CreateUserDto{
			ChatId:    update.CallbackQuery.From.ID,
			FirstName: update.CallbackQuery.From.FirstName,
			LastName:  update.CallbackQuery.From.LastName,
			UserName:  update.CallbackQuery.From.UserName,
			Lang:      "en",
		})

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

	case "setLang:ru":
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		t.userService.Create(ctx, userDto.CreateUserDto{
			ChatId:    update.CallbackQuery.From.ID,
			FirstName: update.CallbackQuery.From.FirstName,
			LastName:  update.CallbackQuery.From.LastName,
			UserName:  update.CallbackQuery.From.UserName,
			Lang:      "ru",
		})

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

func (t *TgClient) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// Specify allowed update types
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := t.client.GetUpdatesChan(u)

	t.logger.Infoln("Listening messages.")

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				t.HandleCommand(t.client, &update)
				continue
			} else {
				t.HandleMessage(t.client, &update, t.forismatic)
			}
			continue
		} else if update.CallbackQuery != nil {
			// this is where we handle the callback query
			t.HandleQuery(t.client, &update)
			continue
		}

	}
}
