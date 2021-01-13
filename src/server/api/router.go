package api

import (
	"github.com/DevSDK/DFD/src/server/api/v1"
	"github.com/gin-gonic/gin"
)

func Initialize(router *gin.RouterGroup) {
	v1API := router.Group("/api")
	v1.Initialize(v1API)
}
