package v1

import (
	"github.com/DevSDK/DFD/src/server/middleware"
	"github.com/gin-gonic/gin"
)

func Initialize(router *gin.RouterGroup) {
	v1router := router.Group("v1")
	v1router.PUT("/user", middleware.JWTAuthMiddleware("user.put"), middleware.JsonParseMiddleware, PutUser)
	v1router.PUT("/user/lol", middleware.JWTAuthMiddleware("user.put"), middleware.JsonParseMiddleware, PutUserLolName)
	v1router.GET("/user/:id", middleware.JWTAuthMiddleware("user.get"), GetUser)
	v1router.GET("/user", middleware.JWTAuthMiddleware("user.get"), GetMe)

	v1router.GET("/images", middleware.JWTAuthMiddleware("imagelist.get"), GetImageList)
	v1router.POST("/image", middleware.JWTAuthMiddleware("image.post"), middleware.JsonParseMiddleware, PostImage)
	v1router.GET("/image/:id", middleware.JWTAuthMiddleware("image.get"), GetImage)
	v1router.DELETE("/image/:id", middleware.JWTAuthMiddleware("image.delete"), DelImage)

	v1router.GET("/states/:id", middleware.JWTAuthMiddleware("states.get"), GetStateHistory)
	v1router.GET("/states", middleware.JWTAuthMiddleware("states.get"), GetOwnStateHistory)
	v1router.POST("/state", middleware.JWTAuthMiddleware("states.post"), middleware.JsonParseMiddleware, PostState)

	v1router.GET("/announces", middleware.JWTAuthMiddleware("announces.get"), GetAnnounceList)
	v1router.GET("/announces/:id", middleware.JWTAuthMiddleware("announces.get"), GetAnnounceList)

	v1router.POST("/announce", middleware.JWTAuthMiddleware("announce.post"), middleware.JsonParseMiddleware, PostAnnounce)
	v1router.GET("/announce", middleware.JWTAuthMiddleware("announce.get"), GetAnnounceListMe)
	v1router.GET("/announce/:id", middleware.JWTAuthMiddleware("announce.get"), GetAnnounce)
	v1router.DELETE("/announce/:id", middleware.JWTAuthMiddleware("announce.delete"), DelAnnounce)
	v1router.PUT("/announce/:id", middleware.JWTAuthMiddleware("announce.put"), middleware.JsonParseMiddleware, PutAnnounce)

	v1router.POST("/application/token", middleware.JWTAuthMiddleware("admin.token.create"), PostAppicationToken)
	v1router.POST("/lol/history/updater", middleware.JsonParseMiddleware, middleware.VerifyApplicationTokenMiddleware, PostLolHistoryUpdate)
	v1router.GET("/lol/historys", middleware.JsonParseMiddleware, GetLolHistoryList)
	v1router.GET("/lol/history/:id", middleware.JsonParseMiddleware, GetLolHistory)
}
