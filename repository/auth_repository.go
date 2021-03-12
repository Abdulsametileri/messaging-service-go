package repository

/*
import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

var authRepo *AuthRepository

type AuthRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) {
	authRepo = &AuthRepository{db: db}
}

func Get() *AuthRepository {
	return authRepo
}

func (p *AuthRepository) CreateUser(user *models.User) error {
	err := p.db.Model(&models.User{}).Create(&user).Error
	return err
}

func (p *AuthRepository) IsUserExist(user *models.User) (bool, error) {
	err := p.db.Model(&models.User{}).
		First(&user, "user_name = ? AND password = ?", user.UserName, user.Password).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}*/
