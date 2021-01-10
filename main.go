package main

import (
	"github.com/DevSDK/DFD/src/server"
	"github.com/DevSDK/DFD/src/server/database"
	"log"
)

func main() {
	if err := database.Initialize(); err != nil {
		log.Fatal("Database Error Occured" + err.Error())
	}

	defer database.Disconnect()
	server.RunServer()
}
