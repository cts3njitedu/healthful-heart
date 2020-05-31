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
	GetWorkoutType(categoryCd string, workoutTypeName string) (models.WorkoutType, error)
	GetCategoryNameFromCode(categoryCd string) (string, error)
	GetCategoryCodeFromName(categoryName string) (string, error)
	GetCategories() (map[string]string, map[string]string)
	GetCategoriesAndWorkoutsMap(catCode string) (map[string]map[int64]models.WorkoutType)
	GetSortedCategoriesAndWorkoutTypes() ([]models.SortedCategoryWorkoutType)
	GetWorkoutTypeByIds(ids []int64) (map[int64]models.WorkoutType)

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

type IGymRepositoryService interface {
	LoadWorkoutDayOriginal(workoutDaysCurrent []models.WorkoutDay) []models.WorkoutDay
	GetWorkoutDaysByParams(options models.QueryOptions) ([]models.WorkoutDay, error)
	SaveWorkoutDayLocation(workDay *models.WorkoutDay) (*models.WorkoutDay, error)
	UpdateAllWorkoutDay(workDays []models.WorkoutDay, ids map[string][]string) error
	GetWorkoutByParams(queryOptions models.QueryOptions) ([]models.Workout, error)
	GetGroupByParams(queryOptions models.QueryOptions) ([]models.Group, error)
	GetLocationsQueryParams(queryOptions models.QueryOptions) ([]models.Location, error)
	GetWorkoutDaysLocationByParams(queryOptions models.QueryOptions) ([]models.WorkoutDay, error)
}

type IEventService interface {
	FindWorkoutDaysAdded(currs [] models.WorkoutDay, origins []models.WorkoutDay, eventDetails *[]models.ModEventDetail)
	FindDeletedIds(origins [] models.WorkoutDay, deletedIds map[string][]string, eventDetails *[]models.ModEventDetail)
	FindWorkoutDaysDifferences(currs [] models.WorkoutDay, origins []models.WorkoutDay, eventDetails *[]models.ModEventDetail) []models.WorkoutDay
	FindWorkoutsDifferences(currs [] models.Workout, origins []models.Workout, eventDetails *[]models.ModEventDetail) []models.Workout
	FindGroupDifferences(currs [] models.Group, origins []models.Group, eventDetails *[]models.ModEventDetail) []models.Group
}