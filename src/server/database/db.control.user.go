package database

import (
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

//UserDB data structure for user
type UserDB struct {
	BaseDB
}

//Register user into user db
func (db *UserDB) Register(userMap bson.M) (models.User, error) {
	tokenString := userMap["tokenString"].(string)
	user := models.User{
		DiscordID:    userMap["id"].(string),
		Username:     userMap["username"].(string),
		Email:        userMap["email"].(string),
		Role:         "guest",
		RefreshToken: string(tokenString),
		Created:      time.Now(),
		Modified:     time.Now(),
	}
	result, err := db.collection.InsertOne(timeoutContext(), user)
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, err
}

//FindByID from user db
func (db *UserDB) FindByID(id primitive.ObjectID) (models.User, error) {
	user := models.User{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

//FindByDiscordID from user db
func (db *UserDB) FindByDiscordID(id string) (models.User, error) {
	user := models.User{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"discord_id": id}).Decode(&user)
	return user, err
}

//FindByEmail from user db
func (db *UserDB) FindByEmail(email string) (models.User, error) {
	user := models.User{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"email": email}).Decode(&user)
	return user, err
}

//UpdateByID from user db
func (db *UserDB) UpdateByID(id primitive.ObjectID, setElement *bson.D) error {
	if len(*setElement) > 0 {
		*setElement = append(*setElement, bson.E{"modified", time.Now()})
	}
	setMap := bson.D{
		{"$set", *setElement},
	}
	_, err := db.collection.UpdateOne(timeoutContext(), bson.M{"_id": id}, setMap)
	return err
}

//GetLoLInfoList get all non guest user league of legends user information from user db
func (db *UserDB) GetLoLInfoList() []bson.M {
	aggregateStage := bson.D{{"$project", bson.D{{"id", "$_id"},
		{"lol_username", "$lol_username"},
		{"lol_account_id", "$lol_account_id"},
		{"lol_puu_id", "$lol_puu_id"},
		{"lol_id", "$lol_id"}}}}
	matchStage := bson.D{{"$match", bson.D{{"role", bson.D{{"$ne", "guest"}}}}}}
	cursor, err := db.collection.Aggregate(timeoutContext(), mongo.Pipeline{matchStage, aggregateStage})
	if err != nil {
		log.Print(err.Error())
	}
	result := []bson.M{}
	if err := cursor.All(timeoutContext(), &result); err != nil {
		log.Print(err.Error())
	}
	return result
}

//GetUserList from get all non guest user
func (db *UserDB) GetUserList() []primitive.ObjectID {
	aggregateStage := bson.D{{"$project", bson.D{{"id", "$_id"}}}}
	matchStage := bson.D{{"$match", bson.D{{"role", bson.D{{"$ne", "guest"}}}}}}
	cursor, err := db.collection.Aggregate(timeoutContext(), mongo.Pipeline{matchStage, aggregateStage})
	resultFromDB := []bson.M{}
	result := []primitive.ObjectID{}
	if err != nil {
		log.Print(err.Error())
		return result
	}

	if err := cursor.All(timeoutContext(), &resultFromDB); err != nil {
		log.Print(err.Error())
		return result
	}

	for _, v := range resultFromDB {
		result = append(result, v["id"].(primitive.ObjectID))
	}
	return result
}
