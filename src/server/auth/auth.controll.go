package auth

import (
	"encoding/json"
	"fmt"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func intenralServerError(c *gin.Context, err error) {
	c.JSON(502, gin.H{"message": "Internal Server Error"})
	log.Print(err.Error())
}

func checkAndInsertUser(userMap map[string]interface{}) models.User {
	user, err := database.Instance.User.FindByEmail(userMap["email"].(string))
	if err != nil {
		user, _ = database.Instance.User.Register(userMap)
	}
	return user
}

func CreateDiscordOauthURI() string {
	DISCORD_CLIENT_ID := os.Getenv("DISCORD_CLIENT_ID")
	DISCORD_REDIRECT_URI := os.Getenv("DISCORD_REDIRECT_URI")
	DISCORD_SCOPES := os.Getenv("DISCORD_SCOPES")
	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s"+
		"&redirect_uri=%s&response_type=code&scope=%s", DISCORD_CLIENT_ID, DISCORD_REDIRECT_URI, DISCORD_SCOPES)
}

func createToken(atClaims jwt.MapClaims) (string, error) {
	DFD_SECRET_CODE := os.Getenv("DFD_SECRET_CODE")
	gen := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := gen.SignedString([]byte(DFD_SECRET_CODE))

	if err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func CreateAccessToken(user models.User) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = user.Id.Hex()
	atClaims["email"] = user.Email
	atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	token, err := createToken(atClaims)
	return token, err
}

func CreateRefreshToken() (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	uuid, _ := uuid.NewUUID()
	atClaims["uuid"] = uuid
	token, err := createToken(atClaims)
	return token, err
}

func Login(c *gin.Context) {
	c.Redirect(http.StatusFound, CreateDiscordOauthURI())
}

func Logout(c *gin.Context) {
	SERVER_URI := os.Getenv("SERVER_URI")
	_, err := c.Cookie("access")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Auth failed"})
		return
	}
	c.SetCookie("access", "", -1, "/", SERVER_URI, false, true)

	refreshToken, err := c.Cookie("refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Auth failed"})
		return
	}
	c.SetCookie("refresh", "", -1, "/", SERVER_URI, false, true)

	//Register refresh blacklist
	database.Instance.Redis.Set(refreshToken, true)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}

func Refresh(c *gin.Context) {
	DFD_SECRET_CODE := os.Getenv("DFD_SECRET_CODE")
	SERVER_URI := os.Getenv("SERVER_URI")
	accessToken, err := c.Cookie("access")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("no access token"))
		return
	}
	refreshToken, err := c.Cookie("refresh")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("no refresh token"))
		return
	}
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, _ := parser.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(DFD_SECRET_CODE), nil
	})
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("no refresh token"))
		return
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	useridHex := claims["id"].(string)
	isBlackListed, err := database.Instance.Redis.Get(refreshToken)

	token, _ = parser.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(DFD_SECRET_CODE), nil
	})

	if isBlackListed == "true" || err == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("refresh token is expired"))
		return
	}

	userid, _ := primitive.ObjectIDFromHex(useridHex)
	user, err := database.Instance.User.FindById(userid)

	newAccessToken, err := CreateAccessToken(user)
	if err != nil {
		intenralServerError(c, err)
		return
	}
	c.SetCookie("access", newAccessToken, 0, "/", SERVER_URI, false, true)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}

func Redirect(c *gin.Context) {
	DISCORD_CLIENT_ID := os.Getenv("DISCORD_CLIENT_ID")
	DISCORD_REDIRECT_URI := os.Getenv("DISCORD_REDIRECT_URI")
	DISCORD_SECRET_ID := os.Getenv("DISCORD_SECRET_ID")
	DISCORD_API_BASE := os.Getenv("DISCORD_API_BASE")
	//Get Access Token from discord server
	resp, _ := http.PostForm(DISCORD_API_BASE+"/oauth2/token",
		url.Values{"code": {c.Query("code")},
			"client_id":     {DISCORD_CLIENT_ID},
			"client_secret": {DISCORD_SECRET_ID},
			"redirect_uri":  {DISCORD_REDIRECT_URI},
			"grant_type":    {"authorization_code"}})

	tokenString, _ := ioutil.ReadAll(resp.Body)
	var accessMap map[string]interface{}
	if err := json.Unmarshal([]byte(tokenString), &accessMap); err != nil {
		intenralServerError(c, err)
		return
	}

	bearer := "Bearer " + accessMap["access_token"].(string)
	//Reqeust user information to discord server
	userInfoRequest, err := http.NewRequest("GET", DISCORD_API_BASE+"/users/@me", nil)
	if err != nil {
		c.JSON(502, gin.H{"message": "Internal Server Error"})
		log.Print(err.Error())
		return
	}
	userInfoRequest.Header.Add("Authorization", bearer)
	resp, err = (&http.Client{}).Do(userInfoRequest)
	if err != nil {
		c.JSON(502, gin.H{"message": "Discord Server Error"})
		log.Print(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var userMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &userMap); err != nil {
		intenralServerError(c, err)
		return
	}

	userMap["tokenString"] = string(accessMap["refresh_token"].(string))
	user := checkAndInsertUser(userMap)
	accessToken, err := CreateAccessToken(user)
	if err != nil {
		intenralServerError(c, err)
		return
	}
	refreshToken, err := CreateRefreshToken()
	if err != nil {
		intenralServerError(c, err)
		return
	}
	database.Instance.Redis.Set(user.Id.Hex(), refreshToken)
	SERVER_URI := os.Getenv("SERVER_URI")
	c.SetCookie("access", accessToken, 0, "/", SERVER_URI, false, true)
	c.SetCookie("refresh", refreshToken, 0, "/", SERVER_URI, false, true)
	c.Redirect(http.StatusFound, "/")
}
