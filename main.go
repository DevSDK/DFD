package main

import (
	"github.com/DevSDK/DFD/docs/v1"
	"github.com/DevSDK/DFD/src/server"
	"github.com/DevSDK/DFD/src/server/database"
	"os"

	"log"
)

// @title DFD API
// @version 1.0
// @description # This is a DFD server
// @description Most of api endpoints aim to restful.
// @description ## Permissions
// @description API endpoints request a permission.
// @description If not described, It is public API.
// @description Permissions are described below table:
// @description | Permission        | Description                         | Role               |
// @description |-------------------|-------------------------------------|--------------------|
// @description | user.patch        | Allows edit user information.       | Admin <br/> User   |
// @description | user.get          | Allows get user information         | Admin <br/> User   |
// @description | imagelist.get     |  Allows get own image list          | Admin <br /> User  |
// @description | image.post        | Allows image upload                 | Admin  <br /> User |
// @description | image.delete      | Allows image delete                 | Admin  <br /> User |
// @description | states.get        | Allows get states history           | Admin  <br /> User |
// @description | states.post       | Allows create states                | Admin  <br /> User |
// @description | announces.get     | Allows get states announce list     | Admin  <br /> User |
// @description | announce.post     | Allows create announce              | Admin  <br /> User |
// @description | announce.get      | Allows get announce                 | Admin  <br /> User |
// @description | announce.delete   | Allows delete my announce           | Admin  <br /> User |
// @description | announce.patch    | Allows patch announce               | Admin  <br /> User |
// @description | admin.token.create| Allows create and patch admin token | Admin              |
// @description ## Authentication
// @description Login is working with discord Oauth2.
// @description After login, the access and refresh token stored in Cookie.
// @description You should get from the tokens via /token endpoint.
// @description However, all request takes access token in **Authorization** header for security.
// @description The access token is JWT statless token.
// @description Refresh token is stored and used cookie.
// @contact.name Seokho Song
// @contact.url http://github.com/devsdkGIN_MODE
// @license.url https://opensource.org/licenses/MIT
// @host devsdk.net/api/dfd
// @BasePath /
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := database.Initialize(); err != nil {
		log.Fatal("Database Error Occured " + err.Error())
	}
	if os.Getenv("GIN_MODE") == "release" {
		v1.SwaggerInfo.Host = "devsdk.net/api/dfd"
	} else {
		v1.SwaggerInfo.Host = "localhost:18020"
	}

	defer database.Disconnect()
	server.RunServer()
}
