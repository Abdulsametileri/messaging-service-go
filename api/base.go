package api

import (
	"encoding/json"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

var logRepo = repository.GetLogRepository()

type Props struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Data(c *gin.Context, code int, data interface{}, message string) {
	props := &Props{
		Code:    code,
		Data:    data,
		Message: message,
	}

	requestBody, _ := c.Get(c.FullPath())
	requestBodyJson, _ := json.Marshal(&requestBody)
	responseBodyJson, _ := json.Marshal(&props)

	_ = logRepo.CreateLog(&models.Log{
		ApiPath:  c.FullPath(),
		Request:  requestBodyJson,
		Response: responseBodyJson,
		Type:     models.LogInfo,
	})

	c.JSON(code, props)
}

func Error(c *gin.Context, code int, err error) {
	props := &Props{
		Code:    code,
		Data:    nil,
		Message: err.Error(),
	}

	responseJson, _ := json.Marshal(&props)

	request, isExist := c.Get(c.FullPath())
	var requestJson datatypes.JSON

	if isExist {
		requestJson, _ = json.Marshal(&request)
	} else {
		requestJson, _ = json.Marshal(err.Error())
	}

	_ = logRepo.CreateLog(&models.Log{
		ApiPath:  c.FullPath(),
		Request:  requestJson,
		Response: responseJson,
		Type:     models.LogError,
	})

	c.AbortWithStatusJSON(code, props)
}
