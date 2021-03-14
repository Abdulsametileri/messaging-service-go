package api

import (
	"errors"
	"fmt"
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

func messages(c *gin.Context) {
	userRepo := repository.GetUserRepository()

	userClaims := getUserClaims(c)

	_, err := userRepo.GetUser(userClaims.Id)
	if err != nil {
		Error(c, http.StatusBadRequest, err, fmt.Sprintf("userClaimsId=%d does not found in the database", userClaims.Id))
		return
	}

	userName := helpers.LowerTrimString(c.Query("user_name"))
	receiverUser, err := userRepo.GetUserByUsername(userName)
	if err != nil {
		Error(c, http.StatusBadRequest,
			errors.New("Receiver user does not exist."),
			fmt.Sprintf("User id=%d send message not existed user!! receiverUsername=%s", userClaims.Id, c.Query("user_name")),
		)
		return
	}

	messageRepo := repository.GetMessageRepository()
	messages, err := messageRepo.GetMessages(userClaims.Id, receiverUser.ID)
	if err != nil {
		logMsg := fmt.Sprintf("Failed to get all messages between sender-receiver=[%d-%d]", userClaims.Id, receiverUser.ID)
		Error(c, http.StatusBadRequest, err, logMsg)
		return
	}

	ret := make([]gin.H, len(messages))
	for i, message := range messages {
		ret[i] = gin.H{
			"date":     message.CreatedAt.Format("03 January 2006 - 15:04"),
			"sender":   message.Sender.UserName,
			"receiver": message.Receiver.UserName,
			"text":     message.Text,
		}
	}

	logMsg := fmt.Sprintf("senderId-receiverId %d-%d got all messages", userClaims.Id, receiverUser.ID)
	Data(c, http.StatusOK, ret, "", logMsg)
}

func sendMessage(c *gin.Context) {
	var vm viewmodels.MessageVm

	if err := c.ShouldBindJSON(&vm); err != nil {
		Error(c, http.StatusBadRequest, err, "payload error"+err.Error())
		return
	}

	userRepo := repository.GetUserRepository()

	userClaims := getUserClaims(c)

	_, err := userRepo.GetUser(userClaims.Id)
	if err != nil {
		Error(c, http.StatusBadRequest, err, fmt.Sprintf("userClaimsId=%d does not found in the database", userClaims.Id))
		return
	}

	userName := helpers.LowerTrimString(c.Query("user_name"))
	receiverUser, err := userRepo.GetUserByUsername(userName)
	if err != nil {
		Error(c, http.StatusBadRequest,
			errors.New("Receiver user does not exist."),
			fmt.Sprintf("User id=%d send message not existed user!! receiverUsername=%s", userClaims.Id, userName),
		)
		return
	}

	if receiverUser.ID == userClaims.Id {
		Error(c, http.StatusBadRequest,
			errors.New("You cannot type message yourself, sorry."),
			fmt.Sprintf("User id=%d try to message yourself!", userClaims.Id))
		return
	}

	if found := helpers.SearchNumberInPgArray(receiverUser.ID, receiverUser.MutedUserIDs); found {
		errMsg := fmt.Sprintf("You cannot message %s Because he/she mutated you!", receiverUser.UserName)
		logMsg := fmt.Sprintf("UserId=%d tried to send message a another user which block the user before receiverId=%d",
			userClaims.Id, receiverUser.ID)
		Error(c, http.StatusBadRequest, errors.New(errMsg), logMsg)
	}

	messageRepo := repository.GetMessageRepository()
	err = messageRepo.CreateMessage(&models.Message{
		SenderId:   userClaims.Id,
		ReceiverId: receiverUser.ID,
		Text:       vm.Text,
	})

	if err != nil {
		Error(c, http.StatusBadRequest,
			errors.New("Error creating message, please try again."),
			"Error create message "+err.Error())
		return
	}

	logMsg := fmt.Sprintf("%d -> %d send message %s", userClaims.Id, receiverUser.ID, vm.Text)
	Data(c, http.StatusOK, nil, "", logMsg)
}
