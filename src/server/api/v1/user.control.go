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

func createCommonUserMap(user models.User) gin.H {
	recent, err := database.Instance.DFDHistory.GetRecent(user.Id)
	state := ""
	if err == nil {
		state = recent["state"].(string)
	}
	return gin.H{
		"id":            user.Id.Hex(),
		"email":         user.Email,
		"state":         state,
		"lol_name":      user.LolUsername,
		"discord_id":    user.DiscordId,
		"username":      user.Username,
		"role":          user.Role,
		"profile_image": user.ProfileImage,
		"created":       user.Created,
	}
}

// @Summary Get User Information
// @Description Get user by user id
// @Description Permission : **user.get**
// @ID get-string-by-string
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} docmodels.ResponseSuccess{user=models.User} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/user
// @Router /api/v1/user/{id} [get]
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

// @Summary Get My User Information
// @Description Get me
// @Description Permission : **user.get**
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{user=models.User} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/user
// @Router /api/v1/user [get]
func GetMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	result := createCommonUserMap(user)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"user": result}))
}

// @Summary Edit user information
// @Description edit userfield. For now, only username can could be patched
// @Description Permission : **user.patch**
// @Accept  json
// @Produce  json
// @Param username body docmodels.RequestEmpty{username=string} true "username"
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/user
// @Router /api/v1/user [patch]
func PatchUser(c *gin.Context) {
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

// @Summary Patch LOL Information
// @Description Update LOL user informations
// @Description Permission : **user.patch**
// @Accept  json
// @Produce  json
// @Param lol_username body docmodels.RequestEmpty{lol_username=string} true "league of legends username"
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/user
// @Router /api/v1/user/lol [patch]
func PatchUserLolName(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	user := c.MustGet("user").(models.User)
	setElement := &bson.D{}

	if bodyMap["lol_username"] != nil {
		username, ok := bodyMap["lol_username"].(string)

		if !ok {
			c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("lol_username must be string"))
			return
		}

		respMap, respCode := utils.RequestToRiotServer("/lol/summoner/v4/summoners/by-name/"+username, nil)

		if respCode != 200 {
			c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
			return
		}

		utils.ApplySetElementString(setElement, respMap, "name", "lol_username")
		utils.ApplySetElementString(setElement, respMap, "id", "lol_id")
		utils.ApplySetElementString(setElement, respMap, "accountId", "lol_account_id")
		utils.ApplySetElementString(setElement, respMap, "puuid", "lol_puu_id")
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
