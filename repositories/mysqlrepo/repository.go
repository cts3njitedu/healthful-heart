package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"time"
)

type IUserRepository interface {
	GetUser(user models.User) (models.User, error)
	CreateUser(user *models.User) error

}

type ITokenRepository interface {
	SaveRefreshToken(token string, expirationTime time.Time, userId string) error 
}