package api

import (
	"github.com/DevSDK/DFD/src/server/api/v1"
	"github.com/DevSDK/DFD/src/server/middleware"
	"github.com/gin-gonic/gin"
)

//Initialize all routes on API
func Initialize(router *gin.RouterGroup) {
	router.Use(middleware.CORSMiddleware)
	v1.Initialize(router)
}
