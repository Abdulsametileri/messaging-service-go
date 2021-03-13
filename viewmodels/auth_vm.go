package viewmodels

import (
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"strings"
)

type AuthVm struct {
	UserName string `json:"user_name" binding:"required,min=3"`
	Password string `json:"password" binding:"required"`
}

func (vm *AuthVm) formatUserName() {
	vm.UserName = strings.ToLower(strings.TrimSpace(vm.UserName))
}

func (vm *AuthVm) formatPassword() {
	vm.Password = strings.ToLower(strings.TrimSpace(vm.Password))
}

func (vm *AuthVm) hashPassword() {
	vm.Password = helpers.Sha256String(vm.Password)
}

func (vm *AuthVm) ToModel() models.User {
	vm.formatUserName()
	vm.formatPassword()
	vm.hashPassword()

	return models.User{UserName: vm.UserName, Password: vm.Password}
}
