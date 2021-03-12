package viewmodels

import (
	"github.com/Abdulsametileri/messaging-service/helpers"
	"github.com/Abdulsametileri/messaging-service/models"
	"strings"
)

type RegisterVm struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (vm *RegisterVm) formatUserName() {
	vm.UserName = strings.ToLower(strings.TrimSpace(vm.UserName))
}

func (vm *RegisterVm) formatPassword() {
	vm.Password = strings.ToLower(strings.TrimSpace(vm.Password))
}

func (vm *RegisterVm) hashPassword() {
	vm.Password = helpers.Sha256String(vm.Password)
}

func (vm *RegisterVm) ToModel() models.User {
	vm.formatUserName()
	vm.formatPassword()
	vm.hashPassword()

	return models.User{UserName: vm.UserName, Password: vm.Password}
}
