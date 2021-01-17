package v1

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func GetStateHistory(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	list, _ := database.Instance.DFDHistory.GetList(id)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"history": list}))
}

func GetOwnStateHistory(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, _ := database.Instance.DFDHistory.GetList(user.Id)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"history": list}))
}

func PostState(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	if bodyMap["state"] == nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("state field is required"))
		return
	}

	stateStromg, ok := bodyMap["state"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("state field must be string"))
		return
	}
	database.Instance.DFDHistory.Push(user.Id, stateStromg)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}
