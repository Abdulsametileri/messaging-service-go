package viewmodels

import "github.com/dgrijalva/jwt-go"

type UserClaim struct {
	Id       int    `json:"id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}
