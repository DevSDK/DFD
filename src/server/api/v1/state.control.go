package v1

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetStateHistory(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	list, _ := database.Instance.DFDHistory.GetList(id)
	c.JSON(200, gin.H{"message": "success", "states": list})
}

func GetOwnStateHistory(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, _ := database.Instance.DFDHistory.GetList(user.Id)
	c.JSON(200, gin.H{"message": "success", "states": list})
}

func PostState(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	if bodyMap["state"] == nil {
		c.JSON(400, gin.H{"message": "wrong body request"})
		return
	}

	database.Instance.DFDHistory.Push(user.Id, bodyMap["state"].(string))
	c.JSON(200, gin.H{"message": "success"})
}
