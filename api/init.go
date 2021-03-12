package api

import "github.com/gin-gonic/gin"

func Setup(r *gin.Engine) {
	v1 := r.Group("api/v1")
	{
		v1.POST("register", register)
	}
}
