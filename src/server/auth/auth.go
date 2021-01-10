package auth

import (
	"github.com/gin-gonic/gin"
)

func Initialize(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	authGroup.GET("/login", Login)
	authGroup.GET("/redirect", Redirect)
}
