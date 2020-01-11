package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/models"
)

type IUserRepository interface {
	GetUser(username string) models.User
	CreateUser(user *models.User) error

}