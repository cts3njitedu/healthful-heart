package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"log"
)

type SignupService struct {
	workflowService IWorkflowService
	userRepository mysqlrepo.IUserRepository
	
}

func NewSignupService(workflowService IWorkflowService,
	userRepository mysqlrepo.IUserRepository) *SignupService {
	return &SignupService{workflowService, userRepository}
}

func (signup * SignupService) SignupService(page models.Page) (models.Page, models.Credentials, error) {
	
	page.PageId = "SIGNUP"
	
	work,err:=signup.workflowService.ExecuteWorkflow(page);

	if err != nil {
		return work.page,models.Credentials{},err
	}
	user:=work.user
	cred:=work.credentials
	err = signup.userRepository.CreateUser(&user)
	
	if err != nil {
		log.Printf("%+v\n",err)
		return models.Page{}, models.Credentials{}, err
	}

	return models.Page{}, cred, nil
}