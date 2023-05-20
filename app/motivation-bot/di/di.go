package di

import (
	"go.uber.org/dig"
	"motivation-bot/common"
	"motivation-bot/config"
	"motivation-bot/integrations/forismatic"
	"motivation-bot/logging"
	"motivation-bot/pkg/mongoClient"
	"motivation-bot/pkg/tgClient"
	"motivation-bot/users"
)

type App struct {
	//Config *config.Env
	//Logger *logging.Logger
	Client *tgClient.Client
	//UsersRepo *users.Repo
}

func NewApp(client *tgClient.Client) *App {
	return &App{
		//Config:    config,
		//Logger:    logger,
		Client: client,
		//UsersRepo: repo,
	}
}

func ProvideMongoConnection(container *dig.Container) {
	err := container.Provide(mongoClient.NewMongoConnection)
	common.CriticErrorHandler(err)
}

func ProvideNewUsersRepository(container *dig.Container) {
	err := container.Provide(users.NewRepository)
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
	ProvideNewUsersRepository(container)

	return container
}
