package v1

import (
	_ "github.com/DevSDK/DFD/src/server/api/v1/docmodels"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

// @Summary Get current announce list
// @Description Get announce list greater then target-date from NOW sorted by creation time
// @Description Permission : **announces.get**
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{list=[]models.Announce} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/announce
// @Router /v1/announces/current [get]
func GetCurrentAnnounceList(c *gin.Context) {
	list, _ := database.Instance.Announce.GetListWithTimestamp(time.Now().Unix())
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"list": list}))
}

// @Summary Get all announce list
// @Description Get all announce list sorted by creation time
// @Description Permission : **announces.get**
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{list=[]models.Announce} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/announce
// @Router /v1/announces/all [get]
func GetAllAnnounceList(c *gin.Context) {
	list, _ := database.Instance.Announce.GetList()
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"list": list}))
}

// @Summary Get announce list by user id
// @Description Get announce list written by user
// @Description Permission : **announces.get**
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} docmodels.ResponseSuccess{list=[]models.Announce} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/announce
// @Router /v1/announces/user/{id} [get]
func GetAnnounceList(c *gin.Context) {
	idString := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idString)
	list, _ := database.Instance.Announce.GetListByAuthorId(id)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"list": list}))
}

// @Summary Get own announce list
// @Description get announce list written by me
// @Description Permission : **announces.get**
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{list=[]models.Announce} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/announce
// @Router /v1/announces/me [get]
func GetAnnounceListMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, _ := database.Instance.Announce.GetListByAuthorId(user.Id)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"list": list}))
}

// @Summary Write announce
// @Description write announce
// @Description Permission : **announce.post**
// @Accept  json
// @Produce  json
// @Param bodymap body docmodels.RequestBodyAnnouncePost true "username"
// @Success 200 {object} docmodels.ResponseSuccess{id=primitive.ObjectId} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/announce
// @Router /v1/announce [post]
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

// @Summary Patch Announce
// @Description Edit spcific announce written by me
// @Description Permission : **announce.get**
// @Accept  json
// @Produce  json
// @Param id path string true "Announce ID"
// @Param body body docmodels.RequestBodyAnnouncePost true "body"
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/announce
// @Router /v1/announce/{id} [patch]
func PatchAnnounce(c *gin.Context) {
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
		utils.ApplySetElementStringSameTarget(setElement, bodyMap, "target_date")
	}
	if err := database.Instance.Announce.UpdateAnnounceById(id, user.Id, setElement); err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("cannot found announce"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))

}

// @Summary Get announce
// @Description Get specific announce by announce id
// @Description Permission : **announce.get**
// @Accept  json
// @Produce  json
// @Param id path string true "Announce ID"
// @Success 200 {object} docmodels.ResponseSuccess{announce=models.Announce} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/announce
// @Router /v1/announce/{id} [get]
func GetAnnounce(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	res, err := database.Instance.Announce.GetAnnounceById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("cannot found announce"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"announce": res}))
}

// @Summary Delete announce
// @Description Delete announce by announce id
// @Description Permission : **announce.delete**
// @Accept  json
// @Produce  json
// @Param id path string true "Announce ID"
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/announce
// @Router /v1/announce/{id} [delete]
func DelAnnounce(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	user := c.MustGet("user").(models.User)
	if err := database.Instance.Announce.DeleteAnnounceById(id, user.Id); err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("cannot found announce"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}
