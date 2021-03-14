package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserController(t *testing.T) {
	as := &authSvc{}
	us := &userSvc{}
	bctl := &baseController{ls: &logSvc{}}

	userCtl := NewUserController(bctl, as, us)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	v1 := router.Group("api/v1")
	{
		v1.POST("register", userCtl.Register)
	}

	t.Run("Register", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			reqBody := gin.H{
				"user_name": "abdulsametnew",
				"password":  "123456",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
			c.Request = request

			userCtl.Register(c)

			assert.Equal(t, http.StatusCreated, w.Code)
		})
		t.Run("Fails to register user", func(t *testing.T) {
			reqBody := gin.H{
				"user_name": "abdulsamet",
				"password":  "123456",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
			c.Request = request

			userCtl.Register(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})

	t.Run("Login", func(t *testing.T) {
		reqBody := map[string]string{
			"user_name": "abdulsamet",
			"password":  "123456",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		payload, _ := json.Marshal(reqBody)
		request := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(payload))
		c.Request = request

		userCtl.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody := Props{}
		json.NewDecoder(w.Body).Decode(&resBody)

		expectedResBody := Props{
			Code:    http.StatusOK,
			Data:    map[string]interface{}{"token": "nice-token"},
			Message: "",
		}

		assert.Equal(t, expectedResBody.Data, resBody.Data)
	})
}
