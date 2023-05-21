package di

import (
	"go.uber.org/dig"
	"motivation-bot/adapters/mongoClient"
	"motivation-bot/adapters/tgClient"
	"motivation-bot/common"
	"motivation-bot/config"
	cronJobs "motivation-bot/cron"
	"motivation-bot/integrations/forismatic"
	"motivation-bot/logging"
	"motivation-bot/users"
)

type App struct {
	//Config *config.Env
	//Logger *logging.Logger
	Client *tgClient.TgClient
	//UsersRepo *users.Repo
	Cron *cronJobs.CronJob
}

func NewApp(client *tgClient.TgClient, cron *cronJobs.CronJob) *App {
	return &App{
		//Config:    config,
		//Logger:    logger,
		Client: client,
		//UsersRepo: repo,
		Cron: cron,
	}
}

func ProvideMongoConnection(container *dig.Container) {
	err := container.Provide(mongoClient.NewMongoConnection)
	common.CriticErrorHandler(err)
}

func ProvideUsersRepository(container *dig.Container) {
	err := container.Provide(users.NewRepository)
	common.CriticErrorHandler(err)
}

func ProvideUsersService(container *dig.Container) {
	err := container.Provide(users.NewService)
	common.CriticErrorHandler(err)
}

func ProvideCronJob(container *dig.Container) {
	err := container.Provide(cronJobs.NewCronJob)
	common.CriticErrorHandler(err)
}

func BuildContainer() *dig.Container {
	container := dig.New()

	// Provide the logger object to the container
	err := container.Provide(logging.GetLogger)
	common.CriticErrorHandler(err)

	// Provide the config object to the container
	err = container.Provide(config.GetConfig)
	common.CriticErrorHandler(err)

	// Provide the App object to the container
	err = container.Provide(NewApp)
	common.CriticErrorHandler(err)

	// Provide the TG client object to the container
	err = container.Provide(forismatic.NewClient)
	common.CriticErrorHandler(err)

	// Provide the TG client object to the container
	err = container.Provide(tgClient.NewClient)
	common.CriticErrorHandler(err)

	ProvideMongoConnection(container) // todo: split all containers like this.
	ProvideUsersRepository(container)
	ProvideUsersService(container)
	ProvideCronJob(container)

	return container
}
