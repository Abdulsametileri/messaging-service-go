package main

import (
	"fmt"
	"github.com/Abdulsametileri/messaging-service/config"
	"github.com/Abdulsametileri/messaging-service/controllers"
	"github.com/Abdulsametileri/messaging-service/database"
	"github.com/Abdulsametileri/messaging-service/infra/redisclient"
	"github.com/Abdulsametileri/messaging-service/middlewares"
	"github.com/Abdulsametileri/messaging-service/repository/logrepo"
	"github.com/Abdulsametileri/messaging-service/repository/messagerepo"
	"github.com/Abdulsametileri/messaging-service/repository/userrepo"
	"github.com/Abdulsametileri/messaging-service/services/authservice"
	"github.com/Abdulsametileri/messaging-service/services/logservice"
	"github.com/Abdulsametileri/messaging-service/services/messageservice"
	"github.com/Abdulsametileri/messaging-service/services/redisservice"
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
		====== Setup infra ==============
	*/
	redisClient := redisclient.NewRedisClient()

	/*
		====== Setup repositories =======
	*/
	userRepo := userrepo.NewUserRepo(db)
	logRepo := logrepo.NewLogRepo(db)
	messageRepo := messagerepo.NewMessageRepo(db)

	/*
		====== Setup services ===========
	*/
	authService := authservice.NewAuthService()
	userService := userservice.NewUserService(userRepo)
	logService := logservice.NewLogService(logRepo)
	messageService := messageservice.NewMessageService(messageRepo)
	redisService := redisservice.NewRedisService(redisClient)

	/*
		====== Setup controllers ========
	*/
	baseCtl := controllers.NewBaseController(logService)
	userCtl := controllers.NewUserController(baseCtl, authService, userService)
	messageCtl := controllers.NewMessageController(baseCtl, userService, messageService)
	websocketCtl := controllers.NewWebSocketController(baseCtl, redisService)

	if !config.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	var router = gin.New()

	if config.IsDebug {
		router.Use(gin.Logger())
	}

	router.Use(middlewares.CustomRecoveryMiddleware(baseCtl))

	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	v1 := router.Group("api/v1")
	{
		v1.POST("register", userCtl.Register)
		v1.POST("login", userCtl.Login)

		v1.GET("mutateUser/:mutateUserId", middlewares.RequireLoggedIn(baseCtl), userCtl.MutateUser)
		v1.GET("users", middlewares.RequireLoggedIn(baseCtl), userCtl.GetUserList)

		v1.GET("messagesWith/:userName", middlewares.RequireLoggedIn(baseCtl), messageCtl.GetMessages)
		v1.POST("sendMessage/:userName", middlewares.RequireLoggedIn(baseCtl), messageCtl.SendMessage)
	}

	router.GET("ws/:chatId", websocketCtl.ServeWs)

	fmt.Println("Server Started")
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
