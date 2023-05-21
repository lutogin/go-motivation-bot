package cronJobs

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"motivation-bot/adapters/tgClient"
	"motivation-bot/config"
	"motivation-bot/logging"
	"motivation-bot/users"
	usersDto "motivation-bot/users/dto"
	"time"
)

type CronJob struct {
	cron         *cron.Cron
	logger       *logging.Logger
	tgClient     *tgClient.TgClient
	usersService *users.Service
	cfg          *config.Env
}

func NewCronJob(logger *logging.Logger, tgClient *tgClient.TgClient, usersService *users.Service, cfg *config.Env) *CronJob {
	logger.Info("Registering cron job.")

	c := cron.New()
	return &CronJob{
		cron:         c,
		logger:       logger,
		tgClient:     tgClient,
		usersService: usersService,
		cfg:          cfg,
	}
}

func (c *CronJob) StartCron() {
	interval := fmt.Sprintf("@every %dm", c.cfg.CronInterval)
	_, err := c.cron.AddFunc(interval, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		c.logger.Infoln("Running job", time.Now())
		from := time.Now().Add(-5 * time.Minute)
		to := time.Now().Add(5 * time.Minute)
		users, _ := c.usersService.GetByAlertingDate(ctx, usersDto.GetUserByAlertingDateDto{From: from, To: to})

		for _, u := range users {
			c.logger.Infoln(u)
		}
	})
	if err != nil {
		c.logger.Error("Cron job is not set.")
		return
	}

	c.logger.Info(fmt.Sprintf("Cron job is set to %s.", interval))

	c.cron.Start()
}
