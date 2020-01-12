package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/models"
)

type IUserRepository interface {
	GetUser(user models.User) models.User
	CreateUser(user *models.User) error

}