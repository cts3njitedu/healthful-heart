package services

import (
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/enrichers"
	"log"
)

type AuthenticationService struct {
	pageRepository mongorepo.IPageRepository
	restructPage IRestructurePageService
	enricherExecutor enrichers.IEnricherExecutor
}

func NewAuthenticationService(pageRepository mongorepo.IPageRepository, 
		restructPage IRestructurePageService, enricherExecutor enrichers.IEnricherExecutor) *AuthenticationService {
			return &AuthenticationService {pageRepository, restructPage, enricherExecutor}
}


func (auth *AuthenticationService) GetAuthenticationPage(authenticationType string) models.Page {
	page:=auth.pageRepository.GetPage("LOGIN_PAGE");
	page.PageId = authenticationType
	log.Println("Enriching sign up/login page...")
	auth.enricherExecutor.Enricher(&page, page)
	// if authenticationType == "LOGIN" {
	// 	return page
	// } else if authenticationType == "SIGNUP" {
	// 	log.Println("Retrieving sign up page...")
	// 	auth.restructPage.RestructureLoginPage(&page)
	// 	return page
	// }
	return page
}