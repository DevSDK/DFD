package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

//ApplicationAuth is model structure for applications
type ApplicationAuth struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

//User is model structure for user
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProfileImage primitive.ObjectID `json:"profile_image" bson:"profile_image_id"`
	DiscordID    string             `json:"discord_id" bson:"discord_id,unique"`
	Username     string             `bson:"username,omitempty"`
	Email        string             `bson:"email,omitempty,unique" format:"email"`
	LolID        string             `bson:"lol_id" swaggerignore:"true"`
	LolPuuID     string             `bson:"lol_puu_id" swaggerignore:"true"`
	LolAccountID string             `bson:"lol_account_id" swaggerignore:"true"`
	LolUsername  string             `bson:"lol_username"`
	RefreshToken string             `bson:"refresh,omitempty" swaggerignore:"true"`
	Role         string             `bson:"role"`
	Created      time.Time          `bson:"created" format:"date-time"`
	Modified     time.Time          `bson:"modified" format:"date-time"`
}

//Role is model structure for role
type Role struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `bson:"name,unique"`
	Description string             `bson:"description"`
	Permissions []string           `bson:"permissions"`
}

//LOLHistory is model structure for lol history
type LOLHistory struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	GameID       string             `bson:"gameid,unique"`
	Game         bson.M             `bson:game`
	Win          bool               `bson:win`
	Timestamp    time.Time          `bson:timestamp`
	QueueID      int64              `bson:queueid`
	Participants []string           `bson:participants`
	Created      time.Time          `bson:"created"`
}

//DFDHistory is model structure for dfd status history
type DFDHistory struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"user_id" swaggerignore:"true"`
	State   string             `bson:"state"`
	Was     string             `bson:"was"`
	Created time.Time          `bson:"created"`
}

//Announce is model structure for announce
type Announce struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"6006bc289b9bdd2381263063"`
	AuthorID    primitive.ObjectID `bson:"author" example:"5c06bc289b9bdd2381263063"`
	Title       string             `bson:"title" example:"awesome title"`
	Description string             `bson:"description" example:"this project is awesome!"`
	TargetDate  time.Time          `bson:"target_date" example:"2021-01-19T11:01:19+00:00"`

	Created  time.Time `bson:"created example:"2021-01-19T11:01:19+00:00"`
	Modified time.Time `bson:modified example:"2021-01-19T11:01:19+00:00"`
}
