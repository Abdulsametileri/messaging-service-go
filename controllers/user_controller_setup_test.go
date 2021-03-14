package controllers

import (
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
)

type authSvc struct{}

func (as *authSvc) CreateJwtToken(user *models.User) (string, error) {
	return "nice-token", nil
}

type logSvc struct{}

func (ls *logSvc) CreateLog(log *models.Log) error {
	return nil
}

type userSvc struct{}

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
	return &models.User{
		BaseModel: models.BaseModel{
			ID: 100,
		},
		UserName: userName,
		Password: helpers.Sha256String(password),
	}, nil
}
