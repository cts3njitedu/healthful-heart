package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/mappers"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/security"
)

type LoginService struct {
	workflowService IWorkflowService
	userRepository mysqlrepo.IUserRepository
	mapperUtil mappers.IMapper
	hasher security.IPasswordHasher

}

func NewLoginService(workflowService IWorkflowService,
	userRepository mysqlrepo.IUserRepository, mapperUtil mappers.IMapper,hasher security.IPasswordHasher) *LoginService {
	return &LoginService{workflowService, userRepository, mapperUtil, hasher}
}


func (login *LoginService) LoginService(page models.Page) (models.Page, models.Credentials, error) {
	page.PageId = "LOGIN"
	
	work,err:=login.workflowService.ExecuteWorkflow(page);

	if err != nil {
		return work.page,models.Credentials{},err
	}

	cred:=work.credentials
	user:=work.user
	hashedUser :=login.userRepository.GetUser(user)
	// fmt.Printf("User password: %+v, Credentials: %+v\n", hashedUser, cred)
	err = login.hasher.CompareHashWithPassword(hashedUser.Password, cred.PasswordText)
	if err != nil {
		return models.Page{}, models.Credentials{}, err
	}
	cred = login.mapperUtil.MapUserToCredentials(user)

	return models.Page{}, cred, nil

}
