package connection

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	var err error
	ctx := context.TODO()

	clientOptions := options.Client().ApplyURI("mongodb://codeshow:codeshow@localhost:27017")
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Falha ao conectar no MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Falha ao fazer ping no MongoDB: %v", err)
	}

	fmt.Println("Conectado ao MongoDB!")
}

// GetClient é uma função de exportação para retornar a instância do cliente MongoDB.
func GetClient() *mongo.Client {
	return client
}
