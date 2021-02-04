package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

//Instance for singleton
var Instance = DBInstance{}

//BaseDB for DB class
type BaseDB struct {
	collection *mongo.Collection
}

//DBInstance data structure for singlethon
type DBInstance struct {
	mongoClient      *mongo.Client
	redisClient      *redis.Client
	database         *mongo.Database
	User             UserDB
	Announce         AnnounceDB
	Redis            RedisStore
	Role             RoleDB
	Image            ImageDB
	DFDHistory       DFDHistoryDB
	LOLHistory       LOLHistoryDB
	ApplicationToken ApplicationTokenDB
}

func timeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func initializeCollections() {
	Instance.User.collection = Instance.database.Collection("User")
	Instance.Role.collection = Instance.database.Collection("Role")
	Instance.LOLHistory.collection = Instance.database.Collection("LOLHistory")
	Instance.ApplicationToken.collection = Instance.database.Collection("ApplicationToken")
	Instance.Announce.collection = Instance.database.Collection("Announce")
	Instance.DFDHistory.collection = Instance.database.Collection("DFDHistory")
	Instance.Image.database = Instance.mongoClient.Database("Images")
	Instance.Image.collection = Instance.mongoClient.Database("Images").Collection("fs.files")
}

func initializeMongoDB() error {
	DB_LOCATION := os.Getenv("DB_LOCATION")
	DB_AUTH_ID := os.Getenv("DB_AUTH_ID")
	DB_AUTH_PASSWORD := os.Getenv("DB_AUTH_PASSWORD")
	DB_CONNECTION_ADDRESS := fmt.Sprintf("mongodb://%s:%s@%s", DB_AUTH_ID, DB_AUTH_PASSWORD, DB_LOCATION)
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(DB_CONNECTION_ADDRESS))
	if err != nil {
		log.Print("DB URL: " + DB_CONNECTION_ADDRESS)
		return err
	}
	if err := mongoClient.Connect(timeoutContext()); err != nil {
		log.Print("DB URL: " + DB_CONNECTION_ADDRESS)
		return err
	}
	if err := mongoClient.Ping(timeoutContext(), nil); err != nil {
		log.Print("DB URL: " + DB_CONNECTION_ADDRESS)
		return err
	}
	Instance.mongoClient = mongoClient
	Instance.database = mongoClient.Database("DFD")
	initializeCollections()
	Instance.User.collection.Indexes().CreateOne(
		timeoutContext(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}, {Key: "discord_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	Instance.Role.collection.Indexes().CreateOne(
		timeoutContext(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	Instance.LOLHistory.collection.Indexes().CreateOne(
		timeoutContext(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "gameid", Value: 1}},
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
	if err != nil {
		log.Print("Redis")
	}
	return err
}

//Initialize database
func Initialize() error {
	if err := initializeMongoDB(); err != nil {
		return err
	}

	if err := initializeRedis(); err != nil {
		return err
	}
	return nil
}

//Disconnect database
func Disconnect() {
	Instance.mongoClient.Disconnect(timeoutContext())
}
