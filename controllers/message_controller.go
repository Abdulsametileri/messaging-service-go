package controllers

import (
	"errors"
	"fmt"
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/services/messageservice"
	"github.com/Abdulsametileri/messaging-service/services/redisservice"
	"github.com/Abdulsametileri/messaging-service/services/userservice"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrReceiverUserDoesNotExist = errors.New("Receiver user does not exist.")
	ErrUserTryToMuteYourself    = errors.New("You cannot type message yourself, sorry.")
)

type MessageInput struct {
	Text string `json:"text" binding:"required"`
}

type MessageController interface {
	GetMessages(*gin.Context)
	SendMessage(*gin.Context)
}

type messageController struct {
	base     BaseController
	us       userservice.UserService
	msgs     messageservice.MessageService
	redissrv redisservice.RedisService
}

func NewMessageController(bctl BaseController, us userservice.UserService, msgs messageservice.MessageService,
	redisService redisservice.RedisService) MessageController {
	return &messageController{
		base:     bctl,
		us:       us,
		msgs:     msgs,
		redissrv: redisService,
	}
}

func (mc *messageController) GetMessages(c *gin.Context) {
	userId := c.GetInt("user_id")

	_, err := mc.us.GetUserByID(userId)
	if err != nil {
		mc.base.Error(c, http.StatusBadRequest, err, fmt.Sprintf("userClaimsId=%d does not found in the database", userId))
		return
	}

	userName := helpers.LowerTrimString(c.Param("userName"))
	receiverUser, err := mc.us.GetUserByUserName(userName)
	if err != nil {
		mc.base.Error(c, http.StatusBadRequest,
			ErrReceiverUserDoesNotExist,
			fmt.Sprintf("User id=%d get messages between not existed user!! receiverUsername=%s", userId, c.Param("userName")),
		)
		return
	}

	messages, err := mc.msgs.GetMessages(userId, receiverUser.ID)
	if err != nil {
		logMsg := fmt.Sprintf("Failed to get all messages between sender-receiver=[%d-%d]", userId, receiverUser.ID)
		mc.base.Error(c, http.StatusBadRequest, err, logMsg)
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

	logMsg := fmt.Sprintf("senderId-receiverId %d-%d got all messages", userId, receiverUser.ID)
	mc.base.Data(c, http.StatusOK, ret, "", logMsg)
}

func (mc *messageController) SendMessage(c *gin.Context) {
	var vm MessageInput

	if err := c.ShouldBindJSON(&vm); err != nil {
		mc.base.Error(c, http.StatusBadRequest, err, "payload error"+err.Error())
		return
	}

	userId := c.GetInt("user_id")

	requesterUser, err := mc.us.GetUserByID(userId)
	if err != nil {
		mc.base.Error(c, http.StatusBadRequest, err, fmt.Sprintf("userClaimsId=%d does not found in the database", userId))
		return
	}

	userName := helpers.LowerTrimString(c.Param("userName"))
	receiverUser, err := mc.us.GetUserByUserName(userName)
	if err != nil {
		mc.base.Error(c, http.StatusBadRequest,
			ErrReceiverUserDoesNotExist,
			fmt.Sprintf("User id=%d send message not existed user!! receiverUsername=%s", userId, userName),
		)
		return
	}

	if receiverUser.ID == userId {
		mc.base.Error(c, http.StatusBadRequest,
			ErrUserTryToMuteYourself,
			fmt.Sprintf("User id=%d try to message yourself!", userId))
		return
	}

	if found := helpers.SearchNumberInPgArray(userId, receiverUser.MutedUserIDs); found {
		errMsg := fmt.Sprintf("You cannot message %s Because he/she mutated you!", receiverUser.UserName)
		logMsg := fmt.Sprintf("UserId=%d tried to send message a another user which block the user before receiverId=%d",
			userId, receiverUser.ID)
		mc.base.Error(c, http.StatusBadRequest, errors.New(errMsg), logMsg)
		return
	}

	if found := helpers.SearchNumberInPgArray(receiverUser.ID, requesterUser.MutedUserIDs); found {
		errMsg := fmt.Sprintf("You cannot message %s because you blocked before", receiverUser.UserName)
		logMsg := fmt.Sprintf("UserId=%d tried to send message a another user which blocked before receiverId=%d",
			userId, receiverUser.ID)
		mc.base.Error(c, http.StatusBadRequest, errors.New(errMsg), logMsg)
		return
	}

	err = mc.msgs.CreateMessage(&models.Message{
		SenderId:   userId,
		ReceiverId: receiverUser.ID,
		Text:       vm.Text,
	})

	if err != nil {
		mc.base.Error(c, http.StatusBadRequest,
			errors.New("Error creating message, please try again."),
			"Error create message "+err.Error())
		return
	}

	logMsg := fmt.Sprintf("%d -> %d send message %s", userId, receiverUser.ID, vm.Text)
	mc.base.Data(c, http.StatusOK, nil, "", logMsg)

	mc.redissrv.PublishMessage("1-2", vm.Text)
}
