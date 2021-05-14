package main

import (
	"context"
	"github.com/Abdulsametileri/messaging-service/config"
	"github.com/Abdulsametileri/messaging-service/controllers"
	"github.com/Abdulsametileri/messaging-service/database"
	"github.com/Abdulsametileri/messaging-service/infra/mailgunclient"
	"github.com/Abdulsametileri/messaging-service/middlewares"
	"github.com/Abdulsametileri/messaging-service/repository/logrepo"
	"github.com/Abdulsametileri/messaging-service/repository/messagerepo"
	"github.com/Abdulsametileri/messaging-service/repository/userrepo"
	"github.com/Abdulsametileri/messaging-service/services/authservice"
	"github.com/Abdulsametileri/messaging-service/services/logservice"
	"github.com/Abdulsametileri/messaging-service/services/messageservice"
	"github.com/Abdulsametileri/messaging-service/services/userservice"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.Setup()

	db := database.Setup()
	database.Migrate()

	/*
		====== Setup infra ==============
	*/
	emailClient := mailgunclient.NewMailgunClient()
	_ = emailClient

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

	/*
		====== Setup controllers ========
	*/
	baseCtl := controllers.NewBaseController(logService)
	userCtl := controllers.NewUserController(baseCtl, authService, userService)
	messageCtl := controllers.NewMessageController(baseCtl, userService, messageService)

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

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	srv.Shutdown(ctx)
}
