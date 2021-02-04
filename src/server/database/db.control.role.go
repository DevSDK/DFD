package database

import (
	"github.com/DevSDK/DFD/src/server/database/models"
	"go.mongodb.org/mongo-driver/bson"
)

//RoleDB data structure for role
type RoleDB struct {
	BaseDB
}

//FindByName from role db
func (db *RoleDB) FindByName(name string) (models.Role, error) {

	role := models.Role{}
	err := db.collection.FindOne(timeoutContext(),
		bson.M{"name": name}).Decode(&role)
	return role, err
}

//AddRole into role db
func (db *RoleDB) AddRole(role models.Role) error {
	_, err := db.collection.
		InsertOne(timeoutContext(), role)
	return err
}
