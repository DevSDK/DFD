package utils

import (
	"encoding/json"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//RequestToRiotServer function request to riot server
func RequestToRiotServer(endpoint string, params bson.M) (bson.M, int) {
	client := &http.Client{}
	RIOT_URL := os.Getenv("RIOT_API_URI")
	requestURI := endpoint
	if params != nil {
		requestURI += "?"
		for k, v := range params {
			requestURI += k + "=" + v.(string) + "&"
		}
	}

	req, err := http.NewRequest("GET", RIOT_URL+requestURI, nil)
	if err != nil {
		log.Print(err.Error())
	}
	riotAccess, _ := database.Instance.Redis.Get("riot-access-token")

	req.Header.Set("X-Riot-Token", riotAccess)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	resultMap := bson.M{}
	responseString, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Print(string(responseString))
	}
	if err := json.Unmarshal(responseString, &resultMap); err != nil {
		log.Fatal(err.Error())
	}

	return resultMap, resp.StatusCode
}

//ApplySetElementStringSameTarget function to create database set field for same body key and target
func ApplySetElementStringSameTarget(setElement *bson.D, updateMap bson.M, target string) bool {
	return ApplySetElementString(setElement, updateMap, target, target)
}

//ApplySetElementString function to create database set field for different body key and target
func ApplySetElementString(setElement *bson.D, updateMap bson.M, updateMapTarget string, dbTarget string) bool {
	if updateMap[updateMapTarget] != nil {
		*setElement = append(*setElement, bson.E{Key: dbTarget, Value: updateMap[updateMapTarget]})
		return true
	}
	return false
}

//CreateSuccessJSONMessage generates success response object
func CreateSuccessJSONMessage(data gin.H) gin.H {
	result := gin.H{"message": "success", "status": http.StatusOK}
	if data != nil {
		for k, v := range data {
			result[k] = v
		}
	}
	return result
}

//CreateBadRequestJSONMessage generates bad request response object
func CreateBadRequestJSONMessage(message string) gin.H {
	result := gin.H{"message": message, "status": http.StatusBadRequest}
	return result
}

//CreateForbbidnJSONMessage generates forbbiden response object
func CreateForbbidnJSONMessage(message string) gin.H {
	result := gin.H{"message": message, "status": http.StatusForbidden}
	return result
}

//CreateUnauthorizedJSONMessage generates unauthorized response object
func CreateUnauthorizedJSONMessage(message string, isExpired bool) gin.H {
	result := gin.H{"message": message, "status": http.StatusUnauthorized, "token_expired": isExpired}
	return result
}

//CreateNotFoundJSONMessage generates not found response object
func CreateNotFoundJSONMessage(message string) gin.H {
	result := gin.H{"message": message, "status": http.StatusNotFound}
	return result
}

//CreateInternalServerErrorJSONMessage generates server error response object
func CreateInternalServerErrorJSONMessage() gin.H {
	result := gin.H{"message": "Internal Server Error", "status": http.StatusInternalServerError}
	return result
}
