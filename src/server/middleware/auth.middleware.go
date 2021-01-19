package middleware

import (
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func AppAndJWTAuthMiddleware(isApplicationAllowed bool, permissions ...string) gin.HandlerFunc {
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
		if isApplicationAllowed && c.Request.Header["Dfd-App-Auth"] != nil {
			VerifyApplicationTokenMiddleware(c)
			return
		}

		DFD_SECRET_CODE := os.Getenv("DFD_SECRET_CODE")
		headerToken := c.Request.Header["Authorization"]
		if headerToken == nil || len(headerToken) == 0 {
			c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("Access token required", false))
			c.Abort()
			return
		}
		accessToken := headerToken[0]

		claims := &jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(DFD_SECRET_CODE), nil
		})

		if err != nil {
			if (err.(*jwt.ValidationError)).Errors == jwt.ValidationErrorExpired {
				c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("token is expired", true))
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("Token is not valid", false))
			c.Abort()
			return
		}
		userId, _ := primitive.ObjectIDFromHex((*claims)["id"].(string))
		user, err := database.Instance.User.FindById(userId)
		if err != nil {
			c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("cannot found user"))
			c.Abort()
			return
		}
		userRole, err := database.Instance.Role.FindByName(user.Role)
		if err != nil {
			c.JSON(http.StatusForbidden, utils.CreateForbbidnJSONMessage("Cannot found role"))
			c.Abort()
			return
		}
		if !contains(permissions, userRole.Permissions) {
			c.JSON(http.StatusForbidden, utils.CreateForbbidnJSONMessage("You don't have permission"))
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
