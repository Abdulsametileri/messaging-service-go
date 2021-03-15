package userrepo

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"gorm.io/gorm"
)

type Repo interface {
	CreateUser(user *models.User) error
	ExistUser(userName string) (models.User, error)
	GetUser(userName, password string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUserName(userName string) (user models.User, err error)
	SaveUser(user *models.User) error
	GetUserList(userId int, mutatedUserIdsCond string) (users []models.User, err error)
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

func (ur *userRepo) GetUserByID(id int) (*models.User, error) {
	var user *models.User
	err := ur.db.Model(&models.User{}).First(&user, id).Error
	return user, err
}

func (ur *userRepo) GetUserByUserName(userName string) (user models.User, err error) {
	err = ur.db.Model(&models.User{}).First(&user, "user_name = ?", userName).Error
	return
}

func (ur *userRepo) SaveUser(user *models.User) error {
	return ur.db.Save(&user).Error
}

func (ur *userRepo) GetUserList(userId int, mutatedUserIdsCond string) (users []models.User, err error) {
	err = ur.db.Model(&models.User{}).
		Where(mutatedUserIdsCond).
		Where("id <> ?", userId).
		Find(&users).Error
	return
}
