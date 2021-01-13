package server

import (
	"github.com/DevSDK/DFD/src/server/api"
	"github.com/DevSDK/DFD/src/server/auth"
	"github.com/gin-gonic/gin"
	"os"
)

func initialize() *gin.Engine {
	engine := gin.Default()
	base := engine.Group("/")
	auth.Initialize(base)
	api.Initialize(base)
	return engine
}

func RunServer() {
	engine := initialize()
	SERVER_PORT := os.Getenv("SERVER_PORT")
	engine.Run(":" + SERVER_PORT)
}
