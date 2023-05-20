package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/dig"
	"log"
	"motivation-bot/common/logging"
	"motivation-bot/config"
	"time"
)

type MongoConnectOpt struct {
	Host      string
	Port      string
	User      string // optional field
	Password  string // optional field
	Database  string
	UriScheme string
}

type MongoConnection struct {
	DB *mongo.Database
}

func NewMongoConnection(config *config.Env, logger *logging.Logger) *MongoConnection {
	// MongoDB connection string
	suffix := "retryWrites=true&w=majority"

	connectionString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		config.MongoUriScheme,
		config.MongoUser,
		config.MongoPassword,
		config.Host,
		config.Port,
		config.MongoDatabase,
		suffix,
	)

	clientOptions := options.Client().ApplyURI(connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		logger.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)

	if err != nil {
		logger.Fatal(err)
	}

	// Get a handle to the respective database
	database := client.Database("your-database-name")

	return &MongoConnection{
		DB: database,
	}
}

func ProvideMongoConnection(container *dig.Container) {
	err := container.Provide(NewMongoConnection)
	if err != nil {
		log.Fatal(err)
	}
}
