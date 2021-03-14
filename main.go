package main

import (
	"github.com/Abdulsametileri/messaging-service/config"
	"github.com/Abdulsametileri/messaging-service/controllers"
	"github.com/Abdulsametileri/messaging-service/database"
	"github.com/Abdulsametileri/messaging-service/repository/logrepo"
	"github.com/Abdulsametileri/messaging-service/repository/userrepo"
	"github.com/Abdulsametileri/messaging-service/services/authservice"
	"github.com/Abdulsametileri/messaging-service/services/logservice"
	"github.com/Abdulsametileri/messaging-service/services/userservice"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	config.Setup()

	db := database.Setup()
	database.Migrate()

	/*
		====== Setup repositories =======
	*/
	userRepo := userrepo.NewUserRepo(db)
	logRepo := logrepo.NewLogRepo(db)

	/*
		====== Setup services ===========
	*/
	authService := authservice.NewAuthService()
	userService := userservice.NewUserService(userRepo)
	logService := logservice.NewLogService(logRepo)

	/*
		====== Setup controllers ========
	*/
	baseCtl := controllers.NewBaseController(logService)
	userCtl := controllers.NewUserController(baseCtl, authService, userService)

	var router = gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	v1 := router.Group("api/v1")
	{
		v1.POST("register", userCtl.Register)
		v1.POST("login", userCtl.Login)
	}

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}

	/*fmt.Println("Bismillah")
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
	}*/
}
