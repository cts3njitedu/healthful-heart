package security

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {}

type IPasswordHasher interface {
    HashPassword(password string) (string, error)
    CompareHashWithPassword(hash, password string ) error
}

type PasswordError struct {
    s string

}


func (p * PasswordError) Error() string {
    return p.s
}
func NewPasswordHasher() *PasswordHasher {
    return &PasswordHasher{}
}

func (hash *PasswordHasher) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func (hash *PasswordHasher) CompareHashWithPassword(pwHash, password string ) error {
    err := bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(password))
    if err != nil {
        return &PasswordError{"Invalid Username or Password"}
    }
    return err
}