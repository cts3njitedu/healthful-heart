package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/validators"
	"github.com/cts3njitedu/healthful-heart/mappers"
	"github.com/cts3njitedu/healthful-heart/enrichers"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"log"
)

type WorkflowService struct {
	pageValidator validators.IPageValidator
	pageRepository mongorepo.IPageRepository
	mapperUtil mappers.IMapper
	enricherExecutor enrichers.IEnricherExecutor
	credEnricher enrichers.ICredentialEnricher
}

type Workflow struct {
	page models.Page
	credentials models.Credentials
	user models.User
	blankPage models.Page
}


func NewWorkflowService(
	pageValidator validators.IPageValidator,
	pageRepository mongorepo.IPageRepository,
	mapperUtil mappers.IMapper,
	enricherExecutor enrichers.IEnricherExecutor,
	credEnricher enrichers.ICredentialEnricher) *WorkflowService {
	return &WorkflowService{pageValidator, 
		pageRepository, mapperUtil, enricherExecutor, credEnricher}
}

func (work *WorkflowService) ExecuteWorkflow(page models.Page) (Workflow, error) {

	dbPage:=work.pageRepository.GetPage("LOGIN_PAGE");

	work.enricherExecutor.Enricher(&page, dbPage)
	log.Println("Finishing enriching ")
	err := work.pageValidator.ValidatePage(&page)

	if err != nil {
		return Workflow{
			page: page,
			credentials: models.Credentials{},
		}, err
	}
	
	cred:=work.mapperUtil.MapPageToCredentials(page);

	work.credEnricher.EnrichCredentials(&cred)
	
	user:=work.mapperUtil.MapCredentialsToUser(cred);

	return Workflow{
		page: page,
		credentials: cred,
		user: user,
		blankPage: dbPage,
	},nil
}