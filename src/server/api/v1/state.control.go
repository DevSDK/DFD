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

// GetStateHistory is handler for endpoint GET /states/{id}
// @Summary Get state history
// @Description Get user's state change history
// @Description Permission : **states.get**
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} docmodels.ResponseSuccess{history=[]models.DFDHistory} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/state
// @Router /v1/states/{id} [get]
func GetStateHistory(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	list, _ := database.Instance.DFDHistory.GetList(id)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"history": list}))
}

// GetOwnStateHistory is handler for endpoint GET /states
// @Summary Get my state history
// @Description Get my state change history
// @Description Permission : **states.get**
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{history=[]models.DFDHistory} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/state
// @Router /v1/states [get]
func GetOwnStateHistory(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, _ := database.Instance.DFDHistory.GetList(user.ID)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"history": list}))
}

// PostState is handler for endpoint POST /state
// @Description Create state
// @Description Permission : **states.post**
// @Accept  json
// @Produce  json
// @Param state body docmodels.RequestBodyStatePost true "Create user state"
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/state
// @Router /v1/state [post]
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
	database.Instance.DFDHistory.Push(user.ID, stateStromg)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}
