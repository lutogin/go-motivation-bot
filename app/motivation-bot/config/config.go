package config

import (
	"motivation-bot/common/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Env struct {
	IsDebug        bool   `yaml:"isDebug" env:"IS_DEBUG" env-default:"true"`
	Host           string `yaml:"host" env:"HOST" env-default:"127.0.0.1"`
	Port           string `yaml:"port" env:"PORT" env-default:"8080"`
	TgApiKey       string `yaml:"tgApiKey" env:"TG_API_KEY"`
	AdminChatId    int64  `yaml:"adminChatId" env:"ADMIN_CHAT_ID" env-default:"479518542"`
	MongoHost      string `yaml:"mongoHost" env:"MONGO_HOST" env-default:"127.0.0.1"`
	MongoPort      string `yaml:"mongoPort" env:"MONGO_PORT" env-default:"27017"`
	MongoUser      string `yaml:"mongoUser" env:"MONGO_USER"`
	MongoPassword  string `yaml:"mongoPassword" env:"MONGO_PASSWORD"`
	MongoDatabase  string `yaml:"mongoDatabase" env:"MONGO_DATABASE"`
	MongoUriScheme string `yaml:"mongoUriScheme" env:"MONGO_URI_SCHEME" env-default:"mongodb"`
}

var (
	instance *Env
	once     sync.Once
)

func GetConfig(logger *logging.Logger) *Env {
	once.Do(func() { // do it once. Singleton pattern
		logger.Infoln("Read application's config.")
		instance = &Env{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatalln(help)
		}
	})

	logger.Infoln("Read application's config successful.")

	return instance
}
