package middlewares

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/controllers"
	"github.com/Abdulsametileri/messaging-service/services/authservice"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func stripBearer(tok string) (string, error) {
	if len(tok) > 6 && strings.ToLower(tok[0:7]) == "bearer " {
		return tok[7:], nil
	}
	return tok, nil
}

func RequireLoggedIn(base controllers.BaseController) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := stripBearer(c.Request.Header.Get("Authorization"))
		if err != nil {
			base.Error(c, http.StatusUnauthorized, errors.New("Unauthorized access"), err.Error())
			return
		}

		claims := &authservice.UserClaim{}

		tokenClaims, err := jwt.ParseWithClaims(
			token,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("JWT_SECRET_KEY")), nil
			},
		)
		if err != nil {
			base.Error(c, http.StatusUnauthorized, errors.New("Unauthorized access"), err.Error())
			return
		}

		if tokenClaims == nil {
			c.Abort()
			return
		}

		if !tokenClaims.Valid {
			base.Error(c, http.StatusUnauthorized, errors.New("Invalid token error"), err.Error())
			return
		}

		if claims.Id == 0 {
			base.Error(c, http.StatusBadRequest, errors.New("Invalid user"), "user id = 0")
			return
		}

		if claims.UserName == "" {
			base.Error(c, http.StatusBadRequest, errors.New("Invalid username"), "empty user name")
			return
		}

		c.Set("user_id", claims.Id)
		c.Set("user_name", claims.UserName)

		c.Next()
	}
}
