package v1

import (
	"github.com/DevSDK/DFD/src/server/middleware"
	"github.com/gin-gonic/gin"
)

func Initialize(router *gin.RouterGroup) {
	v1router := router.Group("v1")
	v1router.PUT("/user", middleware.JWTAuthMiddleware("user.put"), middleware.JsonParseMiddleware, PutUser)
	v1router.GET("/user/:id", middleware.JWTAuthMiddleware("user.get"), GetUser)
	v1router.GET("/user", middleware.JWTAuthMiddleware("user.get"), GetMe)

	v1router.GET("/images", middleware.JWTAuthMiddleware("imagelist.get"), GetImageList)
	v1router.POST("/image", middleware.JWTAuthMiddleware("image.post"), middleware.JsonParseMiddleware, PostImage)
	v1router.GET("/image/:id", middleware.JWTAuthMiddleware("image.get"), GetImage)
	v1router.DELETE("/image/:id", middleware.JWTAuthMiddleware("image.delete"), DelImage)

	v1router.GET("/states/:id", middleware.JWTAuthMiddleware("states.get"), GetStateHistory)
	v1router.GET("/states", middleware.JWTAuthMiddleware("states.get"), GetOwnStateHistory)
	v1router.POST("/state", middleware.JWTAuthMiddleware("states.post"), middleware.JsonParseMiddleware, PostState)

	v1router.GET("/announces", middleware.JWTAuthMiddleware(), GetAnnounceList)
	v1router.GET("/announces/:id", middleware.JWTAuthMiddleware(), GetAnnounceList)

	v1router.POST("/announce", middleware.JWTAuthMiddleware(), middleware.JsonParseMiddleware, PostAnnounce)
	v1router.GET("/announce", middleware.JWTAuthMiddleware(), GetAnnounceListMe)
	v1router.GET("/announce/:id", middleware.JWTAuthMiddleware(), GetAnnounce)
	v1router.DELETE("/announce/:id", middleware.JWTAuthMiddleware(), DelAnnounce)
	v1router.PUT("/announce/:id", middleware.JWTAuthMiddleware(), middleware.JsonParseMiddleware, PutAnnounce)
}
