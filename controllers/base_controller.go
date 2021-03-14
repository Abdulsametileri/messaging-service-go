package controllers

import (
	"encoding/json"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/services/authservice"
	"github.com/Abdulsametileri/messaging-service/services/logservice"
	"github.com/gin-gonic/gin"
)

type Props struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type BaseController interface {
	Data(c *gin.Context, code int, data interface{}, message string, requestDetailForLogPurpose interface{})
	Error(c *gin.Context, code int, friendlyErrorForClient error, requestDetailForLogPurpose interface{})
}

type baseController struct {
	ls logservice.LogService
}

func NewBaseController(ls logservice.LogService) BaseController {
	return &baseController{ls: ls}
}

func (bc *baseController) Data(c *gin.Context, code int, data interface{}, message string, requestDetailForLogPurpose interface{}) {
	props := &Props{
		Code:    code,
		Data:    data,
		Message: message,
	}

	requestDetailJson, _ := json.Marshal(&requestDetailForLogPurpose)
	responseBodyJson, _ := json.Marshal(&props)

	_ = bc.ls.CreateLog(&models.Log{
		ApiPath:  c.FullPath(),
		Request:  requestDetailJson,
		Response: responseBodyJson,
		Type:     models.LogInfo,
	})

	c.JSON(code, props)
}

func (bc *baseController) Error(c *gin.Context, code int, friendlyErrorForClient error, requestDetailForLogPurpose interface{}) {
	props := &Props{
		Code:    code,
		Data:    nil,
		Message: friendlyErrorForClient.Error(),
	}

	requestDetailJson, _ := json.Marshal(&requestDetailForLogPurpose)
	responseJson, _ := json.Marshal(&props)

	_ = bc.ls.CreateLog(&models.Log{
		ApiPath:  c.FullPath(),
		Request:  requestDetailJson,
		Response: responseJson,
		Type:     models.LogError,
	})

	c.AbortWithStatusJSON(code, props)
}

func getUserClaims(c *gin.Context) (claims authservice.UserClaim) {
	cl, ok := c.Get("claims")
	if !ok {
		return claims
	}

	userClaims, ok := cl.(authservice.UserClaim)
	if !ok {
		return claims
	}

	return userClaims
}
