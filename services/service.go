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
	GetCategories() (map[string]string, map[string]string)
	GetCategoriesAndWorkouts() (map[string]map[string]models.WorkoutType)
	GetCategoriesAndWorkoutTypes() (map[string]map[string]string)

}

type ILocationService interface {
	GetLocation(locationId int64) (models.Location, error)
	GetLocations() (map[int64]models.Location, error)
}

type IWorkoutService interface {
	GetWorkoutDays(queryParams models.QueryParams, cred models.Credentials) ([]models.WorkoutDay, error)
	GetWorkoutDaysPage(queryParams models.QueryParams, cred models.Credentials) (models.HeartResponse, error)
	GetWorkoutDaysLocationsView(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error)
	AddWorkoutDateLocation(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error)
	GetWorkoutPageHeader(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error)
	GetWorkouts(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error)
	GetWorkoutDetails(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error)
	GetWorkoutDetailsMetaInfo(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error)
	ActionWorkoutDay(heartRequest models.HeartRequest, cred models.Credentials) (models.HeartResponse, error)
}