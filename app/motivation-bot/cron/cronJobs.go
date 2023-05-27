package cronJobs

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"motivation-bot/adapters/tgClient"
	"motivation-bot/config"
	"motivation-bot/integrations/forismatic"
	"motivation-bot/logging"
	"motivation-bot/users"
	usersDto "motivation-bot/users/dto"
	"time"
)

type CronJob struct {
	cron             *cron.Cron
	logger           *logging.Logger
	tgClient         *tgClient.TgClient
	usersService     *users.Service
	forismaticClient *forismatic.Client
	cfg              *config.Env
}

func NewCronJob(logger *logging.Logger, tgClient *tgClient.TgClient, usersService *users.Service, forismaticClient *forismatic.Client, cfg *config.Env) *CronJob {
	logger.Info("Registering cron job.")

	c := cron.New()
	return &CronJob{
		cron:             c,
		logger:           logger,
		tgClient:         tgClient,
		usersService:     usersService,
		forismaticClient: forismaticClient,
		cfg:              cfg,
	}
}

func (c *CronJob) StartCron() {
	interval := fmt.Sprintf("*/%d * * * *", c.cfg.CronInterval)
	_, err := c.cron.AddFunc(interval, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		c.logger.Infoln("Running job", time.Now())
		users, _ := c.usersService.GetByAlertingDate(ctx, usersDto.GetUserByAlertingTimeDto{Date: time.Now()})

		if len(users) == 0 {
			return
		}

		quotes := map[string]forismatic.GetQuoteResponse{
			"ru": c.forismaticClient.GetQuote("ru"),
			"en": c.forismaticClient.GetQuote("en"),
		}

		for _, u := range users {
			c.logger.Infoln(u)

			c.tgClient.SendQuote(u.ChatId, quotes[u.Lang])
		}
	})

	if err != nil {
		c.logger.Error("Cron job is not set.")
		return
	}

	c.logger.Info(fmt.Sprintf("Cron job is set to %s.", interval))

	c.cron.Start()
}
