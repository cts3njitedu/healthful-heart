package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"mime/multipart"
)

type IRestructurePageService interface {
	RestructureLoginPage(page *models.Page)
}

type IAuthenticationService interface {
	GetAuthenticationPage(authenticationType string) models.Page
}

type ISignupService interface {
	SignupService(page models.Page) (models.Page, models.Credentials, error)
}

type IWorkflowService interface {
	ExecuteWorkflow(page models.Page) (Workflow, error)
}

type ILoginService interface {
	LoginService(page models.Page) (models.Page, models.Credentials, error)
}

type IFileService interface {
	UploadFile(file multipart.File, fileHeader * multipart.FileHeader,metaData models.WorkoutFile, cred models.Credentials) error
}

type IRabbitService interface {
	PushFileMetaDataToQueue(file *models.WorkoutFile) error
	// PullFileMetaDataFromQueue()
}

type IFileProcessorService interface {
	ProcessWorkoutFile(file models.WorkoutFile) (error)
}

type IGroupParserService interface {
	GetGroups(workoutType string, groupText string, categoryCode string) ([]models.Group)
}

type IWorkoutTypeService interface {
	GetWorkoutTypeCode(categoryCd string, workoutTypeName string) string
	GetCategoryNameFromCode(categoryCd string) (string, error)
	GetCategoryCodeFromName(categoryName string) (string, error)

}

type ILocationService interface {
	GetLocation(locationId int64) (models.Location, error)
	GetLocations() ([]models.Location, error)
}

type IWorkoutService interface {
	GetWorkoutDays(queryParams models.QueryParams, cred models.Credentials) ([]models.WorkoutDay, error)
	GetWorkoutDaysPage(queryParams models.QueryParams, cred models.Credentials) (models.HeartResponse, error)
}