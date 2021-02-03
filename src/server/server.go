package server

import (
	_ "github.com/DevSDK/DFD/docs/v1"
	"github.com/DevSDK/DFD/src/server/api"
	"github.com/DevSDK/DFD/src/server/auth"
	"github.com/DevSDK/DFD/src/server/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

func initialize() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.CORSMiddleware)
	base := engine.Group("/")
	//API Document
	url := ginSwagger.URL("http://localhost:18020/docs/v1/doc.json")
	if os.Getenv("GIN_MODE") == "release" {
		url = ginSwagger.URL("https://devsdk.net/api/dfd/docs/v1/doc.json")
	}
	base.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	auth.Initialize(base)
	api.Initialize(base)
	return engine
}

func RunServer() {
	engine := initialize()
	SERVER_PORT := "18020"
	engine.Run("0.0.0.0:" + SERVER_PORT)
}
