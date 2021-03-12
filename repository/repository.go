package repository

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

var Repo *Repository

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) {
	Repo = &Repository{db: db}
}

func Get() *Repository {
	return Repo
}

func (p *Repository) CreateUser(user *models.User) error {
	err := p.db.Model(&models.User{}).Create(&user).Error
	return err
}

func (p *Repository) IsUserExist(user *models.User) (bool, error) {
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
}
