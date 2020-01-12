package mappers

import (
	"github.com/cts3njitedu/healthful-heart/models"

)


type IMapper interface {
	MapPageToCredentials(page models.Page) models.Credentials
	MapCredentialsToUser(cred models.Credentials) models.User
}


type Mapper struct {}
func NewMapper() *Mapper {
	return &Mapper{}
}

func (mapper *Mapper) MapPageToCredentials(page models.Page) models.Credentials {
	var cred models.Credentials
	for _,section := range page.Sections {
		for _, field := range section.Fields {
			switch field.Name {
			case "username":
				cred.Username = field.Value
			case "password":
				cred.Password = field.Value
			case "confirmPassword":
				cred.ConfirmPassword = field.Value
			case "email":
				cred.Email = field.Value
			case "firstname":
				cred.FirstName = field.Value
			case "lastname": 
				cred.LastName = field.Value
			}
		}
	}
	return cred;	

}

func (mapper *Mapper) MapCredentialsToUser(cred models.Credentials) models.User {
	user:=models.User {
		Username: cred.Username,
		Password: cred.Password,
		FirstName: cred.FirstName,
		LastName: cred.LastName,
		Email: cred.Email,
	}
	return user
}