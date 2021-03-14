package userservice

import (
	"errors"
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/Abdulsametileri/messaging-service/repository/userrepo"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(user *models.User) error
	ExistUser(userName string) (bool, error)
	GetUser(userName, password string) (*models.User, error)
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
