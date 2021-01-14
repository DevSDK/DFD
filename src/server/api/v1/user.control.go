package v1

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createCommonUserMap(user models.User) gin.H {
	recent, err := database.Instance.DFDHistory.GetRecent(user.Id)
	state := ""
	if err == nil {
		state = recent["state"].(string)
	}
	return gin.H{
		"id":         user.Id.Hex(),
		"email":      user.Email,
		"state":      state,
		"lol_name":   user.LolUsername,
		"discord_id": user.DiscordId,
		"username":   user.Username,
		"role":       user.Role,
	}
}

func GetUser(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	user, err := database.Instance.User.FindById(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Cannot found user by id"})
		return
	}
	response := createCommonUserMap(user)
	c.JSON(200, response)
}

func GetMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	response := createCommonUserMap(user)
	c.JSON(200, response)
}
func PutUser(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	setElement := &bson.D{}
	utils.ApplySetElementStringByName(setElement, bodyMap, "username")
	if err := database.Instance.User.UpdateById(user.Id, setElement); err != nil {
		c.JSON(501, gin.H{"message": "update failed "})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
