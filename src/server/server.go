package server

import (
	_ "github.com/DevSDK/DFD/docs/v1"
	"github.com/DevSDK/DFD/src/server/api"
	"github.com/DevSDK/DFD/src/server/auth"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

func initialize() *gin.Engine {
	engine := gin.Default()
	base := engine.Group("/")
	//API Document
	url := ginSwagger.URL("http://localhost:8080/docs/v1/doc.json")
	base.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	auth.Initialize(base)
	api.Initialize(base)
	return engine
}

func RunServer() {
	engine := initialize()
	SERVER_PORT := os.Getenv("SERVER_PORT")
	engine.Run(":" + SERVER_PORT)
}
