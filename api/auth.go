package api

import (
	"errors"
	"fmt"
	"github.com/Abdulsametileri/messaging-service/config"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
		c.Set(
			c.FullPath(),
			"Invalid body in register method"+err.Error(),
		)
		Error(c, 400, err)
		return
	}

	userModel := vm.ToModel()

	c.Set(c.FullPath(), userModel)

	authRepo := repository.GetAuthRepository()

	isExist, err := authRepo.ExistUser(userModel.UserName)

	if err == nil && isExist {
		c.Set(c.FullPath(), "User already exist error occured")
		Error(c, http.StatusBadRequest, ErrUserAlreadyExist)
		return
	}

	if err != nil && !c.GetBool("test") {
		c.Set(c.FullPath(), "Error occured in the database "+err.Error())
		Error(c, http.StatusBadRequest, err)
		return
	}

	if errCreateUser := authRepo.CreateUser(&userModel); errCreateUser != nil {
		c.Set(c.FullPath(), "Error creating user in the database "+errCreateUser.Error())
		Error(c, http.StatusBadRequest, errCreateUser)
		return
	}

	c.Set(c.FullPath(), fmt.Sprintf("A User has been registered successfully. ID=%d", userModel.ID))
	Data(c, http.StatusOK, nil, "")
}

func login(c *gin.Context) {
	var vm viewmodels.AuthVm

	if err := c.ShouldBindJSON(&vm); err != nil {
		c.Set(
			c.FullPath(),
			"Invalid body in register method"+err.Error(),
		)
		Error(c, 400, err)
		return
	}

	userModel := vm.ToModel()

	c.Set(c.FullPath(), userModel)

	authRepo := repository.GetAuthRepository()

	user, err := authRepo.GetUser(&userModel)

	if err != nil && !c.GetBool("test") {
		Error(c, http.StatusBadRequest, err)
		return
	}

	if user.ID == 0 {
		Error(c, http.StatusBadRequest, ErrUserDoesNotExist)
		return
	}

	claims := createUserClaim(user)

	t, err := createJwtToken(claims)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	Data(c, http.StatusOK, gin.H{"token": t}, "")
}

func createJwtToken(claims *viewmodels.UserClaim) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(config.JwtSecretKey)
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
