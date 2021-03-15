package userservice

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository/userrepo"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"strings"
)

type UserService interface {
	CreateUser(user *models.User) error
	ExistUser(userName string) (bool, error)
	GetUser(userName, password string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUserName(userName string) (user models.User, err error)
	SaveUser(user *models.User) error
	GetUserList(userId int, mutatedUserIds pq.Int32Array) (users []models.User, err error)
}

type userService struct {
	Repo userrepo.Repo
}

func NewUserService(repo userrepo.Repo) UserService {
	return &userService{
		Repo: repo,
	}
}

func (us *userService) CreateUser(user *models.User) error {
	return us.Repo.CreateUser(user)
}

func (us *userService) ExistUser(userName string) (bool, error) {
	_, err := us.Repo.ExistUser(userName)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (us *userService) GetUser(userName, password string) (*models.User, error) {
	user, err := us.Repo.GetUser(userName, password)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, nil
	}

	if err != nil {
		return user, err
	}

	return user, nil
}

func (us *userService) GetUserByID(id int) (*models.User, error) {
	return us.Repo.GetUserByID(id)
}

func (us *userService) GetUserByUserName(userName string) (user models.User, err error) {
	return us.Repo.GetUserByUserName(userName)
}

func (us *userService) SaveUser(user *models.User) error {
	return us.Repo.SaveUser(user)
}

func (us *userService) GetUserList(userId int, mutatedUserIds pq.Int32Array) (users []models.User, err error) {
	if mutatedUserIds == nil {
		return us.Repo.GetUserList(userId, "")
	}
	v, _ := mutatedUserIds.Value()
	mutatedUserCondition := "NOT id = ANY('{values}')"
	mutatedUserCondition = strings.Replace(mutatedUserCondition, "{values}", v.(string), 1)

	return us.Repo.GetUserList(userId, mutatedUserCondition)
}
