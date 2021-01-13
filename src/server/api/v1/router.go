package v1

import (
	"github.com/DevSDK/DFD/src/server/auth"
	"github.com/gin-gonic/gin"
)

func Initialize(router *gin.RouterGroup) {
	v1router := router.Group("v1")
	v1router.PUT("/user", auth.JWTAuthMiddleware("user.put"), PutUser)
	v1router.GET("/user/:id", auth.JWTAuthMiddleware("user.get"), GetUser)
	v1router.GET("/user", auth.JWTAuthMiddleware("user.get"), GetMe)

	v1router.GET("/images", auth.JWTAuthMiddleware("imagelist.get"), GetImageList)
	v1router.POST("/image", auth.JWTAuthMiddleware("image.post"), PostImage)
	v1router.GET("/image/:id", auth.JWTAuthMiddleware("image.get"), GetImage)
	v1router.DELETE("/image/:id", auth.JWTAuthMiddleware("image.delete"), DelImage)

}
