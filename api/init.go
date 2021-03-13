package api

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/config"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SetupRouter() *gin.Engine {
	if !config.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.POST("register", register)
		v1.POST("login", login)
	}

	return r
}

func jwtMiddleware(c *gin.Context) {
	r := c.Request
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	if token == "" {
		token = c.Query("token")
	}
	if token == "" || token == "Bearer" {
		c.String(http.StatusUnauthorized, "Token is required!")
		c.Abort()
		return
	}

	claims := &viewmodels.UserClaim{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			Error(c, http.StatusUnauthorized, errors.New("Invalid token signature. Please try to log in again."))
			return
		}
		Error(c, http.StatusUnauthorized, errors.New("Session has ended. Please try to log in."))
		return
	}
	if !tkn.Valid {
		Error(c, http.StatusUnauthorized, errors.New("Invalid Token"))
		return
	}
	c.Set("claims", *claims)
	c.Next()
}
