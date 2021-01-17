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
	bodyMap := c.MustGet("bodymap").(bson.M)
	tokenString, ok := bodyMap["token"].(string)
	if bodyMap["token"] == nil || !ok {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("token filed required"))
		c.Abort()
		return
	}
	token, _ := primitive.ObjectIDFromHex(tokenString)
	if !database.Instance.ApplicationToken.Exist(token) {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("invalid token"))
		return
		c.Abort()
	}
	c.Next()
}
