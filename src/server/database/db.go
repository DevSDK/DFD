package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var Instance = DBInstance{}

type DBInstance struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	database    *mongo.Database
	User        UserDB
	Redis       RedisStore
	Role        RoleDB
	Image       ImageDB
}

func timeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func initializeMongoDB() error {
	DB_LOCATION := os.Getenv("DB_LOCATION")
	DB_AUTH_ID := os.Getenv("DB_AUTH_ID")
	DB_AUTH_PASSWORD := os.Getenv("DB_AUTH_PASSWORD")
	DB_CONNECTION_ADDRESS := fmt.Sprintf("mongodb://%s:%s@%s", DB_AUTH_ID, DB_AUTH_PASSWORD, DB_LOCATION)
	mongoClient, _ := mongo.NewClient(options.Client().ApplyURI(DB_CONNECTION_ADDRESS))

	if err := mongoClient.Connect(timeoutContext()); err != nil {
		return err
	}
	if err := mongoClient.Ping(timeoutContext(), nil); err != nil {
		return err
	}
	Instance.mongoClient = mongoClient
	Instance.database = mongoClient.Database("DFD")
	userCollection := Instance.database.Collection("User")
	userCollection.Indexes().CreateOne(
		timeoutContext(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}, {Key: "discord_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	Instance.database.Collection("Role").Indexes().CreateOne(
		timeoutContext(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	return nil
}

func initializeRedis() error {
	REDIS_LOCATION := os.Getenv("REDIS_LOCATION")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")

	Instance.redisClient = redis.NewClient(&redis.Options{
		Addr:     REDIS_LOCATION,
		Password: REDIS_PASSWORD,
		DB:       0,
	})
	_, err := Instance.redisClient.Ping().Result()
	return err
}

func Initialize() error {
	if err := initializeMongoDB(); err != nil {
		return err
	}

	if err := initializeRedis(); err != nil {
		return err
	}
	return nil
}

func Disconnect() {
	Instance.mongoClient.Disconnect(timeoutContext())
}
