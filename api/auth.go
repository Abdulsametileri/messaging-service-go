package api

import (
	"errors"
	"fmt"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var (
	ErrUserAlreadyExist = errors.New("username is already taken")
	ErrUserDoesNotExist = errors.New("user does not exist")
)

func register(c *gin.Context) {
	var vm viewmodels.AuthVm

	if err := c.ShouldBindJSON(&vm); err != nil {
		Error(c, 400, err, "Invalid body in register method"+err.Error())
		return
	}

	userModel := vm.ToModel()

	authRepo := repository.GetAuthRepository()

	isExist, err := authRepo.ExistUser(userModel.UserName)

	if err == nil && isExist {
		Error(c, http.StatusBadRequest, ErrUserAlreadyExist, "User already exist error occured")
		return
	}

	if err != nil && !c.GetBool("test") {
		Error(c, http.StatusBadRequest, err, "Error occured in the database "+err.Error())
		return
	}

	if errCreateUser := authRepo.CreateUser(&userModel); errCreateUser != nil {
		Error(c, http.StatusBadRequest, errCreateUser, "Error creating user in the database "+errCreateUser.Error())
		return
	}

	Data(c, http.StatusCreated, nil, "", fmt.Sprintf("A User has been registered successfully. ID=%d", userModel.ID))
}

func login(c *gin.Context) {
	var vm viewmodels.AuthVm

	if err := c.ShouldBindJSON(&vm); err != nil {
		Error(c, 400, err, "Invalid body"+err.Error())
		return
	}

	userModel := vm.ToModel()

	authRepo := repository.GetAuthRepository()

	user, err := authRepo.GetUser(userModel.UserName, userModel.Password)

	if err != nil && !c.GetBool("test") {
		Error(c, http.StatusBadRequest, err, "Error occured in the database. "+err.Error())
		return
	}

	if user.ID == 0 {
		Error(c, http.StatusBadRequest, ErrUserDoesNotExist, "User does not exist error occured")
		return
	}

	claims := createUserClaim(user)

	t, err := createJwtToken(claims)
	if err != nil {
		Error(c, http.StatusBadRequest, err, "create jwt token error")
		return
	}

	Data(c, http.StatusOK, gin.H{"token": t}, "", fmt.Sprintf("token generated for the user id=%d", user.ID))
}

func createJwtToken(claims *viewmodels.UserClaim) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString([]byte(viper.GetString("JWT_SECRET_KEY")))
	return tokenString, err
}

func createUserClaim(user *models.User) *viewmodels.UserClaim {
	return &viewmodels.UserClaim{
		Id:       user.ID,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365).Unix(),
		},
	}
}
