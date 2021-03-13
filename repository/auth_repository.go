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

func (p *AuthRepository) ExistUser(userName string) (bool, error) {
	var user models.User
	err := p.db.Model(&models.User{}).First(&user, "user_name = ?", userName).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *AuthRepository) GetUser(user *models.User) (*models.User, error) {
	err := p.db.Model(&models.User{}).
		First(&user, "user_name = ? AND password = ?", user.UserName, user.Password).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, nil
	}

	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}
