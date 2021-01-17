package database

import (
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type AnnounceDB struct {
	BaseDB
}

func (db *AnnounceDB) getListFromDB(authorId *primitive.ObjectID, limit int) ([]bson.M, error) {
	aggregateStages := bson.D{{"$project", bson.D{{"id", "$_id"},
		{"_id", 0},
		{"title", "$title"},
		{"description", "$description"},
		{"created", "$created"},
		{"modified", "$modified"},
		{"target_date", "$target_date"}}}}
	sortStage := bson.D{{"$sort", bson.M{"created": -1}}}
	stages := mongo.Pipeline{aggregateStages, sortStage}
	if authorId != nil {
		stages = append(stages, bson.D{{"$match", bson.D{{"_id", *authorId}}}})
	}
	if limit > 0 {
		stages = append(stages, bson.D{{"$limit", limit}})
	}
	cursor, err := db.collection.Aggregate(timeoutContext(), stages)
	if err != nil {
		log.Fatal(err.Error())
	}
	result := []bson.M{}
	if err := cursor.All(timeoutContext(), &result); err != nil {
		log.Fatal(err.Error())
	}
	return result, nil
}
func (db *AnnounceDB) GetList() ([]bson.M, error) {
	return db.getListFromDB(nil, 0)
}

func (db *AnnounceDB) GetListByAuthorId(authorId primitive.ObjectID) ([]bson.M, error) {
	return db.getListFromDB(&authorId, 0)
}

func (db *AnnounceDB) GetAnnounceById(id primitive.ObjectID) (bson.M, error) {
	result := bson.M{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"_id": id}).Decode(&result)
	return result, err
}

func (db *AnnounceDB) AddAnnounce(authorId primitive.ObjectID, announceMap bson.M) (primitive.ObjectID, error) {

	announce := models.Announce{
		AuthorId:    authorId,
		Title:       announceMap["title"].(string),
		Description: announceMap["description"].(string),
		TargetDate:  announceMap["target_date"].(time.Time),
		Created:     time.Now(),
		Modified:    time.Now(),
	}
	res, err := db.collection.InsertOne(timeoutContext(), announce)
	return res.InsertedID.(primitive.ObjectID), err
}

func (db *AnnounceDB) UpdateAnnounceById(id primitive.ObjectID, userId primitive.ObjectID, setElement *bson.D) error {
	if len(*setElement) > 0 {
		*setElement = append(*setElement, bson.E{"modified", time.Now()})
	}
	setMap := bson.D{
		{"$set", *setElement},
	}
	_, err := db.collection.UpdateOne(timeoutContext(), bson.M{"_id": id}, setMap)
	return err
}

func (db *AnnounceDB) DeleteAnnounceById(id primitive.ObjectID, userId primitive.ObjectID) error {
	_, err := db.collection.DeleteOne(timeoutContext(), bson.M{"$and": []bson.M{{"_id": id}, {"author": userId}}})
	return err
}
