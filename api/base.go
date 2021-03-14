package api

import (
	"github.com/Abdulsametileri/messaging-service/repository"
	"github.com/Abdulsametileri/messaging-service/viewmodels"
	"github.com/gin-gonic/gin"
)

var logRepo = repository.GetLogRepository()

type Props struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Data(c *gin.Context, code int, data interface{}, message string, requestDetailForLogPurpose interface{}) {
	props := &Props{
		Code:    code,
		Data:    data,
		Message: message,
	}

	/*requestDetailJson, _ := json.Marshal(&requestDetailForLogPurpose)
	responseBodyJson, _ := json.Marshal(&props)

	_ = logRepo.CreateLog(&models.Log{
		ApiPath:  c.FullPath(),
		Request:  requestDetailJson,
		Response: responseBodyJson,
		Type:     models.LogInfo,
	})*/

	c.JSON(code, props)
}

func Error(c *gin.Context, code int, friendlyErrorForClient error, requestDetailForLogPurpose interface{}) {
	props := &Props{
		Code:    code,
		Data:    nil,
		Message: friendlyErrorForClient.Error(),
	}

	/*requestDetailJson, _ := json.Marshal(&requestDetailForLogPurpose)
	responseJson, _ := json.Marshal(&props)

	_ = logRepo.CreateLog(&models.Log{
		ApiPath:  c.FullPath(),
		Request:  requestDetailJson,
		Response: responseJson,
		Type:     models.LogError,
	})*/

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
