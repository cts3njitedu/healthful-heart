package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/validators"
	"github.com/cts3njitedu/healthful-heart/mappers"
	"github.com/cts3njitedu/healthful-heart/mergers"
	"github.com/cts3njitedu/healthful-heart/enrichers"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
)

type SignupService struct {
	pageValidator validators.IPageValidator
	pageRepository mongorepo.IPageRepository
	mapperUtil mappers.IMapper
	pageMerger mergers.IPageMerger
	userRepository mysqlrepo.IUserRepository
	credEnricher enrichers.ICredentialEnricher
}

func NewSignupService(
	pageValidator validators.IPageValidator,
	pageRepository mongorepo.IPageRepository,
	mapperUtil mappers.IMapper,
	pageMerger mergers.IPageMerger,
	userRepository mysqlrepo.IUserRepository, credEnricher enrichers.ICredentialEnricher) *SignupService {
	return &SignupService{pageValidator, 
		pageRepository, mapperUtil, pageMerger, userRepository, credEnricher}
}

func (signup * SignupService) SignupService(page models.Page) (models.Page, models.Credentials, error) {
	
	dbPage:=signup.pageRepository.GetPage("LOGIN_PAGE");

	signup.pageMerger.MergeRequestPageToPage(&page, dbPage)

	err := signup.pageValidator.ValidatePage(&page)

	if err != nil {
		return page, models.Credentials{}, err
	}
	
	cred:=signup.mapperUtil.MapPageToCredentials(page);

	user:=signup.mapperUtil.MapCredentialsToUser(cred);

	err = signup.userRepository.CreateUser(&user)
	

	return models.Page{}, cred, nil
}