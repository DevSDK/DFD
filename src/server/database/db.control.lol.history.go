package database

import (
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type LOLHistoryDB struct{ BaseDB }

func (db *LOLHistoryDB) GetList() []bson.M {
	aggregateStage := bson.D{{"$project", bson.D{{"id", "$_id"},
		{"_id", 0},
		{"created", "$created"},
		{"win", "$win"},
		{"timestamp", "$timestamp"}}}}
	sortStage := bson.D{{"$sort", bson.M{"timestamp": -1}}}
	result := []bson.M{}
	cursor, err := db.collection.Aggregate(timeoutContext(), mongo.Pipeline{aggregateStage, sortStage})
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := cursor.All(timeoutContext(), &result); err != nil {
		log.Fatal(err.Error())
	}
	return result
}

func (db *LOLHistoryDB) AddLolHistory(dataMap bson.M, win bool, timestamp int64) (primitive.ObjectID, error) {
	game := models.LOLHistory{
		Game:      dataMap,
		Win:       win,
		Timestamp: timestamp,
		Created:   time.Now(),
	}
	res, err := db.collection.InsertOne(timeoutContext(), game)
	return res.InsertedID.(primitive.ObjectID), err
}

func (db *LOLHistoryDB) GetLolHistory(id primitive.ObjectID) (models.LOLHistory, error) {
	history := models.LOLHistory{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"_id": id}).Decode(&history)
	return history, err

}
