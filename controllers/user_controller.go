package controllers

import (
	"errors"
	"fmt"
	"github.com/Abdulsametileri/messaging-service/api"
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/services/authservice"
	"github.com/Abdulsametileri/messaging-service/services/userservice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	ErrUserAlreadyExist = errors.New("username is already taken")
	ErrUserDoesNotExist = errors.New("user does not exist")
)

type UserInput struct {
	UserName string `json:"user_name" binding:"required,min=3"`
	Password string `json:"password" binding:"required"`
}

type UserController interface {
	Register(*gin.Context)
	Login(*gin.Context)
	MutateUser(*gin.Context)
}

type userController struct {
	base BaseController
	auth authservice.AuthService
	us   userservice.UserService
}

func NewUserController(bctl BaseController, auth authservice.AuthService, us userservice.UserService) UserController {
	return &userController{base: bctl, auth: auth, us: us}
}

func (ctl *userController) Register(c *gin.Context) {
	var vm UserInput

	if err := c.ShouldBindJSON(&vm); err != nil {
		ctl.base.Error(c, 400, err, "Invalid body in register method"+err.Error())
		return
	}

	vm.UserName = helpers.LowerTrimString(vm.UserName)
	vm.Password = helpers.LowerTrimString(vm.Password)
	vm.Password = helpers.HashPassword(vm.Password)

	isExist, err := ctl.us.ExistUser(vm.UserName)

	if err == nil && isExist {
		ctl.base.Error(c, http.StatusBadRequest, ErrUserAlreadyExist, "User already exist error occured")
		return
	}

	if err != nil {
		ctl.base.Error(c, http.StatusBadRequest, err, "Error occured in the database "+err.Error())
		return
	}

	if errCreateUser := ctl.us.CreateUser(&models.User{
		UserName: vm.UserName, Password: vm.Password,
	}); errCreateUser != nil {
		ctl.base.Error(c, http.StatusBadRequest, errCreateUser, "Error creating user in the database "+errCreateUser.Error())
		return
	}

	ctl.base.Data(c, http.StatusCreated, nil, "", fmt.Sprintf("A User has been registered successfully. UserName=%s", vm.UserName))
}

func (ctl *userController) Login(c *gin.Context) {
	var vm UserInput

	if err := c.ShouldBindJSON(&vm); err != nil {
		api.Error(c, 400, err, "Invalid body"+err.Error())
		return
	}

	vm.UserName = helpers.LowerTrimString(vm.UserName)
	vm.Password = helpers.LowerTrimString(vm.Password)
	vm.Password = helpers.HashPassword(vm.Password)

	user, err := ctl.us.GetUser(vm.UserName, vm.Password)

	if err != nil {
		ctl.base.Error(c, http.StatusBadRequest, err, "Error occured in the database. "+err.Error())
		return
	}

	if user.ID == 0 {
		ctl.base.Error(c, http.StatusBadRequest, ErrUserDoesNotExist, "User does not exist error occured")
		return
	}

	t, err := ctl.auth.CreateJwtToken(user)
	if err != nil {
		ctl.base.Error(c, http.StatusBadRequest, err, "create jwt token error")
		return
	}

	ctl.base.Data(c, http.StatusOK, gin.H{"token": t}, "", fmt.Sprintf("token generated for the user id=%d", user.ID))
}

func (ctl *userController) MutateUser(c *gin.Context) {
	userClaims := getUserClaims(c)

	user, err := ctl.us.GetUserByID(userClaims.Id)
	if err != nil {
		ctl.base.Error(c, http.StatusBadRequest, err, fmt.Sprintf("userClaimsId=%d does not found in the database", userClaims.Id))
		return
	}

	mutateThisUserId, err := strconv.Atoi(c.Query("mutateUserId"))
	if err != nil {
		logMsg := fmt.Sprintf("User id=%d tried to pass api via invalid param %s",
			userClaims.Id, c.Query("mutateUserId"),
		)
		ctl.base.Error(c, http.StatusBadRequest, errors.New("invalid param"), logMsg)
		return
	}

	if mutateThisUserId == userClaims.Id {
		ctl.base.Error(c, http.StatusBadRequest, errors.New("You cannot mutate yourself, sorry."),
			fmt.Sprintf("User id=%d tried to mutate yourself", userClaims.Id))
		return
	}

	mutatedUser, err := ctl.us.GetUserByID(mutateThisUserId)
	if err != nil {
		logMsg := fmt.Sprintf("User Id=%d and UserName=%s tried to mutate not existed user. mutateUserId=%d",
			user.ID, user.UserName, mutateThisUserId)
		ctl.base.Error(c, http.StatusBadRequest, errors.New("Want to mutated user does not exist."), logMsg)
		return
	}

	found := helpers.SearchNumberInPgArray(mutateThisUserId, user.MutedUserIDs)
	if found {
		logMsg := fmt.Sprintf("User Id=%d and UserName=%s tried to mutate user before he/she mutated. mutateUserId=%d",
			user.ID, user.UserName, mutateThisUserId)
		ctl.base.Error(c, http.StatusBadRequest, errors.New("This user is already mutated."), logMsg)
		return
	}

	unSortedMutatedUserIds := append(user.MutedUserIDs, int32(mutateThisUserId))
	user.MutedUserIDs = helpers.SortPgArrayAscending(unSortedMutatedUserIds)

	err = ctl.us.SaveUser(user)
	if err != nil {
		ctl.base.Error(c, http.StatusBadRequest, err, "Error saving the user"+err.Error())
		return
	}

	logMsg := fmt.Sprintf("%s has been mutated successfully from %s", mutatedUser.UserName, user.UserName)
	ctl.base.Data(c, http.StatusOK, nil, "", logMsg)
}
