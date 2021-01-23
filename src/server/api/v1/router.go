package v1

import (
	"github.com/DevSDK/DFD/src/server/middleware"
	"github.com/gin-gonic/gin"
)

func Initialize(router *gin.RouterGroup) {
	v1router := router.Group("v1")
	v1router.PATCH("/user", middleware.AppAndJWTAuthMiddleware(false, "user.patch"), middleware.JsonParseMiddleware, PatchUser)
	v1router.PATCH("/user/lol", middleware.AppAndJWTAuthMiddleware(false, "user.patch"), middleware.JsonParseMiddleware, PatchUserLolName)
	v1router.GET("/user/:id", middleware.AppAndJWTAuthMiddleware(true, "user.get"), GetUser)
	v1router.GET("/user", middleware.AppAndJWTAuthMiddleware(false, "user.get"), GetMe)

	v1router.GET("/images", middleware.AppAndJWTAuthMiddleware(true, "imagelist.get"), GetImageList)
	v1router.POST("/image", middleware.AppAndJWTAuthMiddleware(false, "image.post"), middleware.JsonParseMiddleware, PostImage)
	v1router.GET("/image/:id", GetImage)
	v1router.DELETE("/image/:id", middleware.AppAndJWTAuthMiddleware(false, "image.delete"), DelImage)

	v1router.GET("/states/:id", middleware.AppAndJWTAuthMiddleware(true, "states.get"), GetStateHistory)
	v1router.GET("/states", middleware.AppAndJWTAuthMiddleware(false, "states.get"), GetOwnStateHistory)
	v1router.POST("/state", middleware.AppAndJWTAuthMiddleware(false, "states.post"), middleware.JsonParseMiddleware, PostState)

	v1router.GET("/announces/all", middleware.AppAndJWTAuthMiddleware(true, "announces.get"), GetAllAnnounceList)
	v1router.GET("/announces/current", middleware.AppAndJWTAuthMiddleware(true, "announces.get"), GetCurrentAnnounceList)
	v1router.GET("/announces/user/:id", middleware.AppAndJWTAuthMiddleware(true, "announces.get"), GetAnnounceList)
	v1router.POST("/announce", middleware.AppAndJWTAuthMiddleware(false, "announce.post"), middleware.JsonParseMiddleware, PostAnnounce)
	v1router.GET("/announces/me", middleware.AppAndJWTAuthMiddleware(false, "announces.get"), GetAnnounceListMe)
	v1router.GET("/announce/:id", middleware.AppAndJWTAuthMiddleware(true, "announce.get"), GetAnnounce)
	v1router.DELETE("/announce/:id", middleware.AppAndJWTAuthMiddleware(false, "announce.delete"), DelAnnounce)
	v1router.PATCH("/announce/:id", middleware.AppAndJWTAuthMiddleware(false, "announce.patch"), middleware.JsonParseMiddleware, PatchAnnounce)

	v1router.POST("/application/token", middleware.AppAndJWTAuthMiddleware(false, "admin.token.create"), PostAppicationToken)
	v1router.PATCH("/application/riot/access", middleware.AppAndJWTAuthMiddleware(false, "admin.token.create"), middleware.JsonParseMiddleware, PatchRiotAccessToken)
	v1router.POST("/lol/history/updater", middleware.VerifyApplicationTokenMiddleware, PostLolHistoryUpdate)
	v1router.GET("/lol/historys", GetLolHistoryList)
	v1router.GET("/lol/history/:id", GetLolHistory)
}
