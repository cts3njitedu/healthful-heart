package services

import (
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"log"
)

type AuthenticationService struct {
	pageRepository mongorepo.IPageRepository
	restructPage IRestructurePageService
}

func NewAuthenticationService(pageRepository mongorepo.IPageRepository, 
		restructPage IRestructurePageService) *AuthenticationService {
			return &AuthenticationService {pageRepository, restructPage}
}


func (auth *AuthenticationService) GetAuthenticationPage(authenticationType string) models.Page {
	page:=auth.pageRepository.GetPage("LOGIN_PAGE");

	if authenticationType == "LOGIN" {
		return page
	} else if authenticationType == "SIGNUP" {
		log.Println("Retrieving sign up page...")
		auth.restructPage.RestructureLoginPage(&page)
		return page
	}
	return models.Page{}
}