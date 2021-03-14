package controllers

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/services/authservice"
	"github.com/lib/pq"
)

type authSvc struct{}

func (a authSvc) CreateJwtToken(user *models.User) (string, error) {
	return "mock-token", nil
}

func (a authSvc) ParseJwtToken(token string) (*authservice.UserClaim, error) {
	return nil, nil
}

// ----------
type logSvc struct{}

func (ls *logSvc) CreateLog(log *models.Log) error {
	return nil
}

// ----------

type userSvc struct{}

func (us *userSvc) GetUserByID(id int) (*models.User, error) {
	if id == -1 {
		return nil, errors.New("non existed user")
	}

	return &models.User{
		BaseModel: models.BaseModel{
			ID: id,
		},
		UserName:     "abdulsamet",
		Password:     "123456",
		MutedUserIDs: pq.Int32Array{1000},
	}, nil
}

func (us *userSvc) SaveUser(user *models.User) error {
	return nil
}

func (us *userSvc) CreateUser(user *models.User) error {
	return nil
}

func (us *userSvc) ExistUser(userName string) (bool, error) {
	if userName == "abdulsamet" {
		return true, nil
	}
	return false, nil
}

func (us *userSvc) GetUser(userName, password string) (*models.User, error) {
	if userName == "abdulsamet" {
		return &models.User{
			BaseModel: models.BaseModel{
				ID: 100,
			},
			UserName: userName,
			Password: helpers.Sha256String(password),
		}, nil
	}
	return &models.User{}, nil
}
