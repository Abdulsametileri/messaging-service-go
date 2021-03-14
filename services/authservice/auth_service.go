package authservice

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type UserClaim struct {
	Id       int    `json:"id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

type AuthService interface {
	CreateJwtToken(user *models.User) (string, error)
	ParseJwtToken(token string) (*UserClaim, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (as *authService) CreateJwtToken(user *models.User) (string, error) {
	claims := &UserClaim{
		Id:       user.ID,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365).Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(viper.GetString("JWT_SECRET_KEY")))
}

func (as *authService) ParseJwtToken(token string) (*UserClaim, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("JWT_SECRET_KEY")), nil
		},
	)

	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*UserClaim)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
