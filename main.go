package main

import (
	"fmt"
	"github.com/Abdulsametileri/messaging-service/api"
	"github.com/Abdulsametileri/messaging-service/config"
	"github.com/Abdulsametileri/messaging-service/database"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	fmt.Println("Bismillah")
	config.Setup()

	db := database.Setup()

	repository.New(db)

	//database.Init()

	if !config.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	api.Setup(r)
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
