package database

import (
	"errors"
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type DFDHistoryDB struct {
	BaseDB
}

func (db *DFDHistoryDB) getList(userId primitive.ObjectID, limit int) ([]bson.M, error) {
	result := []bson.M{}

	matchStage := bson.D{{"$match", bson.D{{"user_id", userId}}}}
	projectStage := bson.D{{"$project", bson.D{{"id", "$_id"}, {"_id", 0}, {"created", "$created"},
		{"state", "$state"}, {"was", "$was"}}}}
	limitStage := bson.D{{"$limit", limit}}
	sortStage := bson.D{{"$sort", bson.D{{"created", -1}}}}
	stages := mongo.Pipeline{matchStage, projectStage, sortStage}
	if limit > 0 {
		stages = append(stages, limitStage)
	}

	cursor, err := db.collection.Aggregate(timeoutContext(), stages)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := cursor.All(timeoutContext(), &result); err != nil {
		log.Fatal(err.Error())
	}
	return result, nil
}

func (db *DFDHistoryDB) GetList(userId primitive.ObjectID) ([]bson.M, error) {
	return db.getList(userId, 0)
}

func (db *DFDHistoryDB) GetRecent(userId primitive.ObjectID) (bson.M, error) {
	result, err := db.getList(userId, 1)
	if err != nil {
		return bson.M{}, err
	}
	if len(result) != 1 {
		return bson.M{}, errors.New("DFDHistory is empty")
	}
	return result[0], nil
}

func (db *DFDHistoryDB) Push(userId primitive.ObjectID, newState string) error {
	recent, err := db.GetRecent(userId)
	wasState := ""
	if err == nil {
		wasState = recent["state"].(string)
	}
	insertData := models.DFDHistory{
		UserId:  userId,
		Was:     wasState,
		State:   newState,
		Created: time.Now(),
	}
	_, err = db.collection.InsertOne(timeoutContext(), insertData)
	return err
}
