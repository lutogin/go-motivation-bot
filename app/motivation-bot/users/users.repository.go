package users

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"motivation-bot/adapters/mongoClient"
	"motivation-bot/logging"
	"motivation-bot/users/dto"
)

type GetUsersByAlertingTimeDto struct {
	HoursFrom   int
	MinutesFrom int
	HoursTo     int
	MinutesTo   int
}

type Repository interface {
	Create(ctx context.Context, payload usersDto.CreateUserDto) (id string, err error)
	GetById(ctx context.Context, payload usersDto.GetUserByIdDto) (user UserEntity, err error)
	GetByFilter(ctx context.Context, payload usersDto.GetUsersDto) (user []UserEntity, err error)
	Update(ctx context.Context, payload usersDto.UpdateUserDto) (err error)
	Upsert(ctx context.Context, payload usersDto.CreateUserDto) (err error)
	Delete(ctx context.Context, payload usersDto.DeleteUserDto) (err error)
}

type Repo struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewRepository(mongoConn *mongoClient.MongoConnection, logger *logging.Logger) *Repo {
	logger.Logger.Infoln("Registering new users repository.")

	return &Repo{
		collection: mongoConn.DB.Collection("users"),
		logger:     logger,
	}
}

func (r *Repo) Create(ctx context.Context, payload usersDto.CreateUserDto) (string, error) {
	result, err := r.collection.InsertOne(ctx, payload)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		errMsg := "error during getting oid"
		r.logger.Errorln(errMsg)
		r.logger.Traceln(payload)
		return "", errors.New(errMsg)
	}

	return oid.Hex(), nil
}

func (r *Repo) GetById(ctx context.Context, payload usersDto.GetUserByIdDto) (UserEntity, error) {
	var user UserEntity
	oid, err := primitive.ObjectIDFromHex(payload.Id)
	if err != nil {
		return user, err
	}
	filter := bson.M{"_id": oid}
	result := r.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return user, errors.New("user not found")
		}
		r.logger.Errorf("Error during looking for users by id: %s \n", payload.Id)
		r.logger.Traceln(payload.Id)
		return user, err
	}
	if err = result.Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (r *Repo) GetByChatId(ctx context.Context, payload usersDto.GetUserByChatIdDto) (UserEntity, error) {
	var user UserEntity
	filter := bson.M{"chatId": payload.ChatId}
	result := r.collection.FindOne(ctx, filter)

	if err := result.Err(); err != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return user, errors.New("user not found")
		}
		r.logger.Errorf("Error during looking for users by chat id: %s \n", payload.ChatId)
		r.logger.Traceln(payload.ChatId)
		return user, err
	}

	if err := result.Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repo) GetByFilter(ctx context.Context, payload usersDto.GetUsersDto) ([]UserEntity, error) {
	var foundUsers []UserEntity
	// Marshal the anonymous struct into BSON bytes
	bsonBytes, err := bson.Marshal(payload)
	if err != nil {
		log.Fatal("Error marshaling usersDto.GetUsersDto struct:", err)
	}

	// Unmarshal the BSON bytes into a bson.M
	var filter bson.M
	err = bson.Unmarshal(bsonBytes, &filter)
	if err != nil {
		log.Fatal("Error unmarshaling BSON bytes:", err)
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("Error during looking for users by filter: %s /n", payload)
		r.logger.Traceln(payload)
		return foundUsers, err
	}
	defer cur.Close(ctx)

	// Iterate over the cursor and decode the results into a User struct
	for cur.Next(context.Background()) {
		var user UserEntity

		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		foundUsers = append(foundUsers, user)
	}

	// Check if there was an error in the cursor iteration
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return foundUsers, nil
}

func (r *Repo) Update(ctx context.Context, payload usersDto.UpdateUserDto) error {
	filter := bson.M{"userName": payload.UserName}
	payloadBytes, err := bson.Marshal(payload)
	if err != nil {
		r.logger.Errorf("Failed to bson marshal: %v", err)
		return err
	}

	var updatePayload bson.M
	err = bson.Unmarshal(payloadBytes, &updatePayload)
	if err != nil {
		r.logger.Errorf("Failed to Unmarshal: %v", err)
		return err
	}

	update := bson.M{"$set": updatePayload}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("Failed to update users: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		r.logger.Warnf("not found at update by UserName: %s", payload.UserName)
		return errors.New("not found user for update alerting time")
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, payload usersDto.DeleteUserDto) error {
	uid, err := primitive.ObjectIDFromHex(payload.Id)
	if err != nil {
		r.logger.Errorf("Failed to ObjectIDFromHex for ID: %s", payload.Id)
		return err
	}
	filter := bson.M{"_id": uid}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Tracef("User ID: %s", payload.Id)
		r.logger.Errorf("Failed to delete users: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		r.logger.Errorf("Not found users with ID: %s", payload.Id)
	}

	return nil
}

func (r *Repo) GetUsersByAlertingTime(ctx context.Context, payload GetUsersByAlertingTimeDto) ([]UserEntity, error) {
	// Build the query
	filter := bson.M{
		"alertingTime.hours": bson.M{
			"$gte": payload.HoursFrom,
			"$lte": payload.HoursTo,
		},
		"alertingTime.minutes": bson.M{
			"$gte": payload.MinutesFrom,
			"$lte": payload.MinutesTo,
		},
	}

	fmt.Println(payload)

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("Failed to get users by alerting date: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []UserEntity
	if err = cursor.All(ctx, &users); err != nil {
		r.logger.Errorf("Failed to load users by alerting date to the entity: %v", err)
		return nil, err
	}

	return users, nil
}

func (r *Repo) Upsert(ctx context.Context, payload usersDto.UpdateUserDto) error {
	// Define the filter for the update operation
	filter := bson.M{"chatId": payload.ChatId}

	payloadBytes, err := bson.Marshal(payload)
	if err != nil {
		r.logger.Errorf("Failed to bson marshal: %v", err)
		return err
	}

	var updatePayload bson.M
	err = bson.Unmarshal(payloadBytes, &updatePayload)
	if err != nil {
		r.logger.Errorf("Failed to bson marshal: %v", err)
		return err
	}

	// Define the update to be performed if a document matching the filter is found
	update := bson.M{
		"$set": updatePayload,
	}

	// Define options to indicate that upsert should be used
	opts := options.Update().SetUpsert(true)

	// Perform the update operation
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
