package auth

import (
	"encoding/json"
	"fmt"
	_ "github.com/DevSDK/DFD/src/server/api/v1/docmodels"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/DevSDK/DFD/src/server/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func checkAndInsertUser(userMap bson.M) models.User {
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

// @Summary Login
// @Description Redirect to discord Oauth2 login page
// @Accept  json
// @Produce  json
// @tags auth/
// @Router /auth/login [get]
func Login(c *gin.Context) {
	url := CreateDiscordOauthURI()
	c.Redirect(http.StatusFound, url)
}

// @Summary Logout
// @Description Delete access token and refresh token from cookie. And register the refresh token into blacklist
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Security ApiKeyAuth
// @tags auth/
// @Router /auth/logout [get]
func Logout(c *gin.Context) {
	SERVER_URI := os.Getenv("SERVER_URI")
	BASE_URL := os.Getenv("BASE_URL")
	headerToken := c.Request.Header["Authorization"]
	if headerToken == nil || len(headerToken) == 0 {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("Access token required", false))
		c.Abort()
		return
	}
	c.SetCookie("access", "", -1, BASE_URL, SERVER_URI, false, true)

	refreshHeaderToken := c.Request.Header["X-Dfd-Refresh"]
	if refreshHeaderToken == nil || len(refreshHeaderToken) == 0 {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("no refresh token"))
		c.Abort()
		return
	}
	refreshToken := refreshHeaderToken[0]

	c.SetCookie("refresh", "", -1, BASE_URL, SERVER_URI, false, true)

	//Register refresh blacklist
	database.Instance.Redis.Set(refreshToken, true)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(nil))
}

// @Summary Refresh token
// @Description Refresh access token token. REQUIRED: access and refresh JWT token in cookie.
// @Accept  json
// @Produce  json
// @Param X-Dfd-Refresh header string true "refresh token"
// @Success 200 {object} docmodels.ResponseSuccess "success"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @Security ApiKeyAuth
// @tags auth/
// @Router /auth/refresh [get]
func Refresh(c *gin.Context) {
	DFD_SECRET_CODE := os.Getenv("DFD_SECRET_CODE")

	headerToken := c.Request.Header["Authorization"]
	if headerToken == nil || len(headerToken) == 0 {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("Access token required", false))
		c.Abort()
		return
	}
	accessToken := headerToken[0]
	refreshHeaderToken := c.Request.Header["X-Dfd-Refresh"]
	if refreshHeaderToken == nil || len(refreshHeaderToken) == 0 {
		c.JSON(http.StatusBadRequest, utils.CreateBadRequestJSONMessage("no refresh token"))
		c.Abort()
		return
	}
	refreshToken := refreshHeaderToken[0]
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, _ := parser.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(DFD_SECRET_CODE), nil
	})
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("no refresh token", false))
		return
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	useridHex := claims["id"].(string)
	isBlackListed, err := database.Instance.Redis.Get(refreshToken)

	token, _ = parser.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(DFD_SECRET_CODE), nil
	})

	if isBlackListed == "true" || err == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("refresh token is expired", true))
		return
	}

	userid, _ := primitive.ObjectIDFromHex(useridHex)
	user, err := database.Instance.User.FindById(userid)

	newAccessToken, err := CreateAccessToken(user)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		return
	}
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"access": newAccessToken}))
}

func Redirect(c *gin.Context) {
	queryError := c.Query("error")

	if queryError != "" {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage(c.Query("error_description"), false))
		return
	}
	DISCORD_CLIENT_ID := os.Getenv("DISCORD_CLIENT_ID")
	DISCORD_REDIRECT_URI := os.Getenv("DISCORD_REDIRECT_URI")
	DISCORD_SECRET_ID := os.Getenv("DISCORD_SECRET_ID")
	DISCORD_API_BASE := os.Getenv("DISCORD_API_BASE")
	//Get Access Token from discord server
	resp, err := http.PostForm(DISCORD_API_BASE+"/oauth2/token",
		url.Values{"code": {c.Query("code")},
			"client_id":     {DISCORD_CLIENT_ID},
			"client_secret": {DISCORD_SECRET_ID},
			"redirect_uri":  {DISCORD_REDIRECT_URI},
			"grant_type":    {"authorization_code"}})
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		return
	}

	tokenString, _ := ioutil.ReadAll(resp.Body)
	accessMap := bson.M{}
	if err := json.Unmarshal([]byte(tokenString), &accessMap); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		return
	}

	bearer := "Bearer " + accessMap["access_token"].(string)
	//Reqeust user information to discord server
	userInfoRequest, err := http.NewRequest("GET", DISCORD_API_BASE+"/users/@me", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		log.Print(err.Error())
		return
	}
	userInfoRequest.Header.Add("Authorization", bearer)
	resp, err = (&http.Client{}).Do(userInfoRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		log.Print(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	userMap := bson.M{}
	if err := json.Unmarshal([]byte(body), &userMap); err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		return
	}
	userMap["tokenString"] = string(accessMap["refresh_token"].(string))

	user := checkAndInsertUser(userMap)
	accessToken, err := CreateAccessToken(user)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		return
	}
	refreshToken, err := CreateRefreshToken()
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, utils.CreateInternalServerErrorJSONMessage())
		return
	}
	database.Instance.Redis.Set(user.Id.Hex(), refreshToken)
	SERVER_URI := os.Getenv("SERVER_URI")
	BASE_URL := os.Getenv("BASE_URL")

	REDIRECT_URL := os.Getenv("REDIRECT_URL")
	c.SetCookie("access", accessToken, 0, BASE_URL, SERVER_URI, false, true)
	c.SetCookie("refresh", refreshToken, 0, BASE_URL, SERVER_URI, false, true)
	c.Redirect(http.StatusFound, REDIRECT_URL)
}

// @Summary Get access token and refresh
// @Description Return access and refresh token from cookie if it is valid.
// @Accept  json
// @Produce  json
// @Success 200 {object} docmodels.ResponseSuccess{access=string,refresh=string} "success"
// @Failure 404 {object} docmodels.ResponseNotFound "Cannt found user"
// @Failure 401 {object} docmodels.ResponseUnauthorized "Unauthorized Request. If token is expired, **token_expired** filed must be set true"
// @tags auth/
// @Router /auth/token [get]
func Token(c *gin.Context) {
	SERVER_URI := os.Getenv("SERVER_URI")
	BASE_URL := os.Getenv("BASE_URL")
	DFD_SECRET_CODE := os.Getenv("DFD_SECRET_CODE")
	accessToken, err := c.Cookie("access")
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("Access token required", false))
		c.Abort()
		return
	}
	refreshToken, err := c.Cookie("refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.CreateUnauthorizedJSONMessage("Refresh token required", false))
		c.Abort()
		return
	}

	claims := &jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
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
	_, err = database.Instance.User.FindById(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateNotFoundJSONMessage("cannot found user"))
		c.Abort()
		return
	}
	c.SetCookie("access", accessToken, -1, BASE_URL, SERVER_URI, false, true)
	c.SetCookie("refresh", refreshToken, -1, BASE_URL, SERVER_URI, false, true)
	c.JSON(http.StatusOK, utils.CreateSuccessJSONMessage(gin.H{"access": accessToken, "refresh": refreshToken}))
}
