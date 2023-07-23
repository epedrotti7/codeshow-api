package repository

import (
	"context"
	"fmt"
	"log"

	connection "github.com/epedrotti7/codeshow-api/database"
	"github.com/epedrotti7/codeshow-api/internal/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateQuestionByUserId(question *structs.Question) *structs.Question {
	ctx := context.TODO()
	collection := connection.GetClient().Database("codeshow").Collection("questions")

	questionSaved, err := collection.InsertOne(ctx, question)

	if err != nil {
		log.Fatal(err)
	}

	return &structs.Question{
		Id:           questionSaved.InsertedID.(primitive.ObjectID).Hex(),
		UserId:       question.UserId,
		Question:     question.Question,
		Answer:       question.Answer,
		Alternatives: question.Alternatives,
	}

}

func UpdateQuestionById(answerReq *structs.Answer, id string) error {
	ctx := context.TODO()
	collection := connection.GetClient().Database("codeshow").Collection("questions")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	filter := bson.M{"_id": objectID}

	update := bson.D{
		{"$set", bson.D{
			{"userAlternative", answerReq.Answer},
		}},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}

	return nil
}

func FindQuestionsByUserId(userId string) ([]*structs.Question, error) {
	ctx := context.TODO()
	collection := connection.GetClient().Database("codeshow").Collection("questions")

	filter := bson.M{"userid": userId}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var questions []*structs.Question
	if err = cursor.All(ctx, &questions); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return questions, nil
}

func FindQuestionById(id string) (*structs.Question, error) {
	ctx := context.TODO()
	collection := connection.GetClient().Database("codeshow").Collection("questions")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var question structs.Question

	err = collection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &question, nil
}
