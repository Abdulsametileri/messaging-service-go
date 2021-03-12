package api

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserAlreadyExist = errors.New("This username is already taken.")
)

func register(c *gin.Context) {
	var vm viewmodels.RegisterVm

	if err := c.ShouldBindJSON(&vm); err != nil {
		Error(c, 400, err)
		return
	}

	repo := repository.Get()

	userModel := vm.ToModel()

	isExist, errUserExist := repo.IsUserExist(&userModel)

	if isExist {
		Error(c, 400, ErrUserAlreadyExist)
		return
	}

	if errUserExist != nil {
		Error(c, 400, errUserExist)
		return
	}

	if errCreateUser := repo.CreateUser(&userModel); errCreateUser != nil {
		Error(c, 400, errCreateUser)
		return
	}

	Data(c, 200, nil, "")
}
