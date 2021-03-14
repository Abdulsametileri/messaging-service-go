package userrepo

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

type Repo interface {
	CreateUser(user *models.User) error
	ExistUser(userName string) (models.User, error)
	GetUser(userName, password string) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) Repo {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) CreateUser(user *models.User) error {
	return ur.db.Model(&models.User{}).Create(&user).Error
}

func (ur *userRepo) ExistUser(userName string) (models.User, error) {
	var user models.User
	err := ur.db.Model(&models.User{}).First(&user, "user_name = ?", userName).Error

	return user, err
}

func (ur *userRepo) GetUser(userName, password string) (*models.User, error) {
	var user *models.User
	err := ur.db.Model(&models.User{}).
		First(&user, "user_name = ? AND password = ?", userName, password).
		Error

	return user, err
}
