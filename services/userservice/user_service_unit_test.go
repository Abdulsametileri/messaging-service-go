package userservice

import (
	"github.com/Abdulsametileri/messaging-service/models"
	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

func (repo repoMock) CreateUser(user *models.User) error {
	args := repo.Called(user)
	return args.Error(0)
}

func (repo repoMock) ExistUser(userName string) (models.User, error) {
	args := repo.Called(userName)
	return models.User{}, args.Error(1)
}

func (repo repoMock) GetUser(userName, password string) (*models.User, error) {
	args := repo.Called(userName, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func (repo repoMock) GetUserByID(id int) (*models.User, error) {
	panic("implement me")
}

func (repo repoMock) GetUserByUserName(userName string) (user models.User, err error) {
	panic("implement me")
}

func (repo repoMock) SaveUser(user *models.User) error {
	panic("implement me")
}

func (repo repoMock) GetUserList(userId int, mutatedUserIdsCond string) (users []models.User, err error) {
	args := repo.Called(userId, mutatedUserIdsCond)
	return args.Get(0).([]models.User), args.Error(1)
}
