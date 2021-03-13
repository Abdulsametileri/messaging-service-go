package api

import (
	"github.com/Abdulsametileri/messaging-service/config"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	isDebug := config.IsDebug

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	if !isDebug {
		r.Use(gin.Logger())
	}

	r.Use(customRecoveryMiddleware())

	v1 := r.Group("api/v1")
	{
		v1.POST("register", register)
		v1.POST("login", login)

		v1.GET("mutateUser", jwtMiddleware, mutateUser)
		v1.GET("users", jwtMiddleware, users)
	}

	return r
}
