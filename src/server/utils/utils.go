package utils

import "go.mongodb.org/mongo-driver/bson"

func ApplySetElementStringByName(setElement *bson.D, updateMap map[string]interface{}, target string) bool {
	if updateMap[target] != nil {
		*setElement = append(*setElement, bson.E{Key: target, Value: updateMap[target]})
		return true
	}
	return false
}
