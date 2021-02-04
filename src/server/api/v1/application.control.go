package v1

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

// PostAppicationToken is handler for endpoint POST /application/token
// @Summary Create new application token
// @Description 1st party application token.
// @Description Permission : **admin.token.create**
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{token=string} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/admin
// @Router /v1/application/token [post]
func PostAppicationToken(c *gin.Context) {
	token := database.Instance.ApplicationToken.Add()
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"token": token.Hex()}))
}

// PatchRiotAccessToken is handler for endpoint PATCH /application/riot/access
// @Summary Patch Riot Token
// @Description Riot token updator (This api only available on "Riot development token". If I get application token, This endpoint will be deleted)
// @Description Permission : **admin.token.create**
// @Accept  json
// @Produce  json
// @Param body body docmodels.RequestBodyToken true "body"
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/admin
// @Router /v1/application/riot/access [patch]
func PatchRiotAccessToken(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	if bodyMap["token"] != nil {
		database.Instance.Redis.Set("riot-access-token", bodyMap["token"])
		c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
		return
	}
	c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("target filed is required"))
}
