package userservice

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestExistUser(t *testing.T) {
	t.Run("Not Exist user", func(t *testing.T) {
		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("ExistUser", "notexisted").
			Return(mock.Anything, gorm.ErrRecordNotFound)

		found, err := us.ExistUser("notexisted")

		assert.False(t, found)
		assert.Nil(t, err)
	})
	t.Run("Exist user false when arbitrary error", func(t *testing.T) {
		retErr := errors.New("db error")
		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("ExistUser", "abdulsamet").Return(mock.Anything, retErr)

		found, err := us.ExistUser("abdulsamet")
		assert.False(t, found)
		assert.ErrorIs(t, err, retErr)
	})
	t.Run("Exist user", func(t *testing.T) {
		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("ExistUser", "abdulsamet").Return(mock.Anything, nil)

		found, err := us.ExistUser("abdulsamet")
		assert.True(t, found)
		assert.Nil(t, err)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Create user", func(t *testing.T) {
		usr := &models.User{
			UserName: "abdulsamet",
			Password: "123456",
		}

		userRepo := new(repoMock)

		us := NewUserService(userRepo)

		userRepo.On("CreateUser", usr).Return(nil)

		result := us.CreateUser(usr)

		assert.Nil(t, result)
	})

	t.Run("Create user fails", func(t *testing.T) {
		err := errors.New("error")
		usr := &models.User{
			UserName: "abdulsamet",
		}

		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("CreateUser", usr).Return(err)
		result := us.CreateUser(usr)

		assert.EqualValues(t, result, err)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Get user does not exist", func(t *testing.T) {
		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("GetUser", "abdulsamet", "123456").
			Return(&models.User{}, gorm.ErrRecordNotFound)

		usr, err := us.GetUser("abdulsamet", "123456")

		assert.EqualValues(t, usr, &models.User{})
		assert.Nil(t, err)
	})
	t.Run("Get user empty when arbitrary error", func(t *testing.T) {
		retErr := errors.New("arbitrary error")
		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("GetUser", "abdulsamet", "123456").
			Return(&models.User{}, retErr)

		usr, err := us.GetUser("abdulsamet", "123456")

		assert.EqualValues(t, usr, &models.User{})
		assert.ErrorIs(t, err, retErr)
	})
	t.Run("Get user successfully", func(t *testing.T) {
		expected := models.User{
			UserName: "abdulsamet",
			Password: "123456",
		}

		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("GetUser", expected.UserName, expected.Password).
			Return(&expected, nil)

		usr, err := us.GetUser(expected.UserName, expected.Password)

		assert.EqualValues(t, usr, &expected)
		assert.Nil(t, err)
	})
}

func TestGetUserList(t *testing.T) {
	t.Run("get user list when mutated user ids nil", func(t *testing.T) {
		userId := 1
		retUserList := make([]models.User, 0)
		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("GetUserList", 1, "").Return(retUserList, nil)

		_, err := us.GetUserList(userId, nil)

		assert.Nil(t, err)
	})
	t.Run("get user list with values of mutated user ids", func(t *testing.T) {
		userId := 1
		retUserList := make([]models.User, 0)
		userRepo := new(repoMock)
		us := NewUserService(userRepo)
		userRepo.On("GetUserList", 1, "NOT id = ANY('{2,3}')").Return(retUserList, nil)

		_, err := us.GetUserList(userId, pq.Int32Array{2, 3})

		assert.Nil(t, err)
	})
}
