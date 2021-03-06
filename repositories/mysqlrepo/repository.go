package mysqlrepo

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"time"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	GetUser(user models.User) (models.User, error)
	CreateUser(user *models.User) error
	GetUserToken(userId string) (models.User, error)

}

type ITokenRepository interface {
	SaveRefreshToken(token string, expirationTime time.Time, userId string) error
}

type IFileRepository interface {
	SaveFile(file *models.WorkoutFile) (error)
	UpdateFileStatus(file *models.WorkoutFile, newStatus string) (models.WorkoutFile, error)
}

type IWorkoutTypeRepository interface {
	GetWorkoutTypes(queryOptions models.QueryOptions) ([]models.WorkoutType, error)
}

type ICategoryRepository interface {
	GetCategories() ([]models.Category, error) 
	GetCategoriesByParams(queryOptions models.QueryOptions) ([]models.Category, error)
	GetCategoriesAndWorkoutTypes(queryOptions models.QueryOptions) ([]models.Category, error)
}

type IWorkoutDayRepository interface {
	SaveWorkoutDay(workDay *models.WorkoutDay, tx *gorm.DB) error
	GetWorkoutDays(userId string) ([]models.WorkoutDay, error)
	GetWorkoutDaysByParams(queryOptions models.QueryOptions) ([]models.WorkoutDay, error)
	GetWorkoutDaysSpecifyColumns(queryOptions models.QueryOptions) ([]models.WorkoutDay, error)
	SaveWorkoutDayLocation(workDay *models.WorkoutDay) (*models.WorkoutDay, error)
	DeleteWorkoutDays(ids map[string][]string, tx *gorm.DB) bool
	UpdateAllWorkoutDay(workDays []models.WorkoutDay, ids map[string][]string) error
	GetWorkoutDaysLocationByParams(queryOptions models.QueryOptions) ([]models.WorkoutDay, error)

}

type IWorkoutRepository interface {
	SaveWorkout(workDay *models.Workout, tx *gorm.DB) error
	GetWorkoutByParams(queryOptions models.QueryOptions) ([]models.Workout, error)
	DeleteWorkouts(ids map[string][]string, tx *gorm.DB) bool
}

type IGroupRepository interface {
	SaveGroup(group *models.Group, tx *gorm.DB) error
	GetGroupByParams(queryOptions models.QueryOptions) ([]models.Group, error)
	DeleteGroups(ids map[string][]string, tx *gorm.DB) bool
}

type ILocationRepository interface {
	GetLocations() ([]models.Location, error)
	GetLocationsQueryParams(queryOptions models.QueryOptions) ([]models.Location, error)
	GetByLocationIds(ids []string) ([]models.Location, error)
}