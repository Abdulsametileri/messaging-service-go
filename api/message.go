package api

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func sendMessage(c *gin.Context) {
	var vm viewmodels.MessageVm

	if err := c.ShouldBindJSON(&vm); err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	userRepo := repository.GetUserRepository()

	userClaims := getUserClaims(c)

	_, err := userRepo.GetUser(userClaims.Id)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	receiverUserId, err := strconv.Atoi(c.Query("receiverUserId"))
	if err != nil {
		Error(c, http.StatusBadRequest, errors.New("invalid param"))
		return
	}

	if receiverUserId == userClaims.Id {
		Error(c, http.StatusBadRequest, errors.New("You cannot type message yourself, sorry."))
		return
	}

	_, err = userRepo.GetUser(receiverUserId)
	if err != nil {
		Error(c, http.StatusBadRequest, errors.New("Receiver user does not exist."))
		return
	}

	messageRepo := repository.GetMessageRepository()
	err = messageRepo.CreateMessage(&models.Message{
		SenderId:   userClaims.Id,
		ReceiverId: receiverUserId,
		Text:       vm.Text,
	})

	if err != nil {
		Error(c, http.StatusBadRequest, errors.New("Error creating message, please try again."))
		return
	}
}
