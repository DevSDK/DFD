package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DiscordId    string             `bson: discord_id,unique`
	Username     string             `bson:"username,omitempty"`
	Email        string             `bson:"email,omitempty,unique"`
	State        string             `bson:"state,omitempty"`
	RefreshToken string             `bson:"refresh,omitempty"`
	Role         string             `bson:"role,omitempty"`

	Created  time.Time `json:"created"`
	Modified time.Time `json:Modified`
}
