package v1

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"strings"
)

func PostImage(c *gin.Context) {
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
	imgString := bodyMap["img"]
	if imgString == nil {
		c.JSON(400, gin.H{"message": "wrong image"})
		return
	}

	parts := strings.Split(imgString.(string), ",")
	if len(parts) != 2 || len(parts[0]) < 5 {

		c.JSON(400, gin.H{"message": "wrong request"})
		return
	}
	unbased, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		c.JSON(400, gin.H{"message": "wrong image"})
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
		database.Instance.Image.Upload(reader, dataType, user.Id)
		c.JSON(200, gin.H{"message": "success"})
		return
	default:
		c.JSON(400, gin.H{"message": "wrong data"})
		return
	}

}

func GetImage(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	buf, err := database.Instance.Image.DownloadById(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "cannot found image"})
		return
	}
	meta, err := database.Instance.Image.GetMetdataById(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "cannot found metadata"})
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
		c.JSON(404, gin.H{"message": "cannot found metadata"})
		return
	}
	uploader := meta["uploader"].(primitive.ObjectID)
	if uploader.Hex() != user.Id.Hex() {
		c.JSON(400, gin.H{"message": "you cannot delete other's"})
		return
	}

	if err := database.Instance.Image.DeleteImageById(id); err != nil {
		c.JSON(400, gin.H{"message": "delete error"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func GetImageList(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	list, err := database.Instance.Image.ImageList(user.Id)

	if err != nil {
		c.JSON(400, gin.H{"message": "cannt found images"})
		return
	}
	c.JSON(200, gin.H{"message": list})
}
