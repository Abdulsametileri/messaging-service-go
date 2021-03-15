package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessageController(t *testing.T) {
	msgs := &messageSvc{}
	us := &userSvc{}
	bctl := &baseController{ls: &logSvc{}}

	messageCtrl := NewMessageController(bctl, us, msgs)
	gin.SetMode(gin.TestMode)

	t.Run("getMessages", func(t *testing.T) {
		t.Run("getting error specified user_id claim does not exist", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", -1)

			messageCtrl.GetMessages(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("getting error when specified username does not exist", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 1)

			c.Params = []gin.Param{
				{
					Key:   "userName",
					Value: "notexisteduser",
				},
			}

			messageCtrl.GetMessages(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			resBody := Props{}
			json.NewDecoder(w.Body).Decode(&resBody)

			assert.EqualValues(t, ErrReceiverUserDoesNotExist.Error(), resBody.Message)
		})
		t.Run("user getting messages successfully", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 1)

			c.Params = []gin.Param{
				{
					Key:   "userName",
					Value: "abc",
				},
			}

			messageCtrl.GetMessages(c)

			assert.Equal(t, http.StatusOK, w.Code)

			resBody := Props{}
			json.NewDecoder(w.Body).Decode(&resBody)

			val, _ := resBody.Data.([]interface{})
			lenItems := len(val)

			assert.EqualValues(t, lenItems, 2)
		})
	})

	t.Run("sendMessage", func(t *testing.T) {
		t.Run("getting error when message text body is empty", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			messageCtrl.SendMessage(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("getting error specified user_id claim does not exist", func(t *testing.T) {
			reqBody := MessageInput{Text: "hello?"}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{
				{
					Key:   "userName",
					Value: "abc",
				},
			}
			c.Set("user_id", -1)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/sendMessage", bytes.NewBuffer(payload))
			c.Request = request

			messageCtrl.SendMessage(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("getting error receiver user does not exist", func(t *testing.T) {
			reqBody := MessageInput{Text: "hello?"}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{
				{
					Key:   "userName",
					Value: "xxx",
				},
			}
			c.Set("user_id", 1)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/sendMessage", bytes.NewBuffer(payload))
			c.Request = request

			messageCtrl.SendMessage(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			resBody := Props{}
			json.NewDecoder(w.Body).Decode(&resBody)

			assert.EqualValues(t, ErrReceiverUserDoesNotExist.Error(), resBody.Message)
		})
		t.Run("getting error user tries to message yourself", func(t *testing.T) {
			reqBody := MessageInput{Text: "hello?"}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{
				{
					Key:   "userName",
					Value: "ddd",
				},
			}
			c.Set("user_id", 1)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/sendMessage", bytes.NewBuffer(payload))
			c.Request = request

			messageCtrl.SendMessage(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			resBody := Props{}
			json.NewDecoder(w.Body).Decode(&resBody)

			assert.EqualValues(t, ErrUserTryToMuteYourself.Error(), resBody.Message)
		})
		t.Run("getting error when User tries to send message to a user which mutates requester user before", func(t *testing.T) {
			reqBody := MessageInput{Text: "hello?"}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 1)
			c.Params = []gin.Param{
				{
					Key:   "userName",
					Value: "mutateruser",
				},
			}

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/sendMessage", bytes.NewBuffer(payload))
			c.Request = request

			messageCtrl.SendMessage(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("gettinge error when user tries to send own blocked user", func(t *testing.T) {
			reqBody := MessageInput{Text: "hello?"}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", mutaterUser.ID)
			c.Params = []gin.Param{
				{
					Key:   "userName",
					Value: ddduser.UserName,
				},
			}

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/sendMessage", bytes.NewBuffer(payload))
			c.Request = request

			messageCtrl.SendMessage(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("user send message successfully", func(t *testing.T) {
			reqBody := MessageInput{Text: "hello?"}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 1)
			c.Params = []gin.Param{
				{
					Key:   "userName",
					Value: "abc",
				},
			}

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/sendMessage", bytes.NewBuffer(payload))
			c.Request = request

			messageCtrl.SendMessage(c)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
