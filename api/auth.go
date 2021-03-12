package api

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserAlreadyExist = errors.New("username is already taken")
)

func register(c *gin.Context) {
	var vm viewmodels.RegisterVm

	if err := c.ShouldBindJSON(&vm); err != nil {
		Error(c, 400, err)
		return
	}

	repo := repository.GetAuthRepository()

	userModel := vm.ToModel()

	user, err := repo.UserExist(&userModel)

	if err == nil && user.ID > 0 {
		Error(c, 400, ErrUserAlreadyExist)
		return
	}

	if err != nil {
		Error(c, 400, err)
		return
	}

	if errCreateUser := repo.CreateUser(&userModel); errCreateUser != nil {
		Error(c, 400, errCreateUser)
		return
	}

	Data(c, 200, nil, "")
}
