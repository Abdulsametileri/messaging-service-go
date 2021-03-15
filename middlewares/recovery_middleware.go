package middlewares

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CustomRecoveryMiddleware(base controllers.BaseController) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				friendlyErrorToClient := errors.New("Error occured in our own server. Sorry")
				if v, ok := err.(error); ok {
					base.Error(c, http.StatusInternalServerError,
						friendlyErrorToClient, v.Error())
				} else {
					base.Error(c, http.StatusInternalServerError,
						friendlyErrorToClient, err.(string))
				}
			}
		}()
		c.Next()
	}
}
