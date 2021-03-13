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

func customRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if v, ok := err.(error); ok {
					Error(c, http.StatusInternalServerError,
						v, "Error occured in the server "+v.Error())
				} else {
					Error(c, http.StatusInternalServerError,
						errors.New(err.(string)), "Error occured in the server "+err.(string))
				}
			}
		}()
		c.Next()
	}
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
		Error(c, http.StatusUnauthorized, errors.New("token is required"), "jwt middleware empty token error")
		return
	}

	claims := &viewmodels.UserClaim{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			Error(c, http.StatusUnauthorized,
				errors.New("Invalid token signature. Please try to log in again."),
				"jwt middleware invalid token signature error")
			return
		}
		Error(c, http.StatusUnauthorized,
			errors.New("Session has ended. Please try to log in."),
			"jwt middleware session ended error")
		return
	}
	if !tkn.Valid {
		Error(c, http.StatusUnauthorized,
			errors.New("Invalid Token"),
			"jwt middleware invalid token error")
		return
	}
	c.Set("claims", *claims)
	c.Next()
}
