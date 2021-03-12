package repository

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

var authRepo *AuthRepository

type AuthRepository struct {
	Repository
}

func (a AuthRepository) Setup(db *gorm.DB) {
	authRepo = &AuthRepository{Repository{db: db}}
}

func GetAuthRepository() *AuthRepository {
	return authRepo
}

func (p *AuthRepository) CreateUser(user *models.User) error {
	err := p.db.Model(&models.User{}).Create(&user).Error
	return err
}

func (p *AuthRepository) UserExist(user *models.User) (*models.User, error) {
	err := p.db.Model(&models.User{}).
		First(&user, "user_name = ? AND password = ?", user.UserName, user.Password).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
