package middleware

import (
	"encoding/json"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"os"
)

func JsonParseMiddleware(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("body is required"))
		c.Abort()
		return
	}
	var bodyMap bson.M
	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("body must be JSON format"))
		c.Abort()
		return
	}
	c.Set("bodymap", bodyMap)
	c.Next()
}

func VerifyApplicationTokenMiddleware(c *gin.Context) {
	tokenString := c.Request.Header["X-Dfd-App-Auth"]
	if tokenString == nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("App-Auth token required"))
		c.Abort()
		return
	}
	token, _ := primitive.ObjectIDFromHex(tokenString[0])
	if !database.Instance.ApplicationToken.Exist(token) {c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("invalid token", false))
		c.Abort()
		return
	}
	c.Next()
}

func CORSMiddleware(c *gin.Context) {
	if os.Getenv("GIN_MODE") == "release" {
		c.Header("Access-Control-Allow-Origin", "https://devsdk.net/dfd")
	} else {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	}
	
	c.Header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS");
	c.Header("Access-Control-Allow-Headers","*")
	c.Header("Access-Control-Allow-Credentials","true")
	if c.Request.Method == "OPTIONS" {
		c.JSON(http.StatusOK,"")
		return
	}
	return
	c.Next()
}