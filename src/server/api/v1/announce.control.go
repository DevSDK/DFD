package v1

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func GetAnnounceList(c *gin.Context) {
	idString := c.Param("id")
	list := []bson.M{}
	if idString == "" {
		list, _ = database.Instance.Announce.GetList()
	} else {
		id, _ := primitive.ObjectIDFromHex(idString)
		list, _ = database.Instance.Announce.GetListByAuthorId(id)
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"list": list}))
}

func GetAnnounceListMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, _ := database.Instance.Announce.GetListByAuthorId(user.Id)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"list": list}))
}

func PostAnnounce(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	if bodyMap["title"] == nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("title field is required"))
		return
	}
	if bodyMap["target_date"] == nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("target_date filed is required"))
		return
	}

	dateString, ok := bodyMap["target_date"].(string)

	if !ok {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("time field format is invalid"))
		return
	}
	t, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("time field format is invalid"))
		return
	}
	bodyMap["target_date"] = t

	if bodyMap["description"] == nil {
		bodyMap["description"] = ""
	}

	id, _ := database.Instance.Announce.AddAnnounce(user.Id, bodyMap)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"id": id}))
}

func PutAnnounce(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	setElement := &bson.D{}

	utils.ApplySetElementStringSameTarget(setElement, bodyMap, "title")
	utils.ApplySetElementStringSameTarget(setElement, bodyMap, "description")

	if bodyMap["target_date"] != nil {
		dateString, ok := bodyMap["target_date"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("time field format is invalid"))
			return
		}
		_, err := time.Parse(time.RFC3339, dateString)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("time field format is invalid"))
			return
		}
	}
	utils.ApplySetElementStringSameTarget(setElement, bodyMap, "target_date")
	if err := database.Instance.Announce.UpdateAnnounceById(id, user.Id, setElement); err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("cannot found announce"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}

func GetAnnounce(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	res, err := database.Instance.Announce.GetAnnounceById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("cannot found announce"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"announce": res}))
}

func DelAnnounce(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	user := c.MustGet("user").(models.User)
	if err := database.Instance.Announce.DeleteAnnounceById(id, user.Id); err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("cannot found announce"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}
