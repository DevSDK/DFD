package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type ApplicationTokenDB struct{ BaseDB }

func (db *ApplicationTokenDB) Add() primitive.ObjectID {
	res, err := db.collection.InsertOne(timeoutContext(), bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	return res.InsertedID.(primitive.ObjectID)
}

func (db *ApplicationTokenDB) Exist(id primitive.ObjectID) bool {
	if db.collection.FindOne(timeoutContext(), bson.M{"_id": id}).Err() != nil {
		return false
	}
	return true
}
