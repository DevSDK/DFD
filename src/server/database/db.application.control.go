package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

//ApplicationTokenDB is a structure for application token data access
type ApplicationTokenDB struct{ BaseDB }

//Add fucntion creates application token with _id
func (db *ApplicationTokenDB) Add() primitive.ObjectID {
	res, err := db.collection.InsertOne(timeoutContext(), bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	return res.InsertedID.(primitive.ObjectID)
}

//Exist function checks given token existed
func (db *ApplicationTokenDB) Exist(id primitive.ObjectID) bool {
	if db.collection.FindOne(timeoutContext(), bson.M{"_id": id}).Err() != nil {
		return false
	}
	return true
}
