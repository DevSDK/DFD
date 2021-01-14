package v1

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GetAnnounceList(c *gin.Context) {
	idString := c.Param("id")
	if idString == "" {
		list, _ := database.Instance.Announce.GetList()
		c.JSON(200, gin.H{"message": "success", "list": list})
	} else {
		id, _ := primitive.ObjectIDFromHex(idString)
		list, _ := database.Instance.Announce.GetListByAuthorId(id)
		c.JSON(200, gin.H{"message": "success", "list": list})
	}
}

func GetAnnounceListMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, _ := database.Instance.Announce.GetListByAuthorId(user.Id)
	c.JSON(200, gin.H{"message": "success", "list": list})
}

func PostAnnounce(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	if bodyMap["title"] == nil {
		c.JSON(400, gin.H{"message": "title filed is required"})
		return
	}
	if bodyMap["target_date"] == nil {
		c.JSON(400, gin.H{"message": "target_date filed is required"})
		return
	} else {
		t, err := time.Parse(time.RFC3339, bodyMap["target_date"].(string))
		if err != nil {
			c.JSON(400, gin.H{"message": "wrong time format"})
			return
		}
		bodyMap["target_date"] = t
	}
	if bodyMap["description"] == nil {
		bodyMap["description"] = ""
	}

	id, _ := database.Instance.Announce.AddAnnounce(user.Id, bodyMap)
	c.JSON(200, gin.H{"message": "success", "id": id})
}

func PutAnnounce(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	setElement := &bson.D{}

	utils.ApplySetElementStringByName(setElement, bodyMap, "title")
	utils.ApplySetElementStringByName(setElement, bodyMap, "description")

	if bodyMap["target_date"] != nil {
		_, err := time.Parse(time.RFC3339, bodyMap["target_date"].(string))
		if err != nil {
			c.JSON(400, gin.H{"message": "time field error"})
			return
		}
	}
	utils.ApplySetElementStringByName(setElement, bodyMap, "target_date")
	if err := database.Instance.Announce.UpdateAnnounceById(id, user.Id, setElement); err != nil {
		c.JSON(400, gin.H{"message": "cannot update"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func GetAnnounce(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	res, err := database.Instance.Announce.GetAnnounceById(id)
	if err != nil {
		c.JSON(400, gin.H{"message": "cannot found announce"})
		return
	}
	c.JSON(200, gin.H{"message": "success", "announce": res})
}

func DelAnnounce(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	user := c.MustGet("user").(models.User)
	if err := database.Instance.Announce.DeleteAnnounceById(id, user.Id); err != nil {
		c.JSON(400, gin.H{"message": "cannot delete announce"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
