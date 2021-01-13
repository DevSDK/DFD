package database

import (
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserDB struct{}

func (c *UserDB) Register(userMap map[string]interface{}) (models.User, error) {
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
	userCollection := Instance.database.Collection("User")
	result, err := userCollection.InsertOne(timeoutContext(), user)
	user.Id = result.InsertedID.(primitive.ObjectID)
	return user, err
}

func (c *UserDB) FindById(id primitive.ObjectID) (models.User, error) {
	userCollection := Instance.database.Collection("User")
	user := models.User{}
	err := userCollection.FindOne(timeoutContext(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (c *UserDB) FindByDiscordId(id string) (models.User, error) {
	userCollection := Instance.database.Collection("User")
	user := models.User{}
	err := userCollection.FindOne(timeoutContext(), bson.M{"discord_id": id}).Decode(&user)
	return user, err
}

func (c *UserDB) FindByEmail(email string) (models.User, error) {
	userCollection := Instance.database.Collection("User")
	user := models.User{}
	err := userCollection.FindOne(timeoutContext(), bson.M{"email": email}).Decode(&user)
	return user, err
}

func (c *UserDB) UpdateById(id primitive.ObjectID, setElement *bson.D) error {

	userCollection := Instance.database.Collection("User")
	if len(*setElement) > 0 {
		*setElement = append(*setElement, bson.E{"modified", time.Now()})
	}
	setMap := bson.D{
		{"$set", *setElement},
	}
	_, err := userCollection.UpdateOne(timeoutContext(), bson.M{"_id": id}, setMap)
	return err
}
