package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

type ApplicationAuth struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProfileImage primitive.ObjectID `bson:"profile_image_id"`
	DiscordId    string             `bson:"discord_id",unique`
	Username     string             `bson:"username,omitempty"`
	Email        string             `bson:"email,omitempty,unique"`
	LolId        string             `bson:"lol_id"`
	LolPuuId     string             `bson:"lol_puu_id"`
	LolAccountId string             `bson:"lol_account_id"`
	LolUsername  string             `bson:"lol_username"`
	RefreshToken string             `bson:"refresh,omitempty"`
	Role         string             `bson:"role"`
	Created      time.Time          `bson:"created"`
	Modified     time.Time          `bson:modified`
}

type Role struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `bson:"name,unique"`
	Description string             `bson:"description"`
	Permissions []string           `bson:"permissions"`
}

type LOLHistory struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Game      bson.M             `bson:game`
	Win       bool               `bson:win`
	Timestamp int64              `bson:timestamp`
	Created   time.Time          `bson:"created"`
}

type DFDHistory struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId  primitive.ObjectID `bson:"user_id"`
	State   string             `bson:"state"`
	Was     string             `bson:"was"`
	Created time.Time          `bson:"created"`
}

type Announce struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthorId    primitive.ObjectID `bson:"author"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	TargetDate  time.Time          `bson:"target_date"`

	Created  time.Time `bson:"created"`
	Modified time.Time `bson:modified`
}
