package api

import (
	"encoding/json"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

var logRepo = repository.GetLogRepository()

type Props struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Data(c *gin.Context, code int, data interface{}, message string, requestDetail interface{}) {
	props := &Props{
		Code:    code,
		Data:    data,
		Message: message,
	}

	requestDetailJson, _ := json.Marshal(&requestDetail)
	responseBodyJson, _ := json.Marshal(&props)

	_ = logRepo.CreateLog(&models.Log{
		ApiPath:  c.FullPath(),
		Request:  requestDetailJson,
		Response: responseBodyJson,
		Type:     models.LogInfo,
	})

	c.JSON(code, props)
}

func Error(c *gin.Context, code int, err error, requestDetail interface{}) {
	props := &Props{
		Code:    code,
		Data:    nil,
		Message: err.Error(),
	}

	requestDetailJson, _ := json.Marshal(&requestDetail)
	responseJson, _ := json.Marshal(&props)

	_ = logRepo.CreateLog(&models.Log{
		ApiPath:  c.FullPath(),
		Request:  requestDetailJson,
		Response: responseJson,
		Type:     models.LogError,
	})

	if code >= http.StatusInternalServerError {
		props.Message = "Error occured in our own server. Sorry"
	}

	c.AbortWithStatusJSON(code, props)
}

func getUserClaims(c *gin.Context) (claims viewmodels.UserClaim) {
	cl, ok := c.Get("claims")
	if !ok {
		return claims
	}

	userClaims, ok := cl.(viewmodels.UserClaim)
	if !ok {
		return claims
	}

	return userClaims
}
