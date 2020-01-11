package security

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {}

type IPasswordHasher interface {
    HashPassword(password string) (string, error)
}

func NewPasswordHasher() *PasswordHasher {
    return &PasswordHasher{}
}

func (hash *PasswordHasher) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}