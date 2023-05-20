package users

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"motivation-bot/common/logging"
	userDto "motivation-bot/users/dto"
)

type Repository interface {
	Create(ctx context.Context, payload userDto.CreateUserDto) (id string, err error)
	GetById(ctx context.Context, payload userDto.GetUserByIdDto) (user UserEntity, err error)
	GetByFilter(ctx context.Context, payload userDto.GetUsersDto) (user []UserEntity, err error)
	Update(ctx context.Context, payload userDto.UpdateUserDto) (err error)
	Delete(ctx context.Context, payload userDto.DeleteUserDto) (err error)
}

type repo struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *repo) Create(ctx context.Context, payload userDto.CreateUserDto) (id string, err error) {
	result, err := d.collection.InsertOne(ctx, payload)
	if err != nil {
		return "", err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		errMsg := "error during getti ng oid"
		d.logger.Errorln(errMsg)
		d.logger.Traceln(payload)
		return "", errors.New(errMsg)
	}
	return oid.Hex(), nil
}

func (d *repo) GetById(ctx context.Context, payload userDto.GetUserByIdDto) (user UserEntity, err error) {
	oid, err := primitive.ObjectIDFromHex(payload.Id)
	if err != nil {
		return user, err
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return user, errors.New("user not found")
		}
		d.logger.Errorf("Error during looking for users by id: %s \n", payload.Id)
		d.logger.Traceln(payload.Id)
		return user, err
	}
	if err = result.Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (d *repo) GetByFilter(ctx context.Context, payload userDto.GetUsersDto) (foundUsers []UserEntity, err error) {
	// Marshal the anonymous struct into BSON bytes
	bsonBytes, err := bson.Marshal(payload)
	if err != nil {
		log.Fatal("Error marshaling userDto.GetUsersDto struct:", err)
	}

	// Unmarshal the BSON bytes into a bson.M
	var filter bson.M
	err = bson.Unmarshal(bsonBytes, &filter)
	if err != nil {
		log.Fatal("Error unmarshaling BSON bytes:", err)
	}

	cur, err := d.collection.Find(ctx, filter)
	if err != nil {
		d.logger.Errorf("Error during looking for users by filter: %s /n", payload)
		d.logger.Traceln(payload)
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

func (d *repo) Update(ctx context.Context, payload userDto.UpdateUserDto) (err error) {
	oid, err := primitive.ObjectIDFromHex(payload.Id)
	if err != nil {
		d.logger.Errorf("Failed to convert ObjectIDFromHex: %v", err)
		return err
	}

	filter := bson.M{"_id": oid}
	payloadBytes, err := bson.Marshal(payload)
	if err != nil {
		d.logger.Errorf("Failed to bson marshal: %v", err)
		return err
	}

	var updatePayload bson.M
	err = bson.Unmarshal(payloadBytes, &updatePayload)
	if err != nil {
		d.logger.Errorf("Failed to Unmarshal: %v", err)
		return err
	}

	delete(updatePayload, "Id") // drop ID from payload to update

	update := bson.M{"$set": updatePayload}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		d.logger.Errorf("Failed to update users: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		d.logger.Tracef("not found at update by ID: %s", payload.Id)
		return errors.New("not found")
	}
	return nil
}

func (d *repo) Delete(ctx context.Context, payload userDto.DeleteUserDto) (err error) {
	uid, err := primitive.ObjectIDFromHex(payload.Id)
	if err != nil {
		d.logger.Errorf("Failed to ObjectIDFromHex for ID: %s", payload.Id)
		return err
	}
	filter := bson.M{"_id": uid}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		d.logger.Tracef("User ID: %s", payload.Id)
		d.logger.Errorf("Failed to delete users: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		d.logger.Errorf("Not found users with ID: %s", payload.Id)
	}

	return nil
}

func NewRepository(database *mongo.Database, collection string, logger *logging.Logger) Repository {
	return &repo{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
