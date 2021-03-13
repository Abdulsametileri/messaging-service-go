package repository

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"strings"
)

var userRepo *UserRepository

type UserRepository struct {
	Repository
}

func (u UserRepository) Setup(db *gorm.DB) {
	userRepo = &UserRepository{Repository{db: db}}
}

func GetUserRepository() *UserRepository {
	return userRepo
}

func (u *UserRepository) SaveUser(user *models.User) error {
	err := u.db.Save(&user).Error
	return err
}

func (u *UserRepository) GetUser(userId int) (user models.User, err error) {
	err = u.db.Model(&models.User{}).First(&user, userId).Error
	return
}

func (u *UserRepository) GetUserByUsername(userName string) (user models.User, err error) {
	err = u.db.Model(&models.User{}).First(&user, "user_name = ?", userName).Error
	return
}

func (u *UserRepository) GetUserList(userId int, mutatedUserIds pq.Int32Array) (users []models.User, err error) {
	v, _ := mutatedUserIds.Value()
	mutatedUserCondition := "NOT id = ANY('{values}')"
	mutatedUserCondition = strings.Replace(mutatedUserCondition, "{values}", v.(string), 1)

	err = u.db.Model(&models.User{}).
		Where(mutatedUserCondition).
		Where("id <> ?", userId).
		Find(&users).Error
	return
}
