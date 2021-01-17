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
	switch dataType := parts[0][5:]; dataType {
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

func GetImageList(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, err := database.Instance.Image.ImageList(user.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("Image not found"))
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"images": list}))
}
