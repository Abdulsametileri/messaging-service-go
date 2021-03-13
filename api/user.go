package api

import (
	"errors"
	"fmt"
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func mutateUser(c *gin.Context) {
	userRepo := repository.GetUserRepository()

	userClaims := getUserClaims(c)

	user, err := userRepo.GetUser(userClaims.Id)
	if err != nil {
		Error(c, http.StatusBadRequest, err, fmt.Sprintf("userClaimsId=%d does not found in the database", userClaims.Id))
		return
	}

	mutateThisUserId, err := strconv.Atoi(c.Query("mutateUserId"))
	if err != nil {
		logMsg := fmt.Sprintf("User id=%d tried to pass api via invalid param %s",
			userClaims.Id, c.Query("mutateUserId"),
		)
		Error(c, http.StatusBadRequest, errors.New("invalid param"), logMsg)
		return
	}

	if mutateThisUserId == userClaims.Id {
		Error(c, http.StatusBadRequest, errors.New("You cannot mutate yourself, sorry."),
			fmt.Sprintf("User id=%d tried to mutate yourself", userClaims.Id))
		return
	}

	mutatedUser, err := userRepo.GetUser(mutateThisUserId)
	if err != nil {
		logMsg := fmt.Sprintf("User Id=%d and UserName=%s tried to mutate not existed user. mutateUserId=%d",
			user.ID, user.UserName, mutateThisUserId)
		Error(c, http.StatusBadRequest, errors.New("Want to mutated user does not exist."), logMsg)
		return
	}

	found := helpers.SearchNumberInPgArray(mutateThisUserId, user.MutedUserIDs)
	if found {
		logMsg := fmt.Sprintf("User Id=%d and UserName=%s tried to mutate user before he/she mutated. mutateUserId=%d",
			user.ID, user.UserName, mutateThisUserId)
		Error(c, http.StatusBadRequest, errors.New("This user is already mutated."), logMsg)
		return
	}

	unSortedMutatedUserIds := append(user.MutedUserIDs, int32(mutateThisUserId))
	user.MutedUserIDs = helpers.SortPgArrayAscending(unSortedMutatedUserIds)

	err = userRepo.SaveUser(&user)
	if err != nil {
		Error(c, http.StatusBadRequest, err, "Error saving the user"+err.Error())
		return
	}

	logMsg := fmt.Sprintf("%s has been mutated successfully from %s", mutatedUser.UserName, user.UserName)
	Data(c, http.StatusOK, nil, "", logMsg)
}

func users(c *gin.Context) {
	userClaims := getUserClaims(c)

	userRepo := repository.GetUserRepository()

	user, err := userRepo.GetUser(userClaims.Id)
	if err != nil {
		Error(c, http.StatusBadRequest, err, fmt.Sprintf("User Claims Id=%d does not found in the database.", userClaims.Id))
		return
	}

	users, err := userRepo.GetUserList(user.ID, user.MutedUserIDs)
	if err != nil {
		Error(c, http.StatusBadRequest, err, fmt.Sprintf("Error occured within GetUserList for the user id = %d.", user.ID))
		return
	}

	ret := make([]gin.H, len(users))
	for i, user := range users {
		ret[i] = gin.H{
			"id":        user.ID,
			"user_name": user.UserName,
		}
	}

	Data(c, 200, ret, "", fmt.Sprintf("User Id=%d and UserName=%s listed users.", user.ID, user.UserName))
}
