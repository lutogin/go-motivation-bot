package mongoClient

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"motivation-bot/config"
	"motivation-bot/logging"
	"time"
)

type MongoConnection struct {
	DB *mongo.Database
}

const CollectionName = "users"

// todo I don't like it. Like a break SOLID. need to rewrite
func makeIndexes(db *mongo.Database) {
	// Get the collection
	collection := db.Collection(CollectionName)

	// Create a model for the unique index
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"userName", 1}}, // Index key
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create the index
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}

}

func NewMongoConnection(config *config.Env, logger *logging.Logger) *MongoConnection {
	logger.Logger.Infoln("Registering mongo connection.")

	suffix := "retryWrites=true&w=majority"
	connectionString := fmt.Sprintf("%s://%s:%s@%s/%s?%s",
		config.MongoUriScheme,
		config.MongoUser,
		config.MongoPassword,
		config.MongoHost,
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
	err = client.Ping(ctx, nil)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Infoln("DB connection is established.")

	// Get a handle to the respective database
	database := client.Database(config.MongoDatabase)

	makeIndexes(database)

	return &MongoConnection{
		DB: database,
	}
}
