package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func JWTAuthMiddleware(c *gin.Context) {
	DFD_SECRET_CODE := os.Getenv("DFD_SECRET_CODE")
	accessToken, err := c.Cookie("access")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Access token required"})
		c.Abort()
		return
	}

	claims := &jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(DFD_SECRET_CODE), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "token is expired"})
			c.Abort()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden, "message": "Auth failed"})
		c.Abort()
		return
	}
	c.Next()
}
