package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty,unique"`
	Email    string             `bson:"email,omitempty,unique"`
	State    string             `bson:"state,omitempty"`
	Token    string             `bson:"token,omitempty"`
	Role     string             `bson:"role,omitempty"`

	Created  time.Time `json:"created"`
	Modified time.Time `json:Modified`
}
