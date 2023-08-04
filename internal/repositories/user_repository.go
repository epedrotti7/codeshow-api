package repository

import (
	"context"
	"log"

	connection "github.com/epedrotti7/codeshow-api/database"
	"github.com/epedrotti7/codeshow-api/internal/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Create(user *structs.User) *structs.User {
	ctx := context.TODO()
	collection := connection.GetClient().Database("codeshow").Collection("users")

	result, err := collection.InsertOne(ctx, user)

	if err != nil {
		log.Fatal(err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		user.ID = &oid
	}

	return user
}

func UpdateUserScore(userId string) error {
	ctx := context.TODO()
	collection := connection.GetClient().Database("codeshow").Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"score": 1}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func FindUserByID(id string) (*structs.User, error) {
	ctx := context.TODO()
	collection := connection.GetClient().Database("codeshow").Collection("users")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var user structs.User

	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &user, nil
}

func FindUserByEmail(email string) (*structs.User, error) {
	ctx := context.TODO()
	collection := connection.GetClient().Database("codeshow").Collection("users")

	filter := bson.M{"email": email}
	var user structs.User

	err := collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &user, nil
}
