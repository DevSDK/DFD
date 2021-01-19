package v1

import (
	"bytes"
	"encoding/base64"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
)

// @Summary Post image
// @Description Upload Image by base64 image. It supports **png, jpeg, gif**
// @Description Permission : **image.post**
// @Accept  json
// @Produce  json
// @Param img body docmodels.RequestBodyImagePost true "base64 image"
// @Success 200 {object} docmodels.ResponseSuccess{id=string} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/image
// @Router /api/v1/image [post]
func PostImage(c *gin.Context) {
	bodyMap := c.MustGet("bodymap").(bson.M)
	if bodyMap["img"] == nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("cannot found announce"))
		return
	}
	imgString, ok := bodyMap["img"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("invalid image"))
		return
	}
	parts := strings.Split(imgString, ",")
	if len(parts) != 2 || len(parts[0]) < 5 {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("image should contain data-link"))
		return
	}
	unbased, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("invalid base64 image"))
		return
	}
	dataType := strings.Split(parts[0][5:], ";")[0]
	switch dataType {
	case "image/png":
		fallthrough
	case "image/jpeg":
		fallthrough
	case "image/gif":
		reader := bytes.NewReader(unbased)
		user := c.MustGet("user").(models.User)
		dataId, err := database.Instance.Image.Upload(reader, dataType, user.Id)
		if err != nil {
			log.Print(err.Error())
			c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
			return
		}
		c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"id": dataId}))
		return
	default:
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage(dataType+" is not supported"))
		return
	}
}

// @Summary Get image by id
// @Description Get Image. When Success it provide image
// @Accept  json
// @Produce  png
// @Produce  jpeg
// @Produce  gif
// @Param id path string true "Image ID"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @tags api/v1/image
// @Router /api/v1/image/{id} [get]
func GetImage(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	buf, err := database.Instance.Image.DownloadById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Image not found"))
		return
	}
	meta, err := database.Instance.Image.GetMetdataById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Image Metadata not found"))
		return
	}
	c.Header("Content-Type", meta["content-type"].(string))
	c.Writer.Write(buf.Bytes())
}

// @Summary Delete image
// @Description Delete image by image id.
// @Description Permission : **image.delete**
// @Accept  json
// @Produce  json
// @Param id path string true "Image ID"
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/image
// @Router /api/v1/image/{id} [delete]
func DelImage(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	user := c.MustGet("user").(models.User)
	meta, err := database.Instance.Image.GetMetdataById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Image Metadata not found"))
		return
	}
	uploader := meta["uploader"].(primitive.ObjectID)
	if uploader.Hex() != user.Id.Hex() {
		c.JSON(http.StatusForbidden, utils.CreateForbbidnJSONMessage("Permission denied"))
		return
	}

	if err := database.Instance.Image.DeleteImageById(id); err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Image not found"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}

// @Summary Get my image list
// @Description Get Image list uploaded from me.
// @Description Permission : **imagelist.get**
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{images=[]docmodels.ResponseImageElement} "success"
// @Failure 500 {object} docmodels.ResponseInternalServerError "Internal Server Error"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 403 {object} docmodels.ResponseNotFound "You don't have permission"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Failure 400 {object} docmodels.ResponseBadRequest "Bad request"
// @Security ApiKeyAuth
// @tags api/v1/image
// @Router /api/v1/images [get]
func GetImageList(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, err := database.Instance.Image.ImageList(user.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Image not found"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"images": list}))
}
