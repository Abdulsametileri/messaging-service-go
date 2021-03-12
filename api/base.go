package api

import "github.com/gin-gonic/gin"

type Props struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

func Data(c *gin.Context, code int, data interface{}, message string) {
	props := &Props{
		Code:    code,
		Data:    data,
		Message: message,
	}
	c.JSON(code, props)
}

func Error(c *gin.Context, code int, err error) {
	props := &Props{
		Code:    code,
		Data:    nil,
		Message: err.Error(),
	}
	c.AbortWithStatusJSON(code, props)
}
