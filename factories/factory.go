package factories

import (
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/enrichers"
	"github.com/cts3njitedu/healthful-heart/mappers"
	"github.com/cts3njitedu/healthful-heart/mergers"
	"github.com/cts3njitedu/healthful-heart/security"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"github.com/cts3njitedu/healthful-heart/repositories/mongorepo"
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/handlers"
	"github.com/cts3njitedu/healthful-heart/validators"
	"github.com/cts3njitedu/healthful-heart/utils"

)

var (
	
	mongoConnection *connections.MongoConnection
	pageRepository *mongorepo.PageRepository
	restructureService *services.RestructurePageService
	authenticationService *services.AuthenticationService
	fieldValidator *validators.FieldValidator
	pageValidator *validators.PageValidator
	mapperUtil *mappers.Mapper
	pageMerger *mergers.PageMerger
	mysqlConnection *connections.MysqlConnection
	userRepository *mysqlrepo.UserRepository
	hasher *security.PasswordHasher
	credentialEnricher *enrichers.CredentialEnricher
	workflowService *services.WorkflowService
	signupService *services.SignupService
	singupEnricher *enrichers.SignupEnrich
	enricherExecutor *enrichers.EnrichExecutor
	loginService *services.LoginService
	environmentUtiliy *utils.EnvironmentUtility

)

func init() {
	environmentUtiliy = utils.NewEnvironmentUtility()
	mongoConnection = connections.NewMongoConnection(environmentUtiliy)
	pageRepository = mongorepo.NewPageRepository(mongoConnection, environmentUtiliy)
	restructureService = services.NewRestructurePageService()
	authenticationService = services.NewAuthenticationService(pageRepository, restructureService)
	fieldValidator = validators.NewFieldValidator()
	pageValidator = validators.NewPageValidator(fieldValidator)
	mapperUtil = mappers.NewMapper()
	pageMerger = mergers.NewPageMerger()
	mysqlConnection = connections.NewMysqlConnection()
	userRepository = mysqlrepo.NewUserRepository(mysqlConnection)
	hasher = security.NewPasswordHasher()
	credentialEnricher = enrichers.NewCredentialEnricher(hasher)
	singupEnricher = enrichers.NewSignupEnrich();
	enr:= []enrichers.IEnricher {singupEnricher}
	enricherExecutor = enrichers.NewEnrichExecutor(enr)
	workflowService = services.NewWorkflowService(pageValidator, pageRepository, mapperUtil,enricherExecutor , credentialEnricher)
	signupService = services.NewSignupService(workflowService,userRepository)
	loginService = services.NewLoginService(workflowService, userRepository, mapperUtil, hasher)
}





func GetLoginHandler() *handlers.LoginHandler {
	
	return handlers.NewLoginHandler(authenticationService, loginService)
}

func GetSignupHandler() *handlers.SignupHandler {
	return handlers.NewSingupHandler(authenticationService, signupService)
}

func GetTokenHandler() *handlers.TokenHandler {
	return handlers.NewTokenHandler(environmentUtiliy)
}