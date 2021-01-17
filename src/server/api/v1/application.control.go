package v1

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostAppicationToken(c *gin.Context) {
	token := database.Instance.ApplicationToken.Add()
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"token": token.Hex()}))
}
