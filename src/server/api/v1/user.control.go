package v1

import (
	"encoding/json"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
)

func GetUser(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	user, err := database.Instance.User.FindById(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Cannot found user by id"})
		return
	}
	c.JSON(200, gin.H{
		"id":         user.Id.Hex(),
		"state":      user.State,
		"email":      user.Email,
		"lol_name":   user.LolUsername,
		"discord_id": user.DiscordId,
		"username":   user.Username,
		"role":       user.Role,
	})
}

func GetMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	c.JSON(200, gin.H{
		"id":         user.Id.Hex(),
		"state":      user.State,
		"email":      user.Email,
		"lol_name":   user.LolUsername,
		"discord_id": user.DiscordId,
		"username":   user.Username,
		"role":       user.Role,
	})
}

func PutUser(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"message": "wrong body request"})
		return
	}
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		c.JSON(400, gin.H{"message": "wrong body request"})
		return
	}
	user := c.MustGet("user").(models.User)
	setElement := &bson.D{}
	utils.ApplySetElementStringByName(setElement, bodyMap, "state")
	utils.ApplySetElementStringByName(setElement, bodyMap, "username")
	if err := database.Instance.User.UpdateById(user.Id, setElement); err != nil {
		c.JSON(501, gin.H{"message": "update failed "})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
