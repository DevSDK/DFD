package database

import (
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (c *DBInstance) RegisterUser(userMap map[string]interface{}) (models.User, error) {
	tokenString := userMap["tokenString"].(string)
	user := models.User{
		DiscordId:    userMap["id"].(string),
		Username:     userMap["username"].(string),
		Email:        userMap["email"].(string),
		RefreshToken: string(tokenString),
		Role:         "Pending",
		Created:      time.Now(),
		Modified:     time.Now(),
	}
	userCollection := Instance.database.Collection("User")
	_, err := userCollection.InsertOne(timeoutContext(), user)
	return user, err
}

func (c *DBInstance) FindUserById(id string) (models.User, error) {
	userCollection := Instance.database.Collection("User")
	user := models.User{}
	userId, _ := primitive.ObjectIDFromHex(id)
	err := userCollection.FindOne(timeoutContext(), bson.M{"_id": userId}).Decode(&user)
	return user, err
}

func (c *DBInstance) FindUserByDiscordId(id string) (models.User, error) {
	userCollection := Instance.database.Collection("User")
	user := models.User{}
	err := userCollection.FindOne(timeoutContext(), bson.M{"discordid": id}).Decode(&user)
	return user, err
}

func (c *DBInstance) FindUserByEmail(email string) (models.User, error) {
	userCollection := Instance.database.Collection("User")
	user := models.User{}
	err := userCollection.FindOne(timeoutContext(), bson.M{"email": email}).Decode(&user)
	return user, err
}
