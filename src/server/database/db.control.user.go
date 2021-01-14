package database

import (
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserDB struct {
	BaseDB
}

func (db *UserDB) Register(userMap map[string]interface{}) (models.User, error) {
	tokenString := userMap["tokenString"].(string)
	user := models.User{
		DiscordId:    userMap["id"].(string),
		Username:     userMap["username"].(string),
		Email:        userMap["email"].(string),
		Role:         "guest",
		RefreshToken: string(tokenString),
		Created:      time.Now(),
		Modified:     time.Now(),
	}
	result, err := db.collection.InsertOne(timeoutContext(), user)
	user.Id = result.InsertedID.(primitive.ObjectID)
	return user, err
}

func (db *UserDB) FindById(id primitive.ObjectID) (models.User, error) {
	user := models.User{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (db *UserDB) FindByDiscordId(id string) (models.User, error) {
	user := models.User{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"discord_id": id}).Decode(&user)
	return user, err
}

func (db *UserDB) FindByEmail(email string) (models.User, error) {
	user := models.User{}
	err := db.collection.FindOne(timeoutContext(), bson.M{"email": email}).Decode(&user)
	return user, err
}

func (db *UserDB) UpdateById(id primitive.ObjectID, setElement *bson.D) error {
	if len(*setElement) > 0 {
		*setElement = append(*setElement, bson.E{"modified", time.Now()})
	}
	setMap := bson.D{
		{"$set", *setElement},
	}
	_, err := db.collection.UpdateOne(timeoutContext(), bson.M{"_id": id}, setMap)
	return err
}
