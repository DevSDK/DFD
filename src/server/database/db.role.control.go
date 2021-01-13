package database

import (
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleDB struct{}

func (c *RoleDB) FindByName(name string) (models.Role, error) {

	role := models.Role{}
	err := Instance.database.Collection("Role").
		FindOne(timeoutContext(), bson.M{"name": name}).Decode(&role)
	return role, err
}

func (c *RoleDB) AddRole(role models.Role) error {
	_, err := Instance.database.Collection("Role").
		InsertOne(timeoutContext(), role)
	return err
}
