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

//AnnounceDB is a structure for announce data access
type AnnounceDB struct {
	BaseDB
}

func (db *AnnounceDB) getListFromDB(authorID *primitive.ObjectID, timestamp int64) ([]bson.M, error) {
	aggregateStages := bson.D{{"$project", bson.D{{"id", "$_id"},
		{"_id", 0},
		{"title", "$title"},
		{"author", "$author"},
		{"description", "$description"},
		{"created", "$created"},
		{"modified", "$modified"},
		{"target_date", "$target_date"}}}}
	matchStage := bson.D{{"$match", bson.M{"target_date": bson.M{"$gt": time.Unix(timestamp, 64)}}}}
	sortStage := bson.D{{"$sort", bson.M{"target_date": -1}}}
	stages := mongo.Pipeline{aggregateStages, sortStage, matchStage}
	log.Print(time.Unix(timestamp, 64))
	if authorID != nil {
		stages = append(stages, bson.D{{"$match", bson.D{{"author", *authorID}}}})
	}

	result := []bson.M{}
	cursor, err := db.collection.Aggregate(timeoutContext(), stages)
	if err != nil {
		log.Print(err.Error())
		return result, err
	}
	if err := cursor.All(timeoutContext(), &result); err != nil {
		log.Print(err.Error())
		return result, err
	}
	return result, nil
}

//GetListWithTimestamp function returns announces after given timestamp
func (db *AnnounceDB) GetListWithTimestamp(timestamp int64) ([]bson.M, error) {
	return db.getListFromDB(nil, timestamp)
}

//GetList returns all anounces
func (db *AnnounceDB) GetList() ([]bson.M, error) {
	return db.getListFromDB(nil, 0)
}

//GetListByAuthorID returns annoucnes written by gicen user id
func (db *AnnounceDB) GetListByAuthorID(authorID primitive.ObjectID) ([]bson.M, error) {
	return db.getListFromDB(&authorID, 0)
}

//GetAnnounceByID returns announce written by gicen announce id
func (db *AnnounceDB) GetAnnounceByID(id primitive.ObjectID) (bson.M, error) {
	result := []bson.M{}
	matchStage := bson.D{{"$match", bson.D{{"_id", id}}}}
	aggregateStages := bson.D{{"$project", bson.D{{"id", "$_id"},
		{"_id", 0},
		{"title", "$title"},
		{"author", "$author"},
		{"description", "$description"},
		{"created", "$created"},
		{"modified", "$modified"},
		{"target_date", "$target_date"}}}}
	limitStage := bson.D{{"$limit", 1}}
	cursor, err := db.collection.Aggregate(timeoutContext(), mongo.Pipeline{matchStage, aggregateStages, limitStage})
	if err != nil {
		log.Fatal(err.Error())
		return bson.M{}, err
	}
	err = cursor.All(timeoutContext(), &result)
	if err != nil {
		return bson.M{}, err
	}
	if len(result) != 1 {
		return bson.M{}, errors.New("Not found announce")
	}
	return result[0], nil
}

//AddAnnounce add announce into database
func (db *AnnounceDB) AddAnnounce(authorID primitive.ObjectID, announceMap bson.M) (primitive.ObjectID, error) {

	announce := models.Announce{
		AuthorID:    authorID,
		Title:       announceMap["title"].(string),
		Description: announceMap["description"].(string),
		TargetDate:  announceMap["target_date"].(time.Time),
		Created:     time.Now(),
		Modified:    time.Now(),
	}
	res, err := db.collection.InsertOne(timeoutContext(), announce)
	return res.InsertedID.(primitive.ObjectID), err
}

//UpdateAnnounceByID update announce by announce id
func (db *AnnounceDB) UpdateAnnounceByID(id primitive.ObjectID, userID primitive.ObjectID, setElement *bson.D) error {
	if len(*setElement) > 0 {
		*setElement = append(*setElement, bson.E{"modified", time.Now()})
	}
	setMap := bson.D{
		{"$set", *setElement},
	}
	_, err := db.collection.UpdateOne(timeoutContext(), bson.M{"_id": id}, setMap)
	return err
}

//DeleteAnnounceByID delete announce by announce id
func (db *AnnounceDB) DeleteAnnounceByID(id primitive.ObjectID, userID primitive.ObjectID) error {
	_, err := db.collection.DeleteOne(timeoutContext(), bson.M{"$and": []bson.M{{"_id": id}, {"author": userID}}})
	return err
}
