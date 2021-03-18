package controllers

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/services/redisservice"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	WsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WebSocketController interface {
	ServeWs(ctx *gin.Context)
}

type webSocketController struct {
	base         BaseController
	redisService redisservice.RedisService
}

func NewWebSocketController(bctl BaseController, redisService redisservice.RedisService) WebSocketController {
	return &webSocketController{
		base:         bctl,
		redisService: redisService,
	}
}

func (w *webSocketController) ServeWs(c *gin.Context) {
	chatId := c.Param("chatId")
	_ = chatId

	conn, err := WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		w.base.Error(c, http.StatusBadRequest, errors.New("Failed to create ws connection"), err.Error())
		return
	}
	defer conn.Close()

	for {
		mt, _, err := conn.ReadMessage()
		if mt == -1 || err != nil {
			break
		}
	}
}
