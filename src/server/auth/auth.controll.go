package auth

import (
	"encoding/json"
	"fmt"
	"github.com/DevSDK/DFD/src/server/database"
	"github.com/DevSDK/DFD/src/server/database/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

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
	atClaims["email"] = user.Email
	atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	token, err := createToken(atClaims)
	return token, err
}

func CreateRefreshToken() (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	token, err := createToken(atClaims)
	return token, err
}

func Login(c *gin.Context) {
	c.Redirect(http.StatusFound, CreateDiscordOauthURI())
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
		c.JSON(502, gin.H{"msg": "Internal Server Error"})
		log.Fatal(err.Error())
		return
	}
	bearer := "Bearer " + accessMap["access_token"].(string)
	//Reqeust user information to discord server
	userInfoRequest, err := http.NewRequest("GET", DISCORD_API_BASE+"/users/@me", nil)
	if err != nil {
		c.JSON(502, gin.H{"msg": "Internal Server Error"})
		log.Print(err.Error())
		return
	}
	userInfoRequest.Header.Add("Authorization", bearer)
	resp, err = (&http.Client{}).Do(userInfoRequest)
	if err != nil {
		c.JSON(502, gin.H{"msg": "Discord Server Error"})
		log.Print(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var userMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &userMap); err != nil {
		c.JSON(502, gin.H{"msg": "Internal Server Error"})
		log.Fatal(err.Error())
		return
	}
	user, err := database.Instance.FindUserByEmail(userMap["email"].(string))
	if err != nil {
		//New registered User
		userMap["tokenString"] = string(tokenString)
		user, _ = database.Instance.RegisterUser(userMap)
	}
	accessToken, err := CreateAccessToken(user)
	if err != nil {
		c.JSON(502, gin.H{"msg": "Internal Server Error"})
		log.Fatal(err.Error())
		return
	}
	refreshToken, err := CreateRefreshToken()
	if err != nil {
		c.JSON(502, gin.H{"msg": "Internal Server Error"})
		log.Fatal(err.Error())
		return
	}
	SERVER_URI := os.Getenv("SERVER_URI")
	c.SetCookie("access", accessToken, 60*60, "/", SERVER_URI, false, true)
	c.SetCookie("refresh", refreshToken, 60*60*24*30, "/", SERVER_URI, false, true)
	c.JSON(200, gin.H{"access": accessToken, "refresh": refreshToken})

}
