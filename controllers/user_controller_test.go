package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestUserController(t *testing.T) {
	us := &userSvc{}
	as := &authSvc{}
	bctl := &baseController{ls: &logSvc{}}

	userCtl := NewUserController(bctl, as, us)
	gin.SetMode(gin.TestMode)

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
		t.Run("Success", func(t *testing.T) {
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

			fmt.Println(resBody.Data)

			res, _ := resBody.Data.(map[string]interface{})

			assert.Equal(t, res["token"], "mock-token")
		})
		t.Run("Failed to login (not existed user)", func(t *testing.T) {
			reqBody := map[string]string{
				"user_name": "huhu",
				"password":  "123456",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(payload))
			c.Request = request

			userCtl.Login(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})

	t.Run("MutateUser", func(t *testing.T) {
		t.Run("Getting invalid user id error", func(t *testing.T) {
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)

			userCtl.MutateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("Getting error for non exist user id in the database", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", -1)

			userCtl.MutateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("Getting error for url parameter error mutateUserId", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 3)

			userCtl.MutateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("Getting error when user try to mutate yourself ", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 3)
			c.Params = []gin.Param{
				{
					Key:   "mutateUserId",
					Value: "3",
				},
			}

			userCtl.MutateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("Getting error when mutated user does not exist", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 3)
			c.Params = []gin.Param{
				{
					Key:   "mutateUserId",
					Value: "-1",
				},
			}

			userCtl.MutateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("Getting error user try to mutate same user again", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 1)
			c.Params = []gin.Param{
				{
					Key:   "mutateUserId",
					Value: "1000",
				},
			}

			userCtl.MutateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
		t.Run("user will mutate another user successfully?", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", 1)
			c.Params = []gin.Param{
				{
					Key:   "mutateUserId",
					Value: "3",
				},
			}

			userCtl.MutateUser(c)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}