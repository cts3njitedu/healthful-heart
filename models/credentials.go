package models

type Credentials struct {
	UserId string `json:"userId"`
	Username string `json:"username"`
	Password string `json:"password"`
	PasswordText string `json:"passwordText"`
	ConfirmPassword string `json:"confirmPassword"`
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}