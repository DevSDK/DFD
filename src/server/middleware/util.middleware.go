package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
)

func JsonParseMiddleware(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"message": "wrong body request"})
		c.Abort()
		return
	}
	var bodyMap bson.M
	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		c.JSON(400, gin.H{"message": "wrong body request"})
		c.Abort()
		return
	}

	c.Set("bodymap", bodyMap)
	c.Next()
}
