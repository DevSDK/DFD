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
		{"participants", "$participants"},
		{"queueid", "$queueid"},
		{"gameid", "$gameid"},
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

func (db *LOLHistoryDB) AddLolHistory(dataMap bson.M, win bool, timestamp int64, gameId string, queueId int64,participants []string) (primitive.ObjectID, error) {
	game := models.LOLHistory{
		Game:      dataMap,
		Win:       win,
		Timestamp: time.Unix(timestamp,0),
		Participants: participants,
		Created:   time.Now(),
		QueueId: 	queueId,
		GameId: gameId,
	}
	res, err := db.collection.InsertOne(timeoutContext(), game)
	return res.InsertedID.(primitive.ObjectID), err
}

func (db *LOLHistoryDB) GetLolHistory(id primitive.ObjectID) (models.LOLHistory, error) {
	history := models.LOLHistory{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"_id": id}).Decode(&history)
	return history, err

}

func (db *LOLHistoryDB) GetCountByDate() []bson.M {
	winStage := bson.D{{"$project", bson.D{
		{"queueid", "$queueid"},
		{"timestamp", "$timestamp"},	
		{"win", bson.D{{"$cond", []interface{}{"$win", 1, 0}}}},	
	}}}
	groupStage := bson.D{{"$group", bson.D{
		{"_id",  bson.D{ 
			{"month", bson.D{{"$month", "$timestamp"}}},
			{"day", bson.D{{"$dayOfMonth", "$timestamp"}}},		
			{"year", bson.D{{"$year", "$timestamp"}}},
			{"queueid", "$queueid"},
		},
	},
	{"count", bson.D{{"$sum",1}}},
	{"win", bson.D{{"$sum","$win"}}},
}}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", -1}}}}
	aggregateStage := bson.D{{"$project", bson.D{{"_id", 0},
	{"queueid", "$_id.queueid"},
	{"month", "$_id.month"},
	{"day", "$_id.day"},
	{"year", "$_id.year"},
	{"count", "$count"},
	{"win", "$win"},
	
	}}}
	result := []bson.M{}
	cursor, err := db.collection.Aggregate(timeoutContext(), mongo.Pipeline{winStage,groupStage,sortStage,aggregateStage})
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := cursor.All(timeoutContext(), &result); err != nil {
		log.Fatal(err.Error())
	}
	return result
}
