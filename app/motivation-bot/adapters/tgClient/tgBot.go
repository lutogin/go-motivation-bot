package tgClient

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"motivation-bot/config"
	"motivation-bot/integrations/forismatic"
	"motivation-bot/logging"
	"motivation-bot/users"
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

func (t *TgClient) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// Specify allowed update types
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := t.client.GetUpdatesChan(u)

	t.logger.Infoln("Listening messages.")

	for update := range updates {
		if update.Message != nil {
			if update.Message.From.IsBot {
				continue
			}

			if update.Message.IsCommand() {
				t.HandleCommand(&update)
				continue
			} else {
				t.HandleMessage(&update)
			}
			continue
		} else if update.CallbackQuery != nil {
			if update.CallbackQuery.From.IsBot {
				continue
			}

			// this is where we handle the callback query
			t.HandleQuery(&update)
			continue
		}

	}
}
