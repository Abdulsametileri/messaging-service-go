package main

import (
	"fmt"
	"github.com/Abdulsametileri/messaging-service/api"
	"github.com/Abdulsametileri/messaging-service/config"
	"github.com/Abdulsametileri/messaging-service/database"
	"github.com/Abdulsametileri/messaging-service/repository"
	"log"
)

func main() {
	fmt.Println("Bismillah")
	config.Setup()

	db := database.Setup()
	database.Migrate()

	repository.Setup(repository.LogRepository{}, db)
	repository.Setup(repository.AuthRepository{}, db)
	repository.Setup(repository.UserRepository{}, db)
	repository.Setup(repository.MessageRepository{}, db)

	router := api.SetupRouter()

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
