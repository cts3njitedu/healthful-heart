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
	signupService *services.SignupService

)

func init() {
	mongoConnection = connections.NewMongoConnection()
	pageRepository = mongorepo.NewPageRepository(mongoConnection)
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
	signupService = services.NewSignupService(pageValidator, pageRepository, mapperUtil, pageMerger, userRepository, credentialEnricher)
}





func GetLoginHandler() *handlers.LoginHandler {
	
	return handlers.NewLoginHandler(authenticationService)
}

func GetSignupHandler() *handlers.SignupHandler {
	return handlers.NewSingupHandler(authenticationService, signupService)
}