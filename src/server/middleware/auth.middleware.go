package middleware

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func JWTAuthMiddleware(permissions ...string) gin.HandlerFunc {
	contains := func(src []string, dst []string) bool {
		ret := true
		for _, s := range src {
			flag := false
			for _, v := range dst {
				if s == v {
					flag = true
					break
				}
			}
			ret = ret && flag
		}
		return ret
	}

	return func(c *gin.Context) {
		DFD_SECRET_CODE := os.Getenv("DFD_SECRET_CODE")
		accessToken, err := c.Cookie("access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Access token required"})
			c.Abort()
			return
		}

		claims := &jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(DFD_SECRET_CODE), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "token is expired"})
				c.Abort()
				return
			}
			c.JSON(http.StatusForbidden, gin.H{"message": "Auth failed"})
			c.Abort()
			return
		}
		userId, _ := primitive.ObjectIDFromHex((*claims)["id"].(string))
		user, err := database.Instance.User.FindById(userId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Cannot found user"})
			c.Abort()
			return
		}
		userRole, err := database.Instance.Role.FindByName(user.Role)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Cannot found role"})
			c.Abort()
			return
		}
		if !contains(permissions, userRole.Permissions) {
			c.JSON(http.StatusForbidden, gin.H{"message": "You don't have permission"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
