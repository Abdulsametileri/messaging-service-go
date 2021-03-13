package api

import (
	"errors"
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
		Error(c, 400, err)
		return
	}

	repo := repository.GetAuthRepository()

	userModel := vm.ToModel()

	c.Set(c.FullPath(), userModel)

	isExist, err := repo.ExistUser(userModel.UserName)

	if err == nil && isExist {
		Error(c, http.StatusBadRequest, ErrUserAlreadyExist)
		return
	}

	if err != nil && !c.GetBool("test") {
		Error(c, http.StatusBadRequest, err)
		return
	}

	if errCreateUser := repo.CreateUser(&userModel); errCreateUser != nil {
		Error(c, http.StatusBadRequest, errCreateUser)
		return
	}

	Data(c, http.StatusOK, nil, "")
}

func login(c *gin.Context) {
	var vm viewmodels.AuthVm

	if err := c.ShouldBindJSON(&vm); err != nil {
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
