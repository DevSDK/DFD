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
	"log"
	"net/http"
	"os"
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
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Cannot found user"))
		return
	}
	result := createCommonUserMap(user)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"user": result}))
}

func GetMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	result := createCommonUserMap(user)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"user": result}))
}
func PutUser(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	setElement := &bson.D{}
	utils.ApplySetElementStringSameTarget(setElement, bodyMap, "username")
	if err := database.Instance.User.UpdateById(user.Id, setElement); err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Cannot found user"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}

func PutUserLolName(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	setElement := &bson.D{}

	if bodyMap["lol_username"] != nil {
		RIOT_URL := os.Getenv("RIOT_API_URI")
		client := &http.Client{}
		req, err := http.NewRequest("GET", RIOT_URL+"/lol/summoner/v4/summoners/by-name/"+bodyMap["lol_username"].(string), nil)
		if err != nil {
			log.Print(err.Error())
			c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		}
		req.Header.Set("X-Riot-Token", os.Getenv("RIOT_API_ACCESS"))
		resp, err := client.Do(req)
		if err != nil {
			log.Print(err.Error())
			c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		}
		responseString, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			log.Print("RIOT response is " + string(responseString))
			c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
			return
		}
		var responseMap bson.M
		if err := json.Unmarshal(responseString, &responseMap); err != nil {
			log.Print(err.Error())
			c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
			return
		}
		utils.ApplySetElementString(setElement, responseMap, "name", "lol_username")
		utils.ApplySetElementString(setElement, responseMap, "id", "lol_id")
		utils.ApplySetElementString(setElement, responseMap, "accountId", "lol_account_id")
		utils.ApplySetElementString(setElement, responseMap, "puuid", "lol_puu_id")
	} else {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("lol_username filed is required"))
		return
	}

	if err := database.Instance.User.UpdateById(user.Id, setElement); err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Cannot found user"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}
