package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var Instance = DBInstance{}

type DBInstance struct {
	client   *mongo.Client
	database *mongo.Database
}

func timeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func Initialize() error {
	DB_LOCATION := os.Getenv("DB_LOCATION")
	DB_AUTH_ID := os.Getenv("DB_AUTH_ID")
	DB_AUTH_PASSWORD := os.Getenv("DB_AUTH_PASSWORD")
	DB_CONNECTION_ADDRESS := fmt.Sprintf("mongodb://%s:%s@%s", DB_AUTH_ID, DB_AUTH_PASSWORD, DB_LOCATION)
	client, _ := mongo.NewClient(options.Client().ApplyURI(DB_CONNECTION_ADDRESS))

	if err := client.Connect(timeoutContext()); err != nil {
		return err
	}
	if err := client.Ping(timeoutContext(), nil); err != nil {
		return err
	}
	Instance.client = client
	Instance.database = client.Database("DFD")
	userCollection := Instance.database.Collection("User")
	userCollection.Indexes().CreateOne(
		timeoutContext(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return nil
}

func Disconnect() {
	Instance.client.Disconnect(timeoutContext())
}
