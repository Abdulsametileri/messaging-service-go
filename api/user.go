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
		c.Set(
			c.FullPath(),
			fmt.Sprintf("userClaimsId=%d does not found in the database", userClaims.Id),
		)
		Error(c, http.StatusBadRequest, err)
		return
	}

	mutateThisUserId, err := strconv.Atoi(c.Query("mutateUserId"))
	if err != nil {
		c.Set(
			c.FullPath(),
			fmt.Sprintf("User id=%d tried to pass api via invalid param %s",
				userClaims.Id, c.Query("mutateUserId"),
			),
		)
		Error(c, http.StatusBadRequest, errors.New("invalid param"))
		return
	}

	if mutateThisUserId == userClaims.Id {
		c.Set(
			c.FullPath(),
			fmt.Sprintf("User id=%d tried to mutate yourself", userClaims.Id),
		)
		Error(c, http.StatusBadRequest, errors.New("You cannot mutate yourself, sorry."))
		return
	}

	mutatedUser, err := userRepo.GetUser(mutateThisUserId)
	if err != nil {
		c.Set(
			c.FullPath(),
			fmt.Sprintf("User Id=%d and UserName=%s tried to mutate not existed user. mutateUserId=%d",
				user.ID, user.UserName, mutateThisUserId),
		)
		Error(c, http.StatusBadRequest, errors.New("Want to mutated user does not exist."))
		return
	}

	found := helpers.SearchNumberInPgArray(mutateThisUserId, user.MutedUserIDs)
	if found {
		c.Set(c.FullPath(),
			fmt.Sprintf("User Id=%d and UserName=%s tried to mutate user before he/she mutated. mutateUserId=%d",
				user.ID, user.UserName, mutateThisUserId),
		)
		Error(c, http.StatusBadRequest, errors.New("This user is already mutated."))
		return
	}

	unSortedMutatedUserIds := append(user.MutedUserIDs, int32(mutateThisUserId))
	user.MutedUserIDs = helpers.SortPgArrayAscending(unSortedMutatedUserIds)

	err = userRepo.SaveUser(&user)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	c.Set(c.FullPath(), fmt.Sprintf("%s has been mutated successfully from %s", mutatedUser.UserName, user.UserName))

	Data(c, http.StatusOK, nil, "")
}

func users(c *gin.Context) {
	userClaims := getUserClaims(c)

	userRepo := repository.GetUserRepository()

	user, err := userRepo.GetUser(userClaims.Id)
	if err != nil {
		c.Set(
			c.FullPath(),
			fmt.Sprintf("User Claims Id=%d does not found in the database.", userClaims.Id),
		)
		Error(c, http.StatusBadRequest, err)
		return
	}

	users, err := userRepo.GetUserList(user.ID, user.MutedUserIDs)
	if err != nil {
		c.Set(
			c.FullPath(),
			fmt.Sprintf("Error occured within GetUserList for the user id = %d.", user.ID),
		)
		Error(c, http.StatusBadRequest, err)
		return
	}

	ret := make([]gin.H, len(users))
	for i, user := range users {
		ret[i] = gin.H{
			"id":        user.ID,
			"user_name": user.UserName,
		}
	}

	c.Set(
		c.FullPath(),
		fmt.Sprintf("User Id=%d and UserName=%s listed users.", user.ID, user.UserName),
	)

	Data(c, 200, ret, "")
}
