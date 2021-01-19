package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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
	req.Header.Set("X-Riot-Token", os.Getenv("RIOT_API_ACCESS"))
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

func ApplySetElementStringSameTarget(setElement *bson.D, updateMap bson.M, target string) bool {
	return ApplySetElementString(setElement, updateMap, target, target)
}

func ApplySetElementString(setElement *bson.D, updateMap bson.M, updateMapTarget string, dbTarget string) bool {
	if updateMap[updateMapTarget] != nil {
		*setElement = append(*setElement, bson.E{Key: dbTarget, Value: updateMap[updateMapTarget]})
		return true
	}
	return false
}

func CreateSuccessJSONMessage(data gin.H) gin.H {
	result := gin.H{"message": "success", "status": http.StatusOK}
	if data != nil {
		for k, v := range data {
			result[k] = v
		}
	}
	return result
}

func CreateBadRequestJSONMessage(message string) gin.H {
	result := gin.H{"message": message, "status": http.StatusBadRequest}
	return result
}

func CreateForbbidnJSONMessage(message string) gin.H {
	result := gin.H{"message": message, "status": http.StatusForbidden}
	return result
}

func CreateUnauthorizedJSONMessage(message string, isExpired bool) gin.H {
	result := gin.H{"message": message, "status": http.StatusUnauthorized, "token_expired": isExpired}
	return result
}

func CreateNotFoundJSONMessage(message string) gin.H {
	result := gin.H{"message": message, "status": http.StatusNotFound}
	return result
}

func CreateInternalServerErrorJSONMessage() gin.H {
	result := gin.H{"message": "Internal Server Error", "status": http.StatusInternalServerError}
	return result
}
